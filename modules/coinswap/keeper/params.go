package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"

	"github.com/irisnet/irismod/modules/coinswap/types"
)

// GetParams sets the coinswap module parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(types.ParamsKey))
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams sets the parameters for the coinswap module.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}
	store.Set([]byte(types.ParamsKey), bz)

	return nil
}

// SetStandardDenom sets the standard denom for the coinswap module.
func (k Keeper) SetStandardDenom(ctx sdk.Context, denom string) {
	store := ctx.KVStore(k.storeKey)
	denomWrap := gogotypes.StringValue{Value: denom}
	bz := k.cdc.MustMarshal(&denomWrap)
	store.Set(types.KeyStandardDenom, bz)
}

// GetStandardDenom returns the standard denom of the coinswap module.
func (k Keeper) GetStandardDenom(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyStandardDenom)

	var denomWrap = gogotypes.StringValue{}
	k.cdc.MustUnmarshal(bz, &denomWrap)
	return denomWrap.Value
}
