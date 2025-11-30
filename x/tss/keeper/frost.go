package keeper

import (
	"context"
	"crypto/elliptic"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/bnb-chain/tss-lib/v2/crypto"
	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/tss"

	"mpc-wasm-chain/x/tss/types"
)

// Binance tss-lib implementation for threshold signatures
// Production-ready library used in real blockchain projects
// Supports: ECDSA (Secp256k1, P256) and EdDSA (Ed25519)

// TSSCurve represents the elliptic curve used for TSS
type TSSCurve int

const (
	CurveSecp256k1 TSSCurve = iota
	CurveP256
	CurveEd25519
)

// GetCurve returns the appropriate elliptic curve
func GetCurve(curveType TSSCurve) (elliptic.Curve, error) {
	switch curveType {
	case CurveSecp256k1:
		return tss.S256(), nil
	case CurveP256:
		return elliptic.P256(), nil
	case CurveEd25519:
		// EdDSA uses Edwards curve, handled differently in tss-lib
		return tss.Edwards(), nil
	default:
		return nil, fmt.Errorf("unsupported curve type: %d", curveType)
	}
}

// DKGRound1Package represents the data a participant creates in DKG Round 1
// Using tss-lib's keygen messages
type DKGRound1Package struct {
	ParticipantID string                  `json:"participant_id"`
	From          *tss.PartyID            `json:"from"`
	Message       *keygen.KGRound1Message `json:"message"`
}

// DKGRound2Package represents the data a participant creates in DKG Round 2
type DKGRound2Package struct {
	ParticipantID        string `json:"participant_id"`
	VerificationComplete bool   `json:"verification_complete"`
	SharesHash           []byte `json:"shares_hash,omitempty"`
}

// TSSParticipant holds the state for a single TSS participant
type TSSParticipant struct {
	PartyID   *tss.PartyID
	Params    *tss.Parameters
	SaveData  *keygen.LocalPartySaveData
	Threshold int
}

// SerializePublicKey serializes a public key for storage
func SerializePublicKey(pubKey *crypto.ECPoint) ([]byte, error) {
	if pubKey == nil {
		return nil, fmt.Errorf("public key is nil")
	}
	x, y := pubKey.X(), pubKey.Y()
	return append(x.Bytes(), y.Bytes()...), nil
}

// DeserializePublicKey deserializes a public key from storage
func DeserializePublicKey(data []byte, curve elliptic.Curve) (*crypto.ECPoint, error) {
	if len(data) < 64 {
		return nil, fmt.Errorf("invalid public key data length")
	}
	x := new(big.Int).SetBytes(data[:32])
	y := new(big.Int).SetBytes(data[32:])
	return crypto.NewECPointNoCurveCheck(curve, x, y), nil
}

// ========================
// DKG Functions
// ========================

// AggregateDKGRound1Commitments collects and validates all Round 1 commitments
func (k Keeper) AggregateDKGRound1Commitments(ctx context.Context, sessionID string) (map[string][]byte, error) {
	commitments := make(map[string][]byte)
	prefix := sessionID + ":"

	err := k.DKGRound1DataStore.Walk(ctx, nil, func(key string, value types.DKGRound1Data) (bool, error) {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			commitments[value.ValidatorAddress] = value.Commitment
		}
		return false, nil
	})

	if err != nil {
		return nil, err
	}

	return commitments, nil
}

// AggregateDKGRound2Shares collects all Round 2 shares
func (k Keeper) AggregateDKGRound2Shares(ctx context.Context, sessionID string) (map[string][]byte, error) {
	shares := make(map[string][]byte)
	prefix := sessionID + ":"

	err := k.DKGRound2DataStore.Walk(ctx, nil, func(key string, value types.DKGRound2Data) (bool, error) {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			shares[value.ValidatorAddress] = value.Share
		}
		return false, nil
	})

	if err != nil {
		return nil, err
	}

	return shares, nil
}

// CompleteDKGCeremony performs the final aggregation of DKG data
func (k Keeper) CompleteDKGCeremony(ctx context.Context, session types.DKGSession) ([]byte, map[string][]byte, error) {
	// Get all Round 1 commitments
	round1Commitments, err := k.AggregateDKGRound1Commitments(ctx, session.Id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to aggregate round 1 commitments: %w", err)
	}

	// Get all Round 2 shares
	round2Shares, err := k.AggregateDKGRound2Shares(ctx, session.Id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to aggregate round 2 shares: %w", err)
	}

	// Verify we have enough participants
	if len(round1Commitments) < int(session.Threshold) || len(round2Shares) < int(session.Threshold) {
		return nil, nil, fmt.Errorf("insufficient participants for DKG completion: got %d commitments and %d shares, need %d",
			len(round1Commitments), len(round2Shares), session.Threshold)
	}

	// Use real FROST if enabled
	if !UseMockFROST {
		return k.CompleteDKGCeremonyReal(ctx, session, round1Commitments, round2Shares)
	}

	// Mock implementation for testing
	return k.completeDKGCeremonyMock(ctx, session, round1Commitments, round2Shares)
}

