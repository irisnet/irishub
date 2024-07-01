package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"mods.irisnet.org/modules/farm/types"
)

// CreatePool creates an new farm pool
func (k Keeper) CreatePool(
	ctx sdk.Context,
	description string,
	lptDenom string,
	startHeight int64,
	rewardPerBlock sdk.Coins,
	totalReward sdk.Coins,
	editable bool,
	creator sdk.AccAddress,
) (*types.FarmPool, error) {
	// deduct the user's fee for creating a farm pool
	if err := k.DeductPoolCreationFee(ctx, creator); err != nil {
		return nil, err
	}
	// Escrow total reward
	if err := k.bk.SendCoinsFromAccountToModule(ctx,
		creator, types.ModuleName, totalReward); err != nil {
		return nil, err
	}
	return k.createPool(ctx, creator, description, startHeight, editable, lptDenom, totalReward, rewardPerBlock)
}

// Destroy destroy an exist farm pool
func (k Keeper) DestroyPool(ctx sdk.Context, poolId string, creator sdk.AccAddress) (sdk.Coins, error) {
	pool, exist := k.GetPool(ctx, poolId)
	if !exist {
		return nil, errorsmod.Wrapf(types.ErrPoolNotFound, poolId)
	}

	if creator.String() != pool.Creator {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "creator [%s] is not the creator of the pool", creator.String())
	}

	if !pool.Editable {
		return nil, errorsmod.Wrapf(
			types.ErrInvalidOperate, "pool [%s] is not editable", poolId)
	}

	if k.Expired(ctx, pool) {
		return nil, errorsmod.Wrapf(types.ErrPoolExpired,
			"pool [%s] has expired at height [%d], current [%d]",
			poolId,
			pool.EndHeight,
			ctx.BlockHeight(),
		)
	}
	return k.Refund(ctx, pool)
}

// AdjustPool adjusts farm pool parameters
func (k Keeper) AdjustPool(
	ctx sdk.Context,
	poolId string,
	reward sdk.Coins,
	rewardPerBlock sdk.Coins,
	creator sdk.AccAddress,
) (err error) {
	pool, exist := k.GetPool(ctx, poolId)
	// check if the liquidity pool exists
	if !exist {
		return errorsmod.Wrapf(types.ErrPoolNotFound, poolId)
	}

	if !pool.Editable {
		return errorsmod.Wrapf(
			types.ErrInvalidOperate, "pool [%s] is not editable", poolId)
	}

	// check permissions
	if creator.String() != pool.Creator {
		return errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "creator [%s] is not the creator of the pool", creator.String())
	}

	// check for expiration
	if k.Expired(ctx, pool) {
		return errorsmod.Wrapf(types.ErrPoolExpired,
			"pool [%s] has expired at height[%d], current [%d]",
			poolId,
			pool.EndHeight,
			ctx.BlockHeight(),
		)
	}

	pool.Rules = k.GetRewardRules(ctx, pool.Id)
	rules := types.RewardRules(pool.Rules)
	if rewardPerBlock != nil && !rewardPerBlock.DenomsSubsetOf(rules.RewardsPerBlock()) {
		return errorsmod.Wrapf(types.ErrInvalidAppend, "rewardPerBlock: %s", rewardPerBlock.String())
	}

	if reward != nil && !rules.Contains(reward) {
		return errorsmod.Wrapf(types.ErrInvalidAppend, reward.String())
	}

	startHeight := pool.StartHeight
	if pool.Started(ctx) {
		startHeight = ctx.BlockHeight()
	}

	// update pool reward shards
	pool, _, err = k.updatePool(ctx, pool, sdk.ZeroInt(), false)
	if err != nil {
		return err
	}

	// update pool TotalReward、RemainingReward
	rules = types.RewardRules(pool.Rules)
	if reward != nil {
		if err := k.bk.SendCoinsFromAccountToModule(ctx,
			creator, types.ModuleName, reward); err != nil {
			return err
		}
		for i := range rules {
			rules[i].TotalReward = rules[i].TotalReward.Add(reward.AmountOf(rules[i].Reward))
			rules[i].RemainingReward = rules[i].RemainingReward.Add(reward.AmountOf(rules[i].Reward))
		}
	}

	// Calculate remaining available reward
	availableReward := rules.TotalReward()
	if pool.Started(ctx) {
		remainingHeight := pool.EndHeight - startHeight
		remainingReward := sdk.NewCoins()
		for i := range rules {
			remainingReward = remainingReward.Add(
				sdk.NewCoin(
					rules[i].Reward,
					rules[i].RewardPerBlock.Mul(sdk.NewInt(remainingHeight)),
				),
			)
		}
		availableReward = remainingReward.Add(reward...)
	}

	pool.Rules = types.RewardRules(rules).UpdateWith(rewardPerBlock)
	k.SetRewardRules(ctx, pool.Id, pool.Rules)

	// expiredHeight = [(srcEndHeight-beginPoint)*srcRewardPerBlock +appendReward]/RewardPerBlock + beginPoint
	rewardsPerBlock := types.RewardRules(pool.Rules).RewardsPerBlock()
	availableHeight := availableReward[0].Amount.Quo(rewardsPerBlock.AmountOf(availableReward[0].Denom)).Int64()
	for _, c := range availableReward[1:] {
		rpb := rewardsPerBlock.AmountOf(c.Denom)
		inteval := c.Amount.Quo(rpb).Int64()
		if availableHeight > inteval {
			availableHeight = inteval
		}
	}
	expiredHeight := startHeight + availableHeight
	// if the expiration height does not change,
	// there is no need to update the pool and the expired queue
	if expiredHeight == pool.EndHeight {
		return nil
	}
	// remove from Expired Pool at old height
	k.DequeueActivePool(ctx, pool.Id, pool.EndHeight)
	pool.EndHeight = expiredHeight
	k.SetPool(ctx, pool)
	// put to expired farm pool queue at new height
	k.EnqueueActivePool(ctx, pool.Id, pool.EndHeight)
	return nil
}

