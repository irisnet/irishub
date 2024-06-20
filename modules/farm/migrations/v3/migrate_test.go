package v3_test

// import (
// 	"testing"

// 	"github.com/stretchr/testify/require"

// 	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

// 	"github.com/irisnet/irismod/simapp"
// 	v3 "github.com/irisnet/irismod/farm/migrations/v3"
// 	farmtypes "github.com/irisnet/irismod/farm/types"
// )

// func TestMigrate(t *testing.T) {
// 	app := simapp.Setup(t, false)
// 	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

// 	legacySubspace := app.GetSubspace(farmtypes.ModuleName)

// 	params := farmtypes.DefaultParams()
// 	legacySubspace.SetParamSet(ctx, &params)

// 	err := v3.Migrate(
// 		ctx,
// 		app.FarmKeeper,
// 		legacySubspace,
// 	)
// 	require.NoError(t, err)

// 	expParams := app.FarmKeeper.GetParams(ctx)
// 	require.Equal(t, expParams, params, "v3.Migrate failed")

// }
