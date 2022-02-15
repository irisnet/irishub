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

	moduleAddress := k.ak.GetModuleAddress(types.ModuleName)
	// Check if the community pool has enough coins to create the farm pool
	err := k.dk.DistributeFromFeePool(ctx, p.TotalRewards, moduleAddress)
	if err != nil {
		return err
	}
	creator := k.dk.GetDistributionAccount(ctx)
	pool, err := k.createPool(ctx, creator.GetAddress(), p.PoolDescription, ctx.BlockHeight(), false, p.LptDenom, p.TotalRewards, p.RewardsPerBlock)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreatePool,
			sdk.NewAttribute(types.AttributeValuePoolId, pool.Id),
			sdk.NewAttribute(types.AttributeValueAmount, sdk.NewCoins(p.TotalRewards...).String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})
	return nil
}
