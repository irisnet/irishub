package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/farm/types"
)

// CreatePool creates an new farm pool
func (k Keeper) CreatePool(ctx sdk.Context, name string,
	description string,
	lpTokenDenom string,
	startHeight uint64,
	rewardPerBlock sdk.Coins,
	totalReward sdk.Coins,
	destructible bool,
	creator sdk.AccAddress,
) error {
	//Escrow total reward
	if err := k.bk.SendCoinsFromAccountToModule(ctx,
		creator, types.ModuleName, totalReward); err != nil {
		return err
	}

	//send CreatePoolFee to feeCollectorName
	if err := k.bk.SendCoinsFromAccountToModule(ctx,
		creator, k.feeCollectorName, sdk.NewCoins(k.CreatePoolFee(ctx))); err != nil {
		return err
	}

	pool := types.FarmPool{
		Name:               name,
		Creator:            creator.String(),
		Description:        description,
		StartHeight:        startHeight,
		Destructible:       destructible,
		TotalLpTokenLocked: sdk.NewCoin(lpTokenDenom, sdk.ZeroInt()),
		Rules:              []types.RewardRule{},
	}

	//save farm rule
	for _, total := range totalReward {
		rewardRule := types.RewardRule{
			Reward:          total.Denom,
			TotalReward:     total.Amount,
			RemainingReward: total.Amount,
			RewardPerBlock:  rewardPerBlock.AmountOf(total.Denom),
			RewardPerShare:  sdk.ZeroDec(),
		}
		k.SetRewardRule(ctx, name, rewardRule)
		pool.Rules = append(pool.Rules, rewardRule)
	}
	pool.EndHeight = pool.ExpiredHeight()
	//save farm pool
	k.SetPool(ctx, pool)
	// put to expired farm pool queue
	k.EnqueueActivePool(ctx, name, pool.EndHeight)
	return nil
}

// Destroy destroy an exist farm pool
func (k Keeper) DestroyPool(ctx sdk.Context, poolName string,
	creator sdk.AccAddress) error {
	pool, exist := k.GetPool(ctx, poolName)
	if !exist {
		return sdkerrors.Wrapf(types.ErrNotExistPool, "not exist pool [%s]", poolName)
	}

	if creator.String() != pool.Creator {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "creator [%s] is not the creator of the pool", creator.String())
	}

	if !pool.Destructible {
		return sdkerrors.Wrapf(
			types.ErrInvalidOperate, "pool [%s] is not destructible", poolName)
	}

	if pool.IsExpired(ctx.BlockHeight()) {
		return sdkerrors.Wrapf(types.ErrExpiredPool,
			"pool [%s] has expired at height[%d], current [%d]",
			poolName,
			pool.EndHeight,
			ctx.BlockHeight(),
		)
	}

	if err := k.Refund(ctx, pool); err != nil {
		return sdkerrors.Wrapf(types.ErrNotExistPool, "not exist pool [%s]", poolName)
	}
	return nil
}

