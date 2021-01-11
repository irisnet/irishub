package coinswap

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/coinswap/keeper"
	"github.com/irisnet/irismod/modules/coinswap/types"
)

// InitGenesis new coinswap genesis
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	if err := types.ValidateGenesis(data); err != nil {
		panic(fmt.Errorf("panic for ValidateGenesis,%v", err))
	}
	k.SetParams(ctx, data.Params)
	k.SetStandardDenom(ctx, data.StandardDenom)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	return types.NewGenesisState(keeper.GetParams(ctx), keeper.GetStandardDenom(ctx))
}
