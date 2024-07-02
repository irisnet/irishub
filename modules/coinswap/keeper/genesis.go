package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"mods.irisnet.org/modules/coinswap/types"
)

// InitGenesis initializes the coinswap module's state from a given genesis state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	if err := types.ValidateGenesis(genState); err != nil {
		panic(fmt.Errorf("panic for ValidateGenesis,%w", err))
	}
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(fmt.Errorf("panic for SetParams,%w", err))
	}
	k.SetStandardDenom(ctx, genState.StandardDenom)
	k.setSequence(ctx, genState.Sequence)
	for _, pool := range genState.Pool {
		poolCopy := pool // Create a copy of the pool variable
		k.setPool(ctx, &poolCopy)
	}
}

// ExportGenesis returns the coinswap module's genesis state.
func (k Keeper) ExportGenesis(ctx sdk.Context) types.GenesisState {
	return types.GenesisState{
		Params:        k.GetParams(ctx),
		StandardDenom: k.GetStandardDenom(ctx),
		Pool:          k.GetAllPools(ctx),
		Sequence:      k.getSequence(ctx),
	}
}
