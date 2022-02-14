package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/farm/types"
)

// HandleCommunityPoolCreateFarmProposal is a handler for executing a passed community spend proposal
func HandleCommunityPoolCreateFarmProposal(ctx sdk.Context,
	k Keeper,
	dk types.DistrKeeper,
	p *types.CommunityPoolCreateFarmProposal) error {
	// Check if farm pool exists
	poolName := types.GenSysPoolName(p.PoolName)
	_, has := k.GetPool(ctx, poolName)
	if has {
		return sdkerrors.Wrapf(types.ErrPoolExist, p.PoolName)
	}

	//check valid lp token denom
	if err := k.validateLPToken(ctx, p.LpTokenDenom); err != nil {
		return sdkerrors.Wrapf(
			types.ErrInvalidLPToken,
			"The lp token denom[%s] is not exist",
			p.LpTokenDenom,
		)
	}

	moduleAddress := k.ak.GetModuleAddress(types.ModuleName)
	// Check if the community pool has enough coins to create the farm pool
	err := dk.DistributeFromFeePool(ctx, p.TotalRewards, moduleAddress)
	if err != nil {
		return err
	}
	creator := dk.GetDistributionAccount(ctx)
	return k.createPool(ctx, poolName, creator.GetAddress(), p.PoolDescription, ctx.BlockHeight(), false, p.LpTokenDenom, p.TotalRewards, p.RewardsPerBlock)
}
