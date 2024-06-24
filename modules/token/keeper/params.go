package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"mods.irisnet.org/token/types"
	v1 "mods.irisnet.org/token/types/v1"
)

// GetParams sets the token module parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params v1.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(types.PrefixParamsKey))
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams sets the token module parameters.
func (k Keeper) SetParams(ctx sdk.Context, params v1.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}
	store.Set(types.PrefixParamsKey, bz)

	return nil
}

// ERC20Enabled returns true if ERC20 is enabled
func (k Keeper) ERC20Enabled(ctx sdk.Context) bool {
	params := k.GetParams(ctx)
	return params.EnableErc20
}
