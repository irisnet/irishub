package mint

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis new mint genesis
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	if err := ValidateGenesis(data); err != nil {
		panic(fmt.Errorf("failed to initialize mint genesis state: %s", err.Error()))
	}
	keeper.SetMinter(ctx, data.Minter)
	keeper.SetParamSet(ctx, data.Params)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	minter := keeper.GetMinter(ctx)
	params := keeper.GetParamSet(ctx)
	return NewGenesisState(minter, params)
}
