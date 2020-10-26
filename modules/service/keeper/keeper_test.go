package keeper_test

import (
	"testing"
	"time"

	gogotypes "github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto/tmhash"
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
	suite.Equal(testAuthor, svcDef.Author)
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
	suite.Equal(testProvider, svcBinding.Provider)
	suite.Equal(testDeposit, svcBinding.Deposit)
	suite.Equal(testPricing, svcBinding.Pricing)
	suite.Equal(testQoS, svcBinding.QoS)
	suite.True(svcBinding.Available)
	suite.True(svcBinding.DisabledTime.IsZero())
	suite.Equal(testOwner, svcBinding.Owner)

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

	err = suite.keeper.UpdateServiceBinding(suite.ctx, svcBinding.ServiceName, svcBinding.Provider, testAddedDeposit, newPricing, newQoS, newOptions, testOwner)
	suite.NoError(err)

	updatedSvcBinding, found := suite.keeper.GetServiceBinding(suite.ctx, svcBinding.ServiceName, svcBinding.Provider)
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
	ctx := suite.ctx.WithBlockHeight(blockHeight).
		WithValue(types.TxHash, tmhash.Sum([]byte("tx_hash"))).
		WithValue(types.MsgIndex, int64(0))

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
	suite.Equal(providers, requestContext.Providers)
	suite.Equal(consumer, requestContext.Consumer)
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
	suite.Equal(providers, requestContext.Providers)
	suite.Equal(consumer, requestContext.Consumer)
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
	ctx := suite.ctx.WithBlockHeight(blockHeight).
		WithValue(types.TxHash, tmhash.Sum([]byte("tx_hash"))).
		WithValue(types.MsgIndex, int64(0))

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

	requestProviders := []sdk.AccAddress{}
	for ; iterator.Valid(); iterator.Next() {
		var requestIDBz gogotypes.BytesValue
		suite.cdc.MustUnmarshalBinaryBare(iterator.Value(), &requestIDBz)

		requestID := requestIDBz.Value
		request, found := suite.keeper.GetRequest(ctx, requestID)
		suite.True(found)

		suite.Equal(requestContext.ServiceName, request.ServiceName)
		suite.Equal(requestContext.Consumer, request.Consumer)

		requestProviders = append(requestProviders, request.Provider)

		suite.Equal(blockHeight, request.RequestHeight)
		suite.Equal(blockHeight+testTimeout, request.ExpirationHeight)
		suite.Equal(requestContext.BatchCounter, request.RequestContextBatchCounter)
		suite.Equal(requestContextID, request.RequestContextId)
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

func (suite *KeeperTestSuite) TestKeeper_Respond_Service() {
	ctx := suite.ctx.WithValue(types.TxHash, tmhash.Sum([]byte("tx_hash")))

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

	suite.Equal(provider, response.Provider)
	suite.Equal(consumer, response.Consumer)
	suite.Equal(requestContextID, response.RequestContextId)
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
	ctx := suite.ctx.WithValue(types.TxHash, tmhash.Sum([]byte("tx_hash")))

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

	requestContextID := types.GenerateRequestContextID(ctx.Value(types.TxHash).([]byte), 0)
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
