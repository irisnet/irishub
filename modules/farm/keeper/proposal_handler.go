package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/farm/types"
)

// HandleCommunityPoolCreateFarmProposal is a handler for executing a passed community spend proposal
func HandleCommunityPoolCreateFarmProposal(ctx sdk.Context,
	k Keeper,
	p *types.CommunityPoolCreateFarmProposal) error {

	//check valid lp token denom
	if err := k.validateLPToken(ctx, p.LptDenom); err != nil {
		return sdkerrors.Wrapf(
			types.ErrInvalidLPToken,
			"The lp token denom[%s] is not exist",
			p.LptDenom,
		)
	}

	// Check if the community pool has enough coins to create the farm pool
	err := k.distributeFromFeePool(ctx, p.TotalReward)
	if err != nil {
		return err
	}
	creator := k.ak.GetModuleAddress(k.communityPoolName)
	pool, err := k.createPool(ctx, creator, p.PoolDescription, ctx.BlockHeight(), false, p.LptDenom, p.TotalReward, p.RewardPerBlock)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreatePool,
			sdk.NewAttribute(types.AttributeValuePoolId, pool.Id),
			sdk.NewAttribute(types.AttributeValueAmount, sdk.NewCoins(p.TotalReward...).String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})
	return nil
}

// distributeFromFeePool distributes funds from the distribution module account to
// farm module address while updating the community pool
func (k Keeper) distributeFromFeePool(ctx sdk.Context, amount sdk.Coins) error {
	feePool := k.dk.GetFeePool(ctx)

	// NOTE the community pool isn't a module account, however its coins
	// are held in the distribution module account. Thus the community pool
	// must be reduced separately from the SendCoinsFromModuleToAccount call
	newPool, negative := feePool.CommunityPool.SafeSub(sdk.NewDecCoinsFromCoins(amount...))
	if negative {
		return types.ErrBadDistribution
	}

	feePool.CommunityPool = newPool
	err := k.bk.SendCoinsFromModuleToModule(ctx, k.communityPoolName, types.ModuleName, amount)
	if err != nil {
		return err
	}

	k.dk.SetFeePool(ctx, feePool)
	return nil
}

// refundToFeePool return the remaining funds of the farm pool to CommunityPool
func (k Keeper) refundToFeePool(ctx sdk.Context, refundTotal sdk.Coins) error {
	//refund the total remaining reward to creator
	if err := k.bk.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.communityPoolName, refundTotal); err != nil {
		return err
	}
	feelPool := k.dk.GetFeePool(ctx)
	feelPool.CommunityPool = feelPool.CommunityPool.Add(sdk.NewDecCoinsFromCoins(sdk.NewCoins(refundTotal...)...)...)
	k.dk.SetFeePool(ctx, feelPool)
	return nil
}
