package v2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/types/exported"

	"github.com/irisnet/irishub/v2/modules/mint/types"
)

// MintKeeper defines a interface for SetParams function
type MintKeeper interface {
	SetParams(ctx sdk.Context, params types.Params) error
}

// Migrate migrate the coinswap params from legacy x/params module to mint module
func Migrate(ctx sdk.Context, k MintKeeper, legacySubspace exported.Subspace) error {
	var params types.Params
	legacySubspace.GetParamSet(ctx, &params)
	return k.SetParams(ctx, params)
}
