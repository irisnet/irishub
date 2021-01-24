package keeper_test

import (
	"fmt"
	"testing"
	"time"

	gogotypes "github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	abci "github.com/tendermint/tendermint/abci/types"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/service/keeper"
	"github.com/irisnet/irismod/modules/service/types"
	"github.com/irisnet/irismod/simapp"
)

var (
	initCoinAmt = sdk.NewInt(100000)
	testCoin1   = sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000))
	testCoin2   = sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))
	testCoin3   = sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(2))

	testDenom1 = "testdenom1"                                                           // testing the normal cases
	testDenom2 = "testdenom2"                                                           // testing the case in which the feed value is 0
	testDenom3 = "testdenom3"                                                           // testing the case in which the feed is not existent
	testDenom4 = "ibc/9ebf7ebe6f8ffd34617809f3cf00e04a10d8b7226048f68866371fb9dad8a25d" // testing the ibc case
	testDenom5 = "peggy/0xdac17f958d2ee523a2206206994597c13d831ec7"                     // testing the ethpeg case

	testAuthor    sdk.AccAddress
	testOwner     sdk.AccAddress
	testProvider  sdk.AccAddress
	testProvider1 sdk.AccAddress
	testConsumer  sdk.AccAddress

	testServiceName = "test-service"
	testServiceDesc = "test-service-desc"
	testServiceTags = []string{"tag1", "tag2"}
	testAuthorDesc  = "test-author-desc"
	testSchemas     = `{"input":{"type":"object"},"output":{"type":"object"}}`

	testDeposit      = sdk.NewCoins(testCoin1)
	testPricing      = `{"price":"2stake","promotions_by_volume":[{"volume":1,"discount":"0.5"}]}`
	testQoS          = uint64(50)
	testOptions      = "{}"
	testWithdrawAddr = sdk.AccAddress([]byte("test-withdrawal-address"))
	testAddedDeposit = sdk.NewCoins(testCoin2)

	testInput         = `{"header":{},"body":{}}`
	testResult        = `{"code":200,"message":""}`
	testOutput        = `{"header":{},"body":{}}`
	testServiceFee    = sdk.NewCoins(testCoin3)
	testServiceFeeCap = sdk.NewCoins(testCoin3)
	testTimeout       = int64(100)
	testRepeatedFreq  = uint64(120)
	testRepeatedTotal = int64(100)

	callbacked = false
)

type KeeperTestSuite struct {
	suite.Suite

	cdc    codec.Marshaler
	ctx    sdk.Context
	keeper *keeper.Keeper
	app    *simapp.SimApp
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	isCheckTx := false
	app := simapp.Setup(isCheckTx)

	suite.cdc = codec.NewAminoCodec(app.LegacyAmino())
	suite.ctx = app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	suite.app = app
	suite.keeper = &app.ServiceKeeper

	suite.keeper.SetParams(suite.ctx, types.DefaultParams())

	suite.setTestAddrs()
}

func (suite *KeeperTestSuite) setTestAddrs() {
	testAddrs := simapp.AddTestAddrs(suite.app, suite.ctx, 5, initCoinAmt)

	testAuthor = testAddrs[0]
	testOwner = testAddrs[1]
	testProvider = testAddrs[2]
	testProvider1 = testAddrs[3]
	testConsumer = testAddrs[4]
}

func (suite *KeeperTestSuite) setServiceDefinition() {
	svcDef := types.NewServiceDefinition(testServiceName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, testSchemas)
	suite.keeper.SetServiceDefinition(suite.ctx, svcDef)
}

func (suite *KeeperTestSuite) setServiceBinding(available bool, disabledTime time.Time, provider, owner sdk.AccAddress) {
	svcBinding := types.NewServiceBinding(testServiceName, provider, testDeposit, testPricing, testQoS, testOptions, available, disabledTime, owner)
	suite.keeper.SetServiceBinding(suite.ctx, svcBinding)
	suite.keeper.SetOwnerServiceBinding(suite.ctx, svcBinding)

	suite.keeper.SetOwner(suite.ctx, provider, owner)
	suite.keeper.SetOwnerProvider(suite.ctx, owner, provider)

	pricing, _ := suite.keeper.ParsePricing(suite.ctx, testPricing)
	suite.keeper.SetPricing(suite.ctx, testServiceName, provider, pricing)
}

