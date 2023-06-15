package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v1 "github.com/irisnet/irismod/modules/token/types/v1"
)

// GetParam returns token params from the global param store
func (k Keeper) GetParam(ctx sdk.Context) v1.Params {
	var p v1.Params
	k.paramSpace.GetParamSet(ctx, &p)
	return p
}

// SetParam sets token params to the global param store
func (k Keeper) SetParam(ctx sdk.Context, params v1.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
