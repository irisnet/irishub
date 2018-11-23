package gov

import (
sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TerminateTagKey		= "terminate_blockchain"
	TerminateTagValue	= "true"
)

var _ Proposal = (*TerminatorProposal)(nil)

type TerminatorProposal struct {
	TextProposal
}

func (sp *TerminatorProposal) Execute(ctx sdk.Context, k Keeper) error {
	//logger := ctx.Logger().With("module", "x/gov")

	return nil
}

// Key for getting a the next available proposalID from the store
var (
	KeyTerminatorHeight = []byte("TerminatorHeight")
)

// Get Proposal from store by ProposalID
func (keeper Keeper) GetTerminatorHeight(ctx sdk.Context) int64 {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyTerminatorHeight)
	if bz == nil {
		return 0
	}

	var height int64
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &height)

	return height
}

// Implements sdk.AccountKeeper.
func (keeper Keeper) SetTerminatorHeight(ctx sdk.Context, height int64) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(height)
	store.Set(KeyTerminatorHeight, bz)
}
