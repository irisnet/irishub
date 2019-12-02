package asset

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/simapp"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestExportGatewayGenesis(t *testing.T) {
	isCheckTx := false
	app := simapp.Setup(isCheckTx)
	keeper := app.AssetKeeper
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})

	// add token
	addr := sdk.AccAddress([]byte("addr1"))
	ft := NewFungibleToken(NATIVE, "", "bch", "bch", 1, "", "satoshi", sdk.NewIntWithDecimal(1, 0), sdk.NewIntWithDecimal(1, 0), true, addr)
	_, _, err := keeper.AddToken(ctx, ft)
	require.NoError(t, err)

	// query all token
	var tokens Tokens
	keeper.IterateTokens(ctx, func(token FungibleToken) (stop bool) {
		tokens = append(tokens, token)
		return false
	})

	require.Equal(t, len(tokens), 1)

	// export gateways
	genesisState := ExportGenesis(ctx, keeper)

	for _, token := range genesisState.Tokens {
		require.Equal(t, token, ft)
	}
}
