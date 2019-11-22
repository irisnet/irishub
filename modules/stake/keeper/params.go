package keeper

import (
	"time"

	"github.com/irisnet/irishub/modules/params"
	"github.com/irisnet/irishub/modules/stake/types"
	sdk "github.com/irisnet/irishub/types"
)

// ParamTable for stake module
func ParamTypeTable() params.TypeTable {
	return params.NewTypeTable().RegisterParamSet(&types.Params{})
}

// UnbondingTime
func (k Keeper) UnbondingTime(ctx sdk.Context) (res time.Duration) {
	k.paramstore.Get(ctx, types.KeyUnbondingTime, &res)
	return
}

// MaxValidators - Maximum number of validators
func (k Keeper) MaxValidators(ctx sdk.Context) (res uint16) {
	k.paramstore.Get(ctx, types.KeyMaxValidators, &res)
	return
}

// Get all parameteras as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (res types.Params) {
	res.UnbondingTime = k.UnbondingTime(ctx)
	res.MaxValidators = k.MaxValidators(ctx)
	return
}

// set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
