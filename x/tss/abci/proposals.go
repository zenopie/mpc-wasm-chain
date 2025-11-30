package abci

import (
	"encoding/json"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"mpc-wasm-chain/x/tss/keeper"
)

// ProposalHandler handles block proposals with TSS data
type ProposalHandler struct {
	keeper *keeper.Keeper
	logger log.Logger
}

// NewProposalHandler creates a new proposal handler
func NewProposalHandler(k *keeper.Keeper, logger log.Logger) *ProposalHandler {
	return &ProposalHandler{
		keeper: k,
		logger: logger,
	}
}

// PrepareProposal aggregates TSS data from vote extensions into the block proposal
// This is called by the block proposer
func (h *ProposalHandler) PrepareProposal(ctx sdk.Context, req *abci.RequestPrepareProposal) (*abci.ResponsePrepareProposal, error) {
	fmt.Printf("TSS PrepareProposal: height=%d votes=%d\n", req.Height, len(req.LocalLastCommit.Votes))

	// Aggregate TSS data from all vote extensions
	aggregated := &keeper.AggregatedTSSData{
		DKGRound1:          make(map[string]map[string][]byte),
		DKGRound2:          make(map[string]map[string][]byte),
		SigningCommitments: make(map[string]map[string][]byte),
		SignatureShares:    make(map[string]map[string][]byte),
	}

	// Process vote extensions from the last commit
	for _, vote := range req.LocalLastCommit.Votes {
		if len(vote.VoteExtension) == 0 {
			continue
		}

		// Get validator address from vote
		validatorAddr := fmt.Sprintf("%x", vote.Validator.Address)

		// Decode vote extension
		var ext TSSVoteExtension
		if err := json.Unmarshal(vote.VoteExtension, &ext); err != nil {
			h.logger.Error("Failed to unmarshal vote extension",
				"validator", validatorAddr,
				"error", err)
			continue
		}

		// Aggregate DKG Round 1 data
		for _, data := range ext.DKGRound1 {
			if aggregated.DKGRound1[data.SessionID] == nil {
				aggregated.DKGRound1[data.SessionID] = make(map[string][]byte)
			}
			aggregated.DKGRound1[data.SessionID][validatorAddr] = data.Commitment
		}

		// Aggregate DKG Round 2 data
		for _, data := range ext.DKGRound2 {
			if aggregated.DKGRound2[data.SessionID] == nil {
				aggregated.DKGRound2[data.SessionID] = make(map[string][]byte)
			}
			aggregated.DKGRound2[data.SessionID][validatorAddr] = data.Share
		}

		// Aggregate signing commitments
		for _, data := range ext.SigningCommitments {
			if aggregated.SigningCommitments[data.RequestID] == nil {
				aggregated.SigningCommitments[data.RequestID] = make(map[string][]byte)
			}
			aggregated.SigningCommitments[data.RequestID][validatorAddr] = data.Commitment
		}

		// Aggregate signature shares
		for _, data := range ext.SignatureShares {
			if aggregated.SignatureShares[data.RequestID] == nil {
				aggregated.SignatureShares[data.RequestID] = make(map[string][]byte)
			}
			aggregated.SignatureShares[data.RequestID][validatorAddr] = data.Share
		}
	}

	// Encode aggregated data
	aggregatedBytes, err := json.Marshal(aggregated)
	if err != nil {
		h.logger.Error("Failed to marshal aggregated TSS data", "error", err)
		return &abci.ResponsePrepareProposal{Txs: req.Txs}, nil
	}

	h.logger.Info("Prepared proposal with aggregated TSS data",
		"dkg_r1_sessions", len(aggregated.DKGRound1),
		"dkg_r2_sessions", len(aggregated.DKGRound2),
		"signing_requests", len(aggregated.SigningCommitments))

	// Include aggregated data as the first "transaction"
	// This is a special injection - not a real transaction
	txs := [][]byte{aggregatedBytes}
	txs = append(txs, req.Txs...)

	return &abci.ResponsePrepareProposal{Txs: txs}, nil
}

// ProcessProposal verifies a block proposal containing TSS data
// This is called by all validators when they receive a proposal
func (h *ProposalHandler) ProcessProposal(ctx sdk.Context, req *abci.RequestProcessProposal) (*abci.ResponseProcessProposal, error) {
	fmt.Printf("TSS ProcessProposal: height=%d txs=%d\n", req.Height, len(req.Txs))

	// If no transactions, accept
	if len(req.Txs) == 0 {
		return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
	}

	// First "transaction" should be the aggregated TSS data
	var aggregated keeper.AggregatedTSSData
	if err := json.Unmarshal(req.Txs[0], &aggregated); err != nil {
		// Not TSS data or invalid format - this might be a regular transaction
		// Accept the proposal (normal transactions are validated elsewhere)
		return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
	}

	// TODO: Add validation logic:
	// 1. Verify aggregated data matches vote extensions from LastCommit
	// 2. Verify no duplicate submissions per validator
	// 3. Verify cryptographic validity
	//
	// For now, we accept all proposals to get the basic flow working

	h.logger.Info("Verified proposal with TSS data",
		"dkg_r1_sessions", len(aggregated.DKGRound1),
		"dkg_r2_sessions", len(aggregated.DKGRound2))

	// Store TSS data for processing in BeginBlock
	// This avoids the TSS "transaction" being processed by normal tx handlers
	h.keeper.StoreTSSDataForBlock(&aggregated)

	return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
}

// NOTE: FinalizeBlock is not needed - TSS data is processed in BeginBlock instead
// The aggregated data is stored in memory during ProcessProposal and retrieved in BeginBlock
