package keeper

import (
	"context"
)

// Vote Extension Helper Methods
// These methods generate TSS data for inclusion in vote extensions
// They are called by the ABCI vote extension handlers

// UseMockFROST controls whether to use mock or real FROST
// Set to false to use real FROST (requires all validators to have state)
var UseMockFROST = false

// GenerateDKGRound1Data creates DKG Round 1 commitment data for this validator
// Returns serialized commitment bytes for inclusion in vote extension
func (k Keeper) GenerateDKGRound1Data(ctx context.Context, sessionID, validatorAddr string) []byte {
	if UseMockFROST {
		return k.generateDKGRound1DataMock(ctx, sessionID, validatorAddr)
	}
	return k.GenerateDKGRound1DataReal(ctx, sessionID, validatorAddr)
}

// GenerateDKGRound2Data creates DKG Round 2 share data for this validator
// Returns serialized share bytes for inclusion in vote extension
func (k Keeper) GenerateDKGRound2Data(ctx context.Context, sessionID, validatorAddr string) []byte {
	if UseMockFROST {
		return k.generateDKGRound2DataMock(ctx, sessionID, validatorAddr)
	}
	return k.GenerateDKGRound2DataReal(ctx, sessionID, validatorAddr)
}

// GenerateSigningCommitment creates signing Round 1 commitment for this validator
// Returns serialized commitment bytes for inclusion in vote extension
func (k Keeper) GenerateSigningCommitment(ctx context.Context, requestID, validatorAddr string) []byte {
	if UseMockFROST {
		return k.generateSigningCommitmentMock(ctx, requestID, validatorAddr)
	}
	return k.GenerateSigningCommitmentReal(ctx, requestID, validatorAddr)
}

// GenerateSignatureShare creates signing Round 2 signature share for this validator
// Returns serialized share bytes for inclusion in vote extension
func (k Keeper) GenerateSignatureShare(ctx context.Context, requestID, validatorAddr string) []byte {
	if UseMockFROST {
		return k.generateSignatureShareMock(ctx, requestID, validatorAddr)
	}
	return k.GenerateSignatureShareReal(ctx, requestID, validatorAddr)
}
