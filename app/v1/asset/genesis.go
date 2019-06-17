package asset

import (
	sdk "github.com/irisnet/irishub/types"
)

// GenesisState - all asset state that must be provided at genesis
type GenesisState struct {
	Params         Params  `json:"params"` // asset params
	FungibleTokens []Token `json:"tokens"` // issued assets
}

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	if err := ValidateGenesis(data); err != nil {
		panic(err.Error())
	}

	k.SetParamSet(ctx, data.Params)

	// TODO: init assets with data.FungibleTokens
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var tokens []Token // TODO: extract existing tokens from app state
	return GenesisState{
		Params:         k.GetParamSet(ctx),
		FungibleTokens: tokens,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:         DefaultParams(),
		FungibleTokens: []Token{},
	}
}

// get raw genesis raw message for testing
func DefaultGenesisStateForTest() GenesisState {
	return GenesisState{
		Params:         DefaultParamsForTest(),
		FungibleTokens: []Token{},
	}
}

// ValidateGenesis validates the provided asset genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	err := validateParams(data.Params)
	if err != nil {
		return err
	}
	return nil
}
