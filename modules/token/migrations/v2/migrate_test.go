package v2_test

// import (
// 	"testing"

// 	"github.com/stretchr/testify/require"

// 	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

// 	"github.com/irisnet/irismod/simapp"
// 	v2 "irismod.io/token/migrations/v2"
// 	tokentypes "irismod.io/token/types"
// 	v1 "irismod.io/token/types/v1"
// )

// func TestMigrate(t *testing.T) {
// 	app := simapp.Setup(t, false)
// 	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

// 	legacySubspace := app.GetSubspace(tokentypes.ModuleName)

// 	params := v1.DefaultParams()
// 	legacySubspace.SetParamSet(ctx, &params)

// 	err := v2.Migrate(
// 		ctx,
// 		app.TokenKeeper,
// 		legacySubspace,
// 	)
// 	require.NoError(t, err)

// 	expParams := app.TokenKeeper.GetParams(ctx)
// 	// compatible with previous logic
// 	expParams.EnableErc20 = true
// 	require.Equal(t, expParams, params, "v2.Migrate failed")

// }
