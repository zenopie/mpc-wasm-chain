package keeper

import (
	"context"
	"encoding/json"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AggregatedTSSData contains all TSS data from vote extensions
// This matches the structure in abci/proposals.go
type AggregatedTSSData struct {
	DKGRound1          map[string]map[string][]byte `json:"dkg_round1"`
	DKGRound2          map[string]map[string][]byte `json:"dkg_round2"`
	SigningCommitments map[string]map[string][]byte `json:"signing_commitments"`
	SignatureShares    map[string]map[string][]byte `json:"signature_shares"`
}

// In-memory storage for TSS data between ProcessProposal and BeginBlock
// This is safe because ProcessProposal and BeginBlock run sequentially on the same node
var (
	pendingTSSData *AggregatedTSSData
	tssDataMutex   sync.RWMutex
)

// StoreTSSDataForBlock stores aggregated TSS data to be processed in BeginBlock
// Called by ProposalHandler during ProcessProposal
func (k Keeper) StoreTSSDataForBlock(data *AggregatedTSSData) {
	tssDataMutex.Lock()
	defer tssDataMutex.Unlock()
	pendingTSSData = data
}

// ProcessTSSDataFromBlock retrieves and processes TSS data stored during ProcessProposal
// Called during BeginBlock
func (k Keeper) ProcessTSSDataFromBlock(ctx context.Context) error {
	tssDataMutex.Lock()
	defer tssDataMutex.Unlock()

	// No pending TSS data - this is normal for blocks without TSS activity
	if pendingTSSData == nil {
		return nil
	}

	// Get data and clear it
	data := pendingTSSData
	pendingTSSData = nil

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	logger := sdkCtx.Logger().With("module", "tss", "phase", "begin_block")

	logger.Info("Processing aggregated TSS data from vote extensions",
		"height", sdkCtx.BlockHeight(),
		"dkg_r1_sessions", len(data.DKGRound1),
		"dkg_r2_sessions", len(data.DKGRound2),
		"signing_commitments", len(data.SigningCommitments),
		"signature_shares", len(data.SignatureShares))

	// Process DKG Round 1 data
	for sessionID, validators := range data.DKGRound1 {
		for validatorAddr, commitment := range validators {
			if err := k.ProcessDKGRound1(ctx, sessionID, validatorAddr, commitment); err != nil {
				logger.Error("Failed to process DKG Round 1 data",
					"session", sessionID,
					"validator", validatorAddr,
					"error", err)
				// Continue processing other data even if one fails
			}
		}
	}

	// Process DKG Round 2 data
	for sessionID, validators := range data.DKGRound2 {
		for validatorAddr, share := range validators {
			if err := k.ProcessDKGRound2(ctx, sessionID, validatorAddr, share); err != nil {
				logger.Error("Failed to process DKG Round 2 data",
					"session", sessionID,
					"validator", validatorAddr,
					"error", err)
			}
		}
	}

	// Process signing commitments
	for requestID, validators := range data.SigningCommitments {
		for validatorAddr, commitment := range validators {
			if err := k.ProcessSigningCommitment(ctx, requestID, validatorAddr, commitment); err != nil {
				logger.Error("Failed to process signing commitment",
					"request", requestID,
					"validator", validatorAddr,
					"error", err)
			}
		}
	}

	// Process signature shares
	for requestID, validators := range data.SignatureShares {
		for validatorAddr, share := range validators {
			if err := k.ProcessSignatureShare(ctx, requestID, validatorAddr, share); err != nil {
				logger.Error("Failed to process signature share",
					"request", requestID,
					"validator", validatorAddr,
					"error", err)
			}
		}
	}

	logger.Info("Finished processing TSS data from vote extensions",
		"height", sdkCtx.BlockHeight())

	return nil
}

// ParseTSSDataFromBytes parses aggregated TSS data from JSON bytes
func ParseTSSDataFromBytes(data []byte) (*AggregatedTSSData, error) {
	var aggregated AggregatedTSSData
	if err := json.Unmarshal(data, &aggregated); err != nil {
		return nil, err
	}
	return &aggregated, nil
}
