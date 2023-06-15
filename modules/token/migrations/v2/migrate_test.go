package v2_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	v2 "github.com/irisnet/irismod/modules/token/migrations/v2"
	tokentypes "github.com/irisnet/irismod/modules/token/types"
	v1 "github.com/irisnet/irismod/modules/token/types/v1"
	"github.com/irisnet/irismod/simapp"
)

func TestMigrate(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	legacySubspace := app.GetSubspace(tokentypes.ModuleName)

	params := v1.DefaultParams()
	legacySubspace.SetParamSet(ctx, &params)

	err := v2.Migrate(
		ctx,
		app.TokenKeeper,
		legacySubspace,
	)
	require.NoError(t, err)

	expParams := app.TokenKeeper.GetParams(ctx)
	require.Equal(t, expParams, params, "v2.Migrate failed")

}
