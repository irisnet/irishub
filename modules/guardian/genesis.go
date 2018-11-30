package guardian

import (
	sdk "github.com/irisnet/irishub/types"
)

// GenesisState - all profiling state that must be provided at genesis
type GenesisState struct {
	Profilers []Profiler `json:"profilers"`
}

func NewGenesisState(profilers []Profiler) GenesisState {
	return GenesisState{
		Profilers: profilers,
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	// Add Profilers
	for _, profiler := range data.Profilers {
		keeper.AddProfiler(ctx, profiler)
	}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	profilersIterator := k.GetProfilers(ctx)
	var profilers []Profiler
	for ; profilersIterator.Valid(); profilersIterator.Next() {
		var profiler Profiler
		k.cdc.MustUnmarshalBinaryLengthPrefixed(profilersIterator.Value(), &profiler)
		profilers = append(profilers, profiler)
	}
	return GenesisState{
		Profilers: profilers,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	profiler := Profiler{Name: "genessis"}
	return GenesisState{
		Profilers: []Profiler{profiler},
	}
}

// get raw genesis raw message for testing
func DefaultGenesisStateForTest() GenesisState {
	profiler := Profiler{Name: "genessis"}
	return GenesisState{
		Profilers: []Profiler{profiler},
	}
}
