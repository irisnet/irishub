package upgrade

import (
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis - build the genesis version For first Version
func InitGenesis(ctx sdk.Context, k Keeper, router bam.Router) {

	RegisterModuleList(router)

	genesisVersion := NewVersion(0, 0, 0, GetModuleListFromBucket(0))
	k.AddNewVersion(ctx, genesisVersion)

	k.SetCurrentProposalAcceptHeight(ctx,-1)
	k.SetCurrentProposalID(ctx,-1)
	InitGenesis_commitID(ctx, k)
}
