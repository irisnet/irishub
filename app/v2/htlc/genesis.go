package htlc

import (
	sdk "github.com/irisnet/irishub/types"
)

// InitGenesis stores genesis data
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {

}

// ExportGenesis outputs genesis data
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return GenesisState{}
}

// DefaultGenesisState gets the default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

// DefaultGenesisStateForTest gets the default genesis state for test
func DefaultGenesisStateForTest() GenesisState {
	return GenesisState{}
}
