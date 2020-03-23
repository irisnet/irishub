package keeper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	cmn "github.com/tendermint/tendermint/libs/common"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v3/service/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

var (
	testChainID = "test-chain"

	testCoin1, _ = sdk.IrisCoinType.ConvertToCoin("10000iris")
	testCoin2, _ = sdk.IrisCoinType.ConvertToCoin("100iris")
	testCoin3, _ = sdk.IrisCoinType.ConvertToCoin("1iris")

	testServiceName = "test-service"
	testServiceDesc = "test-service-desc"
	testServiceTags = []string{"tag1", "tag2"}
	testAuthorDesc  = "test-author-desc"
	testSchemas     = `{"input":{"type":"object"},"output":{"type":"object"}}`

	testDeposit      = sdk.NewCoins(testCoin1)
	testPricing      = `{"price":"1iris","promotions_by_volume":[{"volume":1,"discount":"0.8"}]}`
	testWithdrawAddr = sdk.AccAddress([]byte("test-withdrawal-address"))
	testAddedDeposit = sdk.NewCoins(testCoin2)

	testInput         = `{"pair":"iris-usdt"}`
	testResult        = `{"code":200,"message":""}`
	testOutput        = `{"last":"100"}`
	testServiceFee    = sdk.NewCoins(testCoin3)
	testServiceFeeCap = sdk.NewCoins(testCoin3)
	testTimeout       = int64(100)
	testRepeatedFreq  = uint64(120)
	testRepeatedTotal = int64(100)

	callbacked = false
)

func setServiceDefinition(ctx sdk.Context, k Keeper, author sdk.AccAddress) {
	svcDef := types.NewServiceDefinition(testServiceName, testServiceDesc, testServiceTags, author, testAuthorDesc, testSchemas)
	k.SetServiceDefinition(ctx, svcDef)
}

func setServiceBinding(ctx sdk.Context, k Keeper, provider sdk.AccAddress, available bool, disabledTime time.Time) {
	svcBinding := types.NewServiceBinding(testServiceName, provider, testDeposit, testPricing, available, disabledTime)
	k.SetServiceBinding(ctx, svcBinding)

	pricing, _ := k.ParsePricing(ctx, testPricing)
	k.SetPricing(ctx, testServiceName, provider, pricing)
}

func setRequestContext(
	ctx sdk.Context, k Keeper, consumer sdk.AccAddress,
	providers []sdk.AccAddress, state types.RequestContextState,
	threshold uint16, moduleName string,
) (cmn.HexBytes, types.RequestContext) {
	requestContext := types.NewRequestContext(
		testServiceName, providers, consumer, testInput,
		testServiceFeeCap, testTimeout, false, true, testRepeatedFreq,
		testRepeatedTotal, 0, 0, 0, types.BATCHCOMPLETED, state, threshold, moduleName,
	)

	requestContextID := types.GenerateRequestContextID(ctx.BlockHeight(), 0)
	k.SetRequestContext(ctx, requestContextID, requestContext)

	return requestContextID, requestContext
}

func setRequest(ctx sdk.Context, k Keeper, consumer sdk.AccAddress, provider sdk.AccAddress, requestContextID []byte) cmn.HexBytes {
	requestContext, _ := k.GetRequestContext(ctx, requestContextID)

	_ = k.DeductServiceFees(ctx, consumer, testServiceFee)

	request := types.NewCompactRequest(
		requestContextID, requestContext.BatchCounter, provider,
		testServiceFee, ctx.BlockHeight(),
	)

	requestContext.BatchRequestCount++

	requestID := types.GenerateRequestID(requestContextID, request.RequestContextBatchCounter, int16(requestContext.BatchRequestCount))
	k.SetCompactRequest(ctx, requestID, request)

	requestContext.BatchState = types.BATCHRUNNING
	k.SetRequestContext(ctx, requestContextID, requestContext)

	k.AddActiveRequest(ctx, requestContext.ServiceName, provider, request.RequestHeight+requestContext.Timeout, requestID)
	k.AddActiveRequestByID(ctx, requestID)

	return requestID
}

