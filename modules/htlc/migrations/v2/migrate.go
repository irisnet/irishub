package v2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"mods.irisnet.org/modules/htlc/types"
)

// HTLCKeeper defines a interface for SetParams function
type HTLCKeeper interface {
	SetParams(ctx sdk.Context, params types.Params) error
}

// Migrate migrate the htlc params from legacy x/params module to htlc module
func Migrate(ctx sdk.Context, k HTLCKeeper, legacySubspace types.Subspace) error {
	var params types.Params
	legacySubspace.GetParamSet(ctx, &params)
	return k.SetParams(ctx, params)
}