// AppendReward creates an new farm pool
func (k Keeper) AppendReward(ctx sdk.Context, poolName string,
	reward sdk.Coins,
	creator sdk.AccAddress,
) (remaining sdk.Coins, err error) {
	pool, exist := k.GetPool(ctx, poolName)
	if !exist {
		return remaining, sdkerrors.Wrapf(types.ErrNotExistPool, "not exist pool [%s]", poolName)
	}

	if creator.String() != pool.Creator {
		return remaining, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "creator [%s] is not the creator of the pool", creator.String())
	}

	if pool.IsExpired(ctx.BlockHeight()) {
		return remaining, sdkerrors.Wrapf(types.ErrExpiredPool,
			"pool [%s] has expired at height[%d], current [%d]",
			poolName,
			pool.EndHeight,
			ctx.BlockHeight(),
		)
	}

	rules := k.GetRewardRules(ctx, poolName)
	if !rules.Contains(reward) {
		return remaining, sdkerrors.Wrapf(types.ErrInvalidAppend, reward.String())
	}

	if err := k.bk.SendCoinsFromAccountToModule(ctx,
		creator, types.ModuleName, reward); err != nil {
		return remaining, err
	}

	for i := range rules {
		rules[i].TotalReward = rules[i].TotalReward.Add(reward.AmountOf(rules[i].Reward))
		rules[i].RemainingReward = rules[i].RemainingReward.Add(reward.AmountOf(rules[i].Reward))
		remaining = remaining.Add(sdk.NewCoin(rules[i].Reward, rules[i].RemainingReward))
		k.SetRewardRule(ctx, poolName, rules[i])
	}
	pool.Rules = rules

	//if the expiration height does not change, there is no need to update the pool and the expired queue
	expiredHeight := pool.ExpiredHeight()
	if expiredHeight == pool.EndHeight {
		return remaining, nil
	}

	// remove from Expired Pool at old height
	k.DequeueActivePool(ctx, poolName, pool.EndHeight)
	pool.EndHeight = expiredHeight
	k.SetPool(ctx, pool)
	// put to expired farm pool queue at new height
	k.EnqueueActivePool(ctx, poolName, pool.EndHeight)
	return remaining, nil
}

// Stake is responsible for the user to mortgage the lp token to the system and get back the reward accumulated before then
func (k Keeper) Stake(ctx sdk.Context, poolName string,
	lpToken sdk.Coin,
	sender sdk.AccAddress,
) (reward sdk.Coins, err error) {
	pool, exist := k.GetPool(ctx, poolName)
	if !exist {
		return reward, sdkerrors.Wrapf(types.ErrNotExistPool, "not exist pool [%s]", poolName)
	}

	if pool.StartHeight > uint64(ctx.BlockHeight()) {
		return reward, sdkerrors.Wrapf(types.ErrNotStartPool,
			"farm pool [%s] will start at height[%d], current [%d]",
			poolName,
			pool.StartHeight,
			ctx.BlockHeight(),
		)
	}

	if pool.IsExpired(ctx.BlockHeight()) {
		return reward, sdkerrors.Wrapf(types.ErrExpiredPool,
			"pool [%s] has expired at height[%d], current [%d]",
			poolName,
			pool.EndHeight,
			ctx.BlockHeight(),
		)
	}

	if lpToken.Denom != pool.TotalLpTokenLocked.Denom {
		return reward, sdkerrors.Wrapf(types.ErrNotMatch,
			"pool [%s] only accept [%s] token, but got [%s]",
			poolName, pool.TotalLpTokenLocked.Denom, lpToken.Denom)
	}

	if err := k.bk.SendCoinsFromAccountToModule(ctx,
		sender, types.ModuleName, sdk.NewCoins(lpToken)); err != nil {
		return reward, err
	}

	//update pool reward shards
	pool, _, err = k.UpdatePool(ctx, pool, lpToken.Amount, false)
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
		if err = k.bk.SendCoinsFromModuleToAccount(ctx, types.RewardCollector,
			sender, rewards); err != nil {
			return reward, err
		}
	}

	farmInfo.RewardDebt = rewardDebt
	farmInfo.Locked = farmInfo.Locked.Add(lpToken.Amount)
	k.SetFarmInfo(ctx, farmInfo)
	return rewards, nil
}

