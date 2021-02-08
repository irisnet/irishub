package token

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/token/keeper"
	"github.com/irisnet/irismod/modules/token/types"
)

// InitGenesis stores the genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	if err := types.ValidateGenesis(data); err != nil {
		panic(err.Error())
	}

	k.SetParamSet(ctx, data.Params)

	// init tokens
	for _, token := range data.Tokens {
		if err := k.AddToken(ctx, token); err != nil {
			panic(err.Error())
		}
	}

	for _, coin := range data.BurnedCoins {
		k.AddBurnCoin(ctx, coin)
	}

	// assert the symbol exists
	if !k.HasSymbol(ctx, data.Params.IssueTokenBaseFee.Denom) {
		panic(fmt.Sprintf("Token %s does not exist", data.Params.IssueTokenBaseFee.Denom))
	}

}

// ExportGenesis outputs the genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	var tokens []types.Token
	for _, token := range k.GetTokens(ctx, nil) {
		t := token.(*types.Token)
		tokens = append(tokens, *t)
	}
	return &types.GenesisState{
		Params:      k.GetParamSet(ctx),
		Tokens:      tokens,
		BurnedCoins: k.GetAllBurnCoin(ctx),
	}
}

// DefaultGenesisState returns the default genesis state for testing
func DefaultGenesisState() *types.GenesisState {
	return &types.GenesisState{
		Params: types.DefaultParams(),
		Tokens: []types.Token{types.GetNativeToken()},
	}
}
