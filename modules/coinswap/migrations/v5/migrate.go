package v5

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"mods.irisnet.org/modules/coinswap/types"
)

// CoinswapKeeper defines a interface for SetParams function
type CoinswapKeeper interface {
	SetParams(ctx sdk.Context, params types.Params) error
}

// Migrate migrate the coinswap params from legacy x/params module to coinswap module
func Migrate(ctx sdk.Context, k CoinswapKeeper, legacySubspace types.Subspace) error {
	var params types.Params
	legacySubspace.GetParamSet(ctx, &params)
	return k.SetParams(ctx, params)
}
