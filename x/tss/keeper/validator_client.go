package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/tss"

	"mpc-wasm-chain/x/tss/types"
)

// ValidatorClient handles automatic TSS participation for validators
// This runs in the EndBlocker to detect and participate in DKG/signing sessions

// AutoParticipateDKG checks for DKG sessions and automatically submits data for all validators
// In testnet mode, this runs once and submits for all validators
// In production, each validator node runs this and submits for itself
func (k Keeper) AutoParticipateDKG(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get all active validator addresses
	activeValidators, err := k.GetActiveValidatorAddresses(ctx)
	if err != nil || len(activeValidators) == 0 {
		return nil
	}

	// Walk through all active DKG sessions
	return k.DKGSessionStore.Walk(ctx, nil, func(sessionID string, session types.DKGSession) (bool, error) {
		// Skip completed or failed sessions
		if session.State == types.DKGState_DKG_STATE_COMPLETE || session.State == types.DKGState_DKG_STATE_FAILED {
			return false, nil
		}

		// Check timeout
		if sdkCtx.BlockHeight() >= session.TimeoutHeight {
			return false, nil
		}

		// Submit for each validator that is a participant
		for _, validatorAddr := range session.Participants {
			// Participate based on current state
			switch session.State {
			case types.DKGState_DKG_STATE_ROUND1:
				// Check if already submitted
				key := fmt.Sprintf("%s:%s", sessionID, validatorAddr)
				has, err := k.DKGRound1DataStore.Has(ctx, key)
				if err != nil {
					return true, err
				}
				if has {
					continue // Already submitted
				}

				// Generate and submit Round 1 data
				if err := k.submitDKGRound1(ctx, session, validatorAddr); err != nil {
					// Log error but don't halt the chain
					sdkCtx.Logger().Error("Failed to auto-submit DKG Round 1", "error", err, "session", sessionID, "validator", validatorAddr)
				}

			case types.DKGState_DKG_STATE_ROUND2:
				// Check if already submitted
				key := fmt.Sprintf("%s:%s", sessionID, validatorAddr)
				has, err := k.DKGRound2DataStore.Has(ctx, key)
				if err != nil {
					return true, err
				}
				if has {
					continue // Already submitted
				}

				// Generate and submit Round 2 data
				if err := k.submitDKGRound2(ctx, session, validatorAddr); err != nil {
					// Log error but don't halt the chain
					sdkCtx.Logger().Error("Failed to auto-submit DKG Round 2", "error", err, "session", sessionID, "validator", validatorAddr)
				}
			}
		}

		return false, nil
	})
}

// AutoParticipateSigning checks for signing sessions and automatically submits data for all validators
// In testnet mode, this runs once and submits for all validators
// In production, each validator node runs this and submits for itself
func (k Keeper) AutoParticipateSigning(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get all active validator addresses
	activeValidators, err := k.GetActiveValidatorAddresses(ctx)
	if err != nil || len(activeValidators) == 0 {
		return nil
	}

	// Walk through all active signing sessions
	return k.SigningSessionStore.Walk(ctx, nil, func(requestID string, session types.SigningSession) (bool, error) {
		// Skip completed or failed sessions
		if session.State == types.SigningState_SIGNING_STATE_COMPLETE || session.State == types.SigningState_SIGNING_STATE_FAILED {
			return false, nil
		}

		// Check timeout
		if sdkCtx.BlockHeight() >= session.TimeoutHeight {
			return false, nil
		}

		// Submit for each validator that is a participant
		for _, validatorAddr := range session.Participants {
			// Participate based on current state
			switch session.State {
			case types.SigningState_SIGNING_STATE_ROUND1:
				// Check if already submitted commitment
				key := fmt.Sprintf("%s:%s", requestID, validatorAddr)
				has, err := k.SigningCommitmentStore.Has(ctx, key)
				if err != nil {
					return true, err
				}
				if has {
					continue
				}

				// Submit signing commitment
				if err := k.submitSigningCommitment(ctx, session, validatorAddr); err != nil {
					sdkCtx.Logger().Error("Failed to auto-submit signing commitment", "error", err, "request", requestID, "validator", validatorAddr)
				}

			case types.SigningState_SIGNING_STATE_ROUND2:
				// Check if already submitted share
				key := fmt.Sprintf("%s:%s", requestID, validatorAddr)
				has, err := k.SignatureShareStore.Has(ctx, key)
				if err != nil {
					return true, err
				}
				if has {
					continue
				}

				// Submit signature share
				if err := k.submitSignatureShare(ctx, session, validatorAddr); err != nil {
					sdkCtx.Logger().Error("Failed to auto-submit signature share", "error", err, "request", requestID, "validator", validatorAddr)
				}
			}
		}

		return false, nil
	})
}

