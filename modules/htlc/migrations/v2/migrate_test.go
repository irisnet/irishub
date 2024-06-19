package v2_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cometbft/cometbft/crypto"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/simapp"
	v2 "irismod.io/htlc/migrations/v2"
	htlctypes "irismod.io/htlc/types"
)

func TestMigrate(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	legacySubspace := app.GetSubspace(htlctypes.ModuleName)

	params := htlctypes.Params{
		AssetParams: []htlctypes.AssetParam{
			{
				Denom: "htltbnb",
				SupplyLimit: htlctypes.SupplyLimit{
					Limit:          sdk.NewInt(350000000000000),
					TimeLimited:    false,
					TimeBasedLimit: sdk.ZeroInt(),
					TimePeriod:     time.Hour,
				},
				Active:        true,
				DeputyAddress: sdk.AccAddress(crypto.AddressHash([]byte("TestDeputy"))).String(),
				FixedFee:      sdk.NewInt(1000),
				MinSwapAmount: sdk.OneInt(),
				MaxSwapAmount: sdk.NewInt(1000000000000),
				MinBlockLock:  220,
				MaxBlockLock:  270,
			},
		},
	}
	legacySubspace.SetParamSet(ctx, &params)

	err := v2.Migrate(
		ctx,
		app.HTLCKeeper,
		legacySubspace,
	)
	require.NoError(t, err)

	expParams := app.HTLCKeeper.GetParams(ctx)
	require.Equal(t, expParams, params, "v2.Migrate failed")

}
