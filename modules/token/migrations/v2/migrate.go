package v2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v1 "github.com/irisnet/irismod/modules/token/types/v1"
	"github.com/irisnet/irismod/types/exported"
)

// TokenKeeper defines a interface for SetParams function
type TokenKeeper interface {
	SetParams(ctx sdk.Context, params v1.Params) error
}

// Migrate migrate the service params from legacy x/params module to htlc module
func Migrate(ctx sdk.Context, k TokenKeeper, legacySubspace exported.Subspace) error {
	var params v1.Params
	legacySubspace.GetParamSet(ctx, &params)
	return k.SetParams(ctx, params)
}