func (suite *KeeperTestSuite) TestDefineService() {
	err := suite.keeper.AddServiceDefinition(suite.ctx, testServiceName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, testSchemas)
	suite.NoError(err)

	svcDef, found := suite.keeper.GetServiceDefinition(suite.ctx, testServiceName)
	suite.True(found)

	suite.Equal(testServiceName, svcDef.Name)
	suite.Equal(testServiceDesc, svcDef.Description)
	suite.Equal(testServiceTags, svcDef.Tags)
	suite.Equal(testAuthor.String(), svcDef.Author)
	suite.Equal(testAuthorDesc, svcDef.AuthorDescription)
	suite.Equal(testSchemas, svcDef.Schemas)
}

func (suite *KeeperTestSuite) TestBindService() {
	suite.setServiceDefinition()

	err := suite.keeper.AddServiceBinding(suite.ctx, testServiceName, testProvider, testDeposit, testPricing, testQoS, testOptions, testOwner)
	suite.NoError(err)

	svcBinding, found := suite.keeper.GetServiceBinding(suite.ctx, testServiceName, testProvider)
	suite.True(found)

	suite.Equal(testServiceName, svcBinding.ServiceName)
	suite.Equal(testProvider.String(), svcBinding.Provider)
	suite.Equal(testDeposit, svcBinding.Deposit)
	suite.Equal(testPricing, svcBinding.Pricing)
	suite.Equal(testQoS, svcBinding.QoS)
	suite.True(svcBinding.Available)
	suite.True(svcBinding.DisabledTime.IsZero())
	suite.Equal(testOwner.String(), svcBinding.Owner)

	svcBindings := suite.keeper.GetOwnerServiceBindings(suite.ctx, testOwner, testServiceName)
	suite.Equal(1, len(svcBindings))
	suite.Equal(svcBinding.String(), svcBindings[0].String())

	providerOwner, found := suite.keeper.GetOwner(suite.ctx, testProvider)
	suite.True(found)
	suite.Equal(testOwner, providerOwner)

	iterator := suite.keeper.OwnerProvidersIterator(suite.ctx, testOwner)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		suite.Equal(testProvider, sdk.AccAddress(iterator.Key()[sdk.AddrLen+1:]))
	}

	// update binding
	newPricing := `{"price":"1stake"}`
	newQoS := uint64(80)
	newOptions := "{}"

	provider, _ := sdk.AccAddressFromBech32(svcBinding.Provider)
	err = suite.keeper.UpdateServiceBinding(suite.ctx, svcBinding.ServiceName, provider, testAddedDeposit, newPricing, newQoS, newOptions, testOwner)
	suite.NoError(err)

	updatedSvcBinding, found := suite.keeper.GetServiceBinding(suite.ctx, svcBinding.ServiceName, provider)
	suite.True(found)

	suite.True(updatedSvcBinding.Deposit.IsEqual(svcBinding.Deposit.Add(testAddedDeposit...)))
	suite.Equal(newPricing, updatedSvcBinding.Pricing)
	suite.Equal(newQoS, updatedSvcBinding.QoS)
}

func (suite *KeeperTestSuite) TestSetWithdrawAddress() {
	suite.setServiceBinding(true, time.Time{}, testProvider, testOwner)

	withdrawAddr := suite.keeper.GetWithdrawAddress(suite.ctx, testOwner)
	suite.Equal(testOwner, withdrawAddr)

	suite.keeper.SetWithdrawAddress(suite.ctx, testOwner, testWithdrawAddr)

	withdrawAddr = suite.keeper.GetWithdrawAddress(suite.ctx, testOwner)
	suite.Equal(testWithdrawAddr, withdrawAddr)
}

