package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/taurusgroup/frost-ed25519/pkg/eddsa"
	"github.com/taurusgroup/frost-ed25519/pkg/frost"
	"github.com/taurusgroup/frost-ed25519/pkg/frost/keygen"
	"github.com/taurusgroup/frost-ed25519/pkg/frost/party"
	"github.com/taurusgroup/frost-ed25519/pkg/frost/sign"
	"github.com/taurusgroup/frost-ed25519/pkg/helpers"
	"github.com/taurusgroup/frost-ed25519/pkg/state"

	"mpc-wasm-chain/x/tss/types"
)

// FROSTStateManager manages FROST protocol state for validators
// State is kept in memory since validators need it across blocks
type FROSTStateManager struct {
	mu sync.RWMutex

	// DKG state per session
	dkgStates  map[string]*state.State
	dkgOutputs map[string]*keygen.Output

	// Signing state per request
	signStates  map[string]*state.State
	signOutputs map[string]*sign.Output

	// Stored key shares for signing (indexed by keySetID)
	keyShares    map[string]*eddsa.SecretShare
	publicShares map[string]*eddsa.Public
}

// Global state manager (validators maintain this across blocks)
var frostStateManager = &FROSTStateManager{
	dkgStates:    make(map[string]*state.State),
	dkgOutputs:   make(map[string]*keygen.Output),
	signStates:   make(map[string]*state.State),
	signOutputs:  make(map[string]*sign.Output),
	keyShares:    make(map[string]*eddsa.SecretShare),
	publicShares: make(map[string]*eddsa.Public),
}

// ========================
// DKG Functions
// ========================

// InitDKGState initializes FROST DKG state for this validator
func (k Keeper) InitDKGState(sessionID string, selfIndex int, participantCount int, threshold uint32) error {
	frostStateManager.mu.Lock()
	defer frostStateManager.mu.Unlock()

	// Check if already initialized
	if _, exists := frostStateManager.dkgStates[sessionID]; exists {
		return nil // Already initialized
	}

	// Create party IDs (1-indexed as FROST expects)
	partyIDs := make(party.IDSlice, participantCount)
	for i := 0; i < participantCount; i++ {
		partyIDs[i] = party.ID(i + 1)
	}

	// Our party ID (1-indexed)
	selfID := party.ID(selfIndex + 1)

	// Initialize FROST keygen state
	frostState, output, err := frost.NewKeygenState(selfID, partyIDs, party.Size(threshold), 0)
	if err != nil {
		return fmt.Errorf("failed to init DKG state: %w", err)
	}

	frostStateManager.dkgStates[sessionID] = frostState
	frostStateManager.dkgOutputs[sessionID] = output

	return nil
}

// GenerateDKGRound1Message generates real FROST DKG Round 1 data
func (k Keeper) GenerateDKGRound1Message(ctx context.Context, sessionID, validatorAddr string) ([]byte, error) {
	frostStateManager.mu.Lock()
	defer frostStateManager.mu.Unlock()

	frostState, exists := frostStateManager.dkgStates[sessionID]
	if !exists {
		return nil, fmt.Errorf("DKG state not initialized for session %s", sessionID)
	}

	// Process round 1 (no input messages for first round)
	msgs, err := helpers.PartyRoutine(nil, frostState)
	if err != nil {
		return nil, fmt.Errorf("failed to generate round 1: %w", err)
	}

	// Combine all messages into a single package
	pkg := FROSTDKGRound1Msg{
		SessionID:    sessionID,
		ValidatorAddr: validatorAddr,
		Messages:     msgs,
	}

	return json.Marshal(pkg)
}

