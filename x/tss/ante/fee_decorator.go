package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"mpc-wasm-chain/x/tss/types"
)

// TSSFeeDecorator exempts TSS messages from gas fees
// TSS participation is protocol-required for validators, so they shouldn't pay fees
type TSSFeeDecorator struct{}

func NewTSSFeeDecorator() TSSFeeDecorator {
	return TSSFeeDecorator{}
}

// AnteHandle exempts TSS messages from fee requirements
func (tfd TSSFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// Check if all messages in the transaction are TSS messages
	allTSSMessages := true
	msgs := tx.GetMsgs()

	for _, msg := range msgs {
		switch msg.(type) {
		case *types.MsgSubmitDKGRound1,
			*types.MsgSubmitDKGRound2,
			*types.MsgSubmitCommitment,
			*types.MsgSubmitSignatureShare:
			// TSS message - continue checking
			continue
		default:
			// Non-TSS message found
			allTSSMessages = false
			break
		}
	}

	// If all messages are TSS messages, mark the transaction as fee-exempt
	if allTSSMessages && len(msgs) > 0 {
		// Set zero minimum gas prices for this transaction
		ctx = ctx.WithMinGasPrices(sdk.NewDecCoins())

		// Verify that the transaction has zero fees
		feeTx, ok := tx.(sdk.FeeTx)
		if ok {
			fees := feeTx.GetFee()
			if !fees.IsZero() {
				// TSS messages should have zero fees
				ctx.Logger().Info("TSS transaction has non-zero fees, allowing anyway",
					"fees", fees.String())
			}
		}

		ctx.Logger().Debug("Exempting TSS transaction from gas fees",
			"msg_count", len(msgs))
	}

	return next(ctx, tx, simulate)
}
