package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/service/internal/types"
	"github.com/irisnet/irishub/simapp"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	testChainID     = "test-chain-id"
	testServiceName = "test-service"
	testServiceDesc = "test-service-desc"
	testServiceTags = []string{"tag1", "tag2"}
	testAuthor      = sdk.AccAddress([]byte("test-author"))
	testAuthorDesc  = "test-author-desc"
	testIDLContent  = ""

	testBindingType = types.Global
	testLevel       = types.Level{AvgRspTime: 10000, UsableTime: 9999}
	testProvider    = sdk.AccAddress([]byte("test-provider"))
	testDeposit, _  = sdk.ParseCoins("100iris-atto")
	testPrices      = []sdk.Coin{sdk.NewCoin("iris-atto", sdk.NewInt(50))}

	testConsumer       = sdk.AccAddress([]byte("test-consumer"))
	testMethodID       = int16(0)
	testServiceFees, _ = sdk.ParseCoins("50iris-atto")
	testInput          = []byte{}

	testProviderCoins, _ = sdk.ParseCoins("200iris-atto")
	testConsumerCoins, _ = sdk.ParseCoins("50iris-atto")
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

	app.bankKeeper.SetCoins(suite.ctx, testProvider, testProviderCoins)
	app.bankKeeper.SetCoins(suite.ctx, testConsumer, testConsumerCoins)
}

func (suite *KeeperTestSuite) setServiceDefinition() {
	svc := types.NewSvcDef(
		testServiceName, testChainID, testServiceDesc,
		testServiceTags, testAuthor, testAuthorDesc, testIDLContent,
	)

	suite.app.serviceKeeper.SetServiceDefinition(suite.ctx, svc)
	suite.app.serviceKeeper.AddMethods(suite.ctx, svc)
}

func (suite *KeeperTestSuite) setServiceBinding() {
	svcBinding := types.NewSvcBinding(
		suite.ctx, testChainID, testServiceName, testChainID, testProvider,
		testBindingType, testDeposit, testPrices, testLevel, true,
	)

	suite.app.serviceKeeper.SetServiceBinding(suite.ctx, svcBinding)
}
func (suite *KeeperTestSuite) TestServiceDefinition() {
	_, err := suite.app.serviceKeeper.AddServiceDefinition(
		suite.ctx, testServiceName, testChainID, testServiceDesc,
		testServiceTags, testAuthor, testAuthorDesc, testIDLContent,
	)
	suite.NoError(err)

	svc, found := suite.app.serviceKeeper.GetServiceDefinition(suite.ctx, testChainID, testServiceName)
	suite.True(found)

	expectedSvc := types.NewSvcDef(
		testServiceName, testChainID, testServiceDesc,
		testServiceTags, testAuthor, testAuthorDesc, testIDLContent,
	)
	suite.Equal(expectedSvc, svc)
}

func (suite *KeeperTestSuite) TestServiceBinding() {
	suite.setServiceDefinition()

	err := suite.app.serviceKeeper.AddServiceBinding(
		suite.ctx, testChainID, testServiceName, testChainID, testProvider,
		testBindingType, testDeposit, testPrices, testLevel,
	)
	suite.NoError(err)

	providerCoins1 := suite.app.BankKeeper.GetCoins(suite.ctx, testProvider)
	suite.Equal(testProviderCoins.Sub(testDeposit), providerCoins1)

	depositMaccAddr := suite.app.supplyKeeper.GetModuleAddress(suite.ctx, types.DepositAccName)
	depositMaccCoins1 := suite.app.BankKeeper.GetCoins(suite.ctx, depositMaccAddr)
	suite.Equal(testDeposit, depositMaccCoins1)

	binding, found := suite.app.serviceKeeper.GetServiceBinding(suite.ctx, testChainID, testServiceName, testChainID, testProvider)
	suite.True(found)

	expectedBinding := types.NewSvcBinding(
		suite.ctx, testChainID, testServiceName, testChainID, testProvider,
		testBindingType, testDeposit, testPrices, testLevel, true,
	)
	suite.Equal(expectedBinding, binding)

	_, err = suite.app.serviceKeeper.UpdateServiceBinding(
		suite.ctx, testChainID, testServiceName, testChainID, testProvider,
		testBindingType, testDeposit, testPrices, testLevel,
	)
	suite.NoError(err)

	providerCoins2 := suite.app.BankKeeper.GetCoins(suite.ctx, testProvider)
	suite.Equal(providerCoins1.Sub(testDeposit), providerCoins2)

	depositMaccCoins2 := suite.app.BankKeeper.GetCoins(suite.ctx, depositMaccAddr)
	suite.Equal(depositMaccCoins1.Add(testDeposit), depositMaccCoins2)

	binding, found = suite.app.serviceKeeper.GetServiceBinding(suite.ctx, testChainID, testServiceName, testChainID, testProvider)
	suite.True(found)
	suite.Equal(testDeposit.Add(testDeposit), binding.Deposit)
}

func (suite *KeeperTestSuite) TestServiceRequest() {
	suite.setServiceDefinition()
	suite.setServiceBinding()

	svcReq := types.NewSvcRequest(
		testChainID, testServiceName, testChainID, testChainID, testConsumer,
		testProvider, testMethodID, testInput, testServiceFees, false,
	)

	_, err := suite.app.serviceKeeper.AddRequest(
		suite.ctx, testChainID, testServiceName, testChainID, testChainID,
		testConsumer, testProvider, testMethodID, testInput, testServiceFees, false,
	)
	suite.NoError(err)

	consumerCoins := suite.app.BankKeeper.GetCoins(suite.ctx, testConsumer)
	suite.Equal(consumerCoins.Sub(testServiceFees), consumerCoins)

	requestMaccAddr := suite.app.supplyKeeper.GetModuleAddress(suite.ctx, types.RequestAccName)
	requestMaccCoins := suite.app.BankKeeper.GetCoins(suite.ctx, requestMaccAddr)
	suite.Equal(testServiceFees, requestMaccCoins)

	activeReq, found := suite.app.serviceKeeper.GetActiveRequest(suite.ctx, svcReq.ExpirationHeight, svcReq.RequestHeight, svcReq.RequestIntraTxCounter)
	suite.True(found)
	suite.Equal(svcReq.RequestID(), activeReq.RequestID())
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
