package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/farm/types"
)

// Stake is responsible for the user to mortgage the lp token to the system and get back the reward accumulated before then
func (k Keeper) Stake(
	ctx sdk.Context,
	poolName string,
	lpToken sdk.Coin,
	sender sdk.AccAddress,
) (reward sdk.Coins, err error) {
	pool, exist := k.GetPool(ctx, poolName)
	if !exist {
		return reward, sdkerrors.Wrapf(types.ErrPoolNotFound, poolName)
	}

	if pool.StartHeight > ctx.BlockHeight() {
		return reward, sdkerrors.Wrapf(
			types.ErrPoolNotStart,
			"farm pool [%s] will start at height [%d], current [%d]",
			poolName, pool.StartHeight, ctx.BlockHeight(),
		)
	}

	if k.Expired(ctx, pool) {
		return reward, sdkerrors.Wrapf(
			types.ErrPoolExpired,
			"pool [%s] has expired at height [%d], current [%d]",
			poolName, pool.EndHeight, ctx.BlockHeight(),
		)
	}

	if lpToken.Denom != pool.TotalLptLocked.Denom {
		return reward, sdkerrors.Wrapf(
			types.ErrNotMatch,
			"pool [%s] only accept [%s] token, but got [%s]",
			poolName, pool.TotalLptLocked.Denom, lpToken.Denom,
		)
	}

	if err := k.bk.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(lpToken)); err != nil {
		return reward, err
	}

	//update pool reward shards
	pool, _, err = k.updatePool(ctx, pool, lpToken.Amount, false)
	if err != nil {
		return nil, err
	}

	farmInfo, exist := k.GetFarmInfo(ctx, poolName, sender.String())
	if !exist {
		farmInfo = types.FarmInfo{
			PoolName:   poolName,
			Address:    sender.String(),
			Locked:     sdk.ZeroInt(),
			RewardDebt: sdk.NewCoins(),
		}
	}

	rewards, rewardDebt := pool.CaclRewards(farmInfo, lpToken.Amount)
	//reward users
	if rewards.IsAllPositive() {
		if err = k.bk.SendCoinsFromModuleToAccount(ctx, types.RewardCollector, sender, rewards); err != nil {
			return reward, err
		}
	}

	farmInfo.RewardDebt = rewardDebt
	farmInfo.Locked = farmInfo.Locked.Add(lpToken.Amount)
	k.SetFarmInfo(ctx, farmInfo)
	return rewards, nil
}

// Unstake withdraw lp token from farm pool
func (k Keeper) Unstake(ctx sdk.Context, poolName string, lpToken sdk.Coin, sender sdk.AccAddress) (_ sdk.Coins, err error) {
	pool, exist := k.GetPool(ctx, poolName)
	if !exist {
		return nil, sdkerrors.Wrapf(types.ErrPoolNotFound, poolName)
	}

	//lpToken demon must be same as pool.TotalLptLocked.Denom
	if lpToken.Denom != pool.TotalLptLocked.Denom {
		return nil, sdkerrors.Wrapf(
			types.ErrNotMatch,
			"pool [%s] only accept [%s] token, but got [%s]",
			poolName, pool.TotalLptLocked.Denom, lpToken.Denom,
		)
	}

	//farmInfo must be exist
	farmInfo, exist := k.GetFarmInfo(ctx, poolName, sender.String())
	if !exist {
		return nil, sdkerrors.Wrapf(
			types.ErrFarmerNotFound,
			"farmer [%s] not found in pool [%s]",
			sender.String(), poolName,
		)
	}

	//the lp token unstaked must be less than staked
	if farmInfo.Locked.LT(lpToken.Amount) {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrInsufficientFunds,
			"farmer locked lp token [%s], but unstake [%s]",
			farmInfo.Locked.String(), lpToken.Amount.String(),
		)
	}

	//the lp token unstaked must be less than pool
	if pool.TotalLptLocked.Amount.LT(lpToken.Amount) {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrInsufficientFunds,
			"farmer locked lp token [%s], but farm pool total: [%s]",
			farmInfo.Locked.String(), pool.TotalLptLocked.Amount.String(),
		)
	}

	if k.Expired(ctx, pool) {
		//If the farm has ended, the reward rules cannot be updated
		pool.Rules = k.GetRewardRules(ctx, pool.Name)
		pool.TotalLptLocked = pool.TotalLptLocked.Sub(lpToken)
		k.SetPool(ctx, pool)
	} else {
		//update pool reward shards
		pool, _, err = k.updatePool(ctx, pool, lpToken.Amount.Neg(), false)
		if err != nil {
			return nil, err
		}
	}

	//unstake lpToken to sender account
	if err = k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.NewCoins(lpToken)); err != nil {
		return nil, err
	}

	//compute farmer rewards
	rewards, rewardDebt := pool.CaclRewards(farmInfo, lpToken.Amount.Neg())
	if rewards.IsAllPositive() {
		//distribute reward
		if err = k.bk.SendCoinsFromModuleToAccount(ctx, types.RewardCollector, sender, rewards); err != nil {
			return nil, err
		}
	}

	farmInfo.RewardDebt = rewardDebt
	farmInfo.Locked = farmInfo.Locked.Sub(lpToken.Amount)
	if farmInfo.Locked.IsZero() {
		k.DeleteFarmInfo(ctx, poolName, sender.String())
		return rewards, nil
	}
	k.SetFarmInfo(ctx, farmInfo)
	return rewards, nil
}