// ProcessDKGRound1Messages processes Round 1 messages and generates Round 2
func (k Keeper) ProcessDKGRound1Messages(sessionID string, round1Messages [][]byte) ([]byte, error) {
	frostStateManager.mu.Lock()
	defer frostStateManager.mu.Unlock()

	frostState, exists := frostStateManager.dkgStates[sessionID]
	if !exists {
		return nil, fmt.Errorf("DKG state not initialized for session %s", sessionID)
	}

	// Collect all inner messages
	var allMsgs [][]byte
	for _, msgData := range round1Messages {
		var pkg FROSTDKGRound1Msg
		if err := json.Unmarshal(msgData, &pkg); err != nil {
			continue
		}
		allMsgs = append(allMsgs, pkg.Messages...)
	}

	// Process round 1 messages to generate round 2
	msgs, err := helpers.PartyRoutine(allMsgs, frostState)
	if err != nil {
		return nil, fmt.Errorf("failed to process round 1: %w", err)
	}

	// Package round 2 messages
	pkg := FROSTDKGRound2Msg{
		SessionID: sessionID,
		Messages:  msgs,
	}

	return json.Marshal(pkg)
}

// ProcessDKGRound2Messages processes Round 2 messages and finalizes DKG
func (k Keeper) ProcessDKGRound2Messages(sessionID string, round2Messages [][]byte) (*eddsa.PublicKey, *eddsa.SecretShare, *eddsa.Public, error) {
	frostStateManager.mu.Lock()
	defer frostStateManager.mu.Unlock()

	frostState, exists := frostStateManager.dkgStates[sessionID]
	if !exists {
		return nil, nil, nil, fmt.Errorf("DKG state not initialized for session %s", sessionID)
	}

	output, exists := frostStateManager.dkgOutputs[sessionID]
	if !exists {
		return nil, nil, nil, fmt.Errorf("DKG output not initialized for session %s", sessionID)
	}

	// Collect all inner messages
	var allMsgs [][]byte
	for _, msgData := range round2Messages {
		var pkg FROSTDKGRound2Msg
		if err := json.Unmarshal(msgData, &pkg); err != nil {
			continue
		}
		allMsgs = append(allMsgs, pkg.Messages...)
	}

	// Process round 2 messages to finalize
	_, err := helpers.PartyRoutine(allMsgs, frostState)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to process round 2: %w", err)
	}

	// Wait for completion
	if err := frostState.WaitForError(); err != nil {
		return nil, nil, nil, fmt.Errorf("DKG failed: %w", err)
	}

	// Get results
	groupKey := output.Public.GroupKey
	secretShare := output.SecretKey
	publicShares := output.Public

	return groupKey, secretShare, publicShares, nil
}

// StoreFROSTKeyShare stores the FROST key share for future signing
// It saves to both memory and disk for persistence across restarts
func (k Keeper) StoreFROSTKeyShare(keySetID string, secretShare *eddsa.SecretShare, publicShares *eddsa.Public) {
	frostStateManager.mu.Lock()
	defer frostStateManager.mu.Unlock()

	frostStateManager.keyShares[keySetID] = secretShare
	frostStateManager.publicShares[keySetID] = publicShares

	// Persist to disk for restart recovery
	if err := saveFROSTKeyShareToDisk(keySetID, secretShare, publicShares); err != nil {
		fmt.Printf("Warning: failed to persist FROST key share to disk: %v\n", err)
	} else {
		fmt.Printf("FROST key share saved to disk for keyset: %s\n", keySetID)
	}
}

// CleanupDKGState removes DKG state after completion
func (k Keeper) CleanupDKGState(sessionID string) {
	frostStateManager.mu.Lock()
	defer frostStateManager.mu.Unlock()

	delete(frostStateManager.dkgStates, sessionID)
	delete(frostStateManager.dkgOutputs, sessionID)
}

// ========================
// Signing Functions
// ========================

