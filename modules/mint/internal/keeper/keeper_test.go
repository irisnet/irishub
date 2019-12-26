package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/irisnet/irishub/modules/mint/internal/types"
	"github.com/irisnet/irishub/simapp"
)

type KeeperTestSuite struct {
	suite.Suite

	cdc *codec.Codec
	ctx sdk.Context
	app *simapp.SimApp
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.Setup(false)

	suite.cdc = app.Codec()
	suite.ctx = app.BaseApp.NewContext(false, abci.Header{})
	suite.app = app

	app.MintKeeper.SetParamSet(suite.ctx, types.DefaultParams())
	app.MintKeeper.SetMinter(suite.ctx, types.DefaultMinter())
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestSetGetMinter() {
	minter := types.NewMinter(time.Now().UTC(), sdk.NewInt(100000))
	suite.app.MintKeeper.SetMinter(suite.ctx, minter)
	expMinter := suite.app.MintKeeper.GetMinter(suite.ctx)

	require.Equal(suite.T(), minter, expMinter)
}

func (suite *KeeperTestSuite) TestSetGetParamSet() {
	suite.app.MintKeeper.SetParamSet(suite.ctx, types.DefaultParams())
	expParamSet := suite.app.MintKeeper.GetParamSet(suite.ctx)

	require.Equal(suite.T(), types.DefaultParams(), expParamSet)
}

func (suite *KeeperTestSuite) TestMintCoins() {
	suite.app.SupplyKeeper.SetSupply(suite.ctx, supply.Supply{})

	mintCoins := sdk.NewCoins(sdk.NewCoin("iris", sdk.NewInt(1000)))
	err := suite.app.MintKeeper.MintCoins(suite.ctx, mintCoins)
	require.NoError(suite.T(), err)

	acc := suite.app.SupplyKeeper.GetModuleAccount(suite.ctx, types.ModuleName)
	require.Equal(suite.T(), acc.GetCoins(), mintCoins)
}

func (suite *KeeperTestSuite) TestAddCollectedFees() {
	suite.app.SupplyKeeper.SetSupply(suite.ctx, supply.Supply{})

	mintCoins := sdk.NewCoins(sdk.NewCoin("iris", sdk.NewInt(1000)))

	err := suite.app.MintKeeper.MintCoins(suite.ctx, mintCoins)
	require.NoError(suite.T(), err)

	acc := suite.app.SupplyKeeper.GetModuleAccount(suite.ctx, types.ModuleName)
	require.Equal(suite.T(), acc.GetCoins(), mintCoins)

	err = suite.app.MintKeeper.AddCollectedFees(suite.ctx, mintCoins)
	require.NoError(suite.T(), err)

	acc = suite.app.SupplyKeeper.GetModuleAccount(suite.ctx, types.ModuleName)
	require.True(suite.T(), acc.GetCoins().Empty())

	acc1 := suite.app.SupplyKeeper.GetModuleAccount(suite.ctx, "fee_collector")
	require.Equal(suite.T(), acc1.GetCoins(), mintCoins)

}