func (suite *KeeperTestSuite) TestDisableServiceBinding() {
	suite.setServiceBinding(true, time.Time{}, testProvider, testOwner)

	currentTime := time.Now().UTC()
	suite.ctx = suite.ctx.WithBlockTime(currentTime)

	err := suite.keeper.DisableServiceBinding(suite.ctx, testServiceName, testProvider, testOwner)
	suite.NoError(err)

	svcBinding, found := suite.keeper.GetServiceBinding(suite.ctx, testServiceName, testProvider)
	suite.True(found)

	suite.False(svcBinding.Available)
	suite.Equal(currentTime, svcBinding.DisabledTime)
}

func (suite *KeeperTestSuite) TestEnableServiceBinding() {
	disabledTime := time.Now().UTC()
	suite.setServiceBinding(false, disabledTime, testProvider, testOwner)

	err := suite.keeper.EnableServiceBinding(suite.ctx, testServiceName, testProvider, nil, testOwner)
	suite.NoError(err)

	svcBinding, found := suite.keeper.GetServiceBinding(suite.ctx, testServiceName, testProvider)
	suite.True(found)

	suite.True(svcBinding.Available)
	suite.True(svcBinding.DisabledTime.IsZero())
}

func (suite *KeeperTestSuite) TestRefundDeposit() {
	disabledTime := time.Now().UTC()
	suite.setServiceBinding(false, disabledTime, testProvider, testOwner)

	err := suite.app.BankKeeper.AddCoins(suite.ctx, suite.keeper.GetServiceDepositAccount(suite.ctx).GetAddress(), testDeposit)
	suite.NoError(err)

	params := suite.keeper.GetParams(suite.ctx)
	blockTime := disabledTime.Add(params.ArbitrationTimeLimit).Add(params.ComplaintRetrospect)
	suite.ctx = suite.ctx.WithBlockTime(blockTime)

	err = suite.keeper.RefundDeposit(suite.ctx, testServiceName, testProvider, testOwner)
	suite.NoError(err)

	svcBinding, found := suite.keeper.GetServiceBinding(suite.ctx, testServiceName, testProvider)
	suite.True(found)

	suite.Equal(sdk.Coins(nil), svcBinding.Deposit)
}

func (suite *KeeperTestSuite) TestRegisterCallback() {
	moduleName := "test-module"

	err := suite.keeper.RegisterResponseCallback(moduleName, callback)
	suite.NoError(err)

	_, err = suite.keeper.GetResponseCallback(moduleName)
	suite.NoError(err)

	err = suite.keeper.RegisterResponseCallback(moduleName, callback)
	suite.Error(err, "module already registered")
}