// GetActiveValidatorAddresses returns all active validator consensus addresses
// Uses staking module to get validators, so it works during transaction execution
func (k Keeper) GetActiveValidatorAddresses(ctx context.Context) ([]string, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get all validators from staking module
	validators, err := k.stakingKeeper.GetAllValidators(ctx)
	if err != nil {
		sdkCtx.Logger().Error("Failed to get all validators", "error", err)
		return nil, fmt.Errorf("failed to get validators: %w", err)
	}

	sdkCtx.Logger().Info("GetActiveValidatorAddresses", "total_validators", len(validators))

	addresses := make([]string, 0)
	for _, val := range validators {
		// Only include bonded, non-jailed validators
		if !val.IsBonded() || val.IsJailed() {
			sdkCtx.Logger().Debug("Skipping validator", "operator", val.GetOperator(), "bonded", val.IsBonded(), "jailed", val.IsJailed())
			continue
		}

		// Get consensus public key
		consPubKey, err := val.ConsPubKey()
		if err != nil {
			sdkCtx.Logger().Error("Failed to get consensus pubkey", "operator", val.GetOperator(), "error", err)
			continue // Skip validators with invalid pubkey
		}

		// Convert to consensus address
		consAddr := sdk.ConsAddress(consPubKey.Address())

		// Convert to hex string for consistency
		hexAddr := fmt.Sprintf("%x", consAddr.Bytes())
		addresses = append(addresses, hexAddr)
		sdkCtx.Logger().Info("Added validator", "operator", val.GetOperator(), "consensus_addr", hexAddr)
	}

	sdkCtx.Logger().Info("GetActiveValidatorAddresses result", "count", len(addresses), "addresses", addresses)
	return addresses, nil
}

// GetValidatorAddress returns this validator's consensus address
// Returns the address set at startup from priv_validator_key.json
func (k Keeper) GetValidatorAddress(ctx context.Context) (string, error) {
	// Return the address that was set at app startup
	return k.ValidatorConsensusAddress, nil
}

// submitDKGRound1 generates and submits DKG Round 1 data for this validator
func (k Keeper) submitDKGRound1(ctx context.Context, session types.DKGSession, validatorAddr string) error {
	// Find participant index
	participantIndex := -1
	for i, addr := range session.Participants {
		if addr == validatorAddr {
			participantIndex = i
			break
		}
	}
	if participantIndex < 0 {
		return fmt.Errorf("validator not found in participants")
	}

	// Create party ID for this validator
	partyID := tss.NewPartyID(
		fmt.Sprintf("validator-%d", participantIndex),
		"",
		new(big.Int).SetInt64(int64(participantIndex)),
	)

	// In production, this would:
	// 1. Run tss-lib keygen Round 1 locally
	// 2. Generate commitments to secret shares
	// 3. Store the local state securely
	//
	// For now, we generate a placeholder package
	pkg := DKGRound1Package{
		ParticipantID: fmt.Sprintf("validator-%d", participantIndex),
		From:          partyID,
		Message:       &keygen.KGRound1Message{}, // Placeholder
	}

	pkgBytes, err := json.Marshal(pkg)
	if err != nil {
		return fmt.Errorf("failed to marshal DKG Round 1 package: %w", err)
	}

	// Store Round 1 data
	return k.ProcessDKGRound1(ctx, session.Id, validatorAddr, pkgBytes)
}

// submitDKGRound2 generates and submits DKG Round 2 data for this validator
func (k Keeper) submitDKGRound2(ctx context.Context, session types.DKGSession, validatorAddr string) error {
	// Find participant index
	participantIndex := -1
	for i, addr := range session.Participants {
		if addr == validatorAddr {
			participantIndex = i
			break
		}
	}
	if participantIndex < 0 {
		return fmt.Errorf("validator not found in participants")
	}

	// In production, this would:
	// 1. Verify Round 1 commitments from other participants
	// 2. Run tss-lib keygen Round 2
	// 3. Generate and encrypt secret shares for each participant
	// 4. Compute verification data
	//
	// For now, we generate placeholder verification data
	sharesHash := k.generateDKGRound2SharesHash(session.Id, validatorAddr, participantIndex)

	// Create the Round 2 package
	pkg := DKGRound2Package{
		ParticipantID:        fmt.Sprintf("validator-%d", participantIndex),
		VerificationComplete: true,
		SharesHash:           sharesHash,
	}

	pkgBytes, err := json.Marshal(pkg)
	if err != nil {
		return fmt.Errorf("failed to marshal DKG Round 2 package: %w", err)
	}

	// Store Round 2 data
	return k.ProcessDKGRound2(ctx, session.Id, validatorAddr, pkgBytes)
}

