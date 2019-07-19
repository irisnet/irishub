package rand

import (
	sdk "github.com/irisnet/irishub/types"
)

// InitGenesis stores genesis parameters
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	if err := ValidateGenesis(data); err != nil {
		panic(err.Error())
	}

	// TODO
}

// ExportGenesis outputs genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	// TODO

	return GenesisState{}
}

// DefaultGenesisState gets the default genesis state
func DefaultGenesisState() GenesisState {
	// TODO

	return GenesisState{}
}

// DefaultGenesisStateForTest gets the default genesis state for test
func DefaultGenesisStateForTest() GenesisState {
	// TODO

	return GenesisState{}
}

// ValidateGenesis validates the provided rand genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	// TODO

	return nil
}
