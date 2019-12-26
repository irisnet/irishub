package token_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	token "github.com/irisnet/irishub/modules/asset/01-token"
	"github.com/irisnet/irishub/simapp"
)

func TestExportGatewayGenesis(t *testing.T) {
	isCheckTx := false
	app := simapp.Setup(isCheckTx)
	keeper := app.AssetKeeper.TokenKeeper
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})

	token.InitGenesis(ctx, keeper, token.DefaultGenesisState())

	// add token
	addr := sdk.AccAddress([]byte("addr1"))
	ft := token.NewFungibleToken(token.NATIVE, "bch", "bch", 1, "", "satoshi", sdk.NewIntWithDecimal(1, 0), sdk.NewIntWithDecimal(1, 0), true, addr)
	_, _, err := keeper.AddToken(ctx, ft)
	require.NoError(t, err)

	// query all token
	var tokens token.Tokens
	keeper.GetAllTokens(ctx, func(token token.FungibleToken) (stop bool) {
		tokens = append(tokens, token)
		return false
	})

	require.Equal(t, len(tokens), 1)

	// export gateways
	genesisState := token.ExportGenesis(ctx, keeper)

	for _, token := range genesisState.Tokens {
		require.Equal(t, token, ft)
	}
}
