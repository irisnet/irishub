package service

import (
	"testing"

	sdk "github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_service_Definition(t *testing.T) {
	mapp, keeper, _, addrs, _, _ := getMockApp(t, 3)
	SortAddresses(addrs)
	mapp.BeginBlock(abci.RequestBeginBlock{})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})
	keeper.ck.AddCoins(ctx, addrs[1], sdk.Coins{sdk.NewCoin("iris", sdk.NewInt(1100))})

	serviceDef := NewSvcDef("myService",
		"testnet",
		"the service for unit test",
		[]string{"test", "tutorial"},
		addrs[0],
		"unit test author",
		idlContent)

	keeper.AddServiceDefinition(ctx, serviceDef)
	serviceDefB, _ := keeper.GetServiceDefinition(ctx, "testnet", "myService")

	require.Equal(t, serviceDefB.IDLContent, idlContent)
	require.Equal(t, serviceDefB.Name, "myService")

	// test methods
	keeper.AddMethods(ctx, serviceDef)
	iterator := keeper.GetMethods(ctx, "testnet", "myService")
	require.True(t, iterator.Valid())
	for ; ; iterator.Next() {
		var method MethodProperty
		if !iterator.Valid() {
			break
		}
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &method)
		require.Equal(t, method.Name, "SayHello")
		require.Equal(t, method.Description, "sayHello")
		require.Equal(t, method.OutputPrivacy.String(), "NoPrivacy")
		require.Equal(t, method.OutputCached.String(), "NoCached")
	}

	// test binding
	svcBinding := NewSvcBinding(ctx, "testnet", "myService", "testnet",
		addrs[1], Global, sdk.Coins{sdk.NewCoin("iris", sdk.NewInt(1000))}, []sdk.Coin{{"iris", sdk.NewInt(1)}},
		Level{AvgRspTime: 10000, UsableTime: 9999}, true)
	err, _ := keeper.AddServiceBinding(ctx, svcBinding)
	require.NoError(t, err)

	require.True(t, keeper.ck.HasCoins(ctx, addrs[1], sdk.Coins{sdk.NewCoin("iris", sdk.NewInt(100))}))

	gotSvcBinding, found := keeper.GetServiceBinding(ctx, svcBinding.DefChainID, svcBinding.DefName, svcBinding.BindChainID, svcBinding.Provider)
	require.True(t, found)
	require.True(t, SvcBindingEqual(svcBinding, gotSvcBinding))

	// test binding update
	svcBindingUpdate := NewSvcBinding(ctx, "testnet", "myService", "testnet",
		addrs[1], Global, sdk.Coins{sdk.NewCoin("iris", sdk.NewInt(100))}, []sdk.Coin{{"iris", sdk.NewInt(1)}},
		Level{AvgRspTime: 10000, UsableTime: 9999}, true)
	err, _ = keeper.UpdateServiceBinding(ctx, svcBindingUpdate)
	require.NoError(t, err)

	require.True(t, keeper.ck.HasCoins(ctx, addrs[1], sdk.Coins{sdk.NewCoin("iris", sdk.NewInt(0))}))

	upSvcBinding, found := keeper.GetServiceBinding(ctx, svcBinding.DefChainID, svcBinding.DefName, svcBinding.BindChainID, svcBinding.Provider)
	require.True(t, found)
	require.True(t, upSvcBinding.Deposit.IsEqual(gotSvcBinding.Deposit.Plus(svcBindingUpdate.Deposit)))
}

func TestKeeper_service_Call(t *testing.T) {
	mapp, keeper, _, addrs, _, _ := getMockApp(t, 3)
	SortAddresses(addrs)
	mapp.BeginBlock(abci.RequestBeginBlock{})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})
	keeper.ck.AddCoins(ctx, addrs[1], sdk.Coins{sdk.NewCoin("iris", sdk.NewInt(1100))})
	keeper.ck.AddCoins(ctx, addrs[2], sdk.Coins{sdk.NewCoin("iris", sdk.NewInt(1100))})

	serviceDef := NewSvcDef("myService",
		"testnet",
		"the service for unit test",
		[]string{"test", "tutorial"},
		addrs[0],
		"unit test author",
		idlContent)

	keeper.AddServiceDefinition(ctx, serviceDef)

	svcBinding := NewSvcBinding(ctx, "testnet", "myService", "testnet",
		addrs[1], Global, sdk.Coins{sdk.NewCoin("iris", sdk.NewInt(1000))}, []sdk.Coin{{"iris", sdk.NewInt(1)}},
		Level{AvgRspTime: 10000, UsableTime: 9999}, true)
	keeper.AddServiceBinding(ctx, svcBinding)

	// service request
	svcRequest := NewSvcRequest("testnet", "myService", "testnet", "testnet",
		addrs[2], addrs[1], 1, []byte("1234"), sdk.Coins{sdk.NewCoin("iris", sdk.NewInt(1))}, false)
	svcRequest, err := keeper.AddRequest(ctx, svcRequest)
	require.NoError(t, err)

	svcRequest1, found := keeper.GetActiveRequest(ctx, svcRequest.ExpirationHeight, svcRequest.RequestHeight, svcRequest.RequestIntraTxCounter)
	require.True(t, found)
	require.Equal(t, svcRequest.RequestID(), svcRequest1.RequestID())

	iterator := keeper.ActiveRequestQueueIterator(ctx, svcRequest.ExpirationHeight)
	require.True(t, iterator.Valid())
	for ; ; iterator.Next() {
		var req SvcRequest
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
