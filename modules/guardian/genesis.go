package guardian

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/guardian/keeper"
	"github.com/irisnet/irishub/modules/guardian/types"
)

// InitGenesis stores genesis data
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) {
	// Add supers
	for _, super := range data.Supers {
		keeper.AddSuper(ctx, super)
	}
}

// ExportGenesis outputs genesis data
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	var supers []types.Super
	k.IterateSupers(
		ctx,
		func(super types.Super) bool {
			supers = append(supers, super)
			return false
		},
	)

	return types.NewGenesisState(supers)
}
