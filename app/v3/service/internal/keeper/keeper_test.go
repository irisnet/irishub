package keeper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/irisnet/irishub/app/v3/service/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

var (
	testChainID = "test-chain"

	testServiceName = "test-service"
	testServiceDesc = "test-service-desc"
	testServiceTags = []string{"tag1", "tag2"}
	testAuthorDesc  = "test-author-desc"
	testSchemas     = `{"input":{"type":"object"},"output":{"type":"object"},"error":{"type":"object"}}`

	testDeposit, _      = sdk.IrisCoinType.ConvertToMinDenomCoin("1000iris")
	testPricing         = `{"price":"1iris"}`
	testAddedDeposit, _ = sdk.IrisCoinType.ConvertToMinDenomCoin("1000iris")
	testWithdrawAddr    = sdk.AccAddress{}

	testMethodID       = int16(1)
	testServiceFees, _ = sdk.IrisCoinType.ConvertToMinDenomCoin("1iris")
	testInput          = []byte{}
)

func setServiceDefinition(ctx sdk.Context, k Keeper, author sdk.AccAddress) {
	svcDef := types.NewServiceDefinition(testServiceName, testServiceDesc, testServiceTags, author, testAuthorDesc, testSchemas)
	k.SetServiceDefinition(ctx, svcDef)
}

func setServiceBinding(ctx sdk.Context, k Keeper, provider sdk.AccAddress, available bool, disabledTime time.Time) {
	svcBinding := types.NewServiceBinding(testServiceName, provider, sdk.NewCoins(testDeposit), testPricing, testWithdrawAddr, available, disabledTime)
	k.SetServiceBinding(ctx, svcBinding)
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
	ctx, keeper, accs := createTestInput(t, sdk.NewIntWithDecimal(2000, 18), 2)

	author := accs[0].GetAddress()
	provider := accs[1].GetAddress()

	setServiceDefinition(ctx, keeper, author)

	err := keeper.AddServiceBinding(ctx, testServiceName, provider, sdk.NewCoins(testDeposit), testPricing, testWithdrawAddr)
	require.NoError(t, err)

	svcBinding, found := keeper.GetServiceBinding(ctx, testServiceName, provider)
	require.True(t, found)

	require.Equal(t, testServiceName, svcBinding.ServiceName)
	require.Equal(t, provider, svcBinding.Provider)
	require.Equal(t, sdk.NewCoins(testDeposit), svcBinding.Deposit)
	require.Equal(t, testPricing, svcBinding.Pricing)
	require.Equal(t, provider, svcBinding.WithdrawAddress)
	require.True(t, svcBinding.Available)
	require.Equal(t, time.Time{}, svcBinding.DisabledTime)

	// update binding
	err = keeper.UpdateServiceBinding(ctx, svcBinding.ServiceName, svcBinding.Provider, sdk.NewCoins(testAddedDeposit), testPricing)
	require.NoError(t, err)

	updatedSvcBinding, found := keeper.GetServiceBinding(ctx, svcBinding.ServiceName, svcBinding.Provider)
	require.True(t, found)

	require.True(t, updatedSvcBinding.Deposit.IsEqual(svcBinding.Deposit.Add(sdk.NewCoins(testAddedDeposit))))
}

func TestKeeper_Set_Withdraw_Address(t *testing.T) {
	ctx, keeper, accs := createTestInput(t, sdk.NewIntWithDecimal(2000, 18), 2)

	provider := accs[0].GetAddress()
	withdrawAddr := accs[1].GetAddress()

	setServiceBinding(ctx, keeper, provider, true, time.Time{})

	err := keeper.SetWithdrawAddress(ctx, testServiceName, provider, withdrawAddr)
	require.NoError(t, err)

	svcBinding, found := keeper.GetServiceBinding(ctx, testServiceName, provider)
	require.True(t, found)

	require.Equal(t, withdrawAddr, svcBinding.WithdrawAddress)
}

func TestKeeper_Disable_Service(t *testing.T) {
	ctx, keeper, accs := createTestInput(t, sdk.NewIntWithDecimal(2000, 18), 1)

	provider := accs[0].GetAddress()
	setServiceBinding(ctx, keeper, provider, true, time.Time{})

	currentTime := time.Now()
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
	setServiceBinding(ctx, keeper, provider, false, time.Now())

	err := keeper.EnableService(ctx, testServiceName, provider, nil)
	require.NoError(t, err)

	svcBinding, found := keeper.GetServiceBinding(ctx, testServiceName, provider)
	require.True(t, found)

	require.True(t, svcBinding.Available)
	require.Equal(t, time.Time{}, svcBinding.DisabledTime)
}

func TestKeeper_Refund_Deposit(t *testing.T) {
	ctx, keeper, accs := createTestInput(t, sdk.NewIntWithDecimal(2000, 18), 1)

	provider := accs[0].GetAddress()

	disabledTime := time.Now()
	setServiceBinding(ctx, keeper, provider, false, disabledTime)

	params := keeper.GetParamSet(ctx)
	blockTime := disabledTime.Add(params.ArbitrationTimeLimit).Add(params.ComplaintRetrospect)
	ctx = ctx.WithBlockTime(blockTime)

	err := keeper.RefundDeposit(ctx, testServiceName, provider)
	require.NoError(t, err)

	svcBinding, found := keeper.GetServiceBinding(ctx, testServiceName, provider)
	require.True(t, found)

	require.Equal(t, sdk.Coins{}, svcBinding.Deposit)
}

func TestKeeper_Call_Service(t *testing.T) {
	ctx, keeper, accs := createTestInput(t, sdk.NewIntWithDecimal(2000, 18), 3)

	author := accs[0].GetAddress()
	provider := accs[1].GetAddress()
	consumer := accs[2].GetAddress()

	_ = keeper.AddServiceDefinition(ctx, testServiceName, testServiceDesc, testServiceTags, author, testAuthorDesc, testSchemas)
	_ = keeper.AddServiceBinding(ctx, testServiceName, provider, sdk.NewCoins(testDeposit), testPricing, sdk.AccAddress{})

	svcRequest, err := keeper.AddRequest(ctx, testChainID, testServiceName, testChainID, testChainID, consumer, provider, testMethodID, testInput, sdk.NewCoins(testServiceFees), false)
	require.NoError(t, err)

	storedSvcRequest, found := keeper.GetActiveRequest(ctx, svcRequest.ExpirationHeight, svcRequest.RequestHeight, svcRequest.RequestIntraTxCounter)
	require.True(t, found)
	require.Equal(t, svcRequest.RequestID(), storedSvcRequest.RequestID())

	iterator := keeper.ActiveRequestQueueIterator(ctx, svcRequest.ExpirationHeight)
	defer iterator.Close()

	require.True(t, iterator.Valid())
	for ; iterator.Valid(); iterator.Next() {
		var req types.SvcRequest
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &req)

		require.Equal(t, svcRequest.RequestID(), req.RequestID())
	}
}
