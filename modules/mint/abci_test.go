package mint_test

import (
	"cosmossdk.io/math"
	"github.com/irisnet/irishub/v4/app/keepers"
	"testing"

	"github.com/stretchr/testify/require"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/irisnet/irishub/v4/modules/mint"
	"github.com/irisnet/irishub/v4/modules/mint/types"
	apptestutil "github.com/irisnet/irishub/v4/testutil"
)

func TestBeginBlocker(t *testing.T) {
	app, ctx := createTestApp(t, true)

	mint.BeginBlocker(ctx, app.MintKeeper)
	minter := app.MintKeeper.GetMinter(ctx)
	param := app.MintKeeper.GetParams(ctx)
	mintCoins := minter.BlockProvision(param)

	acc1 := app.AccountKeeper.GetModuleAccount(ctx, "fee_collector")
	mintedCoins := app.BankKeeper.GetAllBalances(ctx, acc1.GetAddress())
	require.Equal(t, mintedCoins, sdk.NewCoins(mintCoins))
}

// returns context and an app with updated mint keeper
func createTestApp(t *testing.T, isCheckTx bool) (*apptestutil.AppWrapper, sdk.Context) {
	app := apptestutil.CreateApp(t)

	ctx := app.BaseApp.NewContextLegacy(isCheckTx, tmproto.Header{Height: 2})
	app.MintKeeper.SetParams(ctx, types.NewParams(
		sdk.DefaultBondDenom,
		math.LegacyNewDecWithPrec(4, 2),
	))
	app.MintKeeper.SetMinter(ctx, types.DefaultMinter())
	err := keepers.NewDistrKeeperAdapter(app.DistrKeeper).SetFeePool(ctx, distributiontypes.InitialFeePool())
	require.NoError(t, err)

	return app, ctx
}
