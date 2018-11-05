package baseapp

import (
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/irisnet/irishub/types"
	"runtime/debug"
)

var (
	nativeFeeTokenKey          = []byte("feeTokenNative")
	nativeGasPriceThresholdKey = []byte("feeTokenGasPriceThreshold")
	//	FeeExchangeRatePrefix = "feeToken/exchangeRate/"	//  key = gov/feeToken/exchangeRate/<denomination>, rate = BigInt(value)/10^9
	//	RatePrecision = int64(1000000000) //10^9
)

// NewFeePreprocessHandler creates a fee token preprocess handler
func NewFeePreprocessHandler(fm FeeManager) types.FeePreprocessHandler {
	return func(ctx sdk.Context, tx sdk.Tx) error {
		stdTx, ok := tx.(auth.StdTx)
		if !ok {
			return sdk.ErrInternal("tx must be StdTx")
		}
		fee := auth.StdFee{
			Gas:    stdTx.Fee.Gas,
			Amount: sdk.Coins{fm.getNativeFeeToken(ctx, stdTx.Fee.Amount)},
		}
		return fm.feePreprocess(ctx, fee.Amount, fee.Gas)
	}
}

// NewFeePreprocessHandler creates a fee token refund handler
func NewFeeRefundHandler(am auth.AccountKeeper, fck auth.FeeCollectionKeeper, fm FeeManager) types.FeeRefundHandler {
	return func(ctx sdk.Context, tx sdk.Tx, txResult sdk.Result) (refundResult sdk.Coin, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("encountered panic error during fee refund, recovered: %v\nstack:\n%v", r, string(debug.Stack()))
			}
		}()

		txAccounts := auth.GetSigners(ctx)
		// If this tx failed in anteHandler, txAccount length will be less than 1
		if len(txAccounts) < 1 {
			panic("invalid transaction, should not reach here")
		}
		firstAccount := txAccounts[0]

		stdTx, ok := tx.(auth.StdTx)
		if !ok {
			return sdk.Coin{}, errors.New("transaction is not Stdtx")
		}
		// Refund process will also cost gas, but this is compensation for previous fee deduction.
		// It is not reasonable to consume users' gas. So the context gas is reset to transaction gas
		ctx = ctx.WithGasMeter(sdk.NewInfiniteGasMeter())

		fee := auth.StdFee{
			Gas:    stdTx.Fee.Gas,
			Amount: sdk.Coins{fm.getNativeFeeToken(ctx, stdTx.Fee.Amount)}, // consume gas
		}

		//If all gas has been consumed, then there is no necessary to run fee refund process
		if txResult.GasWanted <= txResult.GasUsed {
			refundResult = sdk.Coin{
				Denom:  fee.Amount[0].Denom,
				Amount: fee.Amount[0].Amount,
			}
			return refundResult, nil
		}

		unusedGas := txResult.GasWanted - txResult.GasUsed
		var refundCoins sdk.Coins
		for _, coin := range fee.Amount {
			newCoin := sdk.Coin{
				Denom:  coin.Denom,
				Amount: coin.Amount.Mul(sdk.NewInt(unusedGas)).Div(sdk.NewInt(txResult.GasWanted)),
			}
			refundCoins = append(refundCoins, newCoin)
		}
		coins := am.GetAccount(ctx, firstAccount.GetAddress()).GetCoins() // consume gas
		err = firstAccount.SetCoins(coins.Plus(refundCoins))
		if err != nil {
			return sdk.Coin{}, err
		}

		am.SetAccount(ctx, firstAccount)          // consume gas
		fck.RefundCollectedFees(ctx, refundCoins) // consume gas
		// There must be just one fee token
		refundResult = sdk.Coin{
			Denom:  fee.Amount[0].Denom,
			Amount: fee.Amount[0].Amount.Mul(sdk.NewInt(txResult.GasUsed)).Div(sdk.NewInt(txResult.GasWanted)),
		}

		return refundResult, nil
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

func (fck FeeManager) feePreprocess(ctx sdk.Context, coins sdk.Coins, gasLimit int64) sdk.Error {
	if gasLimit <= 0 {
		return sdk.ErrInternal(fmt.Sprintf("gaslimit %d should be larger than 0", gasLimit))
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
		return sdk.ErrInvalidCoins(fmt.Sprintf("no native fee token, expected native token is %s", nativeFeeToken))
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
	gasPrice := equivalentTotalFee.Div(sdk.NewInt(gasLimit))
	if gasPrice.LT(threshold) {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("equivalent gas price (%s%s) is less than threshold (%s%s)", gasPrice.String(), nativeFeeToken, threshold.String(), nativeFeeToken))
	}
	return nil
}

type FeeGenesisStateConfig struct {
	FeeTokenNative    string `json:"fee_token_native"`
	GasPriceThreshold int64  `json:"gas_price_threshold"`
}

func InitGenesis(ctx sdk.Context, ps FeeManager, data FeeGenesisStateConfig) {
	ps.paramSpace.Set(ctx, nativeFeeTokenKey, data.FeeTokenNative)
	ps.paramSpace.Set(ctx, nativeGasPriceThresholdKey, sdk.NewInt(data.GasPriceThreshold).String())
}