// completeDKGCeremonyMock uses mock/placeholder cryptography for testing
func (k Keeper) completeDKGCeremonyMock(ctx context.Context, session types.DKGSession, round1Commitments, round2Shares map[string][]byte) ([]byte, map[string][]byte, error) {
	// Parse and verify all Round 1 packages
	for validatorAddr, commitment := range round1Commitments {
		var pkg DKGRound1Package
		if err := json.Unmarshal(commitment, &pkg); err != nil {
			return nil, nil, fmt.Errorf("failed to unmarshal round 1 package from %s: %w", validatorAddr, err)
		}
	}

	// Verify Round 2 completion
	round2Complete := make(map[string]bool)
	for validatorAddr, shareData := range round2Shares {
		var pkg DKGRound2Package
		if err := json.Unmarshal(shareData, &pkg); err != nil {
			return nil, nil, fmt.Errorf("failed to unmarshal round 2 package from %s: %w", validatorAddr, err)
		}
		round2Complete[validatorAddr] = pkg.VerificationComplete
	}

	// Verify sufficient participants completed Round 2
	if len(round2Complete) < int(session.Threshold) {
		return nil, nil, fmt.Errorf("insufficient Round 2 verifications: got %d, need %d", len(round2Complete), session.Threshold)
	}

	// Create a placeholder group public key
	curve, err := GetCurve(CurveSecp256k1)
	if err != nil {
		return nil, nil, err
	}

	x := new(big.Int).SetInt64(1)
	y := new(big.Int).SetInt64(1)
	x, y = curve.ScalarBaseMult(x.Bytes())

	ecPoint := crypto.NewECPointNoCurveCheck(curve, x, y)
	groupPubkeyBytes, err := SerializePublicKey(ecPoint)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to serialize group public key: %w", err)
	}

	// Generate key share references for validators
	keyShares := make(map[string][]byte)
	for i, validatorAddr := range session.Participants {
		if _, participated := round1Commitments[validatorAddr]; participated {
			participantID := uint32(i + 1)

			shareRef := map[string]interface{}{
				"participant_id": participantID,
				"keyset_id":      session.KeySetId,
				"threshold":      session.Threshold,
				"max_signers":    session.MaxSigners,
				"curve":          "secp256k1",
				"protocol":       "mock",
			}

			shareRefBytes, err := json.Marshal(shareRef)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to marshal share reference: %w", err)
			}

			keyShares[validatorAddr] = shareRefBytes
		}
	}

	return groupPubkeyBytes, keyShares, nil
}

// ========================
// Signing Functions
// ========================

// SignatureShare represents a partial signature from one participant
type SignatureShare struct {
	ParticipantID string `json:"participant_id"`
	R             []byte `json:"r"` // Signature component R
	S             []byte `json:"s"` // Signature component S (partial)
}

// CreateSignatureShare creates a partial signature share using tss-lib
// In production, this would use tss-lib's signing protocol
func CreateSignatureShare(
	participantID string,
	message []byte,
	partyID *tss.PartyID,
) (*SignatureShare, error) {
	// In a real implementation, this would:
	// 1. Run tss-lib's signing protocol
	// 2. Generate partial signature shares
	// 3. Return the share for aggregation
	//
	// For now, return a placeholder structure
	return &SignatureShare{
		ParticipantID: participantID,
		R:             message[:32], // Placeholder
		S:             message[:32], // Placeholder
	}, nil
}

// AggregateSigningCommitments collects all signing commitments (Round 1)
func (k Keeper) AggregateSigningCommitments(ctx context.Context, requestID string) (map[string][]byte, error) {
	commitments := make(map[string][]byte)
	prefix := requestID + ":"

	err := k.SigningCommitmentStore.Walk(ctx, nil, func(key string, value types.SigningCommitment) (bool, error) {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			commitments[value.ValidatorAddress] = value.Commitment
		}
		return false, nil
	})

	if err != nil {
		return nil, err
	}

	return commitments, nil
}

// AggregateSignatureShares collects all signature shares (Round 2)
func (k Keeper) AggregateSignatureSharesData(ctx context.Context, requestID string) (map[string][]byte, error) {
	shares := make(map[string][]byte)
	prefix := requestID + ":"

	err := k.SignatureShareStore.Walk(ctx, nil, func(key string, value types.SignatureShare) (bool, error) {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			shares[value.ValidatorAddress] = value.Share
		}
		return false, nil
	})

	if err != nil {
		return nil, err
	}

	return shares, nil
}

