package token

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis - store genesis parameters and tokens
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	if err := ValidateGenesis(data); err != nil {
		panic(fmt.Errorf("failed to initialize asset genesis state: %s", err.Error()))
	}
	k.SetParamSet(ctx, data.Params)
	//init tokens
	for _, token := range data.Tokens {
		if _, _, err := k.AddToken(ctx, token); err != nil {
			panic(err.Error())
		}
	}
}

// ExportGenesis - output genesis parameters and tokens
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
