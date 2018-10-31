package iservice

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_IService_Definition(t *testing.T) {
	ctx, keeper := createTestInput(t)

	serviceDef := NewSvcDef("myService",
		"testnet",
		"the iservice for unit test",
		[]string{"test", "tutorial"},
		addrs[0],
		"unit test author",
		idlContent,
		Unicast)

	keeper.AddServiceDefinition(ctx, serviceDef)
	serviceDefB, _ := keeper.GetServiceDefinition(ctx, "testnet", "myService")

	require.Equal(t, serviceDefB.IDLContent, idlContent)
	require.Equal(t, serviceDefB.Name, "myService")
	require.Equal(t, serviceDefB.Messaging, Unicast)

	// test methods
	keeper.AddMethods(ctx, serviceDef)
	iterator := keeper.GetMethods(ctx, "testnet", "myService")
	require.True(t, iterator.Valid())
	for ; ; iterator.Next() {
		var method MethodProperty
		if !iterator.Valid() {
			break
		}
		keeper.cdc.MustUnmarshalBinary(iterator.Value(), &method)
		require.Equal(t, method.Name, "SayHello")
		require.Equal(t, method.Description, "sayHello")
		require.Equal(t, method.OutputPrivacy.String(), "NoPrivacy")
		require.Equal(t, method.OutputCached.String(), "NoCached")
	}
}

func TestKeeper_IService_Binding(t *testing.T) {
	ctx, keeper := createTestInput(t)

	// test binding
	svcBinding := NewSvcBinding("testnet", "myService", "testnet",
		addrs[1], Local, sdk.Coins{sdk.NewCoin("iris", sdk.NewInt(100))}, []sdk.Coin{{"iris", sdk.NewInt(100)}},
		[]int{1}, 1000)
	keeper.AddServiceBinding(ctx, svcBinding)
	gotSvcBinding, found := keeper.GetServiceBinding(ctx, svcBinding.DefChainID, svcBinding.DefName, svcBinding.BindChainID, svcBinding.Provider)
	require.True(t, found)
	require.True(t, SvcBindingEqual(svcBinding, gotSvcBinding))
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
