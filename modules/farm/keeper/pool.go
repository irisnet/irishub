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
	startHeight int64,
	rewardPerBlock sdk.Coins,
	totalReward sdk.Coins,
	editable bool,
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
		Editable:           editable,
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

	endHeight, err := pool.ExpiredHeight()
	if err != nil {
		return err
	}
	//save farm pool
	pool.EndHeight = endHeight
	k.SetPool(ctx, pool)
	// put to expired farm pool queue
	k.EnqueueActivePool(ctx, name, pool.EndHeight)
	return nil
}

// Destroy destroy an exist farm pool
func (k Keeper) DestroyPool(ctx sdk.Context, poolName string,
	creator sdk.AccAddress) (sdk.Coins, error) {
	pool, exist := k.GetPool(ctx, poolName)
	if !exist {
		return nil, sdkerrors.Wrapf(types.ErrPoolNotFound, poolName)
	}

	if creator.String() != pool.Creator {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "creator [%s] is not the creator of the pool", creator.String())
	}

	if !pool.Editable {
		return nil, sdkerrors.Wrapf(
			types.ErrInvalidOperate, "pool [%s] is not editable", poolName)
	}

	if k.Expired(ctx, pool) {
		return nil, sdkerrors.Wrapf(types.ErrPoolExpired,
			"pool [%s] has expired at height [%d], current [%d]",
			poolName,
			pool.EndHeight,
			ctx.BlockHeight(),
		)
	}
	return k.Refund(ctx, pool)
}

// AdjustPool adjusts farm pool parameters
func (k Keeper) AdjustPool(ctx sdk.Context,
	poolName string,
	reward sdk.Coins,
	rewardPerBlock sdk.Coins,
	creator sdk.AccAddress,
) (err error) {
	pool, exist := k.GetPool(ctx, poolName)
	//check if the liquidity pool exists
	if !exist {
		return sdkerrors.Wrapf(types.ErrPoolNotFound, poolName)
	}

	if !pool.Editable {
		return sdkerrors.Wrapf(
			types.ErrInvalidOperate, "pool [%s] is not editable", poolName)
	}

	//check permissions
	if creator.String() != pool.Creator {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "creator [%s] is not the creator of the pool", creator.String())
	}

	//check for expiration
	if k.Expired(ctx, pool) {
		return sdkerrors.Wrapf(types.ErrPoolExpired,
			"pool [%s] has expired at height[%d], current [%d]",
			poolName,
			pool.EndHeight,
			ctx.BlockHeight(),
		)
	}

	//update pool reward shards
	pool, _, err = k.updatePool(ctx, pool, sdk.ZeroInt(), false)
	if err != nil {
		return err
	}

	rules := types.RewardRules(pool.Rules)
	if rewardPerBlock != nil && !rewardPerBlock.DenomsSubsetOf(rules.RewardsPerBlock()) {
		return sdkerrors.Wrapf(types.ErrInvalidAppend, "rewardPerBlock: %s", rewardPerBlock.String())
	}

	availableReward := sdk.NewCoins()
	remainingHeight := pool.EndHeight - ctx.BlockHeight()

	if reward != nil {
		if !rules.Contains(reward) {
			return sdkerrors.Wrapf(types.ErrInvalidAppend, reward.String())
		}

		if err := k.bk.SendCoinsFromAccountToModule(ctx,
			creator, types.ModuleName, reward); err != nil {
			return err
		}

		for i := range rules {
			coin := sdk.NewCoin(rules[i].Reward, rules[i].RewardPerBlock.Mul(sdk.NewInt(remainingHeight)))
			availableReward = availableReward.Add(coin)
			rules[i].TotalReward = rules[i].TotalReward.Add(reward.AmountOf(rules[i].Reward))
			rules[i].RemainingReward = rules[i].RemainingReward.Add(reward.AmountOf(rules[i].Reward))
		}
	}

	if rewardPerBlock != nil {
		pool.Rules = types.RewardRules(rules).UpdateWith(rewardPerBlock)
	}
	k.SetRewardRules(ctx, pool.Name, pool.Rules)

	//expiredHeight = [(srcEndHeight-curHeight)*srcRewardPerBlock +appendReward]/RewardPerBlock + curHeight
	availableReward = availableReward.Add(reward...)
	rewardsPerBlock := types.RewardRules(pool.Rules).RewardsPerBlock()
	availableHeight := availableReward[0].Amount.
		Quo(rewardsPerBlock.AmountOf(availableReward[0].Denom)).Int64()
	for _, c := range availableReward[1:] {
		rpb := rewardsPerBlock.AmountOf(c.Denom)
		inteval := c.Amount.Quo(rpb).Int64()
		if availableHeight > inteval {
			availableHeight = inteval
		}
	}
	expiredHeight := ctx.BlockHeight() + availableHeight
	//if the expiration height does not change,
	// there is no need to update the pool and the expired queue
	if expiredHeight == pool.EndHeight {
		return nil
	}
	// remove from Expired Pool at old height
	k.DequeueActivePool(ctx, pool.Name, pool.EndHeight)
	pool.EndHeight = expiredHeight
	k.SetPool(ctx, pool)
	// put to expired farm pool queue at new height
	k.EnqueueActivePool(ctx, pool.Name, pool.EndHeight)
	return nil
}

// updatePool is responsible for updating the status information of the farm pool, including the total accumulated bonus from the last time the bonus was distributed to the present, the current remaining bonus in the farm pool, and the ratio of the current liquidity token to the bonus.

// Note that when multiple transactions at the same block height trigger the farm pool update at the same time, only the first transaction will trigger the `RewardPerShare` update operation

// updatePool returns the updated farm pool and the reward collected in this period
func (k Keeper) updatePool(ctx sdk.Context,
	pool types.FarmPool,
	amount sdk.Int,
	isDestroy bool,
) (types.FarmPool, sdk.Coins, error) {
	height := ctx.BlockHeight()
	if height < pool.LastHeightDistrRewards {
		return pool, nil, sdkerrors.Wrapf(types.ErrExpiredHeight,
			"invalid height: [%d], last distribution height: [%d]",
			height,
			pool.LastHeightDistrRewards,
		)
	}

	rules := k.GetRewardRules(ctx, pool.Name)
	if len(rules) == 0 {
		return pool, nil, sdkerrors.Wrapf(types.ErrPoolNotFound, pool.Name)
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
					"the remaining reward of the pool [%s] is [%s], but got [%s]",
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
	pool.LastHeightDistrRewards = ctx.BlockHeight()
	if isDestroy {
		pool.EndHeight = ctx.BlockHeight()
		if pool.StartHeight > pool.EndHeight {
			pool.StartHeight = pool.EndHeight
		}
	}
	pool.Rules = rules
	k.SetPool(ctx, pool)
	return pool, rewardTotal, nil
}
