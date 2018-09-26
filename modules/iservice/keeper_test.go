package iservice

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const idl = "message SearchRequest { required string query = 1; optional int32 page_number = 2; optional int32 result_per_page = 3; }"

func TestKeeper_IService_Definition(t *testing.T) {
	ctx, keeper := createTestInput(t)

	serviceDef := NewMsgSvcDef("myService",
		"testnet",
		"the iservice for unit test",
		[]string{"test", "tutorial"},
		addrs[0],
		"unit test author",
		idl,
		Broadcast)

	keeper.AddServiceDefinition(ctx, serviceDef)

	serviceDefB := keeper.GetServiceDefinition(ctx, "testnet", "myService")

	require.Equal(t, serviceDefB.IDLContent, idl)
	require.Equal(t, serviceDefB.Name, "myService")
	require.Equal(t, serviceDefB.Broadcast, Broadcast)
}
