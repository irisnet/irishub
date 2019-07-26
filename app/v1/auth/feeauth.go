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

func ValidateFee(auth FeeAuth, collectedFee sdk.Coins) error {
	if !collectedFee.IsValid() {
		return sdk.ErrInvalidCoins("")
	}
	if len(auth.NativeFeeDenom) == 0 {
		return sdk.ErrInvalidFeeDenom("")
	}
	return nil
}

// NewFeePreprocessHandler creates a fee token preprocesser
func NewFeePreprocessHandler(fk FeeKeeper) types.FeePreprocessHandler {
	return func(ctx sdk.Context, tx sdk.Tx) sdk.Error {
		stdTx, ok := tx.(StdTx)
		if !ok {
			return sdk.ErrInternal("tx must be StdTx")
		}

		fa := fk.GetFeeAuth(ctx)
		feeParams := fk.GetParamSet(ctx)

		totalNativeFee := fa.getNativeFeeToken(ctx, stdTx.Fee.Amount)

		return fa.feePreprocess(ctx, feeParams, sdk.Coins{totalNativeFee}, stdTx.Fee.Gas)
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

		fm := fk.GetFeeAuth(ctx)
		totalNativeFee := fm.getNativeFeeToken(ctx, stdTx.Fee.Amount)

		//If all gas has been consumed, then there is no necessary to run fee refund process
		if txResult.GasWanted <= txResult.GasUsed {
			actualCostFee = totalNativeFee
			return actualCostFee, nil
		}

		unusedGas := txResult.GasWanted - txResult.GasUsed
		refundCoin := sdk.NewCoin(totalNativeFee.Denom,
			totalNativeFee.Amount.Mul(sdk.NewInt(int64(unusedGas))).Div(sdk.NewInt(int64(txResult.GasWanted))))

		coins := am.GetAccount(ctx, firstAccount.GetAddress()).GetCoins() // consume gas
		err = firstAccount.SetCoins(coins.Plus(sdk.Coins{refundCoin}))
		if err != nil {
			return sdk.Coin{}, err
		}

		am.SetAccount(ctx, firstAccount)
		fk.RefundCollectedFees(ctx, sdk.Coins{refundCoin})

		actualCostFee = sdk.NewCoin(totalNativeFee.Denom, totalNativeFee.Amount.Sub(refundCoin.Amount))
		return actualCostFee, nil
	}
}

func (fa FeeAuth) getNativeFeeToken(ctx sdk.Context, coins sdk.Coins) sdk.Coin {
	nativeFeeToken := fa.NativeFeeDenom
	for _, coin := range coins {
		if coin.Denom == nativeFeeToken {
			if coin.Amount.BigInt() == nil {
				return sdk.NewCoin(coin.Denom, sdk.ZeroInt())
			}
			return coin
		}
	}
	return sdk.NewCoin("", sdk.ZeroInt())
}

func (fa FeeAuth) feePreprocess(ctx sdk.Context, params Params, coins sdk.Coins, gasLimit uint64) sdk.Error {
	if gasLimit == 0 || int64(gasLimit) < 0 {
		return sdk.ErrInvalidGas(fmt.Sprintf("gaslimit %d should be positive and no more than %d", gasLimit, math.MaxInt64))
	}
	nativeFeeToken := fa.NativeFeeDenom
	threshold := params.GasPriceThreshold

	if len(coins) < 1 || coins[0].Denom != nativeFeeToken {
		return sdk.ErrInvalidTxFee(fmt.Sprintf("fee is required to be specified with %s", nativeFeeToken))
	}

	equivalentTotalFee := coins[0].Amount
	gasPrice := equivalentTotalFee.Div(sdk.NewInt(int64(gasLimit)))
	if gasPrice.LT(threshold) {
		recommendFee := (sdk.NewInt(int64(gasLimit))).Mul(threshold)
		return sdk.ErrGasPriceTooLow(fmt.Sprintf("insufficient fee, gasPrice = fee / gasLimit(default 50000). The gasPrice(%s%s/Gas) cannot be less than %s%s. Recommended fee: %s%s", gasPrice.String(), nativeFeeToken, threshold.String(), nativeFeeToken, recommendFee, nativeFeeToken))
	}
	return nil
}
