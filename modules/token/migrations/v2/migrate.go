package v2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"mods.irisnet.org/token/types"
	v1 "mods.irisnet.org/token/types/v1"
)

// TokenKeeper defines a interface for SetParams function
type TokenKeeper interface {
	SetParams(ctx sdk.Context, params v1.Params) error
}

// Migrate migrate the service params from legacy x/params module to htlc module
func Migrate(ctx sdk.Context, k TokenKeeper, legacySubspace types.Subspace) error {
	var params v1.Params
	legacySubspace.GetParamSet(ctx, &params)
	return k.SetParams(ctx, params)
}
