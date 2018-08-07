package upgrade

import (
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
	"strconv"
)

var (
	Inited              bool
	ModuleListBucket    map[int64]ModuleLifeTimeList
)

func RegisterModuleList(router bam.Router) {
	if Inited {
		return
	}

	moduleList := NewModuleLifeTimeList()
	handlerList := router.RouteTable()

	for _, handler := range handlerList {
		hs := strings.Split(handler, "/")

		stores := strings.Split(hs[1], ":")
		moduleList = moduleList.BuildModuleLifeTime(0, hs[0], stores)
	}

	buildModuleListBucket(moduleList)
}

func buildModuleListBucket(moduleList ModuleLifeTimeList) {
	ModuleListBucket = make(map[int64]ModuleLifeTimeList)

	for _, module := range moduleList { // bucket the module list by the introduced version id
		verstr := strings.Split(module.Handler, "-")
		var ver int
		var err error
		if len(verstr) == 1 {
			ver = 0
		} else {
			ver, err = strconv.Atoi(verstr[1])
			if err != nil {
				panic(err)
			}
		}

		bucket, ok := ModuleListBucket[int64(ver)]
		if ok {
			ModuleListBucket[int64(ver)] = append(bucket, module)
		} else {
			modulelist := NewModuleLifeTimeList()
			ModuleListBucket[int64(ver)] = append(modulelist, module)
		}
	}

	for version := 1; ; version++ {
		bucket, ok := ModuleListBucket[int64(version)]
		if !ok {
			break
		}

		modules := make(map[string]bool) // current module set(only include new version module)
		for _, module := range bucket {
			verstr := strings.Split(module.Handler, "-")
			modules[verstr[0]] = true
		}

		preBucket := ModuleListBucket[int64(version)-1]
		for _, module := range preBucket { // reuse the pre version module if no update in the new version
			verstr := strings.Split(module.Handler, "-")
			if _, ok := modules[verstr[0]]; !ok {
				bucket = append(bucket, module)
			}
		}

		ModuleListBucket[int64(version)] = bucket
	}

	Inited = true
}

func GetModuleListFromBucket(verId int64) ModuleLifeTimeList {
	moduleList, ok := ModuleListBucket[verId]
	if !ok {
		panic("No info found for the specified version module list")
	}

	return moduleList
}

func GetModuleFromBucket(verId int64, handler string) ModuleLifeTime {
	moduleList := GetModuleListFromBucket(verId)
	for _, module := range moduleList {
		verstr := strings.Split(module.Handler, "-")
		if verstr[0] == handler {
			return module
		}
	}

	return ModuleLifeTime{}
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
	currentVersion := k.GetCurrentVersion(ctx)
	if currentVersion == nil {
		panic("No current version info found")
	}

	VersionToBeSwitched := NewVersion(currentVersion.Id+1, 0, 0, GetModuleListFromBucket(currentVersion.Id+1))
	VersionToBeSwitched.ProposalID = k.GetCurrentProposalID(ctx)
	VersionToBeSwitched.Start = ctx.BlockHeight()

	k.AddNewVersion(ctx, VersionToBeSwitched)

	k.SetDoingSwitch(ctx, false)
	k.SetCurrentProposalID(ctx, -1)
	k.SetKVStoreKeylist(ctx)
}
