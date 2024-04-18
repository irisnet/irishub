package token

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/token/keeper"
	v1 "github.com/irisnet/irismod/modules/token/types/v1"
)

// InitGenesis stores the genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data v1.GenesisState) {
	if err := v1.ValidateGenesis(data); err != nil {
		panic(err.Error())
	}

	if err := k.SetParams(ctx, data.Params); err != nil {
		panic(err.Error())
	}

	// init tokens
	for _, token := range data.Tokens {
		if err := k.AddToken(ctx, token, false); err != nil {
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
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *v1.GenesisState {
	var tokens []v1.Token
	for _, token := range k.GetTokens(ctx, nil) {
		t := token.(*v1.Token)
		tokens = append(tokens, *t)
	}
	return &v1.GenesisState{
		Params:      k.GetParams(ctx),
		Tokens:      tokens,
		BurnedCoins: k.GetAllBurnCoin(ctx),
	}
}

// DefaultGenesisState returns the default genesis state for testing
func DefaultGenesisState() *v1.GenesisState {
	return &v1.GenesisState{
		Params: v1.DefaultParams(),
		Tokens: []v1.Token{v1.GetNativeToken()},
	}
}
