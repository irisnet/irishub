package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/farm/types"
)

// CreatePoolFee returns the create pool fee
func (k Keeper) CreatePoolFee(ctx sdk.Context) sdk.Coin {
	return k.GetParams(ctx).PoolCreationFee
}

// MaxRewardCategories returns the maxRewardCategories
func (k Keeper) MaxRewardCategories(ctx sdk.Context) uint32 {
	return k.GetParams(ctx).MaxRewardCategories
}

// MaxRewardCategories returns the maxRewardCategories
func (k Keeper) TaxRate(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).TaxRate
}

// GetParams sets the farm module parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(types.ParamsKey))
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams sets the farm module parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}
	store.Set(types.ParamsKey, bz)

	return nil
}
