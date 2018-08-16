package upgrade

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ParamStoreKeyCurrentProposalAcceptHeight = "upgrade/proposalAcceptHeight"
	ParamStoreKeyCurrentProposalId           = "upgrade/proposalId"
)

func GetCurrentProposalAcceptHeightKey() string {
	return ParamStoreKeyCurrentProposalAcceptHeight
}

func GetCurrentProposalIdKey() string {
	return ParamStoreKeyCurrentProposalId
}

func (k Keeper) GetCurrentProposalID(ctx sdk.Context) int64 {
	var proposalID int64
	err := k.ps.GovGetter().Get(ctx, GetCurrentProposalIdKey(), &proposalID)
	if err != nil {
		panic(err)
	}
	return proposalID
}

func (k Keeper) GetCurrentProposalAcceptHeight(ctx sdk.Context) int64 {
	var height int64
	err := k.ps.GovGetter().Get(ctx, GetCurrentProposalAcceptHeightKey(), &height)
	if err != nil {
		panic(err)
	}
	return height
}

func (k Keeper) SetCurrentProposalID(ctx sdk.Context, proposalID int64) {
	k.ps.GovSetter().Set(ctx, GetCurrentProposalIdKey(), proposalID)
}

func (k Keeper) SetCurrentProposalAcceptHeight(ctx sdk.Context, height int64) {
	k.ps.GovSetter().Set(ctx, GetCurrentProposalAcceptHeightKey(), height)
}
