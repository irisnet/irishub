package asset

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	token "github.com/irisnet/irishub/modules/asset/01-token"
)

// GenesisState - all asset state that must be provided at genesis
type GenesisState struct {
	TokenGenesisState token.GenesisState
}

//NewGenesisState creates a new genesis state.
func NewGenesisState(tGenesisState token.GenesisState) GenesisState {
	return GenesisState{TokenGenesisState: tGenesisState}
}

// InitGenesis - store genesis parameters and tokens
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	token.InitGenesis(ctx, k.TokenKeeper, data.TokenGenesisState)
}

// ExportGenesis - output genesis
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	// export token genesisState
	tokenGenesisState := token.ExportGenesis(ctx, k.TokenKeeper)
	return GenesisState{
		TokenGenesisState: tokenGenesisState,
	}

}

// DefaultGenesisState return the default asset genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(token.DefaultGenesisState())
}

// ValidateGenesis validates the provided asset genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	// validate tokens
	if err := token.ValidateGenesis(data.TokenGenesisState); err != nil {
		return err
	}

	return nil
}
