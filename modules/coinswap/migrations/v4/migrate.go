package v4

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/irisnet/irismod/modules/coinswap/types"
)

// CoinswapKeeper defines a interface for SetParams function
type CoinswapKeeper interface {
	SetParams(ctx sdk.Context, params types.Params) error
}

// Migrate migrate the coinswap params from legacy x/params module to coinswap module
func Migrate(ctx sdk.Context, k CoinswapKeeper, legacySubspace paramstypes.Subspace) error {
	var params types.Params
	legacySubspace.GetParamSet(ctx, &params)
	return k.SetParams(ctx, params)
}
