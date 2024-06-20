package v3

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/farm/types"
	"github.com/irisnet/irismod/types/exported"
)

// FarmKeeper defines a interface for SetParams function
type FarmKeeper interface {
	SetParams(ctx sdk.Context, params types.Params) error
}

// Migrate migrate the coinswap params from legacy x/params module to coinswap module
func Migrate(ctx sdk.Context, k FarmKeeper, legacySubspace exported.Subspace) error {
	var params types.Params
	legacySubspace.GetParamSet(ctx, &params)
	return k.SetParams(ctx, params)
}
