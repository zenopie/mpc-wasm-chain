package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/tss"
)

// Mock FROST implementations for testing
// These generate deterministic placeholder data

// generateDKGRound1DataMock creates mock DKG Round 1 data
func (k Keeper) generateDKGRound1DataMock(ctx context.Context, sessionID, validatorAddr string) []byte {
	session, err := k.GetDKGSession(ctx, sessionID)
	if err != nil {
		return nil
	}

	participantIndex := -1
	for i, addr := range session.Participants {
		if addr == validatorAddr {
			participantIndex = i
			break
		}
	}
	if participantIndex < 0 {
		return nil
	}

	partyID := tss.NewPartyID(
		fmt.Sprintf("validator-%d", participantIndex),
		"",
		new(big.Int).SetInt64(int64(participantIndex)),
	)

	pkg := DKGRound1Package{
		ParticipantID: fmt.Sprintf("validator-%d", participantIndex),
		From:          partyID,
		Message:       &keygen.KGRound1Message{},
	}

	pkgBytes, err := json.Marshal(pkg)
	if err != nil {
		return nil
	}

	return pkgBytes
}

// generateDKGRound2DataMock creates mock DKG Round 2 data
func (k Keeper) generateDKGRound2DataMock(ctx context.Context, sessionID, validatorAddr string) []byte {
	session, err := k.GetDKGSession(ctx, sessionID)
	if err != nil {
		return nil
	}

	participantIndex := -1
	for i, addr := range session.Participants {
		if addr == validatorAddr {
			participantIndex = i
			break
		}
	}
	if participantIndex < 0 {
		return nil
	}

	sharesHash := k.generateDKGRound2SharesHash(sessionID, validatorAddr, participantIndex)

	pkg := DKGRound2Package{
		ParticipantID:        fmt.Sprintf("validator-%d", participantIndex),
		VerificationComplete: true,
		SharesHash:           sharesHash,
	}

	pkgBytes, err := json.Marshal(pkg)
	if err != nil {
		return nil
	}

	return pkgBytes
}

// generateSigningCommitmentMock creates mock signing commitment
func (k Keeper) generateSigningCommitmentMock(ctx context.Context, requestID, validatorAddr string) []byte {
	request, err := k.GetSigningRequest(ctx, requestID)
	if err != nil {
		return nil
	}

	session, err := k.SigningSessionStore.Get(ctx, requestID)
	if err != nil {
		return nil
	}

	participantIndex := -1
	for i, addr := range session.Participants {
		if addr == validatorAddr {
			participantIndex = i
			break
		}
	}
	if participantIndex < 0 {
		return nil
	}

	commitment := k.generateSigningCommitment(requestID, validatorAddr, request.MessageHash)

	pkg := SigningCommitmentPackage{
		ParticipantID: fmt.Sprintf("validator-%d", participantIndex),
		Commitment:    commitment,
	}

	pkgBytes, err := json.Marshal(pkg)
	if err != nil {
		return nil
	}

	return pkgBytes
}

// generateSignatureShareMock creates mock signature share
func (k Keeper) generateSignatureShareMock(ctx context.Context, requestID, validatorAddr string) []byte {
	request, err := k.GetSigningRequest(ctx, requestID)
	if err != nil {
		return nil
	}

	session, err := k.SigningSessionStore.Get(ctx, requestID)
	if err != nil {
		return nil
	}

	participantIndex := -1
	for i, addr := range session.Participants {
		if addr == validatorAddr {
			participantIndex = i
			break
		}
	}
	if participantIndex < 0 {
		return nil
	}

	r, s := k.generateSignatureShare(requestID, validatorAddr, request.MessageHash, participantIndex)

	pkg := SignatureSharePackage{
		ParticipantID: fmt.Sprintf("validator-%d", participantIndex),
		R:             r,
		S:             s,
	}

	pkgBytes, err := json.Marshal(pkg)
	if err != nil {
		return nil
	}

	return pkgBytes
}
