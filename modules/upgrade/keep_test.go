package upgrade

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestUpdateKeeper(t *testing.T) {
	ctx, keeper := createTestInput(t)
	router := baseapp.NewRouter()
	router.AddRoute("main", []*sdk.KVStoreKey{sdk.NewKVStoreKey("main")}, nil)
	router.AddRoute("acc", []*sdk.KVStoreKey{sdk.NewKVStoreKey("acc")}, nil)
	router.AddRoute("ibc", []*sdk.KVStoreKey{sdk.NewKVStoreKey("ibc")}, nil)
	router.AddRoute("stake", []*sdk.KVStoreKey{sdk.NewKVStoreKey("stake")}, nil)
	router.AddRoute("upgrade", []*sdk.KVStoreKey{sdk.NewKVStoreKey("upgrade")}, nil)

	moduleList := getModuleList(router)

	genesisVersion := NewVersion(0, 10, 0, moduleList)
	keeper.AddNewVersion(ctx, genesisVersion)

	version := keeper.GetCurrentVersion(ctx)
	if version == nil || version.Id != genesisVersion.Id {
		t.FailNow()
	}

	router.AddRoute("slashing", []*sdk.KVStoreKey{sdk.NewKVStoreKey("slashing")}, nil)
	moduleList = getModuleList(router)
	version1 := NewVersion(0, 15, 1000, moduleList)
	keeper.AddNewVersion(ctx, version1)

	version = keeper.GetCurrentVersion(ctx)
	if version == nil || version.Id != 1 {
		t.FailNow()
	}

	version = keeper.GetVersionByProposalId(ctx, 15)
	if version == nil || version.Id != 1 {
		t.FailNow()
	}

	versionList := keeper.GetVersionList(ctx)
	if versionList == nil || len(versionList) != 2 {
		t.FailNow()
	}

	version = keeper.GetVersionByVersionId(ctx, 1)
	if version == nil || version.Id != 1 {
		t.FailNow()
	}

	version = keeper.GetVersionByHeight(ctx, 1000)
	if version == nil || version.Id != 1 || version.Start > 1000 {
		t.FailNow()
	}

	version = keeper.GetVersionByHeight(ctx, 1001)
	if version == nil || version.Id != 1 || version.Start > 10011 {
		t.FailNow()
	}

	router.AddRoute("gov", []*sdk.KVStoreKey{sdk.NewKVStoreKey("gov")}, nil)
	router.AddRoute("fee", []*sdk.KVStoreKey{sdk.NewKVStoreKey("fee")}, nil)
	moduleList = getModuleList(router)
	version2 := NewVersion(0, 24, 2000, moduleList)
	keeper.AddNewVersion(ctx, version2)
	version = keeper.GetCurrentVersion(ctx)
	if version == nil && version.Start != 2000 {
		t.FailNow()
	}
	versionList = keeper.GetVersionList(ctx)
	if versionList == nil || len(versionList) != 3 {
		t.FailNow()
	}

}

func TestSwitchKeeper(t *testing.T) {
	//ctx, keeper := createTestInput(t)

	router := baseapp.NewRouter()
	router.AddRoute("main-0", []*sdk.KVStoreKey{sdk.NewKVStoreKey("main")}, nil)
	router.AddRoute("acc-0", []*sdk.KVStoreKey{sdk.NewKVStoreKey("acc")}, nil)
	router.AddRoute("gov-0", []*sdk.KVStoreKey{sdk.NewKVStoreKey("gov")}, nil)
	router.AddRoute("stake-0", []*sdk.KVStoreKey{sdk.NewKVStoreKey("stake")}, nil)
	router.AddRoute("upgrade-0", []*sdk.KVStoreKey{sdk.NewKVStoreKey("upgrade")}, nil)

	RegisterModuleList(router)
	require.Equal(t, len(ModuleList), 5)
	require.Equal(t, len(ModuleListBucket[0]), 5)

	router.AddRoute("gov-1", []*sdk.KVStoreKey{sdk.NewKVStoreKey("gov")}, nil)
	RegisterModuleList(router)
	require.Equal(t, len(ModuleList), 6)
	require.Equal(t, len(ModuleListBucket[1]), 1)

	router.AddRoute("gov-2", []*sdk.KVStoreKey{sdk.NewKVStoreKey("gov")}, nil)
	RegisterModuleList(router)
	require.Equal(t, len(ModuleList), 7)
	require.Equal(t, len(ModuleListBucket[2]), 1)

	router.AddRoute("gov-3", []*sdk.KVStoreKey{sdk.NewKVStoreKey("gov")}, nil)
	RegisterModuleList(router)
	require.Equal(t, len(ModuleList), 8)
	require.Equal(t, len(ModuleListBucket[3]), 1)

	router.AddRoute("stake-3", []*sdk.KVStoreKey{sdk.NewKVStoreKey("stake")}, nil)
	RegisterModuleList(router)
	require.Equal(t, len(ModuleList), 9)
	require.Equal(t, len(ModuleListBucket[3]), 2)
}

func getModuleList(router baseapp.Router) ModuleLifeTimeList {

	modulelist := NewModuleLifeTimeList()
	handlerList := router.RouteTable()

	for _, handler := range handlerList {
		hs := strings.Split(handler, "/")

		stores := strings.Split(hs[1], ":")
		modulelist = modulelist.BuildModuleLifeTime(0, hs[0], stores)
	}

	return modulelist
}

func TestKeeper_InitGenesis_commidID(t *testing.T) {
	ctx, keeper := createTestInput(t)
	router := baseapp.NewRouter()
	router.AddRoute("main", []*sdk.KVStoreKey{sdk.NewKVStoreKey("main")}, nil)
	router.AddRoute("acc", []*sdk.KVStoreKey{sdk.NewKVStoreKey("acc")}, nil)
	router.AddRoute("ibc", []*sdk.KVStoreKey{sdk.NewKVStoreKey("ibc")}, nil)
	router.AddRoute("stake", []*sdk.KVStoreKey{sdk.NewKVStoreKey("stake")}, nil)
	router.AddRoute("upgrade", []*sdk.KVStoreKey{sdk.NewKVStoreKey("upgrade")}, nil)
	router.AddRoute("upgradeI", []*sdk.KVStoreKey{sdk.NewKVStoreKey("upgrade")}, nil)

	moduleList := getModuleList(router)

	genesisVersion := NewVersion(0, 10, 0, moduleList)
	keeper.AddNewVersion(ctx, genesisVersion)
	InitGenesis_commitID(ctx, keeper)
	fmt.Println(keeper.GetKVStoreKeylist(ctx))

	keeper.SetCurrentProposalAcceptHeight(ctx, 1234234000)
	fmt.Println(keeper.GetCurrentProposalAcceptHeight(ctx))
}
