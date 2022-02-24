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
	k.paramSpace.Get(ctx, types.KeyPoolCreationFee, &fee)
	return
}

// MaxRewardCategories returns the maxRewardCategories
func (k Keeper) MaxRewardCategories(ctx sdk.Context) (maxRewardCategories uint32) {
	k.paramSpace.Get(ctx, types.KeyMaxRewardCategories, &maxRewardCategories)
	return
}

// MaxRewardCategories returns the maxRewardCategories
func (k Keeper) TaxRate(ctx sdk.Context) (taxRate sdk.Dec) {
	k.paramSpace.Get(ctx, types.KeyTaxRate, &taxRate)
	return
}

// SetParams sets the params to the store
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// GetParams gets all parameteras as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.CreatePoolFee(ctx),
		k.MaxRewardCategories(ctx),
		k.TaxRate(ctx),
	)
}
