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
		if err := k.AddToken(ctx, token); err != nil {
			panic(err.Error())
		}
	}
}

// ExportGenesis - output genesis parameters and tokens
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	// export created token
	return GenesisState{
		Params: k.GetParamSet(ctx),
		Tokens: k.GetAllTokens(ctx),
	}
}
