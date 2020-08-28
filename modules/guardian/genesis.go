package guardian

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/guardian/keeper"
	"github.com/irisnet/irishub/modules/guardian/types"
)

// InitGenesis stores genesis data
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) {
	// Add profilers
	for _, profiler := range data.Profilers {
		keeper.AddProfiler(ctx, profiler)
	}
	// Add trustees
	for _, trustee := range data.Trustees {
		keeper.AddTrustee(ctx, trustee)
	}
}

// ExportGenesis outputs genesis data
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	var profilers []types.Guardian
	k.IterateProfilers(
		ctx,
		func(profiler types.Guardian) bool {
			profilers = append(profilers, profiler)
			return false
		},
	)
	var trustees []types.Guardian
	k.IterateTrustees(
		ctx,
		func(trustee types.Guardian) bool {
			trustees = append(trustees, trustee)
			return false
		},
	)

	return types.NewGenesisState(profilers, trustees)
}
