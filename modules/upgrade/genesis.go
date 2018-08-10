package upgrade

import (
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis - build the genesis version For first Version
func InitGenesis(ctx sdk.Context, k Keeper, router bam.Router) {

	RegisterModuleList(router)

	moduleList, found := GetModuleListFromBucket(0)
	if !found {
		panic("No module list info found for genesis version")
	}

	genesisVersion := NewVersion(0, 0, 0, moduleList)
	k.AddNewVersion(ctx, genesisVersion)

	k.SetCurrentProposalAcceptHeight(ctx,-1)
	k.SetCurrentProposalID(ctx,-1)
	InitGenesis_commitID(ctx, k)
}
