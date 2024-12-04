package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/v4/modules/mint/types"
	apptestutil "github.com/irisnet/irishub/v4/testutil"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context
	app *apptestutil.AppWrapper
}

func (suite *KeeperTestSuite) SetupTest() {
	app := apptestutil.CreateApp(suite.T())

	suite.ctx = app.BaseApp.NewContextLegacy(false, tmproto.Header{})
	suite.app = app

	err := app.MintKeeper.SetParams(suite.ctx, types.DefaultParams())
	require.NoError(suite.T(), err)
	app.MintKeeper.SetMinter(suite.ctx, types.DefaultMinter())
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestSetGetMinter() {
	minter := types.NewMinter(time.Now().UTC(), math.NewInt(100000))
	suite.app.MintKeeper.SetMinter(suite.ctx, minter)
	expMinter := suite.app.MintKeeper.GetMinter(suite.ctx)

	require.Equal(suite.T(), minter, expMinter)
}

func (suite *KeeperTestSuite) TestSetGetParamSet() {
	err := suite.app.MintKeeper.SetParams(suite.ctx, types.DefaultParams())
	require.NoError(suite.T(), err)

	expParamSet := suite.app.MintKeeper.GetParams(suite.ctx)
	require.Equal(suite.T(), types.DefaultParams(), expParamSet)
}

func (suite *KeeperTestSuite) TestMintCoins() {

	mintCoins := sdk.NewCoins(sdk.NewCoin("iris", math.NewInt(1000)))
	err := suite.app.MintKeeper.MintCoins(suite.ctx, mintCoins)
	require.NoError(suite.T(), err)

	acc := suite.app.AccountKeeper.GetModuleAccount(suite.ctx, types.ModuleName)
	coins := suite.app.BankKeeper.GetAllBalances(suite.ctx, acc.GetAddress())
	require.Equal(suite.T(), coins, mintCoins)
}

func (suite *KeeperTestSuite) TestAddCollectedFees() {
	mintCoins := sdk.NewCoins(sdk.NewCoin("iris", math.NewInt(1000)))

	feeCollector := suite.app.AccountKeeper.GetModuleAccount(suite.ctx, "fee_collector")
	feeCollectorBalance := suite.app.BankKeeper.GetAllBalances(suite.ctx, feeCollector.GetAddress())

	err := suite.app.MintKeeper.MintCoins(suite.ctx, mintCoins)
	require.NoError(suite.T(), err)

	acc := suite.app.AccountKeeper.GetModuleAccount(suite.ctx, types.ModuleName)
	coins := suite.app.BankKeeper.GetAllBalances(suite.ctx, acc.GetAddress())
	require.Equal(suite.T(), coins, mintCoins)

	err = suite.app.MintKeeper.AddCollectedFees(suite.ctx, mintCoins)
	require.NoError(suite.T(), err)

	acc = suite.app.AccountKeeper.GetModuleAccount(suite.ctx, types.ModuleName)
	coins = suite.app.BankKeeper.GetAllBalances(suite.ctx, acc.GetAddress())
	require.True(suite.T(), coins.Empty())

	feeCollectorTotalBalance := suite.app.BankKeeper.GetAllBalances(
		suite.ctx,
		feeCollector.GetAddress(),
	)
	expectedCollectedFees := feeCollectorTotalBalance.Sub(feeCollectorBalance...)
	require.Equal(suite.T(), expectedCollectedFees, mintCoins)

}
