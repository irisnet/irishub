package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/irisnet/irishub/app/v3/service/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

var (
	testChainID     = "test-chain"

	testServiceName = "test-service"
	testServiceDesc = "test-service-desc"
	testServiceTags = []string{"tag1", "tag2"}
	testAuthorDesc  = "test-author-desc"
	testSchemas     = `{"input":{"type":"object"},"output":{"type":"object"},"error":{"type":"object"}}`

	testBindingType     = types.Global
	testLevel           = types.Level{AvgRspTime: 10000, UsableTime: 9999}
	testDeposit, _      = sdk.IrisCoinType.ConvertToMinDenomCoin("1000iris")
	testPrices          = []sdk.Coin{sdk.NewCoin(sdk.IrisAtto, sdk.NewIntWithDecimal(1, 18))}
	testAddedDeposit, _ = sdk.IrisCoinType.ConvertToMinDenomCoin("1000iris")

	testMethodID       = int16(1)
	testServiceFees, _ = sdk.IrisCoinType.ConvertToMinDenomCoin("1iris")
	testInput          = []byte{}
)

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

	_ = keeper.AddServiceDefinition(ctx, testServiceName, testServiceDesc, testServiceTags, author, testAuthorDesc, testSchemas)

	err := keeper.AddServiceBinding(ctx, testChainID, testServiceName, testChainID, provider, testBindingType, sdk.NewCoins(testDeposit), testPrices, testLevel)
	require.NoError(t, err)

	svcBinding, found := keeper.GetServiceBinding(ctx, testChainID, testServiceName, testChainID, provider)
	require.True(t, found)

	require.Equal(t, testServiceName, svcBinding.DefName)
	require.Equal(t, testBindingType, svcBinding.BindingType)
	require.Equal(t, sdk.NewCoins(testDeposit), svcBinding.Deposit)
	require.Equal(t, testPrices, svcBinding.Prices)

	// test binding update
	_, err = keeper.UpdateServiceBinding(ctx, testChainID, testServiceName, testChainID, provider, testBindingType, sdk.NewCoins(testAddedDeposit), testPrices, testLevel)
	require.NoError(t, err)

	updatedSvcBinding, found := keeper.GetServiceBinding(ctx, svcBinding.DefChainID, svcBinding.DefName, svcBinding.BindChainID, svcBinding.Provider)
	require.True(t, found)

	require.True(t, updatedSvcBinding.Deposit.IsEqual(svcBinding.Deposit.Add(sdk.NewCoins(testAddedDeposit))))
}

func TestKeeper_Call_Service(t *testing.T) {
	ctx, keeper, accs := createTestInput(t, sdk.NewIntWithDecimal(2000, 18), 3)

	author := accs[0].GetAddress()
	provider := accs[1].GetAddress()
	consumer := accs[2].GetAddress()

	_ = keeper.AddServiceDefinition(ctx, testServiceName, testServiceDesc, testServiceTags, author, testAuthorDesc, testSchemas)
	_ = keeper.AddServiceBinding(ctx, testChainID, testServiceName, testChainID, provider, testBindingType, sdk.NewCoins(testDeposit), testPrices, testLevel)

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
