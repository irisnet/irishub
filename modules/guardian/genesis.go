package guardian

import (
	sdk "github.com/irisnet/irishub/types"
)

// GenesisState - all profiling state that must be provided at genesis
type GenesisState struct {
	Profilers []Profiler `json:"profilers"`
	Trustees  []Trustee  `json:"trustees"`
}

func NewGenesisState(profilers []Profiler, trustees []Trustee) GenesisState {
	return GenesisState{
		Profilers: profilers,
		Trustees:  trustees,
	}
}

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
	profilersIterator := k.GetProfilers(ctx)
	var profilers []Profiler
	for ; profilersIterator.Valid(); profilersIterator.Next() {
		var profiler Profiler
		k.cdc.MustUnmarshalBinaryLengthPrefixed(profilersIterator.Value(), &profiler)
		profilers = append(profilers, profiler)
	}

	trusteesIterator := k.GetTrustees(ctx)
	var trustees []Trustee
	for ; trusteesIterator.Valid(); trusteesIterator.Next() {
		var trustee Trustee
		k.cdc.MustUnmarshalBinaryLengthPrefixed(trusteesIterator.Value(), &trustee)
		trustees = append(trustees, trustee)
	}
	return GenesisState{
		Profilers: profilers,
		Trustees:  trustees,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	profiler := Profiler{Name: "genessis"}
	trustee := Trustee{}
	return GenesisState{
		Profilers: []Profiler{profiler},
		Trustees:  []Trustee{trustee},
	}
}

// get raw genesis raw message for testing
func DefaultGenesisStateForTest() GenesisState {
	return DefaultGenesisState()
}
