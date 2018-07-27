package upgrade

import (
	"testing"
	"strings"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/baseapp"
	)


func TestUpdateKeeper(t *testing.T) {
	ctx, keeper := createTestInput(t)
	router := baseapp.NewRouter()
	router.AddRoute("main",sdk.NewKVStoreKey("main"),nil)
	router.AddRoute("acc",sdk.NewKVStoreKey("acc"),nil)
	router.AddRoute("ibc",sdk.NewKVStoreKey("ibc"),nil)
	router.AddRoute("stake",sdk.NewKVStoreKey("stake"),nil)
	router.AddRoute("upgrade",sdk.NewKVStoreKey("upgrade"),nil)

	moduleList := getModuleList(router)

	genesisVersion := NewVersion(0, 10,0, moduleList)
	keeper.AddNewVersion(ctx, genesisVersion)

	version := keeper.GetCurrentVersion(ctx)
	if version == nil || version.Id != genesisVersion.Id {
		t.FailNow()
	}

	router.AddRoute("slashing",sdk.NewKVStoreKey("slashing"),nil)
	moduleList = getModuleList(router)
	version1 := NewVersion(0, 15,1000, moduleList)
	keeper.AddNewVersion(ctx, version1)

	version = keeper.GetCurrentVersion(ctx)
	if version == nil || version.Id != 1 {
		t.FailNow()
	}

	version = keeper.GetVersionByProposalId(ctx,15)
	if version == nil || version.Id != 1 {
		t.FailNow()
	}

	versionList := keeper.GetVersionList(ctx)
	if versionList == nil || len(versionList) != 2 {
		t.FailNow()
	}

	version = keeper.GetVersionByVersionId(ctx,1)
	if version == nil || version.Id != 1 {
		t.FailNow()
	}

	version = keeper.GetVersionByHeight(ctx,1000)
	if version == nil || version.Id != 1 || version.Start > 1000 {
		t.FailNow()
	}

	version = keeper.GetVersionByHeight(ctx,1001)
	if version == nil || version.Id != 1 || version.Start > 10011 {
		t.FailNow()
	}

	router.AddRoute("gov",sdk.NewKVStoreKey("gov"),nil)
	router.AddRoute("fee",sdk.NewKVStoreKey("fee"),nil)
	moduleList = getModuleList(router)
	version2 := NewVersion(0, 24,2000, moduleList)
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

func getModuleList(router baseapp.Router) ModuleLifeTimeList {
	modulelist := NewModuleLifeTimeList()
	handlerList := router.RouteTable()

	for _, handler := range handlerList {
		hs := strings.Split(handler, "/")
		modulelist = modulelist.BuildModuleLifeTime(0, hs[0], hs[1])
	}
	return modulelist
}