// Unstake withdraw lp token from farm pool
func (k Keeper) Unstake(ctx sdk.Context, poolName string,
	lpToken sdk.Coin,
	sender sdk.AccAddress) (_ sdk.Coins, err error) {
	pool, exist := k.GetPool(ctx, poolName)
	if !exist {
		return nil, sdkerrors.Wrapf(types.ErrNotExistPool, poolName)
	}

	//lpToken demon must be same as pool.TotalLpTokenLocked.Denom
	if lpToken.Denom != pool.TotalLpTokenLocked.Denom {
		return nil, sdkerrors.Wrapf(types.ErrNotMatch,
			"pool [%s] only accept [%s] token, but got [%s]",
			poolName, pool.TotalLpTokenLocked.Denom, lpToken.Denom)
	}

	//farmInfo must be exist
	farmInfo, exist := k.GetFarmInfo(ctx, poolName, sender.String())
	if !exist {
		return nil, sdkerrors.Wrapf(types.ErrNotExistFarmer,
			"farmer [%s] not found in pool[%s]",
			sender.String(),
			poolName,
		)
	}

	//the lp token unstaked must be less than staked
	if farmInfo.Locked.LT(lpToken.Amount) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
			"farmer locked lp token %s, but unstake %s",
			farmInfo.Locked.String(),
			lpToken.Amount.String(),
		)
	}

	//the lp token unstaked must be less than pool
	if pool.TotalLpTokenLocked.Amount.LT(lpToken.Amount) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
			"farmer locked lp token %s, but farm pool total: %s",
			farmInfo.Locked.String(),
			pool.TotalLpTokenLocked.Amount.String(),
		)
	}

	if pool.IsExpired(ctx.BlockHeight()) {
		//If the farm has ended, the reward rules cannot be updated
		pool.Rules = k.GetRewardRules(ctx, pool.Name)
		pool.TotalLpTokenLocked = pool.TotalLpTokenLocked.Sub(lpToken)
		k.SetPool(ctx, pool)
	} else {
		//update pool reward shards
		pool, _, err = k.UpdatePool(ctx, pool, lpToken.Amount.Neg(), false)
		if err != nil {
			return nil, err
		}
	}

	//unstake lpToken to sender account
	if err = k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName,
		sender, sdk.NewCoins(lpToken)); err != nil {
		return nil, err
	}

	//compute farmer rewards
	rewards, rewardDebt := pool.CaclRewards(farmInfo, lpToken.Amount.Neg())
	if rewards.IsAllPositive() {
		//distribute reward
		if err = k.bk.SendCoinsFromModuleToAccount(ctx, types.RewardCollector,
			sender, rewards); err != nil {
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
func (k Keeper) Harvest(ctx sdk.Context, poolName string,
	sender sdk.AccAddress) (sdk.Coins, error) {
	pool, exist := k.GetPool(ctx, poolName)
	if !exist {
		return nil, sdkerrors.Wrapf(types.ErrNotExistPool, "not exist pool [%s]", poolName)
	}

	if pool.IsExpired(ctx.BlockHeight()) {
		return nil, sdkerrors.Wrapf(types.ErrExpiredPool,
			"pool [%s] has expired at height[%d], current [%d]",
			poolName,
			pool.EndHeight,
			ctx.BlockHeight(),
		)
	}

	farmInfo, exist := k.GetFarmInfo(ctx, poolName, sender.String())
	if !exist {
		return nil, sdkerrors.Wrapf(types.ErrNotExistFarmer,
			"farmer [%s] not found in pool[%s]",
			sender.String(),
			poolName,
		)
	}

	amtAdded := sdk.ZeroInt()
	//update pool reward shards
	pool, _, err := k.UpdatePool(ctx, pool, amtAdded, false)
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
func (k Keeper) Refund(ctx sdk.Context, pool types.FarmPool) error {
	pool, _, err := k.UpdatePool(ctx, pool, sdk.ZeroInt(), true)
	if err != nil {
		return err
	}

	creator, err := sdk.AccAddressFromBech32(pool.Creator)
	if err != nil {
		return err
	}

	var refundTotal sdk.Coins
	for _, r := range pool.Rules {
		refundTotal = refundTotal.Add(sdk.NewCoin(r.Reward, r.RemainingReward))

		r.RemainingReward = sdk.ZeroInt()
		k.SetRewardRule(ctx, pool.Name, r)
	}

	if refundTotal.IsAllPositive() {
		//refund the total remaining reward to creator
		if err := k.bk.SendCoinsFromModuleToAccount(ctx,
			types.ModuleName, creator, refundTotal); err != nil {
			return err
		}
	}
	//remove record
	k.DequeueActivePool(ctx, pool.Name, pool.EndHeight)
	return nil
}

// UpdatePool is responsible for updating the status information of the farm pool, including the total accumulated bonus from the last time the bonus was distributed to the present, the current remaining bonus in the farm pool, and the ratio of the current liquidity token to the bonus.

// Note that when multiple transactions at the same block height trigger the farm pool update at the same time, only the first transaction will trigger the `RewardPerShare` update operation

// This method returns the updated farm pool and the bonuses collected in this period
func (k Keeper) UpdatePool(ctx sdk.Context,
	pool types.FarmPool,
	amount sdk.Int,
	isDestroy bool,
) (types.FarmPool, sdk.Coins, error) {
	height := uint64(ctx.BlockHeight())
	if height < pool.LastHeightDistrRewards {
		return pool, nil, sdkerrors.Wrapf(types.ErrExpiredHeight,
			"invalid height: %d, last distribution height: %d",
			height,
			pool.LastHeightDistrRewards,
		)
	}

	rules := k.GetRewardRules(ctx, pool.Name)
	if len(rules) == 0 {
		return pool, nil, sdkerrors.Wrapf(types.ErrNotExistPool,
			"the rule of the farm pool[%s] not exist", pool.Name)
	}
	var rewardTotal sdk.Coins
	//when there are multiple farm operations in the same block, the value needs to be updated once
	if height > pool.LastHeightDistrRewards &&
		pool.TotalLpTokenLocked.Amount.GT(sdk.ZeroInt()) {
		blockInterval := height - pool.LastHeightDistrRewards
		for i := range rules {
			rewardCollected := rules[i].RewardPerBlock.MulRaw(int64(blockInterval))
			coinCollected := sdk.NewCoin(rules[i].Reward, rewardCollected)
			if rules[i].RemainingReward.LT(rewardCollected) {
				k.Logger(ctx).Error("The remaining amount is not enough to pay the bonus",
					"poolName", pool.Name,
					"remainingReward", rules[i].RemainingReward.String(),
					"rewardCollected", rewardCollected.String(),
				)
				return pool, nil, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
					"the remaining reward of the pool[%s] is %s, but got %s",
					pool.Name,
					sdk.NewCoin(rules[i].Reward, rules[i].RemainingReward).String(),
					coinCollected,
				)
			}
			newRewardPerShare := sdk.NewDecFromInt(rewardCollected).
				QuoInt(pool.TotalLpTokenLocked.Amount)
			rules[i].RewardPerShare = rules[i].RewardPerShare.Add(newRewardPerShare)
			rules[i].RemainingReward = rules[i].RemainingReward.Sub(rewardCollected)

			rewardTotal = rewardTotal.Add(coinCollected)
			k.SetRewardRule(ctx, pool.Name, rules[i])
		}
	}

	//escrow the collected rewards to the `RewardCollector` account
	if rewardTotal.IsAllPositive() {
		if err := k.bk.SendCoinsFromModuleToModule(ctx,
			types.ModuleName, types.RewardCollector, rewardTotal); err != nil {
			return pool, rewardTotal, err
		}
	}

	pool.TotalLpTokenLocked = sdk.NewCoin(
		pool.TotalLpTokenLocked.Denom,
		pool.TotalLpTokenLocked.Amount.Add(amount),
	)
	pool.LastHeightDistrRewards = uint64(ctx.BlockHeight())
	if isDestroy {
		pool.EndHeight = uint64(ctx.BlockHeight())
	}
	pool.Rules = rules
	k.SetPool(ctx, pool)
	return pool, rewardTotal, nil
}
