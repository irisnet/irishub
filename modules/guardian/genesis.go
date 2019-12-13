package guardian

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	profilersIterator := k.ProfilersIterator(ctx)
	defer profilersIterator.Close()
	var profilers []Guardian
	for ; profilersIterator.Valid(); profilersIterator.Next() {
		var profiler Guardian
		ModuleCdc.MustUnmarshalBinaryLengthPrefixed(profilersIterator.Value(), &profiler)
		profilers = append(profilers, profiler)
	}

	trusteesIterator := k.TrusteesIterator(ctx)
	defer trusteesIterator.Close()
	var trustees []Guardian
	for ; trusteesIterator.Valid(); trusteesIterator.Next() {
		var trustee Guardian
		ModuleCdc.MustUnmarshalBinaryLengthPrefixed(trusteesIterator.Value(), &trustee)
		trustees = append(trustees, trustee)
	}
	return NewGenesisState(profilers, trustees)
}
