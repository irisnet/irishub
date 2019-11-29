package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/irisnet/irishub/modules/mint/internal/types"
	"github.com/stretchr/testify/require"
)

func TestSetGetMinter(t *testing.T) {
	app, ctx := createTestApp(true)

	minter := types.NewMinter(time.Now().UTC(), sdk.NewInt(100000))
	app.MintKeeper.SetMinter(ctx, minter)
	expMinter := app.MintKeeper.GetMinter(ctx)

	require.Equal(t, minter, expMinter)
}

func TestSetGetParamSet(t *testing.T) {
	app, ctx := createTestApp(true)

	app.MintKeeper.SetParamSet(ctx, types.DefaultParams())
	expParamSet := app.MintKeeper.GetParamSet(ctx)

	require.Equal(t, types.DefaultParams(), expParamSet)
}

func TestMintCoins(t *testing.T) {
	app, ctx := createTestApp(true)

	app.SupplyKeeper.SetSupply(ctx, supply.Supply{})

	mintCoins := sdk.NewCoins(sdk.NewCoin("iris", sdk.NewInt(1000)))
	err := app.MintKeeper.MintCoins(ctx, mintCoins)
	require.NoError(t, err)

	acc := app.SupplyKeeper.GetModuleAccount(ctx, types.ModuleName)
	require.Equal(t, acc.GetCoins(), mintCoins)
}

func TestAddCollectedFees(t *testing.T) {
	app, ctx := createTestApp(true)

	app.SupplyKeeper.SetSupply(ctx, supply.Supply{})

	mintCoins := sdk.NewCoins(sdk.NewCoin("iris", sdk.NewInt(1000)))

	err := app.MintKeeper.MintCoins(ctx, mintCoins)
	require.NoError(t, err)

	acc := app.SupplyKeeper.GetModuleAccount(ctx, types.ModuleName)
	require.Equal(t, acc.GetCoins(), mintCoins)

	err = app.MintKeeper.AddCollectedFees(ctx, mintCoins)
	require.NoError(t, err)

	acc = app.SupplyKeeper.GetModuleAccount(ctx, types.ModuleName)
	require.True(t, acc.GetCoins().Empty())

	acc1 := app.SupplyKeeper.GetModuleAccount(ctx, "fee_collector")
	require.Equal(t, acc1.GetCoins(), mintCoins)

}