func (suite *KeeperTestSuite) TestKeeperRequestContext() {
	consumer := testConsumer
	providers := []sdk.AccAddress{testProvider}

	suite.setServiceDefinition()

	blockHeight := int64(1000)
	ctx := suite.ctx.WithBlockHeight(blockHeight)
	suite.app.BeginBlocker(ctx, abci.RequestBeginBlock{})

	// create
	requestContextID, err := suite.keeper.CreateRequestContext(
		ctx, testServiceName, providers, consumer, testInput,
		testServiceFeeCap, testTimeout, false, true,
		testRepeatedFreq, testRepeatedTotal, types.RUNNING, 0, "",
	)
	suite.NoError(err)

	suite.True(suite.keeper.HasNewRequestBatch(ctx, requestContextID))

	requestContext, found := suite.keeper.GetRequestContext(ctx, requestContextID)
	suite.True(found)

	suite.Equal(testServiceName, requestContext.ServiceName)
	pds := make([]string, len(providers))
	for i, provider := range providers {
		pds[i] = provider.String()
	}
	suite.Equal(pds, requestContext.Providers)
	suite.Equal(consumer.String(), requestContext.Consumer)
	suite.Equal(testServiceFeeCap, requestContext.ServiceFeeCap)
	suite.Equal(testTimeout, requestContext.Timeout)
	suite.True(requestContext.Repeated)
	suite.Equal(testRepeatedFreq, requestContext.RepeatedFrequency)
	suite.Equal(testRepeatedTotal, requestContext.RepeatedTotal)
	suite.Equal(uint64(0), requestContext.BatchCounter)
	suite.Equal(types.RUNNING, requestContext.State)

	// update
	newServiceFeeCap := sdk.NewCoins(testCoin1)
	newTimeout := testTimeout - 10
	newRepeatedFreq := testRepeatedFreq + 10
	newRepeatedTotal := int64(-1)

	err = suite.keeper.UpdateRequestContext(ctx, requestContextID, nil, 0, newServiceFeeCap, newTimeout, newRepeatedFreq, newRepeatedTotal, consumer)
	suite.NoError(err)

	requestContext, found = suite.keeper.GetRequestContext(ctx, requestContextID)
	suite.True(found)

	suite.Equal(testServiceName, requestContext.ServiceName)
	suite.Equal(pds, requestContext.Providers)
	suite.Equal(consumer.String(), requestContext.Consumer)
	suite.Equal(newServiceFeeCap, requestContext.ServiceFeeCap)
	suite.Equal(newTimeout, requestContext.Timeout)
	suite.True(requestContext.Repeated)
	suite.Equal(newRepeatedFreq, requestContext.RepeatedFrequency)
	suite.Equal(newRepeatedTotal, requestContext.RepeatedTotal)
	suite.Equal(uint64(0), requestContext.BatchCounter)
	suite.Equal(types.RUNNING, requestContext.State)

	// pause
	err = suite.keeper.PauseRequestContext(ctx, requestContextID, consumer)
	suite.NoError(err)

	requestContext, found = suite.keeper.GetRequestContext(ctx, requestContextID)
	suite.True(found)

	suite.Equal(types.PAUSED, requestContext.State)

	// start
	err = suite.keeper.StartRequestContext(ctx, requestContextID, consumer)
	suite.NoError(err)

	requestContext, found = suite.keeper.GetRequestContext(ctx, requestContextID)
	suite.True(found)

	suite.Equal(types.RUNNING, requestContext.State)

	// kill
	err = suite.keeper.KillRequestContext(ctx, requestContextID, consumer)
	suite.NoError(err)

	requestContext, found = suite.keeper.GetRequestContext(ctx, requestContextID)
	suite.True(found)
	suite.Equal(types.COMPLETED, requestContext.State)
}

func (suite *KeeperTestSuite) TestKeeperRequestService() {
	providers := []sdk.AccAddress{testProvider, testProvider1}
	consumer := testConsumer

	suite.setServiceDefinition()

	for _, provider := range providers {
		suite.setServiceBinding(true, time.Time{}, provider, testOwner)
	}

	blockHeight := int64(1000)
	ctx := suite.ctx.WithBlockHeight(blockHeight)
	suite.app.BeginBlocker(ctx, abci.RequestBeginBlock{})

	requestContextID, requestContext := suite.setRequestContext(ctx, consumer, providers, types.RUNNING, 0, "")

	newProviders, totalServiceFees, _, _ := suite.keeper.FilterServiceProviders(ctx, testServiceName, providers, testTimeout, testServiceFeeCap, consumer)
	suite.Equal(providers, newProviders)
	suite.Equal("4stake", totalServiceFees.String())

	err := suite.keeper.DeductServiceFees(ctx, consumer, totalServiceFees)
	suite.NoError(err)

	requestContext.BatchCounter++
	suite.keeper.SetRequestContext(ctx, requestContextID, requestContext)

	providerRequests := make(map[string][]string)
	suite.keeper.InitiateRequests(ctx, requestContextID, newProviders, providerRequests)

	requestContext, _ = suite.keeper.GetRequestContext(ctx, requestContextID)
	suite.Equal(len(newProviders), int(requestContext.BatchRequestCount))
	suite.Equal(types.BATCHRUNNING, requestContext.BatchState)

	iterator := suite.keeper.ActiveRequestsIteratorByReqCtx(ctx, requestContextID, requestContext.BatchCounter)
	defer iterator.Close()

	suite.True(iterator.Valid())

	var requestProviders []sdk.AccAddress
	for ; iterator.Valid(); iterator.Next() {
		var requestIDBz gogotypes.BytesValue
		suite.cdc.MustUnmarshalBinaryBare(iterator.Value(), &requestIDBz)

		requestID := requestIDBz.Value
		request, found := suite.keeper.GetRequest(ctx, requestID)
		suite.True(found)

		suite.Equal(requestContext.ServiceName, request.ServiceName)
		suite.Equal(requestContext.Consumer, request.Consumer)

		provider, _ := sdk.AccAddressFromBech32(request.Provider)
		requestProviders = append(requestProviders, provider)

		suite.Equal(blockHeight, request.RequestHeight)
		suite.Equal(blockHeight+testTimeout, request.ExpirationHeight)
		suite.Equal(requestContext.BatchCounter, request.RequestContextBatchCounter)
		suite.Equal(requestContextID.String(), request.RequestContextId)
	}

	suite.Equal(newProviders, requestProviders)

	// increase volume
	suite.keeper.SetRequestVolume(ctx, consumer, testServiceName, testProvider, 1)
	suite.keeper.SetRequestVolume(ctx, consumer, testServiceName, testProvider1, 1)

	// service fees will change due to the increased volume
	_, totalServiceFees, _, _ = suite.keeper.FilterServiceProviders(ctx, testServiceName, providers, testTimeout, testServiceFeeCap, consumer)
	suite.Equal("4stake", totalServiceFees.String())

	// satifying providers will change due to the condition changed
	newTimeout := int64(40)

	newProviders, _, _, _ = suite.keeper.FilterServiceProviders(ctx, testServiceName, providers, newTimeout, testServiceFeeCap, consumer)
	suite.Equal(0, len(newProviders))
}

