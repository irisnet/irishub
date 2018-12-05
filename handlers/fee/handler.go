package fee

import (
	"errors"
	"fmt"
	"github.com/irisnet/irishub/handlers/auth"
	keeper "github.com/irisnet/irishub/keepers/bank"
	"github.com/irisnet/irishub/keepers/fee"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/types/bank"
	"runtime/debug"
)

// NewFeePreprocessHandler creates a fee token preprocesser
func NewFeePreprocessHandler(fm fee.FeeManager) sdk.FeePreprocessHandler {
	return func(ctx sdk.Context, tx sdk.Tx) error {
		stdTx, ok := tx.(bank.StdTx)
		if !ok {
			return sdk.ErrInternal("tx must be StdTx")
		}
		totalNativeFee := fm.GetNativeFeeToken(ctx, stdTx.Fee.Amount)
		return fm.FeePreprocess(ctx, sdk.Coins{totalNativeFee}, stdTx.Fee.Gas)
	}
}

// NewFeePreprocessHandler creates a fee token refund handler
func NewFeeRefundHandler(am keeper.AccountKeeper, fck fee.FeeCollectionKeeper, fm fee.FeeManager) sdk.FeeRefundHandler {
	return func(ctx sdk.Context, tx sdk.Tx, txResult sdk.Result) (actualCostFee sdk.Coin, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("encountered panic error during fee refund, recovered: %v\nstack:\n%v", r, string(debug.Stack()))
			}
		}()

		txAccounts := auth.GetSigners(ctx)
		// If this tx failed in anteHandler, txAccount length will be less than 1
		if len(txAccounts) < 1 {
			//panic("invalid transaction, should not reach here")
			return sdk.Coin{}, nil
		}
		firstAccount := txAccounts[0]

		stdTx, ok := tx.(bank.StdTx)
		if !ok {
			return sdk.Coin{}, errors.New("transaction is not Stdtx")
		}
		// Refund process will also cost gas, but this is compensation for previous fee deduction.
		// It is not reasonable to consume users' gas. So the context gas is reset to transaction gas
		ctx = ctx.WithGasMeter(sdk.NewInfiniteGasMeter())

		totalNativeFee := fm.GetNativeFeeToken(ctx, stdTx.Fee.Amount)

		//If all gas has been consumed, then there is no necessary to run fee refund process
		if txResult.GasWanted <= txResult.GasUsed {
			actualCostFee = totalNativeFee
			return actualCostFee, nil
		}

		unusedGas := txResult.GasWanted - txResult.GasUsed
		refundCoin := sdk.Coin{
			Denom:  totalNativeFee.Denom,
			Amount: totalNativeFee.Amount.Mul(sdk.NewInt(unusedGas)).Div(sdk.NewInt(txResult.GasWanted)),
		}
		coins := am.GetAccount(ctx, firstAccount.GetAddress()).GetCoins() // consume gas
		err = firstAccount.SetCoins(coins.Plus(sdk.Coins{refundCoin}))
		if err != nil {
			return sdk.Coin{}, err
		}

		am.SetAccount(ctx, firstAccount)
		fck.RefundCollectedFees(ctx, sdk.Coins{refundCoin})

		actualCostFee = sdk.Coin{
			Denom:  totalNativeFee.Denom,
			Amount: totalNativeFee.Amount.Sub(refundCoin.Amount),
		}
		return actualCostFee, nil
	}
}
