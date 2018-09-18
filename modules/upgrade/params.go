package upgrade

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ParamStoreKeyCurrentProposalAcceptHeight = "upgrade/proposalAcceptHeight"
)

func GetCurrentProposalAcceptHeightKey() string {
	return ParamStoreKeyCurrentProposalAcceptHeight
}

func (k Keeper) GetCurrentProposalAcceptHeight(ctx sdk.Context) int64 {
	var height int64
	err := k.params.Get(ctx, GetCurrentProposalAcceptHeightKey(), &height)
	if err != nil {
		panic(err)
	}
	return height
}

func (k Keeper) SetCurrentProposalAcceptHeight(ctx sdk.Context, height int64) {
	k.params.Set(ctx, GetCurrentProposalAcceptHeightKey(), height)
}