func (suite *KeeperTestSuite) TestKeeperRespondService() {
	ctx := suite.ctx

	provider := testProvider
	consumer := testConsumer

	suite.setServiceDefinition()
	suite.keeper.SetOwner(suite.ctx, provider, testOwner)

	blockHeight := int64(1000)
	ctx = ctx.WithBlockHeight(blockHeight)

	requestContextID, requestContext := suite.setRequestContext(ctx, consumer, []sdk.AccAddress{provider}, types.RUNNING, 0, "")

	requestContext.BatchCounter++
	suite.keeper.SetRequestContext(ctx, requestContextID, requestContext)

	requestID1 := suite.setRequest(ctx, consumer, provider, requestContextID)
	requestID2 := suite.setRequest(ctx, consumer, provider, requestContextID)

	// respond request 1
	_, _, err := suite.keeper.AddResponse(ctx, requestID1, provider, testResult, testOutput)
	suite.NoError(err)

	requestContext, _ = suite.keeper.GetRequestContext(ctx, requestContextID)
	suite.Equal(uint32(1), requestContext.BatchResponseCount)
	suite.Equal(types.BATCHRUNNING, requestContext.BatchState)

	response, found := suite.keeper.GetResponse(ctx, requestID1)
	suite.True(found)

	suite.Equal(provider.String(), response.Provider)
	suite.Equal(consumer.String(), response.Consumer)
	suite.Equal(requestContextID.String(), response.RequestContextId)
	suite.Equal(requestContext.BatchCounter, response.RequestContextBatchCounter)

	volume := suite.keeper.GetRequestVolume(ctx, consumer, requestContext.ServiceName, provider)
	suite.Equal(uint64(1), volume)

	// respond request 2
	_, _, err = suite.keeper.AddResponse(ctx, requestID2, provider, testResult, testOutput)
	suite.NoError(err)

	requestContext, _ = suite.keeper.GetRequestContext(ctx, requestContextID)
	suite.Equal(uint32(2), requestContext.BatchResponseCount)
	suite.Equal(types.BATCHCOMPLETED, requestContext.BatchState)

	_, found = suite.keeper.GetResponse(ctx, requestID2)
	suite.True(found)

	volume = suite.keeper.GetRequestVolume(ctx, consumer, requestContext.ServiceName, provider)
	suite.Equal(uint64(2), volume)

	earnedFees, found := suite.keeper.GetEarnedFees(ctx, provider)
	suite.True(found)
	suite.False(earnedFees.Empty())

	ownerEarnedFees, found := suite.keeper.GetOwnerEarnedFees(ctx, testOwner)
	suite.True(found)
	suite.Equal(earnedFees, ownerEarnedFees)

	suite.False(suite.keeper.IsRequestActive(ctx, requestID1))
	suite.False(suite.keeper.IsRequestActive(ctx, requestID2))
}

