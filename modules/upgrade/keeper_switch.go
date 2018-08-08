package upgrade

import (
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
	"strconv"
)

var (
	VersionToBeSwitched Version
	ModuleList          ModuleLifeTimeList
	ModuleListBucket	map[int] ModuleLifeTimeList
)

func (keeper Keeper) GetVersionToBeSwitched() *Version {
	return &VersionToBeSwitched
}

func RegisterModuleList(router bam.Router) {
	ModuleList = NewModuleLifeTimeList()
	handlerList := router.RouteTable()

	for _, handler := range handlerList {
		hs := strings.Split(handler, "/")

		stores := strings.Split(hs[1], ":")
		ModuleList = ModuleList.BuildModuleLifeTime(0, hs[0], stores)
	}

	buildModuleListBucket()
}

func buildModuleListBucket() {

	for _, module := range ModuleList {
		verstr := strings.Split(module.Handler, "-")
		ver, err := strconv.Atoi(verstr[1])
		if err != nil {
			panic(err)
		}

		bucket, ok := ModuleListBucket[ver]
		if ok {
			ModuleListBucket[ver] = bucket.BuildModuleLifeTime(0, verstr[0], module.Store)
		} else {
			modulelist := NewModuleLifeTimeList()
			ModuleListBucket[ver] = modulelist.BuildModuleLifeTime(0, verstr[0], module.Store)
		}
	}
}

func (keeper Keeper) RegisterVersionToBeSwitched(store sdk.KVStore, router bam.Router) {
	currentVersion := keeper.GetCurrentVersionByStore(store)

	if currentVersion == nil { // waiting to create the genesis version
		return
	}

	modulelist := NewModuleLifeTimeList()
	handlerList := router.RouteTable()

	for _, handler := range handlerList {
		hs := strings.Split(handler, "/")

		stores := strings.Split(hs[1], ":")
		modulelist = modulelist.BuildModuleLifeTime(0, hs[0], stores)
	}

	VersionToBeSwitched = NewVersion(currentVersion.Id+1, 0, 0, modulelist)
}

func (k Keeper) SetDoingSwitch(ctx sdk.Context, doing bool) {
	kvStore := ctx.KVStore(k.storeKey)

	var bytes []byte
	if doing {
		bytes = []byte{byte(1)}
	} else {
		bytes = []byte{byte(0)}
	}
	kvStore.Set(GetDoingSwitchKey(), bytes)
}

func (k Keeper) GetDoingSwitch(ctx sdk.Context) bool {
	kvStore := ctx.KVStore(k.storeKey)

	bytes := kvStore.Get(GetDoingSwitchKey())
	if len(bytes) == 1 {
		return bytes[0] == byte(1)
	}

	return false
}

func (k Keeper) DoSwitchBegin(ctx sdk.Context) {
	k.SetDoingSwitch(ctx, true)
}

func (k Keeper) DoSwitchEnd(ctx sdk.Context) {
	VersionToBeSwitched.ProposalID = k.GetCurrentProposalID(ctx)
	VersionToBeSwitched.Start = ctx.BlockHeight()

	k.AddNewVersion(ctx, VersionToBeSwitched)

	k.SetDoingSwitch(ctx, false)
	k.SetCurrentProposalID(ctx, -1)
	k.SetKVStoreKeylist(ctx)
}