func (k Keeper) createPool(
	ctx sdk.Context,
	creator sdk.AccAddress,
	description string,
	startHeight int64,
	editable bool,
	lptDenom string,
	totalReward sdk.Coins,
	rewardPerBlock sdk.Coins,
) (*types.FarmPool, error) {
	pool := types.FarmPool{
		Id:             k.genPoolId(ctx),
		Creator:        creator.String(),
		Description:    description,
		StartHeight:    startHeight,
		Editable:       editable,
		TotalLptLocked: sdk.NewCoin(lptDenom, sdk.ZeroInt()),
		Rules:          []types.RewardRule{},
	}

	for _, total := range totalReward {
		rewardRule := types.RewardRule{
			Reward:          total.Denom,
			TotalReward:     total.Amount,
			RemainingReward: total.Amount,
			RewardPerBlock:  rewardPerBlock.AmountOf(total.Denom),
			RewardPerShare:  sdk.ZeroDec(),
		}
		k.SetRewardRule(ctx, pool.Id, rewardRule)
		pool.Rules = append(pool.Rules, rewardRule)
	}

	endHeight, err := pool.ExpiredHeight()
	if err != nil {
		return nil, err
	}

	pool.EndHeight = endHeight
	k.SetPool(ctx, pool)

	k.EnqueueActivePool(ctx, pool.Id, pool.EndHeight)
	return &pool, nil
}

// updatePool is responsible for updating the status information of the farm pool, including the total accumulated bonus from the last time the bonus was distributed to the present, the current remaining bonus in the farm pool, and the ratio of the current liquidity token to the bonus.

// Note that when multiple transactions at the same block height trigger the farm pool update at the same time, only the first transaction will trigger the `RewardPerShare` update operation

// updatePool returns the updated farm pool and the reward collected in this period
func (k Keeper) updatePool(
	ctx sdk.Context,
	pool types.FarmPool,
	amount sdk.Int,
	isDestroy bool,
) (types.FarmPool, sdk.Coins, error) {
	height := ctx.BlockHeight()
	if height < pool.LastHeightDistrRewards {
		return pool, nil, errorsmod.Wrapf(
			types.ErrExpiredHeight,
			"invalid height: [%d], last distribution height: [%d]",
			height, pool.LastHeightDistrRewards,
		)
	}

	rules := k.GetRewardRules(ctx, pool.Id)
	if len(rules) == 0 {
		return pool, nil, errorsmod.Wrapf(types.ErrPoolNotFound, pool.Id)
	}
	var rewardTotal sdk.Coins
	// when there are multiple farm operations in the same block, the value needs to be updated once
	if height > pool.LastHeightDistrRewards &&
		pool.TotalLptLocked.Amount.GT(sdk.ZeroInt()) {
		blockInterval := height - pool.LastHeightDistrRewards
		for i := range rules {
			rewardCollected := rules[i].RewardPerBlock.MulRaw(int64(blockInterval))
			coinCollected := sdk.NewCoin(rules[i].Reward, rewardCollected)
			if rules[i].RemainingReward.LT(rewardCollected) {
				k.Logger(ctx).Error(
					"The remaining amount is not enough to pay the bonus",
					"poolId", pool.Id,
					"remainingReward", rules[i].RemainingReward.String(),
					"rewardCollected", rewardCollected.String(),
				)
				return pool, nil, errorsmod.Wrapf(
					sdkerrors.ErrInsufficientFunds,
					"the remaining reward of the pool [%s] is [%s], but got [%s]",
					pool.Id, sdk.NewCoin(rules[i].Reward, rules[i].RemainingReward).String(), coinCollected,
				)
			}
			newRewardPerShare := sdk.NewDecFromInt(rewardCollected).QuoInt(pool.TotalLptLocked.Amount)
			rules[i].RewardPerShare = rules[i].RewardPerShare.Add(newRewardPerShare)
			rules[i].RemainingReward = rules[i].RemainingReward.Sub(rewardCollected)

			rewardTotal = rewardTotal.Add(coinCollected)
			k.SetRewardRule(ctx, pool.Id, rules[i])
		}
	}

	// escrow the collected rewards to the `RewardCollector` account
	if rewardTotal.IsAllPositive() {
		if err := k.bk.SendCoinsFromModuleToModule(ctx, types.ModuleName, types.RewardCollector, rewardTotal); err != nil {
			return pool, rewardTotal, err
		}
	}

	pool.TotalLptLocked = sdk.NewCoin(
		pool.TotalLptLocked.Denom,
		pool.TotalLptLocked.Amount.Add(amount),
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

func (k Keeper) genPoolId(ctx sdk.Context) string {
	seq := k.GetSequence(ctx) + 1
	k.SetSequence(ctx, seq)
	return fmt.Sprintf("%s-%d", types.PrefixFarmPool, seq)
}
