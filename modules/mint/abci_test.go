package mint_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/irisnet/irishub/modules/mint"
	"github.com/irisnet/irishub/simapp"
)

func TestBeginBlocker(t *testing.T) {
	app, ctx := createTestApp(true)

	mint.BeginBlocker(ctx, app.MintKeeper)
	minter := app.MintKeeper.GetMinter(ctx)
	param := app.MintKeeper.GetParamSet(ctx)
	mintCoins := minter.BlockProvision(param)

	acc1 := app.SupplyKeeper.GetModuleAccount(ctx, "fee_collector")
	require.Equal(t, acc1.GetCoins(), sdk.NewCoins(mintCoins))
}

// returns context and an app with updated mint keeper
func createTestApp(isCheckTx bool) (*simapp.SimApp, sdk.Context) {
	app := simapp.Setup(isCheckTx)

	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{Height: 2})
	app.MintKeeper.SetParamSet(ctx, mint.Params{
		Inflation: sdk.NewDecWithPrec(4, 2),
		MintDenom: sdk.DefaultBondDenom,
	})
	app.MintKeeper.SetMinter(ctx, mint.DefaultMinter())
	app.SupplyKeeper.SetSupply(ctx, supply.Supply{})
	app.DistrKeeper.SetFeePool(ctx, distribution.InitialFeePool())
	return app, ctx
}
