package v2_test

// import (
// 	"testing"

// 	"github.com/stretchr/testify/require"

// 	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

// 	"mods.irisnet.org/simapp"
// 	v2 "mods.irisnet.org/service/migrations/v2"
// 	servicetypes "mods.irisnet.org/service/types"
// )

// func TestMigrate(t *testing.T) {
// 	app := simapp.Setup(t, false)
// 	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

// 	legacySubspace := app.GetSubspace(servicetypes.ModuleName)

// 	params := servicetypes.DefaultParams()
// 	legacySubspace.SetParamSet(ctx, &params)

// 	err := v2.Migrate(
// 		ctx,
// 		app.ServiceKeeper,
// 		legacySubspace,
// 	)
// 	require.NoError(t, err)

// 	expParams := app.ServiceKeeper.GetParams(ctx)
// 	require.Equal(t, expParams, params, "v2.Migrate failed")

// }
