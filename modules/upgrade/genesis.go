package upgrade

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"strings"
)

// InitGenesis - build the genesis version For first Version
func InitGenesis(ctx sdk.Context, k Keeper, router bam.Router) {
	modulelist := NewModuleLifeTimeList()
	handlerList := router.RouteTable()

	for _, handler := range handlerList {
		hs := strings.Split(handler, "/")

		modulelist = modulelist.BuildModuleLifeTime(0, hs[0], hs[1])
	}

	genesisVersion := NewVersion(0, 0,0, modulelist)
	k.AddNewVersion(ctx, genesisVersion)

}
