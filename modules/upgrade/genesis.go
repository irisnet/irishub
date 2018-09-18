package upgrade

import (
	bam "github.com/irisnet/irishub/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"fmt"
	"github.com/irisnet/irishub/modules/parameter"
	"github.com/irisnet/irishub/modules/upgrade/params"
)

// InitGenesis - build the genesis version For first Version
func InitGenesis(ctx sdk.Context, k Keeper, router bam.Router) {

	RegisterModuleList(router)

	moduleList, found := GetModuleListFromBucket(0)
	fmt.Println(moduleList)
	if !found {
		panic("No module list info found for genesis version")
	}

	genesisVersion := NewVersion(0, 0, 0, moduleList)
	k.AddNewVersion(ctx, genesisVersion)

	parameter.InitGenesisParameter(&upgradeparams.ProposalAcceptHeightParameter, ctx, -1)
	parameter.InitGenesisParameter(&upgradeparams.CurrentUpgradeProposalIdParameter, ctx, -1)
	InitGenesis_commitID(ctx, k)
}