func TestKeeper_Define_Service(t *testing.T) {
	ctx, keeper, accs := createTestInput(t, sdk.NewIntWithDecimal(1, 18), 1)

	author := accs[0].GetAddress()

	err := keeper.AddServiceDefinition(ctx, testServiceName, testServiceDesc, testServiceTags, author, testAuthorDesc, testSchemas)
	require.NoError(t, err)

	svcDef, found := keeper.GetServiceDefinition(ctx, testServiceName)
	require.True(t, found)

	require.Equal(t, testServiceName, svcDef.Name)
	require.Equal(t, testServiceTags, svcDef.Tags)
	require.Equal(t, author, svcDef.Author)
	require.Equal(t, testSchemas, svcDef.Schemas)
}

func TestKeeper_Bind_Service(t *testing.T) {
	ctx, keeper, accs := createTestInput(t, sdk.NewIntWithDecimal(20000, 18), 2)

	author := accs[0].GetAddress()
	provider := accs[1].GetAddress()

	setServiceDefinition(ctx, keeper, author)

	err := keeper.AddServiceBinding(ctx, testServiceName, provider, testDeposit, testPricing)
	require.NoError(t, err)

	svcBinding, found := keeper.GetServiceBinding(ctx, testServiceName, provider)
	require.True(t, found)

	require.Equal(t, testServiceName, svcBinding.ServiceName)
	require.Equal(t, provider, svcBinding.Provider)
	require.Equal(t, testDeposit, svcBinding.Deposit)
	require.Equal(t, testPricing, svcBinding.Pricing)
	require.True(t, svcBinding.Available)
	require.True(t, svcBinding.DisabledTime.IsZero())

	// update binding
	err = keeper.UpdateServiceBinding(ctx, svcBinding.ServiceName, svcBinding.Provider, testAddedDeposit, testPricing)
	require.NoError(t, err)

	updatedSvcBinding, found := keeper.GetServiceBinding(ctx, svcBinding.ServiceName, svcBinding.Provider)
	require.True(t, found)

	require.True(t, updatedSvcBinding.Deposit.IsEqual(svcBinding.Deposit.Add(testAddedDeposit)))
}

func TestKeeper_Set_Withdraw_Address(t *testing.T) {
	ctx, keeper, accs := createTestInput(t, sdk.NewIntWithDecimal(2000, 18), 2)

	provider := accs[0].GetAddress()

	setServiceBinding(ctx, keeper, provider, true, time.Time{})

	withdrawAddr := keeper.GetWithdrawAddress(ctx, provider)
	require.Equal(t, provider, withdrawAddr)

	keeper.SetWithdrawAddress(ctx, provider, testWithdrawAddr)

	withdrawAddr = keeper.GetWithdrawAddress(ctx, provider)
	require.Equal(t, testWithdrawAddr, withdrawAddr)
}

func TestKeeper_Disable_Service(t *testing.T) {
	ctx, keeper, accs := createTestInput(t, sdk.NewIntWithDecimal(2000, 18), 1)

	provider := accs[0].GetAddress()
	setServiceBinding(ctx, keeper, provider, true, time.Time{})

	currentTime := time.Now().UTC()
	ctx = ctx.WithBlockTime(currentTime)

	err := keeper.DisableService(ctx, testServiceName, provider)
	require.NoError(t, err)

	svcBinding, found := keeper.GetServiceBinding(ctx, testServiceName, provider)
	require.True(t, found)

	require.False(t, svcBinding.Available)
	require.Equal(t, currentTime, svcBinding.DisabledTime)
}

func TestKeeper_Enable_Service(t *testing.T) {
	ctx, keeper, accs := createTestInput(t, sdk.NewIntWithDecimal(2000, 18), 1)

	provider := accs[0].GetAddress()

	disabledTime := time.Now().UTC()
	setServiceBinding(ctx, keeper, provider, false, disabledTime)

	err := keeper.EnableService(ctx, testServiceName, provider, nil)
	require.NoError(t, err)

	svcBinding, found := keeper.GetServiceBinding(ctx, testServiceName, provider)
	require.True(t, found)

	require.True(t, svcBinding.Available)
	require.True(t, svcBinding.DisabledTime.IsZero())
}

