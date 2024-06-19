package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"irismod.io/farm/types"
)

// HandleCreateFarmProposal is a handler for executing a passed community spend proposal
func (k Keeper) HandleCreateFarmProposal(ctx sdk.Context, p *types.CommunityPoolCreateFarmProposal) error {
	total := sdk.NewCoins(p.FundApplied...).Add(p.FundSelfBond...)
	// thansfer the coins to ModuleName account from EscrowCollector
	if err := k.bk.SendCoinsFromModuleToModule(ctx, types.EscrowCollector, types.ModuleName, total); err != nil {
		return err
	}

	creator := k.ak.GetModuleAddress(k.communityPoolName)
	pool, err := k.createPool(ctx, creator, p.PoolDescription, ctx.BlockHeight(), false, p.LptDenom, total, p.RewardPerBlock)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreatePool,
			sdk.NewAttribute(types.AttributeValuePoolId, pool.Id),
			sdk.NewAttribute(types.AttributeValueAmount, sdk.NewCoins(total...).String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})
	return nil
}

// escrowFromFeePool distributes funds from the distribution module account to
// farm module address while updating the community pool
func (k Keeper) escrowFromFeePool(ctx sdk.Context, amount sdk.Coins) error {
	feePool := k.dk.GetFeePool(ctx)

	// NOTE the community pool isn't a module account, however its coins
	// are held in the distribution module account. Thus the community pool
	// must be reduced separately from the SendCoinsFromModuleToAccount call
	newPool, negative := feePool.CommunityPool.SafeSub(sdk.NewDecCoinsFromCoins(amount...))
	if negative {
		return types.ErrBadDistribution
	}

	feePool.CommunityPool = newPool
	err := k.bk.SendCoinsFromModuleToModule(ctx, k.communityPoolName, types.EscrowCollector, amount)
	if err != nil {
		return err
	}

	k.dk.SetFeePool(ctx, feePool)
	return nil
}

// refundToFeePool return the remaining funds of the farm pool to CommunityPool
func (k Keeper) refundToFeePool(ctx sdk.Context, fromModule string, refundTotal sdk.Coins) error {
	//refund the total remaining reward to creator
	if err := k.bk.SendCoinsFromModuleToModule(ctx, fromModule, k.communityPoolName, refundTotal); err != nil {
		return err
	}
	feelPool := k.dk.GetFeePool(ctx)
	feelPool.CommunityPool = feelPool.CommunityPool.Add(sdk.NewDecCoinsFromCoins(sdk.NewCoins(refundTotal...)...)...)
	k.dk.SetFeePool(ctx, feelPool)
	return nil
}

func (k Keeper) refundEscrow(ctx sdk.Context, info types.EscrowInfo) {
	proposer, err := sdk.AccAddressFromBech32(info.Proposer)
	if err != nil {
		return
	}
	//refund the amount locked by the user
	if err := k.bk.SendCoinsFromModuleToAccount(ctx,
		types.EscrowCollector, proposer, info.FundSelfBond); err != nil {
		return
	}

	//refund the amount locked by the CommunityPool
	if err := k.refundToFeePool(ctx, types.EscrowCollector, sdk.NewCoins(info.FundApplied...)); err != nil {
		return
	}
	k.deleteEscrowInfo(ctx, info.ProposalId)
	k.Logger(ctx).Info("execute refundEscrow",
		"proposalID", info.ProposalId,
		"proposer", info.Proposer,
		"fundSelfBond", info.FundSelfBond,
		"communityPool", k.communityPoolName,
		"fundApplied", info.FundApplied,
	)
}

func (k Keeper) SetEscrowInfo(ctx sdk.Context, info types.EscrowInfo) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&info)
	store.Set(types.KeyEscrowInfo(info.ProposalId), bz)
}

func (k Keeper) GetEscrowInfo(ctx sdk.Context, proposalId uint64) (types.EscrowInfo, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyEscrowInfo(proposalId))
	if bz == nil {
		return types.EscrowInfo{}, false
	}
	var info types.EscrowInfo
	k.cdc.MustUnmarshal(bz, &info)
	return info, true
}

func (k Keeper) GetAllEscrowInfo(ctx sdk.Context) (infos []types.EscrowInfo) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.EscrowInfoKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var info types.EscrowInfo
		k.cdc.MustUnmarshal(iterator.Value(), &info)
		infos = append(infos, info)
	}
	return
}

func (k Keeper) deleteEscrowInfo(ctx sdk.Context, proposalId uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyEscrowInfo(proposalId))
}
