package asset

import (
	sdk "github.com/irisnet/irishub/types"
)

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	if err := ValidateGenesis(data); err != nil {
		panic(err.Error())
	}

	k.SetParamSet(ctx, data.Params)

	//init tokens
	for _, token := range data.Tokens {
		if err := k.AddToken(ctx, token); err != nil {
			panic(err.Error())
		}
	}
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	// export created token
	var tokens Tokens
	k.IterateTokens(ctx, func(token FungibleToken) (stop bool) {
		tokens = append(tokens, token)
		return false
	})
	return GenesisState{
		Params: k.GetParamSet(ctx),
		Tokens: tokens,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
		Tokens: []FungibleToken{},
	}
}

// get raw genesis raw message for testing
func DefaultGenesisStateForTest() GenesisState {
	return GenesisState{
		Params: DefaultParamsForTest(),
		Tokens: []FungibleToken{},
	}
}

// ValidateGenesis validates the provided asset genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	if err := ValidateParams(data.Params); err != nil {
		return err
	}

	// validate tokens
	if err := data.Tokens.Validate(); err != nil {
		return err
	}

	return nil
}
