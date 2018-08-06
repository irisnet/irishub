package upgrade

import (
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

// InitGenesis - build the genesis version For first Version
func InitGenesis(ctx sdk.Context, k Keeper, router bam.Router) {

	RegisterModuleList(router)

	modulelist := NewModuleLifeTimeList()
	handlerList := router.RouteTable()

	for _, handler := range handlerList {
		hs := strings.Split(handler, "/")
		stores := strings.Split(hs[1], ":")
		modulelist = modulelist.BuildModuleLifeTime(0, hs[0], stores)
	}

	genesisVersion := NewVersion(0, 0, 0, modulelist)
	k.AddNewVersion(ctx, genesisVersion)

	InitGenesis_commitID(ctx, k)
}
