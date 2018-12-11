package gov

import (
	"fmt"
	sdk "github.com/irisnet/irishub/types"
)

var _ Proposal = (*HaltProposal)(nil)

type HaltProposal struct {
	TextProposal
}

func (sp *HaltProposal) Execute(ctx sdk.Context, k Keeper) error {
	logger := ctx.Logger().With("module", "x/gov")

	if k.GetTerminatorHeight(ctx) == -1 {
		k.SetTerminatorHeight(ctx, ctx.BlockHeight()+k.GetTerminatorPeriod(ctx))
		logger.Info("Execute TerminatorProposal begin", "info", fmt.Sprintf("Terminator height:%d", k.GetTerminatorHeight(ctx)))
	} else {
		logger.Info("Terminator Period is in process.", "info", fmt.Sprintf("Terminator height:%d", k.GetTerminatorHeight(ctx)))

	}
	return nil
}

// Key for getting a the next available proposalID from the store
var (
	KeyTerminatorHeight = []byte("TerminatorHeight")
	KeyTerminatorPeriod = []byte("TerminatorPeriod")
)

func (keeper Keeper) GetTerminatorHeight(ctx sdk.Context) int64 {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyTerminatorHeight)
	if bz == nil {
		return -1
	}
	var height int64
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &height)

	return height
}

func (keeper Keeper) SetTerminatorHeight(ctx sdk.Context, height int64) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(height)
	store.Set(KeyTerminatorHeight, bz)
}

func (keeper Keeper) GetTerminatorPeriod(ctx sdk.Context) int64 {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyTerminatorPeriod)
	if bz == nil {
		return -1
	}
	var height int64
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &height)

	return height
}

func (keeper Keeper) SetTerminatorPeriod(ctx sdk.Context, height int64) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(height)
	store.Set(KeyTerminatorPeriod, bz)
}