func TestKeeper_Refund_Deposit(t *testing.T) {
	ctx, keeper, accs := createTestInput(t, sdk.NewIntWithDecimal(20000, 18), 1)

	provider := accs[0].GetAddress()

	disabledTime := time.Now().UTC()
	setServiceBinding(ctx, keeper, provider, false, disabledTime)

	_, err := keeper.bk.SendCoins(ctx, provider, auth.ServiceDepositCoinsAccAddr, testDeposit)
	require.NoError(t, err)

	params := keeper.GetParamSet(ctx)
	blockTime := disabledTime.Add(params.ArbitrationTimeLimit).Add(params.ComplaintRetrospect)
	ctx = ctx.WithBlockTime(blockTime)

	err = keeper.RefundDeposit(ctx, testServiceName, provider)
	require.NoError(t, err)

	svcBinding, found := keeper.GetServiceBinding(ctx, testServiceName, provider)
	require.True(t, found)

	require.Equal(t, sdk.Coins(nil), svcBinding.Deposit)
}

func TestKeeper_Register_Callback(t *testing.T) {
	_, keeper, _ := createTestInput(t, sdk.NewIntWithDecimal(2000, 18), 1)

	moduleName := "test-module"

	err := keeper.RegisterResponseCallback(moduleName, callback)
	require.NoError(t, err)

	_, err = keeper.GetResponseCallback(moduleName)
	require.NoError(t, err)

	err = keeper.RegisterResponseCallback(moduleName, callback)
	require.NotNil(t, err, "module already registered")
}

func TestKeeper_Request_Context(t *testing.T) {
	ctx, keeper, accs := createTestInput(t, sdk.NewIntWithDecimal(2000, 18), 4)

	author := accs[0].GetAddress()
	consumer := accs[0].GetAddress()
	providers := []sdk.AccAddress{accs[0].GetAddress(), accs[1].GetAddress()}

	setServiceDefinition(ctx, keeper, author)

	blockHeight := int64(1000)
	ctx = ctx.WithBlockHeight(blockHeight)

	// create
	requestContextID, err := keeper.CreateRequestContext(
		ctx, testServiceName, providers, consumer, testInput,
		testServiceFeeCap, testTimeout, false, true,
		testRepeatedFreq, testRepeatedTotal, types.RUNNING, 0, "",
	)
	require.NoError(t, err)

	require.True(t, keeper.HasNewRequestBatch(ctx, requestContextID, ctx.BlockHeight()))

	requestContext, found := keeper.GetRequestContext(ctx, requestContextID)
	require.True(t, found)

	require.Equal(t, testServiceName, requestContext.ServiceName)
	require.Equal(t, providers, requestContext.Providers)
	require.Equal(t, consumer, requestContext.Consumer)
	require.Equal(t, testServiceFeeCap, requestContext.ServiceFeeCap)
	require.Equal(t, testTimeout, requestContext.Timeout)
	require.True(t, requestContext.Repeated)
	require.Equal(t, testRepeatedFreq, requestContext.RepeatedFrequency)
	require.Equal(t, testRepeatedTotal, requestContext.RepeatedTotal)
	require.Equal(t, uint64(0), requestContext.BatchCounter)
	require.Equal(t, types.RUNNING, requestContext.State)

	// update
	newServiceFeeCap := sdk.NewCoins(testCoin1)
	newTimeout := testTimeout - 10
	newRepeatedFreq := testRepeatedFreq + 10
	newRepeatedTotal := int64(-1)

	err = keeper.UpdateRequestContext(ctx, requestContextID, nil, newServiceFeeCap, newTimeout, newRepeatedFreq, newRepeatedTotal, consumer)
	require.NoError(t, err)

	requestContext, found = keeper.GetRequestContext(ctx, requestContextID)
	require.True(t, found)

	require.Equal(t, testServiceName, requestContext.ServiceName)
	require.Equal(t, providers, requestContext.Providers)
	require.Equal(t, consumer, requestContext.Consumer)
	require.Equal(t, newServiceFeeCap, requestContext.ServiceFeeCap)
	require.Equal(t, newTimeout, requestContext.Timeout)
	require.True(t, requestContext.Repeated)
	require.Equal(t, newRepeatedFreq, requestContext.RepeatedFrequency)
	require.Equal(t, newRepeatedTotal, requestContext.RepeatedTotal)
	require.Equal(t, uint64(0), requestContext.BatchCounter)
	require.Equal(t, types.RUNNING, requestContext.State)

	// pause
	err = keeper.PauseRequestContext(ctx, requestContextID, consumer)
	require.NoError(t, err)

	requestContext, found = keeper.GetRequestContext(ctx, requestContextID)
	require.True(t, found)

	require.Equal(t, types.PAUSED, requestContext.State)

	// start
	err = keeper.StartRequestContext(ctx, requestContextID, consumer)
	require.NoError(t, err)

	requestContext, found = keeper.GetRequestContext(ctx, requestContextID)
	require.True(t, found)

	require.Equal(t, types.RUNNING, requestContext.State)

	// kill
	err = keeper.KillRequestContext(ctx, requestContextID, consumer)
	require.NoError(t, err)

	requestContext, found = keeper.GetRequestContext(ctx, requestContextID)
	require.True(t, found)

	require.Equal(t, types.COMPLETED, requestContext.State)
}

