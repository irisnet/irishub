package guardian

import (
	sdk "github.com/irisnet/irishub/types"
)

// GenesisState - all guardian state that must be provided at genesis
type GenesisState struct {
	Profilers []Guardian `json:"profilers"`
	Trustees  []Guardian `json:"trustees"`
}

func NewGenesisState(profilers, trustees []Guardian) GenesisState {
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
	profilersIterator := k.ProfilersIterator(ctx)
	defer profilersIterator.Close()
	var profilers []Guardian
	for ; profilersIterator.Valid(); profilersIterator.Next() {
		var profiler Guardian
		k.cdc.MustUnmarshalBinaryLengthPrefixed(profilersIterator.Value(), &profiler)
		profilers = append(profilers, profiler)
	}

	trusteesIterator := k.TrusteesIterator(ctx)
	defer trusteesIterator.Close()
	var trustees []Guardian
	for ; trusteesIterator.Valid(); trusteesIterator.Next() {
		var trustee Guardian
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
	profiler := Guardian{Description: "genessis"}
	trustee := Guardian{Description: "genessis"}
	return GenesisState{
		Profilers: []Guardian{profiler},
		Trustees:  []Guardian{trustee},
	}
}

// get raw genesis raw message for testing
func DefaultGenesisStateForTest() GenesisState {
	return DefaultGenesisState()
}