// InitSignState initializes FROST signing state for this validator
func (k Keeper) InitSignState(requestID, keySetID string, selfIndex int, signerIndices []int, message []byte) error {
	frostStateManager.mu.Lock()
	defer frostStateManager.mu.Unlock()

	// Check if already initialized
	if _, exists := frostStateManager.signStates[requestID]; exists {
		return nil // Already initialized
	}

	// Get stored key shares
	secretShare, exists := frostStateManager.keyShares[keySetID]
	if !exists {
		return fmt.Errorf("no key share found for keyset %s", keySetID)
	}

	publicShares, exists := frostStateManager.publicShares[keySetID]
	if !exists {
		return fmt.Errorf("no public shares found for keyset %s", keySetID)
	}

	// Create signer party IDs (1-indexed)
	signerIDs := make(party.IDSlice, len(signerIndices))
	for i, idx := range signerIndices {
		signerIDs[i] = party.ID(idx + 1)
	}

	// Initialize FROST sign state
	signState, signOutput, err := frost.NewSignState(signerIDs, secretShare, publicShares, message, 0)
	if err != nil {
		return fmt.Errorf("failed to init sign state: %w", err)
	}

	frostStateManager.signStates[requestID] = signState
	frostStateManager.signOutputs[requestID] = signOutput

	return nil
}

// GenerateSigningRound1Message generates real FROST signing commitment
func (k Keeper) GenerateSigningRound1Message(requestID, validatorAddr string) ([]byte, error) {
	frostStateManager.mu.Lock()
	defer frostStateManager.mu.Unlock()

	signState, exists := frostStateManager.signStates[requestID]
	if !exists {
		return nil, fmt.Errorf("sign state not initialized for request %s", requestID)
	}

	// Process round 1 (no input for first round)
	msgs, err := helpers.PartyRoutine(nil, signState)
	if err != nil {
		return nil, fmt.Errorf("failed to generate signing round 1: %w", err)
	}

	pkg := FROSTSignRound1Msg{
		RequestID:    requestID,
		ValidatorAddr: validatorAddr,
		Messages:     msgs,
	}

	return json.Marshal(pkg)
}

// ProcessSigningRound1Messages processes commitments and generates signature shares
func (k Keeper) ProcessSigningRound1Messages(requestID string, round1Messages [][]byte) ([]byte, error) {
	frostStateManager.mu.Lock()
	defer frostStateManager.mu.Unlock()

	signState, exists := frostStateManager.signStates[requestID]
	if !exists {
		return nil, fmt.Errorf("sign state not initialized for request %s", requestID)
	}

	// Collect all inner messages
	var allMsgs [][]byte
	for _, msgData := range round1Messages {
		var pkg FROSTSignRound1Msg
		if err := json.Unmarshal(msgData, &pkg); err != nil {
			continue
		}
		allMsgs = append(allMsgs, pkg.Messages...)
	}

	// Process round 1 messages to generate round 2 (signature shares)
	msgs, err := helpers.PartyRoutine(allMsgs, signState)
	if err != nil {
		return nil, fmt.Errorf("failed to process signing round 1: %w", err)
	}

	pkg := FROSTSignRound2Msg{
		RequestID: requestID,
		Messages:  msgs,
	}

	return json.Marshal(pkg)
}

// ProcessSigningRound2Messages processes signature shares and produces final signature
func (k Keeper) ProcessSigningRound2Messages(requestID string, round2Messages [][]byte) ([]byte, error) {
	frostStateManager.mu.Lock()
	defer frostStateManager.mu.Unlock()

	signState, exists := frostStateManager.signStates[requestID]
	if !exists {
		return nil, fmt.Errorf("sign state not initialized for request %s", requestID)
	}

	signOutput, exists := frostStateManager.signOutputs[requestID]
	if !exists {
		return nil, fmt.Errorf("sign output not initialized for request %s", requestID)
	}

	// Collect all inner messages
	var allMsgs [][]byte
	for _, msgData := range round2Messages {
		var pkg FROSTSignRound2Msg
		if err := json.Unmarshal(msgData, &pkg); err != nil {
			continue
		}
		allMsgs = append(allMsgs, pkg.Messages...)
	}

	// Process round 2 to finalize signature
	_, err := helpers.PartyRoutine(allMsgs, signState)
	if err != nil {
		return nil, fmt.Errorf("failed to process signing round 2: %w", err)
	}

	// Wait for completion
	if err := signState.WaitForError(); err != nil {
		return nil, fmt.Errorf("signing failed: %w", err)
	}

	// Get signature
	sig := signOutput.Signature
	if sig == nil {
		return nil, fmt.Errorf("signature is nil")
	}

	// Marshal to bytes (Ed25519 format - 64 bytes)
	sigBytes := sig.ToEd25519()

	return sigBytes, nil
}

