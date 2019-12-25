package asset_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/asset"
	"github.com/irisnet/irishub/simapp"
)

func TestExportGatewayGenesis(t *testing.T) {
	isCheckTx := false
	app := simapp.Setup(isCheckTx)
	keeper := app.AssetKeeper
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})

	asset.InitGenesis(ctx, keeper, asset.DefaultGenesisState())

	// add token
	addr := sdk.AccAddress([]byte("addr1"))
	ft := asset.NewFungibleToken(asset.NATIVE, "bch", "bch", 1, "", "satoshi", sdk.NewIntWithDecimal(1, 0), sdk.NewIntWithDecimal(1, 0), true, addr)
	_, _, err := keeper.AddToken(ctx, ft)
	require.NoError(t, err)

	// query all token
	var tokens asset.Tokens
	keeper.IterateTokens(ctx, func(token asset.FungibleToken) (stop bool) {
		tokens = append(tokens, token)
		return false
	})

	require.Equal(t, len(tokens), 1)

	// export gateways
	genesisState := asset.ExportGenesis(ctx, keeper)

	for _, token := range genesisState.Tokens {
		require.Equal(t, token, ft)
	}
}