func (suite *KeeperTestSuite) TestRequestServiceFromModule() {
	ctx := suite.ctx

	provider1 := testProvider
	provider2 := testProvider1
	providers := []sdk.AccAddress{provider1, provider2}
	consumer := testConsumer

	suite.setServiceDefinition()

	moduleName := "other-module"
	respThreshold := uint32(2)

	err := suite.keeper.RegisterResponseCallback(moduleName, callback)
	suite.NoError(err)

	blockHeight := int64(1000)
	ctx = ctx.WithBlockHeight(blockHeight)

	requestContextID, requestContext := suite.setRequestContext(ctx, consumer, providers, types.RUNNING, respThreshold, moduleName)

	requestContext.BatchCounter++
	suite.keeper.SetRequestContext(ctx, requestContextID, requestContext)

	requestID1 := suite.setRequest(ctx, consumer, provider1, requestContextID)
	requestID2 := suite.setRequest(ctx, consumer, provider2, requestContextID)

	_, _, err = suite.keeper.AddResponse(ctx, requestID1, provider1, testResult, testOutput)
	suite.NoError(err)

	requestContext, _ = suite.keeper.GetRequestContext(ctx, requestContextID)
	suite.Equal(uint32(1), requestContext.BatchResponseCount)
	suite.Equal(types.BATCHRUNNING, requestContext.BatchState)

	// callback has not occurred due to insufficient responses
	suite.False(callbacked)

	_, _, err = suite.keeper.AddResponse(ctx, requestID2, provider2, testResult, testOutput)
	suite.NoError(err)

	requestContext, _ = suite.keeper.GetRequestContext(ctx, requestContextID)
	suite.Equal(uint32(2), requestContext.BatchResponseCount)
	suite.Equal(types.BATCHCOMPLETED, requestContext.BatchState)

	// callback has occurred because the response count reaches the threshold
	suite.True(callbacked)
}

func (suite *KeeperTestSuite) TestGetMinDeposit() {
	oracleService := MockOracleService{
		feeds: map[string]string{
			fmt.Sprintf("%s-%s", testDenom1, sdk.DefaultBondDenom): "0.5",
			fmt.Sprintf("%s-%s", testDenom2, sdk.DefaultBondDenom): "0",
			fmt.Sprintf("%s-%s", testDenom4, sdk.DefaultBondDenom): "50",
			fmt.Sprintf("%s-%s", testDenom5, sdk.DefaultBondDenom): "20",
		},
	}

	suite.keeper.SetModuleService(types.RegisterModuleName, &types.ModuleService{
		ReuquestService: oracleService.GetExchangeRate,
	})

	testPricing1 := fmt.Sprintf(`{"price":"100%s"}`, testDenom1)
	testPricing2 := fmt.Sprintf(`{"price":"1%s"}`, testDenom1)
	testPricing3 := fmt.Sprintf(`{"price":"0%s"}`, testDenom1)
	testPricing4 := fmt.Sprintf(`{"price":"1%s"}`, testDenom2)
	testPricing5 := fmt.Sprintf(`{"price":"1%s"}`, testDenom3)
	testPricing6 := fmt.Sprintf(`{"price":"10%s"}`, testDenom4)
	testPricing7 := fmt.Sprintf(`{"price":"5%s"}`, testDenom5)

	pricing, err := suite.keeper.ParsePricing(suite.ctx, testPricing)
	suite.NoError(err)

	minDeposit, err := suite.keeper.GetMinDeposit(suite.ctx, pricing)
	suite.Equal(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(5000))), minDeposit)

	pricing1, err := suite.keeper.ParsePricing(suite.ctx, testPricing1)
	suite.NoError(err)

	minDeposit, err = suite.keeper.GetMinDeposit(suite.ctx, pricing1)
	suite.Equal(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(50000))), minDeposit)

	pricing2, err := suite.keeper.ParsePricing(suite.ctx, testPricing2)
	suite.NoError(err)

	minDeposit, err = suite.keeper.GetMinDeposit(suite.ctx, pricing2)
	suite.Equal(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(5000))), minDeposit)

	pricing3, err := suite.keeper.ParsePricing(suite.ctx, testPricing3)
	suite.NoError(err)

	minDeposit, err = suite.keeper.GetMinDeposit(suite.ctx, pricing3)
	suite.Equal(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(0))), minDeposit)

	pricing4, err := suite.keeper.ParsePricing(suite.ctx, testPricing4)
	suite.NoError(err)

	minDeposit, err = suite.keeper.GetMinDeposit(suite.ctx, pricing4)
	suite.NotNil(err, "should error when the exchange rate is zero")

	pricing5, err := suite.keeper.ParsePricing(suite.ctx, testPricing5)
	suite.NoError(err)

	minDeposit, err = suite.keeper.GetMinDeposit(suite.ctx, pricing5)
	suite.NotNil(err, "should error when the feed does not exist")

	pricing6, err := suite.keeper.ParsePricing(suite.ctx, testPricing6)
	suite.NoError(err)

	minDeposit, err = suite.keeper.GetMinDeposit(suite.ctx, pricing6)
	suite.Equal(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(500000))), minDeposit)

	pricing7, err := suite.keeper.ParsePricing(suite.ctx, testPricing7)
	suite.NoError(err)

	minDeposit, err = suite.keeper.GetMinDeposit(suite.ctx, pricing7)
	suite.Equal(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100000))), minDeposit)
}