func TestKeeper_Request_Service(t *testing.T) {
	ctx, keeper, accs := createTestInput(t, sdk.NewIntWithDecimal(2000, 18), 4)

	author := accs[0].GetAddress()
	provider1 := accs[1].GetAddress()
	provider2 := accs[2].GetAddress()
	providers := []sdk.AccAddress{provider1, provider2}
	consumer := accs[3].GetAddress()

	setServiceDefinition(ctx, keeper, author)

	for _, provider := range providers {
		setServiceBinding(ctx, keeper, provider, true, time.Time{})
	}

	blockHeight := int64(1000)
	ctx = ctx.WithBlockHeight(blockHeight)

	requestContextID, requestContext := setRequestContext(ctx, keeper, consumer, providers, types.RUNNING, 0, "")

	newProviders, totalServiceFees := keeper.FilterServiceProviders(ctx, testServiceName, providers, testServiceFeeCap, consumer)
	require.Equal(t, providers, newProviders)
	require.Equal(t, "2iris", totalServiceFees.MainUnitString())

	err := keeper.DeductServiceFees(ctx, consumer, totalServiceFees)
	require.NoError(t, err)

	requestContext.BatchCounter++
	keeper.SetRequestContext(ctx, requestContextID, requestContext)

	keeper.InitiateRequests(ctx, requestContextID, newProviders)

	requestContext, _ = keeper.GetRequestContext(ctx, requestContextID)
	require.Equal(t, len(newProviders), int(requestContext.BatchRequestCount))
	require.Equal(t, types.BATCHRUNNING, requestContext.BatchState)

	iterator := keeper.ActiveRequestsIteratorByReqCtx(ctx, requestContextID, requestContext.BatchCounter)
	defer iterator.Close()

	require.True(t, iterator.Valid())

	requestProviders := []sdk.AccAddress{}
	for ; iterator.Valid(); iterator.Next() {
		var requestID []byte
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &requestID)

		request, found := keeper.GetRequest(ctx, requestID)
		require.True(t, found)

		require.Equal(t, requestContext.ServiceName, request.ServiceName)
		require.Equal(t, requestContext.Consumer, request.Consumer)

		requestProviders = append(requestProviders, request.Provider)

		require.Equal(t, blockHeight, request.RequestHeight)
		require.Equal(t, blockHeight+testTimeout, request.ExpirationHeight)
		require.Equal(t, requestContext.BatchCounter, request.RequestContextBatchCounter)
		require.Equal(t, requestContextID, request.RequestContextID)
	}

	require.Equal(t, newProviders, requestProviders)

	// increase volume
	keeper.SetRequestVolume(ctx, consumer, testServiceName, provider1, 1)
	keeper.SetRequestVolume(ctx, consumer, testServiceName, provider2, 1)

	// service fees will change due to the increased volume
	_, totalServiceFees = keeper.FilterServiceProviders(ctx, testServiceName, providers, testServiceFeeCap, consumer)
	require.Equal(t, "1.6iris", totalServiceFees.MainUnitString())
}

