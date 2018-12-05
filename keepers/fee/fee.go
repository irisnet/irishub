package fee

import (
	"errors"
	"fmt"
	codec "github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/keepers/params"
	sdk "github.com/irisnet/irishub/types"
)

var (
	collectedFeesKey           = []byte("collectedFees")
	nativeFeeTokenKey          = []byte("feeTokenNative")
	nativeGasPriceThresholdKey = []byte("feeTokenGasPriceThreshold")
	//	FeeExchangeRatePrefix = "feeToken/exchangeRate/"	//  key = gov/feeToken/exchangeRate/<denomination>, rate = BigInt(value)/10^9
	//	RatePrecision = int64(1000000000) //10^9
)

// This FeeCollectionKeeper handles collection of fees in the anteHandler
// and setting of MinFees for different fee tokens
type FeeCollectionKeeper struct {

	// The (unexposed) key used to access the fee store from the Context.
	key sdk.StoreKey

	// The codec codec for binary encoding/decoding of accounts.
	cdc *codec.Codec
}

func NewFeeCollectionKeeper(cdc *codec.Codec, key sdk.StoreKey) FeeCollectionKeeper {
	return FeeCollectionKeeper{
		key: key,
		cdc: cdc,
	}
}

// retrieves the collected fee pool
func (fck FeeCollectionKeeper) GetCollectedFees(ctx sdk.Context) sdk.Coins {
	store := ctx.KVStore(fck.key)
	bz := store.Get(collectedFeesKey)
	if bz == nil {
		return sdk.Coins{}
	}

	feePool := &(sdk.Coins{})
	fck.cdc.MustUnmarshalBinaryLengthPrefixed(bz, feePool)
	return *feePool
}

func (fck FeeCollectionKeeper) setCollectedFees(ctx sdk.Context, coins sdk.Coins) {
	bz := fck.cdc.MustMarshalBinaryLengthPrefixed(coins)
	store := ctx.KVStore(fck.key)
	store.Set(collectedFeesKey, bz)
}

// add to the fee pool
func (fck FeeCollectionKeeper) AddCollectedFees(ctx sdk.Context, coins sdk.Coins) sdk.Coins {
	newCoins := fck.GetCollectedFees(ctx).Plus(coins)
	fck.setCollectedFees(ctx, newCoins)

	return newCoins
}

// RefundCollectedFees deducts fees from fee collector
func (fck FeeCollectionKeeper) RefundCollectedFees(ctx sdk.Context, coins sdk.Coins) sdk.Coins {
	newCoins := fck.GetCollectedFees(ctx).Minus(coins)
	if !newCoins.IsNotNegative() {
		panic("fee collector contains negative coins")
	}
	fck.setCollectedFees(ctx, newCoins)
	return newCoins
}

// clear the fee pool
func (fck FeeCollectionKeeper) ClearCollectedFees(ctx sdk.Context) {
	fck.setCollectedFees(ctx, sdk.Coins{})
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

func (fck FeeManager) GetNativeFeeToken(ctx sdk.Context, coins sdk.Coins) sdk.Coin {
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

func (fck FeeManager) FeePreprocess(ctx sdk.Context, coins sdk.Coins, gasLimit int64) sdk.Error {
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