// Harvest creates an new farm pool
func (k Keeper) Harvest(ctx sdk.Context, poolName string, sender sdk.AccAddress) (sdk.Coins, error) {
	pool, exist := k.GetPool(ctx, poolName)
	if !exist {
		return nil, sdkerrors.Wrapf(types.ErrPoolNotFound, poolName)
	}

	if k.Expired(ctx, pool) {
		return nil, sdkerrors.Wrapf(
			types.ErrPoolExpired,
			"pool [%s] has expired at height [%d], current [%d]",
			poolName, pool.EndHeight, ctx.BlockHeight(),
		)
	}

	farmInfo, exist := k.GetFarmInfo(ctx, poolName, sender.String())
	if !exist {
		return nil, sdkerrors.Wrapf(
			types.ErrFarmerNotFound,
			"farmer [%s] not found in pool [%s]",
			sender.String(), poolName,
		)
	}

	amtAdded := sdk.ZeroInt()
	//update pool reward shards
	pool, _, err := k.updatePool(ctx, pool, amtAdded, false)
	if err != nil {
		return nil, err
	}

	rewards, rewardDebt := pool.CaclRewards(farmInfo, amtAdded)
	//reward users
	if rewards.IsAllPositive() {
		if err = k.bk.SendCoinsFromModuleToAccount(ctx, types.RewardCollector, sender, rewards); err != nil {
			return nil, err
		}
	}

	farmInfo.RewardDebt = rewardDebt
	k.SetFarmInfo(ctx, farmInfo)
	return rewards, nil
}

// Refund refund the remaining reward to pool creator
func (k Keeper) Refund(ctx sdk.Context, pool types.FarmPool) (sdk.Coins, error) {
	//remove from active Pool
	k.DequeueActivePool(ctx, pool.Name, pool.EndHeight)
	pool, _, err := k.updatePool(ctx, pool, sdk.ZeroInt(), true)
	if err != nil {
		return nil, err
	}

	creator, err := sdk.AccAddressFromBech32(pool.Creator)
	if err != nil {
		return nil, err
	}

	var refundTotal sdk.Coins
	for _, r := range pool.Rules {
		refundTotal = refundTotal.Add(sdk.NewCoin(r.Reward, r.RemainingReward))
		r.RemainingReward = sdk.ZeroInt()
		k.SetRewardRule(ctx, pool.Name, r)
	}

	if refundTotal.IsAllPositive() {
		//refund the total remaining reward to creator
		if err := k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creator, refundTotal); err != nil {
			return nil, err
		}
	}
	return refundTotal, nil
}