// submitSigningCommitment generates and submits signing Round 1 commitment
func (k Keeper) submitSigningCommitment(ctx context.Context, session types.SigningSession, validatorAddr string) error {
	// Get the signing request
	request, err := k.GetSigningRequest(ctx, session.RequestId)
	if err != nil {
		return err
	}

	// Find participant index
	participantIndex := -1
	for i, addr := range session.Participants {
		if addr == validatorAddr {
			participantIndex = i
			break
		}
	}
	if participantIndex < 0 {
		return fmt.Errorf("validator not found in participants")
	}

	// In production, this would:
	// 1. Load this validator's key share from secure storage
	// 2. Run tss-lib signing Round 1 (nonce generation)
	// 3. Create commitment to nonce
	//
	// For now, generate placeholder commitment
	commitment := k.generateSigningCommitment(session.RequestId, validatorAddr, request.MessageHash)

	pkg := SigningCommitmentPackage{
		ParticipantID: fmt.Sprintf("validator-%d", participantIndex),
		Commitment:    commitment,
	}

	pkgBytes, err := json.Marshal(pkg)
	if err != nil {
		return fmt.Errorf("failed to marshal signing commitment: %w", err)
	}

	// Store commitment
	return k.ProcessSigningCommitment(ctx, session.RequestId, validatorAddr, pkgBytes)
}

// submitSignatureShare generates and submits signature share
func (k Keeper) submitSignatureShare(ctx context.Context, session types.SigningSession, validatorAddr string) error {
	// Get the signing request
	request, err := k.GetSigningRequest(ctx, session.RequestId)
	if err != nil {
		return err
	}

	// Find participant index
	participantIndex := -1
	for i, addr := range session.Participants {
		if addr == validatorAddr {
			participantIndex = i
			break
		}
	}
	if participantIndex < 0 {
		return fmt.Errorf("validator not found in participants")
	}

	// In production, this would:
	// 1. Verify commitments from other participants
	// 2. Run tss-lib signing Round 2
	// 3. Generate partial signature share
	//
	// For now, generate placeholder signature share
	r, s := k.generateSignatureShare(session.RequestId, validatorAddr, request.MessageHash, participantIndex)

	pkg := SignatureSharePackage{
		ParticipantID: fmt.Sprintf("validator-%d", participantIndex),
		R:             r,
		S:             s,
	}

	pkgBytes, err := json.Marshal(pkg)
	if err != nil {
		return fmt.Errorf("failed to marshal signature share: %w", err)
	}

	// Store share
	return k.ProcessSignatureShare(ctx, session.RequestId, validatorAddr, pkgBytes)
}

// Helper functions to generate deterministic placeholder data
// In production, these would be replaced with real tss-lib operations

func (k Keeper) generateDKGRound1Commitment(sessionID, validatorAddr string, index int) []byte {
	h := sha256.New()
	h.Write([]byte(sessionID))
	h.Write([]byte(validatorAddr))
	h.Write([]byte(fmt.Sprintf("round1-%d", index)))
	return h.Sum(nil)
}

func (k Keeper) generateDKGRound2SharesHash(sessionID, validatorAddr string, index int) []byte {
	h := sha256.New()
	h.Write([]byte(sessionID))
	h.Write([]byte(validatorAddr))
	h.Write([]byte(fmt.Sprintf("round2-%d", index)))
	return h.Sum(nil)
}

func (k Keeper) generateSigningCommitment(requestID, validatorAddr string, messageHash []byte) []byte {
	h := sha256.New()
	h.Write([]byte(requestID))
	h.Write([]byte(validatorAddr))
	h.Write(messageHash)
	h.Write([]byte("commitment"))
	return h.Sum(nil)
}

func (k Keeper) generateSignatureShare(requestID, validatorAddr string, messageHash []byte, index int) ([]byte, []byte) {
	// Generate R component
	hR := sha256.New()
	hR.Write([]byte(requestID))
	hR.Write([]byte(validatorAddr))
	hR.Write(messageHash)
	hR.Write([]byte(fmt.Sprintf("R-%d", index)))
	r := hR.Sum(nil)

	// Generate S component
	hS := sha256.New()
	hS.Write([]byte(requestID))
	hS.Write([]byte(validatorAddr))
	hR.Write(messageHash)
	hS.Write([]byte(fmt.Sprintf("S-%d", index)))
	s := hS.Sum(nil)

	return r, s
}
