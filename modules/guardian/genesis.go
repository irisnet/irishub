package guardian

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis stores genesis data
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
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
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	profilersIterator := k.ProfilersIterator(ctx)
	defer profilersIterator.Close()

	var profilers []Guardian
	var profiler Guardian
	for ; profilersIterator.Valid(); profilersIterator.Next() {
		ModuleCdc.MustUnmarshalBinaryLengthPrefixed(profilersIterator.Value(), &profiler)
		profilers = append(profilers, profiler)
	}

	trusteesIterator := k.TrusteesIterator(ctx)
	defer trusteesIterator.Close()

	var trustees []Guardian
	var trustee Guardian
	for ; trusteesIterator.Valid(); trusteesIterator.Next() {
		ModuleCdc.MustUnmarshalBinaryLengthPrefixed(trusteesIterator.Value(), &trustee)
		trustees = append(trustees, trustee)
	}

	return NewGenesisState(profilers, trustees)
}
