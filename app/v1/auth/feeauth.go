package auth

import (
	"errors"
	"fmt"
	"math"

	"github.com/irisnet/irishub/types"
	sdk "github.com/irisnet/irishub/types"
)

type FeeAuth struct {
	NativeFeeDenom string `json:"native_fee_denom"`
}

func NewFeeAuth(nativeFeeDenom string) FeeAuth {
	return FeeAuth{
		NativeFeeDenom: nativeFeeDenom,
	}
}

func InitialFeeAuth() FeeAuth {
	return NewFeeAuth(sdk.IrisAtto)
}

// NewFeePreprocessHandler creates a fee token preprocesser
func NewFeePreprocessHandler(fk FeeKeeper) types.FeePreprocessHandler {
	return func(ctx sdk.Context, tx sdk.Tx) sdk.Error {
		stdTx, ok := tx.(StdTx)
		if !ok {
			return sdk.ErrInternal("tx must be StdTx")
		}

		fee := getFee(stdTx.Fee.Amount)

		feeParams := fk.GetParamSet(ctx)

		return checkFee(feeParams, sdk.Coins{fee}, stdTx.Fee.Gas)
	}
}

// NewFeePreprocessHandler creates a fee token refund handler
func NewFeeRefundHandler(am AccountKeeper, fk FeeKeeper) types.FeeRefundHandler {
	return func(ctx sdk.Context, tx sdk.Tx, txResult sdk.Result) (actualCostFee sdk.Coin, err error) {
		txAccounts := GetSigners(ctx)
		// If this tx failed in anteHandler, txAccount length will be less than 1
		if len(txAccounts) < 1 {
			//panic("invalid transaction, should not reach here")
			return sdk.Coin{}, nil
		}
		firstAccount := txAccounts[0]

		stdTx, ok := tx.(StdTx)
		if !ok {
			return sdk.Coin{}, errors.New("transaction is not Stdtx")
		}
		// Refund process will also cost gas, but this is compensation for previous fee deduction.
		// It is not reasonable to consume users' gas. So the context gas is reset to transaction gas
		ctx = ctx.WithGasMeter(sdk.NewInfiniteGasMeter())

		fee := getFee(stdTx.Fee.Amount)

		//If all gas has been consumed, then there is no necessary to run fee refund process
		if txResult.GasWanted <= txResult.GasUsed {
			actualCostFee = fee
			return actualCostFee, nil
		}

		unusedGas := txResult.GasWanted - txResult.GasUsed
		refundCoin := sdk.NewCoin(fee.Denom,
			fee.Amount.Mul(sdk.NewInt(int64(unusedGas))).Div(sdk.NewInt(int64(txResult.GasWanted))))

		coins := am.GetAccount(ctx, firstAccount.GetAddress()).GetCoins() // consume gas
		err = firstAccount.SetCoins(coins.Add(sdk.Coins{refundCoin}))
		if err != nil {
			return sdk.Coin{}, err
		}

		// set mem regexp
		regexp := am.GetAccount(ctx, firstAccount.GetAddress()).GetMemoRegexp()
		firstAccount.SetMemoRegexp(regexp)

		am.SetAccount(ctx, firstAccount)
		fk.RefundCollectedFees(ctx, sdk.Coins{refundCoin})

		actualCostFee = sdk.NewCoin(fee.Denom, fee.Amount.Sub(refundCoin.Amount))
		return actualCostFee, nil
	}
}

func getFee(coins sdk.Coins) sdk.Coin {
	if coins == nil || coins.Empty() {
		return sdk.NewCoin(sdk.IrisAtto, sdk.ZeroInt())
	}
	return sdk.NewCoin(sdk.IrisAtto, coins.AmountOf(sdk.IrisAtto))
}

func checkFee(params Params, coins sdk.Coins, gasLimit uint64) sdk.Error {
	if gasLimit == 0 || int64(gasLimit) < 0 {
		return sdk.ErrInvalidGas(fmt.Sprintf("gaslimit %d should be positive and no more than %d", gasLimit, math.MaxInt64))
	}

	threshold := params.GasPriceThreshold
	equivalentTotalFee := coins[0].Amount
	gasPrice := equivalentTotalFee.Div(sdk.NewInt(int64(gasLimit)))

	if gasPrice.LT(threshold) {
		recommendFee := (sdk.NewInt(int64(gasLimit))).Mul(threshold)
		return sdk.ErrGasPriceTooLow(fmt.Sprintf("insufficient fee, gasPrice = fee / gasLimit(default 50000). The gasPrice(%s%s/Gas) cannot be less than %s%s. Recommended fee: %s%s", gasPrice.String(), sdk.IrisAtto, threshold.String(), sdk.IrisAtto, recommendFee, sdk.IrisAtto))
	}

	return nil
}
