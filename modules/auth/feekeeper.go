package auth

import (
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/params"
	sdk "github.com/irisnet/irishub/types"
)

var (
	collectedFeesKey = []byte("collectedFees")
	feeAuthKey       = []byte("feeAuth")
)

// This FeeKeeper handles collection of fees in the anteHandler
// and setting of MinFees for different fee tokens
type FeeKeeper struct {

	// The (unexposed) key used to access the fee store from the Context.
	storeKey sdk.StoreKey

	// The codec codec for binary encoding/decoding of accounts.
	cdc *codec.Codec

	paramSpace params.Subspace
}

func NewFeeKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace) FeeKeeper {
	return FeeKeeper{
		storeKey:   key,
		cdc:        cdc,
		paramSpace: paramSpace,
	}
}

// retrieves the collected fee pool
func (fk FeeKeeper) GetCollectedFees(ctx sdk.Context) sdk.Coins {
	store := ctx.KVStore(fk.storeKey)
	bz := store.Get(collectedFeesKey)
	if bz == nil {
		return sdk.Coins{}
	}

	feePool := &(sdk.Coins{})
	fk.cdc.MustUnmarshalBinaryLengthPrefixed(bz, feePool)
	return *feePool
}

func (fk FeeKeeper) setCollectedFees(ctx sdk.Context, coins sdk.Coins) {
	bz := fk.cdc.MustMarshalBinaryLengthPrefixed(coins)
	store := ctx.KVStore(fk.storeKey)
	store.Set(collectedFeesKey, bz)
}

// add to the fee pool
func (fk FeeKeeper) AddCollectedFees(ctx sdk.Context, coins sdk.Coins) sdk.Coins {
	newCoins := fk.GetCollectedFees(ctx).Plus(coins)
	fk.setCollectedFees(ctx, newCoins)

	return newCoins
}

// RefundCollectedFees deducts fees from fee collector
func (fk FeeKeeper) RefundCollectedFees(ctx sdk.Context, coins sdk.Coins) sdk.Coins {
	newCoins := fk.GetCollectedFees(ctx).Minus(coins)
	if !newCoins.IsNotNegative() {
		panic("fee collector contains negative coins")
	}
	fk.setCollectedFees(ctx, newCoins)
	return newCoins
}

func (fk FeeKeeper) ClearCollectedFees(ctx sdk.Context) {
	fk.setCollectedFees(ctx, sdk.Coins{})
}

func (fk FeeKeeper) GetFeeAuth(ctx sdk.Context) (feeAuth FeeAuth) {
	store := ctx.KVStore(fk.storeKey)
	b := store.Get(feeAuthKey)
	if b == nil {
		panic("Stored fee pool should not have been nil")
	}
	fk.cdc.MustUnmarshalBinaryLengthPrefixed(b, &feeAuth)
	return
}

func (fk FeeKeeper) SetFeeAuth(ctx sdk.Context, feeAuth FeeAuth) {
	store := ctx.KVStore(fk.storeKey)
	b := fk.cdc.MustMarshalBinaryLengthPrefixed(feeAuth)
	store.Set(feeAuthKey, b)
}

func (fk FeeKeeper) GetParamSet(ctx sdk.Context) Params {
	var feeParams Params
	fk.paramSpace.GetParamSet(ctx, &feeParams)
	return feeParams
}

func (fk FeeKeeper) SetParamSet(ctx sdk.Context, feeParams Params) {
	fk.paramSpace.SetParamSet(ctx, &feeParams)
}
