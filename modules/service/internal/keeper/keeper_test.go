package keeper_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/service/internal/types"
	"github.com/irisnet/irishub/simapp"
)

var (
	testChainID     = "test-chain-id"
	testServiceName = "test-service"
	testServiceDesc = "test-service-desc"
	testServiceTags = []string{"tag1", "tag2"}
	testAuthor      = sdk.AccAddress([]byte("test-author"))
	testAuthorDesc  = "test-author-desc"
	testSchemas     = `{"input":{"type":"object"},"output":{"type":"object"},"error":{"type":"object"}}`

	testBindingType = types.Global
	testLevel       = types.Level{AvgRspTime: 10000, UsableTime: 9999}
	testProvider    = sdk.AccAddress([]byte("test-provider"))
	testPrices      = []sdk.Coin{sdk.NewCoin("iris", sdk.NewInt(1))}
	testDeposit, _  = sdk.ParseCoins("10000iris") // testPrices * 1000

	testConsumer       = sdk.AccAddress([]byte("test-consumer"))
	testMethodID       = int16(1)
	testServiceFees, _ = sdk.ParseCoins("50iris")
	testInput          = []byte{}

	testProviderCoins, _ = sdk.ParseCoins("50000iris")
	testConsumerCoins, _ = sdk.ParseCoins("10000iris")
)

type KeeperTestSuite struct {
	suite.Suite

	cdc *codec.Codec
	ctx sdk.Context
	app *simapp.SimApp
}

func (suite *KeeperTestSuite) SetupTest() {
	isCheckTx := false
	app := simapp.Setup(isCheckTx)

	suite.cdc = app.Codec()
	suite.ctx = app.BaseApp.NewContext(isCheckTx, abci.Header{})
	suite.app = app

	app.BankKeeper.SetCoins(suite.ctx, testProvider, testProviderCoins)
	app.BankKeeper.SetCoins(suite.ctx, testConsumer, testConsumerCoins)
}

func (suite *KeeperTestSuite) setServiceDefinition() {
	svcDef := types.NewServiceDefinition(
		testServiceName, testServiceDesc,
		testServiceTags, testAuthor, testAuthorDesc, testSchemas,
	)

	suite.app.ServiceKeeper.SetServiceDefinition(suite.ctx, svcDef)
}

func (suite *KeeperTestSuite) setServiceBinding() {
	svcBinding := types.NewSvcBinding(
		suite.ctx, testChainID, testServiceName, testChainID, testProvider,
		testBindingType, testDeposit, testPrices, testLevel, true,
	)

	suite.app.ServiceKeeper.SetServiceBinding(suite.ctx, svcBinding)
}
func (suite *KeeperTestSuite) TestServiceDefinition() {
	err := suite.app.ServiceKeeper.AddServiceDefinition(
		suite.ctx, testServiceName, testServiceDesc,
		testServiceTags, testAuthor, testAuthorDesc, testSchemas,
	)
	suite.NoError(err)

	svcDef, found := suite.app.ServiceKeeper.GetServiceDefinition(suite.ctx, testServiceName)
	suite.True(found)

	expectedSvcDef := types.NewServiceDefinition(
		testServiceName, testServiceDesc,
		testServiceTags, testAuthor, testAuthorDesc, testSchemas,
	)
	suite.Equal(expectedSvcDef, svcDef)
}

func (suite *KeeperTestSuite) TestServiceBinding() {
	suite.setServiceDefinition()

	err := suite.app.ServiceKeeper.AddServiceBinding(
		suite.ctx, testChainID, testServiceName, testChainID, testProvider,
		testBindingType, testDeposit, testPrices, testLevel,
	)
	suite.NoError(err)

	providerCoins1 := suite.app.BankKeeper.GetCoins(suite.ctx, testProvider)
	suite.Equal(testProviderCoins.Sub(testDeposit), providerCoins1)

	depositMaccAddr := suite.app.SupplyKeeper.GetModuleAddress(types.DepositAccName)
	depositMaccCoins1 := suite.app.BankKeeper.GetCoins(suite.ctx, depositMaccAddr)
	suite.Equal(testDeposit, depositMaccCoins1)

	binding, found := suite.app.ServiceKeeper.GetServiceBinding(suite.ctx, testChainID, testServiceName, testChainID, testProvider)
	suite.True(found)

	expectedBinding := types.NewSvcBinding(
		suite.ctx, testChainID, testServiceName, testChainID, testProvider,
		testBindingType, testDeposit, testPrices, testLevel, true,
	)
	suite.Equal(expectedBinding, binding)

	_, err = suite.app.ServiceKeeper.UpdateServiceBinding(
		suite.ctx, testChainID, testServiceName, testChainID, testProvider,
		testBindingType, testDeposit, testPrices, testLevel,
	)
	suite.NoError(err)

	providerCoins2 := suite.app.BankKeeper.GetCoins(suite.ctx, testProvider)
	suite.Equal(providerCoins1.Sub(testDeposit), providerCoins2)

	depositMaccCoins2 := suite.app.BankKeeper.GetCoins(suite.ctx, depositMaccAddr)
	suite.Equal(depositMaccCoins1.Add(testDeposit...), depositMaccCoins2)

	binding, found = suite.app.ServiceKeeper.GetServiceBinding(suite.ctx, testChainID, testServiceName, testChainID, testProvider)
	suite.True(found)
	suite.Equal(testDeposit.Add(testDeposit...), binding.Deposit)
}

func (suite *KeeperTestSuite) TestServiceRequest() {
	suite.setServiceDefinition()
	suite.setServiceBinding()

	svcReq := types.NewSvcRequest(
		testChainID, testServiceName, testChainID, testChainID, testConsumer,
		testProvider, testMethodID, testInput, testServiceFees, false,
	)

	binding, found := suite.app.ServiceKeeper.GetServiceBinding(suite.ctx, testChainID, testServiceName, testChainID, testProvider)
	suite.True(found)

	_, err := suite.app.ServiceKeeper.AddRequest(
		suite.ctx, testChainID, testServiceName, testChainID, testChainID,
		testConsumer, testProvider, testMethodID, testInput, testServiceFees, false,
	)
	suite.NoError(err)

	serviceFees := sdk.NewCoins(binding.Prices[testMethodID-1])
	consumerCoins := suite.app.BankKeeper.GetCoins(suite.ctx, testConsumer)
	suite.Equal(testConsumerCoins.Sub(serviceFees), consumerCoins)

	requestMaccAddr := suite.app.SupplyKeeper.GetModuleAddress(types.RequestAccName)
	requestMaccCoins := suite.app.BankKeeper.GetCoins(suite.ctx, requestMaccAddr)
	suite.Equal(serviceFees, requestMaccCoins)

	requestHeight := suite.ctx.BlockHeight()
	expirationHeight := requestHeight + types.DefaultMaxRequestTimeout

	activeReq, found := suite.app.ServiceKeeper.GetActiveRequest(suite.ctx, expirationHeight, requestHeight, svcReq.RequestIntraTxCounter)
	suite.True(found)
	suite.Equal(fmt.Sprintf("%d-%d-%d", expirationHeight, requestHeight, 0), activeReq.RequestID())
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