func callback(ctx sdk.Context, requestContextID tmbytes.HexBytes, responses []string, err error) {
	callbacked = true
}

func (suite *KeeperTestSuite) setRequestContext(
	ctx sdk.Context, consumer sdk.AccAddress,
	providers []sdk.AccAddress, state types.RequestContextState,
	threshold uint32, moduleName string,
) (tmbytes.HexBytes, types.RequestContext) {
	requestContext := types.NewRequestContext(
		testServiceName, providers, consumer, testInput,
		testServiceFeeCap, testTimeout, false, true, testRepeatedFreq,
		testRepeatedTotal, 0, 0, 0, threshold, types.BATCHCOMPLETED,
		state, threshold, moduleName,
	)

	requestContextID := types.GenerateRequestContextID(keeper.TxHash(ctx), 0)
	suite.keeper.SetRequestContext(ctx, requestContextID, requestContext)

	return requestContextID, requestContext
}

func (suite *KeeperTestSuite) setRequest(ctx sdk.Context, consumer sdk.AccAddress, provider sdk.AccAddress, requestContextID []byte) tmbytes.HexBytes {
	requestContext, _ := suite.keeper.GetRequestContext(ctx, requestContextID)

	_ = suite.keeper.DeductServiceFees(ctx, consumer, testServiceFee)

	request := types.NewCompactRequest(
		requestContextID, requestContext.BatchCounter, provider,
		testServiceFee, ctx.BlockHeight(), requestContext.Timeout,
	)

	requestContext.BatchRequestCount++

	requestID := types.GenerateRequestID(requestContextID, request.RequestContextBatchCounter, ctx.BlockHeight(), int16(requestContext.BatchRequestCount))
	suite.keeper.SetCompactRequest(ctx, requestID, request)

	requestContext.BatchState = types.BATCHRUNNING
	suite.keeper.SetRequestContext(ctx, requestContextID, requestContext)

	suite.keeper.AddActiveRequest(ctx, requestContext.ServiceName, provider, request.RequestHeight+requestContext.Timeout, requestID)
	suite.keeper.AddActiveRequestByID(ctx, requestID)

	return requestID
}

// MockOracleService defines a mock oracle service for exchange rate
type MockOracleService struct {
	feeds map[string]string
}

func (m MockOracleService) GetExchangeRate(ctx sdk.Context, input string) (result string, output string) {
	feedName := gjson.Get(input, "body").Get("pair").String()

	value, ok := m.feeds[feedName]
	if !ok {
		result = `{"code":"400","message":"feed not found"}`
		return
	}

	result = `{"code":"200","message":""}`
	output = fmt.Sprintf(`{"header":{},"body":{"rate":"%s"}}`, value)

	return
}
