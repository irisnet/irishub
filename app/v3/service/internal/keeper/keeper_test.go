package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/irisnet/irishub/app/v3/service/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

func TestKeeper_Service_Definition_Binding(t *testing.T) {
	ctx, keeper, accs := createTestInput(t, sdk.NewIntWithDecimal(2000, 18), 2)

	author := accs[0].GetAddress()
	provider := accs[1].GetAddress()

	serviceName := "myService"
	chainID := "testnet"
	serviceDescription := "the service for unit test"
	tags := []string{"test", "tutorial"}
	authorDescription := "unit test author"

	err := keeper.AddServiceDefinition(ctx, serviceName, chainID, serviceDescription, tags, author, authorDescription, idlContent)
	require.NoError(t, err)

	serviceDef, found := keeper.GetServiceDefinition(ctx, chainID, serviceName)

	require.True(t, found)
	require.Equal(t, idlContent, serviceDef.IDLContent)
	require.Equal(t, serviceName, serviceDef.Name)

	// test methods
	err = keeper.AddMethods(ctx, serviceDef)
	require.NoError(t, err)

	iterator := keeper.GetMethods(ctx, chainID, serviceName)
	defer iterator.Close()

	require.True(t, iterator.Valid())
	for ; ; iterator.Next() {
		var method types.MethodProperty
		if !iterator.Valid() {
			break
		}

		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &method)

		require.Equal(t, "SayHello", method.Name)
		require.Equal(t, "sayHello", method.Description)
		require.Equal(t, "NoPrivacy", method.OutputPrivacy.String())
		require.Equal(t, "NoCached", method.OutputCached.String())
	}

	// test binding
	deposit, _ := sdk.IrisCoinType.ConvertToMinDenomCoin("1000iris")
	price, _ := sdk.IrisCoinType.ConvertToMinDenomCoin("1iris")
	bindingType := types.Global
	level := types.Level{AvgRspTime: 10000, UsableTime: 9999}

	err = keeper.AddServiceBinding(ctx, chainID, serviceName, chainID, provider, bindingType, sdk.NewCoins(deposit), []sdk.Coin{price}, level)
	require.NoError(t, err)

	svcBinding, found := keeper.GetServiceBinding(ctx, chainID, serviceName, chainID, provider)
	require.True(t, found)
	require.Equal(t, sdk.NewCoins(deposit), svcBinding.Deposit)
	require.Equal(t, bindingType, svcBinding.BindingType)

	// test binding update
	addedDeposit, _ := sdk.IrisCoinType.ConvertToMinDenomCoin("100iris")
	_, err = keeper.UpdateServiceBinding(ctx, chainID, serviceName, chainID, provider, bindingType, sdk.NewCoins(addedDeposit), []sdk.Coin{price}, level)
	require.NoError(t, err)

	updatedSvcBinding, found := keeper.GetServiceBinding(ctx, svcBinding.DefChainID, svcBinding.DefName, svcBinding.BindChainID, svcBinding.Provider)
	require.True(t, found)
	require.True(t, updatedSvcBinding.Deposit.IsEqual(svcBinding.Deposit.Add(sdk.NewCoins(addedDeposit))))
}

func TestKeeper_Service_Call(t *testing.T) {
	ctx, keeper, accs := createTestInput(t, sdk.NewIntWithDecimal(2000, 18), 2)

	author := accs[0].GetAddress()
	provider := accs[1].GetAddress()
	consumer := author

	serviceName := "myService"
	chainID := "testnet"
	serviceDescription := "the service for unit test"
	tags := []string{"test", "tutorial"}
	authorDescription := "unit test author"

	err := keeper.AddServiceDefinition(ctx, serviceName, chainID, serviceDescription, tags, author, authorDescription, idlContent)
	require.NoError(t, err)

	deposit, _ := sdk.IrisCoinType.ConvertToMinDenomCoin("1000iris")
	price, _ := sdk.IrisCoinType.ConvertToMinDenomCoin("1iris")
	bindingType := types.Global
	level := types.Level{AvgRspTime: 10000, UsableTime: 9999}

	err = keeper.AddServiceBinding(ctx, chainID, serviceName, chainID, provider, bindingType, sdk.NewCoins(deposit), []sdk.Coin{price}, level)
	require.NoError(t, err)

	// service request
	input := []byte("1234")

	svcRequest, err := keeper.AddRequest(ctx, chainID, serviceName, chainID, chainID, consumer, provider, 1, input, sdk.NewCoins(price), false)
	require.NoError(t, err)

	storedSvcRequest, found := keeper.GetActiveRequest(ctx, svcRequest.ExpirationHeight, svcRequest.RequestHeight, svcRequest.RequestIntraTxCounter)
	require.True(t, found)
	require.Equal(t, svcRequest.RequestID(), storedSvcRequest.RequestID())

	iterator := keeper.ActiveRequestQueueIterator(ctx, svcRequest.ExpirationHeight)
	defer iterator.Close()

	require.True(t, iterator.Valid())
	for ; ; iterator.Next() {
		var req types.SvcRequest
		if !iterator.Valid() {
			break
		}

		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &req)
		require.Equal(t, svcRequest.RequestID(), req.RequestID())
	}
}

const idlContent = `
	syntax = "proto3";

	// The greeting service definition.
	service Greeter {
		//@Attribute description:sayHello
		//@Attribute output_privacy:NoPrivacy
		//@Attribute output_cached:NoCached
		rpc SayHello (HelloRequest) returns (HelloReply) {}
	}

	// The request message containing the user's name.
	message HelloRequest {
		string name = 1;
	}

	// The response message containing the greetings
	message HelloReply {
		string message = 1;
	}`
