package auth

import (
	"errors"
	"fmt"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/params"
	"github.com/irisnet/irishub/types"
	"math"
)

var (
	nativeFeeTokenKey          = []byte("feeTokenNative")
	nativeGasPriceThresholdKey = []byte("feeTokenGasPriceThreshold")
	//	FeeExchangeRatePrefix = "feeToken/exchangeRate/"	//  key = gov/feeToken/exchangeRate/<denomination>, rate = BigInt(value)/10^9
	//	RatePrecision = int64(1000000000) //10^9
)

// NewFeePreprocessHandler creates a fee token preprocesser
func NewFeePreprocessHandler(fm FeeManager) types.FeePreprocessHandler {
	return func(ctx sdk.Context, tx sdk.Tx) error {
		stdTx, ok := tx.(StdTx)
		if !ok {
			return sdk.ErrInternal("tx must be StdTx")
		}
		totalNativeFee := fm.getNativeFeeToken(ctx, stdTx.Fee.Amount)
		return fm.feePreprocess(ctx, sdk.Coins{totalNativeFee}, stdTx.Fee.Gas)
	}
}

// NewFeePreprocessHandler creates a fee token refund handler
func NewFeeRefundHandler(am AccountKeeper, fck FeeCollectionKeeper, fm FeeManager) types.FeeRefundHandler {
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

		totalNativeFee := fm.getNativeFeeToken(ctx, stdTx.Fee.Amount)

		//If all gas has been consumed, then there is no necessary to run fee refund process
		if txResult.GasWanted <= txResult.GasUsed {
			actualCostFee = totalNativeFee
			return actualCostFee, nil
		}

		unusedGas := txResult.GasWanted - txResult.GasUsed
		refundCoin := sdk.Coin{
			Denom:  totalNativeFee.Denom,
			Amount: totalNativeFee.Amount.Mul(sdk.NewInt(int64(unusedGas))).Div(sdk.NewInt(int64(txResult.GasWanted))),
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

// Type declaration for parameters
func ParamTypeTable() params.TypeTable {
	return params.NewTypeTable(
		nativeFeeTokenKey, "",
		nativeGasPriceThresholdKey, "",
	)
}

// FeeManager do fee tokens preprocess according to fee token configuration
type FeeManager struct {
	// The reference to the Paramstore to get and set gov specific params
	paramSpace params.Subspace
}

func NewFeeManager(paramSpace params.Subspace) FeeManager {
	return FeeManager{
		paramSpace: paramSpace.WithTypeTable(ParamTypeTable()),
	}
}

func (fck FeeManager) getNativeFeeToken(ctx sdk.Context, coins sdk.Coins) sdk.Coin {
	var nativeFeeToken string
	fck.paramSpace.Get(ctx, nativeFeeTokenKey, &nativeFeeToken)
	for _, coin := range coins {
		if coin.Denom == nativeFeeToken {
			if coin.Amount.BigInt() == nil {
				return sdk.Coin{
					Denom:  coin.Denom,
					Amount: sdk.ZeroInt(),
				}
			}
			return coin
		}
	}
	return sdk.Coin{
		Denom:  "",
		Amount: sdk.ZeroInt(),
	}
}

func (fck FeeManager) feePreprocess(ctx sdk.Context, coins sdk.Coins, gasLimit uint64) sdk.Error {
	if gasLimit == 0 || int64(gasLimit) < 0 {
		return sdk.ErrInvalidGas(fmt.Sprintf("gaslimit %d should be positive and no more than %d", gasLimit, math.MaxInt64))
	}
	var nativeFeeToken string
	fck.paramSpace.Get(ctx, nativeFeeTokenKey, &nativeFeeToken)

	var nativeGasPriceThreshold string
	fck.paramSpace.Get(ctx, nativeGasPriceThresholdKey, &nativeGasPriceThreshold)

	threshold, ok := sdk.NewIntFromString(nativeGasPriceThreshold)
	if !ok {
		panic(errors.New("failed to parse gas price from string"))
	}

	if len(coins) < 1 || coins[0].Denom != nativeFeeToken {
		return sdk.ErrInvalidTxFee(fmt.Sprintf("no native fee token, expected native token is %s", nativeFeeToken))
	}
	/*
		equivalentTotalFee := sdk.NewInt(0)
		for _,coin := range coins {
			if coin.Denom != nativeFeeToken {
				exchangeRateKey := FeeExchangeRatePrefix + coin.Denom
				rateString, err := fck.getter.GetString(ctx, exchangeRateKey)
				if err != nil {
					continue
				}
				rate, ok := sdk.NewIntFromString(rateString)
				if !ok {
					panic(errors.New("failed to parse rate from string"))
				}
				equivalentFee := rate.Mul(coin.Amount).Div(sdk.NewInt(RatePrecision))
				equivalentTotalFee = equivalentTotalFee.Add(equivalentFee)

			} else {
				equivalentTotalFee = equivalentTotalFee.Add(coin.Amount)
			}
		}
	*/
	equivalentTotalFee := coins[0].Amount
	gasPrice := equivalentTotalFee.Div(sdk.NewInt(int64(gasLimit)))
	if gasPrice.LT(threshold) {
		return sdk.ErrGasPriceTooLow(fmt.Sprintf("equivalent gas price (%s%s) is less than threshold (%s%s)", gasPrice.String(), nativeFeeToken, threshold.String(), nativeFeeToken))
	}
	return nil
}
