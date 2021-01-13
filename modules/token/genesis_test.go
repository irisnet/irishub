package token_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/token"
	"github.com/irisnet/irismod/modules/token/types"
	"github.com/irisnet/irismod/simapp"
)

func TestExportGenesis(t *testing.T) {
	app := simapp.Setup(false)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	// export genesis
	genesisState := token.ExportGenesis(ctx, app.TokenKeeper)

	require.Equal(t, types.DefaultParams(), genesisState.Params)
	for _, token := range genesisState.Tokens {
		require.Equal(t, token, types.GetNativeToken())
	}
}

func TestInitGenesis(t *testing.T) {
	app := simapp.Setup(false)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	// add token
	addr := sdk.AccAddress(tmhash.SumTruncated([]byte("addr1")))
	ft := types.NewToken("btc", "Bitcoin Network", "satoshi", 1, 1, 1, true, addr)

	burnCoins := []sdk.Coin{
		{Denom: ft.MinUnit, Amount: sdk.NewInt(1000)},
	}
	genesis := types.GenesisState{
		Params:    types.DefaultParams(),
		Tokens:    []types.Token{ft},
		BurnCoins: burnCoins,
	}

	// initialize genesis
	token.InitGenesis(ctx, app.TokenKeeper, genesis)

	// query all tokens
	var tokens = app.TokenKeeper.GetTokens(ctx, nil)
	require.Equal(t, len(tokens), 2)
	require.Equal(t, tokens[0], &ft)

	var coins = app.TokenKeeper.GetAllBurnCoin(ctx)
	require.Equal(t, burnCoins, coins)
}