// CleanupSignState removes signing state after completion
func (k Keeper) CleanupSignState(requestID string) {
	frostStateManager.mu.Lock()
	defer frostStateManager.mu.Unlock()

	delete(frostStateManager.signStates, requestID)
	delete(frostStateManager.signOutputs, requestID)
}

// ========================
// Message Types
// ========================

// FROSTDKGRound1Msg wraps DKG Round 1 messages
type FROSTDKGRound1Msg struct {
	SessionID     string   `json:"session_id"`
	ValidatorAddr string   `json:"validator_addr"`
	Messages      [][]byte `json:"messages"`
}

// FROSTDKGRound2Msg wraps DKG Round 2 messages
type FROSTDKGRound2Msg struct {
	SessionID string   `json:"session_id"`
	Messages  [][]byte `json:"messages"`
}

// FROSTSignRound1Msg wraps signing Round 1 messages
type FROSTSignRound1Msg struct {
	RequestID     string   `json:"request_id"`
	ValidatorAddr string   `json:"validator_addr"`
	Messages      [][]byte `json:"messages"`
}

// FROSTSignRound2Msg wraps signing Round 2 messages
type FROSTSignRound2Msg struct {
	RequestID string   `json:"request_id"`
	Messages  [][]byte `json:"messages"`
}

// ========================
// Integration with existing code
// ========================

// CompleteDKGCeremonyReal performs DKG using real FROST
func (k Keeper) CompleteDKGCeremonyReal(ctx context.Context, session types.DKGSession, round1Data, round2Data map[string][]byte) ([]byte, map[string][]byte, error) {
	sessionID := session.Id

	// Collect round 1 messages
	var round1Messages [][]byte
	for _, data := range round1Data {
		round1Messages = append(round1Messages, data)
	}

	// Collect round 2 messages
	var round2Messages [][]byte
	for _, data := range round2Data {
		round2Messages = append(round2Messages, data)
	}

	// Finalize DKG
	groupKey, secretShare, publicShares, err := k.ProcessDKGRound2Messages(sessionID, round2Messages)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to complete DKG: %w", err)
	}

	// Store key shares for future signing
	k.StoreFROSTKeyShare(session.KeySetId, secretShare, publicShares)

	// Serialize group public key
	groupPubkeyBytes := groupKey.ToEd25519()

	// Create key share references for participants
	keyShares := make(map[string][]byte)
	for _, validatorAddr := range session.Participants {
		shareRef := map[string]interface{}{
			"keyset_id": session.KeySetId,
			"threshold": session.Threshold,
			"curve":     "ed25519",
			"protocol":  "frost",
		}
		shareRefBytes, _ := json.Marshal(shareRef)
		keyShares[validatorAddr] = shareRefBytes
	}

	// Cleanup DKG state
	k.CleanupDKGState(sessionID)

	return groupPubkeyBytes, keyShares, nil
}

// AggregateSignatureReal performs signature aggregation using real FROST
func (k Keeper) AggregateSignatureReal(ctx context.Context, request types.SigningRequest, round2Data map[string][]byte) ([]byte, error) {
	requestID := request.Id

	// Collect round 2 messages
	var round2Messages [][]byte
	for _, data := range round2Data {
		round2Messages = append(round2Messages, data)
	}

	// Finalize signature
	signature, err := k.ProcessSigningRound2Messages(requestID, round2Messages)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate signature: %w", err)
	}

	// Cleanup sign state
	k.CleanupSignState(requestID)

	return signature, nil
}
