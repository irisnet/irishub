package iservice

import (
	"github.com/stretchr/testify/require"
	"testing"
)

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

func TestKeeper_IService_Definition(t *testing.T) {
	ctx, keeper := createTestInput(t)

	serviceDef := NewMsgSvcDef("myService",
		"testnet",
		"the iservice for unit test",
		[]string{"test", "tutorial"},
		addrs[0],
		"unit test author",
		idlContent,
		Broadcast)

	keeper.AddServiceDefinition(ctx, serviceDef)
	serviceDefB, _ := keeper.GetServiceDefinition(ctx, "testnet", "myService")

	require.Equal(t, serviceDefB.IDLContent, idlContent)
	require.Equal(t, serviceDefB.Name, "myService")
	require.Equal(t, serviceDefB.Broadcast, Broadcast)

	keeper.AddMethods(ctx, serviceDef)
	iterator := keeper.GetMethods(ctx, "testnet", "myService")
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

}
