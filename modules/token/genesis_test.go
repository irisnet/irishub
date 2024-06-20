package token_test

// import (
// 	"testing"

// 	"github.com/stretchr/testify/require"

// 	"github.com/cometbft/cometbft/crypto/tmhash"
// 	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

// 	sdk "github.com/cosmos/cosmos-sdk/types"

// 	"github.com/irisnet/irismod/modules/token"
// 	"github.com/irisnet/irismod/simapp"
// 	v1 "github.com/irisnet/irismod/token/types/v1"
// )

// func TestExportGenesis(t *testing.T) {
// 	app := simapp.Setup(t, false)

// 	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

// 	// export genesis
// 	genesisState := token.ExportGenesis(ctx, app.TokenKeeper)

// 	require.Equal(t, v1.DefaultParams(), genesisState.Params)
// 	for _, token := range genesisState.Tokens {
// 		require.Equal(t, token, v1.GetNativeToken())
// 	}
// }

// func TestInitGenesis(t *testing.T) {
// 	app := simapp.Setup(t, false)

// 	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

// 	// add token
// 	addr := sdk.AccAddress(tmhash.SumTruncated([]byte("addr1")))
// 	ft := v1.NewToken("btc", "Bitcoin Network", "satoshi", 1, 1, 1, true, addr)

// 	burnCoins := []sdk.Coin{
// 		{Denom: ft.MinUnit, Amount: sdk.NewInt(1000)},
// 	}
// 	genesis := v1.GenesisState{
// 		Params:      v1.DefaultParams(),
// 		Tokens:      []v1.Token{ft},
// 		BurnedCoins: burnCoins,
// 	}

// 	// initialize genesis
// 	token.InitGenesis(ctx, app.TokenKeeper, genesis)

// 	// query all tokens
// 	var tokens = app.TokenKeeper.GetTokens(ctx, nil)
// 	require.Equal(t, len(tokens), 2)
// 	require.Equal(t, tokens[0], &ft)

// 	var coins = app.TokenKeeper.GetAllBurnCoin(ctx)
// 	require.Equal(t, burnCoins, coins)
// }