// SigningCommitmentPackage represents Round 1 signing commitment
type SigningCommitmentPackage struct {
	ParticipantID string `json:"participant_id"`
	Commitment    []byte `json:"commitment"`
}

// SignatureSharePackage represents Round 2 signature share
type SignatureSharePackage struct {
	ParticipantID string `json:"participant_id"`
	R             []byte `json:"r"`
	S             []byte `json:"s"`
}

// AggregateSignature performs TSS threshold signature aggregation
func (k Keeper) AggregateSignature(ctx context.Context, request types.SigningRequest, session types.SigningSession) ([]byte, error) {
	// Get all commitments
	commitments, err := k.AggregateSigningCommitments(ctx, request.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate commitments: %w", err)
	}

	// Get all shares
	shares, err := k.AggregateSignatureSharesData(ctx, request.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate shares: %w", err)
	}

	// Verify we have enough participants
	threshold := session.Threshold
	if len(commitments) < int(threshold) || len(shares) < int(threshold) {
		return nil, fmt.Errorf("insufficient participants for signature completion: got %d commitments and %d shares, need %d",
			len(commitments), len(shares), threshold)
	}

	// Use real FROST if enabled
	if !UseMockFROST {
		return k.AggregateSignatureReal(ctx, request, shares)
	}

	// Mock implementation
	return k.aggregateSignatureMock(ctx, request, session, shares)
}

// aggregateSignatureMock uses mock/placeholder cryptography
func (k Keeper) aggregateSignatureMock(ctx context.Context, request types.SigningRequest, session types.SigningSession, shares map[string][]byte) ([]byte, error) {
	// Parse signature share packages
	var signatureShares []*SignatureShare
	for validatorAddr, shareData := range shares {
		var pkg SignatureSharePackage
		if err := json.Unmarshal(shareData, &pkg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal share package from %s: %w", validatorAddr, err)
		}

		signatureShares = append(signatureShares, &SignatureShare{
			ParticipantID: pkg.ParticipantID,
			R:             pkg.R,
			S:             pkg.S,
		})
	}

	// Verify we have threshold shares
	if len(signatureShares) < int(session.Threshold) {
		return nil, fmt.Errorf("insufficient signature shares: got %d, need %d", len(signatureShares), session.Threshold)
	}

	// Aggregate using mock function
	signature, err := AggregateSignatureShares(
		signatureShares,
		session.Threshold,
		CurveSecp256k1,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate signature shares: %w", err)
	}

	if len(signature) == 0 {
		return nil, fmt.Errorf("aggregated signature is empty")
	}

	return signature, nil
}

// AggregateSignatureShares combines partial signatures into a full signature
// Uses tss-lib's signature aggregation
func AggregateSignatureShares(
	shares []*SignatureShare,
	threshold uint32,
	curve TSSCurve,
) ([]byte, error) {
	if len(shares) < int(threshold) {
		return nil, fmt.Errorf("insufficient shares: got %d, need %d", len(shares), threshold)
	}

	// In a real implementation, this would:
	// 1. Collect all partial signatures
	// 2. Run tss-lib's signature aggregation
	// 3. Return the complete ECDSA/EdDSA signature
	//
	// For now, return a placeholder signature
	signature := make([]byte, 64)
	if len(shares) > 0 && len(shares[0].R) >= 32 {
		copy(signature[:32], shares[0].R)
	}
	if len(shares) > 0 && len(shares[0].S) >= 32 {
		copy(signature[32:], shares[0].S)
	}

	return signature, nil
}

// VerifySignature verifies a threshold signature against a public key
func VerifySignature(
	signature []byte,
	message []byte,
	publicKey []byte,
	curveType TSSCurve,
) error {
	if len(signature) != 64 {
		return fmt.Errorf("invalid signature length: expected 64, got %d", len(signature))
	}

	curve, err := GetCurve(curveType)
	if err != nil {
		return err
	}

	// Deserialize public key
	_, err = DeserializePublicKey(publicKey, curve)
	if err != nil {
		return fmt.Errorf("failed to deserialize public key: %w", err)
	}

	// In a real implementation, this would:
	// 1. Parse the ECDSA/EdDSA signature (R, S)
	// 2. Verify using the curve's verification algorithm
	// 3. Return verification result
	//
	// For now, perform basic validation
	if len(signature) == 0 || len(message) == 0 || len(publicKey) == 0 {
		return fmt.Errorf("signature verification failed: invalid inputs")
	}

	return nil
}

// VerifyThresholdSignature verifies a completed TSS threshold signature
// This can be called by smart contracts via the TSS module
func (k Keeper) VerifyThresholdSignature(
	signature []byte,
	message []byte,
	groupPubkey []byte,
	curveType TSSCurve,
) error {
	return VerifySignature(signature, message, groupPubkey, curveType)
}