func TestKeeper_Respond_Service(t *testing.T) {
	ctx, keeper, accs := createTestInput(t, sdk.NewIntWithDecimal(2000, 18), 3)

	author := accs[0].GetAddress()
	provider := accs[1].GetAddress()
	consumer := accs[2].GetAddress()

	setServiceDefinition(ctx, keeper, author)

	blockHeight := int64(1000)
	ctx = ctx.WithBlockHeight(blockHeight)

	requestContextID, requestContext := setRequestContext(ctx, keeper, consumer, []sdk.AccAddress{provider}, types.RUNNING, 0, "")

	requestContext.BatchCounter++
	keeper.SetRequestContext(ctx, requestContextID, requestContext)

	requestID := setRequest(ctx, keeper, consumer, provider, requestContextID)

	requestIDStr := requestID.String()

	_, _, _, err := keeper.AddResponse(ctx, requestIDStr, provider, testResult, testOutput)
	require.NoError(t, err)

	requestContext, _ = keeper.GetRequestContext(ctx, requestContextID)
	require.Equal(t, uint16(1), requestContext.BatchResponseCount)
	require.Equal(t, types.BATCHCOMPLETED, requestContext.BatchState)

	response, found := keeper.GetResponse(ctx, requestID)
	require.True(t, found)

	require.Equal(t, provider, response.Provider)
	require.Equal(t, consumer, response.Consumer)
	require.Equal(t, requestContextID, response.RequestContextID)
	require.Equal(t, requestContext.BatchCounter, response.RequestContextBatchCounter)

	volume := keeper.GetRequestVolume(ctx, consumer, requestContext.ServiceName, provider)
	require.Equal(t, uint64(1), volume)

	earnedFees, found := keeper.GetEarnedFees(ctx, provider)
	require.True(t, found)
	require.True(t, !earnedFees.Coins.Empty())

	require.False(t, keeper.IsRequestActive(ctx, requestID))
}

func TestKeeper_Request_Service_From_Module(t *testing.T) {
	ctx, keeper, accs := createTestInput(t, sdk.NewIntWithDecimal(2000, 18), 4)

	author := accs[0].GetAddress()
	provider1 := accs[1].GetAddress()
	provider2 := accs[2].GetAddress()
	providers := []sdk.AccAddress{provider1, provider2}
	consumer := accs[3].GetAddress()

	setServiceDefinition(ctx, keeper, author)

	moduleName := "oracle"
	respThreshold := uint16(2)

	err := keeper.RegisterResponseCallback(moduleName, callback)
	require.NoError(t, err)

	blockHeight := int64(1000)
	ctx = ctx.WithBlockHeight(blockHeight)

	requestContextID, requestContext := setRequestContext(ctx, keeper, consumer, providers, types.RUNNING, respThreshold, moduleName)

	requestContext.BatchCounter++
	keeper.SetRequestContext(ctx, requestContextID, requestContext)

	requestID1 := setRequest(ctx, keeper, consumer, provider1, requestContextID)
	requestID2 := setRequest(ctx, keeper, consumer, provider2, requestContextID)

	requestIDStr1 := requestID1.String()
	requestIDStr2 := requestID2.String()

	_, _, _, err = keeper.AddResponse(ctx, requestIDStr1, provider1, testResult, testOutput)
	require.NoError(t, err)

	requestContext, _ = keeper.GetRequestContext(ctx, requestContextID)
	require.Equal(t, uint16(1), requestContext.BatchResponseCount)
	require.Equal(t, types.BATCHRUNNING, requestContext.BatchState)

	// callback has not occurred due to insufficient responses
	require.False(t, callbacked)

	_, _, _, err = keeper.AddResponse(ctx, requestIDStr2, provider2, testResult, testOutput)
	require.NoError(t, err)

	requestContext, _ = keeper.GetRequestContext(ctx, requestContextID)
	require.Equal(t, uint16(2), requestContext.BatchResponseCount)
	require.Equal(t, types.BATCHCOMPLETED, requestContext.BatchState)

	// callback has occurred because the response count reaches the threshold
	require.True(t, callbacked)
}

func callback(ctx sdk.Context, requestContextID cmn.HexBytes, responses []string, err error) {
	callbacked = true
}
