package upgrade

import (
	"fmt"
	"github.com/irisnet/irishub/baseapp"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"github.com/irisnet/irishub/modules/upgrade/params"
	"github.com/irisnet/irishub/modules/params"
)

func TestUpdateKeeper(t *testing.T) {
	ctx, keeper, _ := createTestInput(t)
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

	kvStore := ctx.KVStore(keeper.storeKey)
	keeper.RefreshVersionList(kvStore)

	version = keeper.GetVersionByVersionId(1)
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
	router.AddRoute("main", []*sdk.KVStoreKey{sdk.NewKVStoreKey("main")}, nil)
	router.AddRoute("acc", []*sdk.KVStoreKey{sdk.NewKVStoreKey("acc")}, nil)
	router.AddRoute("gov", []*sdk.KVStoreKey{sdk.NewKVStoreKey("gov")}, nil)
	router.AddRoute("stake", []*sdk.KVStoreKey{sdk.NewKVStoreKey("stake")}, nil)
	router.AddRoute("upgrade", []*sdk.KVStoreKey{sdk.NewKVStoreKey("upgrade")}, nil)

	RegisterModuleList(router)
	require.Equal(t, len(ModuleListBucket[0]), 5)
	require.Equal(t, "gov", GetModuleFromBucket(0, "gov").Handler)
	require.Equal(t, "stake", GetModuleFromBucket(0, "stake").Handler)

	router.AddRoute("gov-1", []*sdk.KVStoreKey{sdk.NewKVStoreKey("gov")}, nil)
	Inited = false
	RegisterModuleList(router)
	require.Equal(t, len(ModuleListBucket[1]), 5)
	require.Equal(t, GetModuleFromBucket(1, "gov").Handler, "gov-1")
	require.Equal(t, GetModuleFromBucket(1, "stake").Handler, "stake")

	router.AddRoute("gov-2", []*sdk.KVStoreKey{sdk.NewKVStoreKey("gov")}, nil)
	Inited = false
	RegisterModuleList(router)
	require.Equal(t, len(ModuleListBucket[2]), 5)
	require.Equal(t, GetModuleFromBucket(2, "gov").Handler, "gov-2")
	require.Equal(t, GetModuleFromBucket(2, "stake").Handler, "stake")

	router.AddRoute("gov-3", []*sdk.KVStoreKey{sdk.NewKVStoreKey("gov")}, nil)
	Inited = false
	RegisterModuleList(router)
	require.Equal(t, len(ModuleListBucket[3]), 5)
	require.Equal(t, GetModuleFromBucket(3, "gov").Handler, "gov-3")
	require.Equal(t, GetModuleFromBucket(3, "stake").Handler, "stake")

	router.AddRoute("stake-3", []*sdk.KVStoreKey{sdk.NewKVStoreKey("stake")}, nil)
	Inited = false
	RegisterModuleList(router)
	require.Equal(t, len(ModuleListBucket[3]), 5)
	require.Equal(t, GetModuleFromBucket(3, "gov").Handler, "gov-3")
	require.Equal(t, GetModuleFromBucket(3, "stake").Handler, "stake-3")
}

func TestSetKVStoreKeylist(t *testing.T) {
	ctx, keeper, paramKeeper := createTestInput(t)

	router := baseapp.NewRouter()
	router.AddRoute("main-0", []*sdk.KVStoreKey{sdk.NewKVStoreKey("main")}, nil)
	router.AddRoute("acc-0", []*sdk.KVStoreKey{sdk.NewKVStoreKey("acc")}, nil)
	router.AddRoute("gov-0", []*sdk.KVStoreKey{sdk.NewKVStoreKey("gov")}, nil)
	router.AddRoute("stake-0", []*sdk.KVStoreKey{sdk.NewKVStoreKey("stake")}, nil)
	router.AddRoute("upgrade-0", []*sdk.KVStoreKey{sdk.NewKVStoreKey("upgrade")}, nil)

	subspace := paramKeeper.Subspace("Sig").WithTypeTable(params.NewTypeTable(
		upgradeparams.CurrentUpgradeProposalIdParameter.GetStoreKey(), uint64((0)),
		upgradeparams.ProposalAcceptHeightParameter.GetStoreKey(), int64(0),
		upgradeparams.SwitchPeriodParameter.GetStoreKey(), int64(0),
	))

	upgradeparams.ProposalAcceptHeightParameter.SetReadWriter(subspace)
	upgradeparams.CurrentUpgradeProposalIdParameter.SetReadWriter(subspace)
	upgradeparams.SwitchPeriodParameter.SetReadWriter(subspace)

	InitGenesis(ctx, keeper, router, DefaultGenesisStateForTest())
	keeper.SetKVStoreKeylist(ctx)
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
	ctx, keeper, paramKeeper := createTestInput(t)
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

	subspace := paramKeeper.Subspace("Sig").WithTypeTable(params.NewTypeTable(
		upgradeparams.CurrentUpgradeProposalIdParameter.GetStoreKey(), uint64((0)),
		upgradeparams.ProposalAcceptHeightParameter.GetStoreKey(), int64(0),
		upgradeparams.SwitchPeriodParameter.GetStoreKey(), int64(0),
	))

	upgradeparams.ProposalAcceptHeightParameter.SetReadWriter(subspace)

	upgradeparams.SetProposalAcceptHeight(ctx, 1234234000)
	fmt.Println(upgradeparams.GetProposalAcceptHeight(ctx))
}