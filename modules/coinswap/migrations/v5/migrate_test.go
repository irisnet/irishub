package v5_test

// import (
// 	"testing"

// 	"github.com/stretchr/testify/require"

// 	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

// 	"mods.irisnet.org/simapp"
// 	v5 "mods.irisnet.org/coinswap/migrations/v5"
// 	coinswaptypes "mods.irisnet.org/coinswap/types"
// )

// func TestMigrate(t *testing.T) {
// 	app := simapp.Setup(t, false)
// 	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

// 	legacySubspace := app.GetSubspace(coinswaptypes.ModuleName)

// 	params := coinswaptypes.DefaultParams()
// 	legacySubspace.SetParamSet(ctx, &params)

// 	err := v5.Migrate(
// 		ctx,
// 		app.CoinswapKeeper,
// 		legacySubspace,
// 	)
// 	require.NoError(t, err)

// 	expParams := app.CoinswapKeeper.GetParams(ctx)
// 	require.Equal(t, expParams, params, "v4.Migrate failed")

// }
