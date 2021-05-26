package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/irisnet/irismod/modules/farm/types"
)

// ParamKeyTable for farm module
func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&types.Params{})
}

// CreatePoolFee returns the create pool fee
func (k Keeper) CreatePoolFee(ctx sdk.Context) (fee sdk.Coin) {
	k.paramSpace.Get(ctx, types.KeyCreatePoolFee, &fee)
	return
}

// MaxRewardCategoryN returns the maxRewardCategoryN
func (k Keeper) MaxRewardCategoryN(ctx sdk.Context) (maxRewardCategoryN uint32) {
	k.paramSpace.Get(ctx, types.KeyMaxRewardCategoryN, &maxRewardCategoryN)
	return
}

// SetParams sets the params to the store
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) (fee sdk.Coin) {
	k.paramSpace.SetParamSet(ctx, &params)
	return
}

// GetParams gets all parameteras as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.CreatePoolFee(ctx),
		k.MaxRewardCategoryN(ctx),
	)
}
