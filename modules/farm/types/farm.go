package types

import (
	math "math"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (pool FarmPool) Started(ctx sdk.Context) bool {
	return pool.StartHeight <= ctx.BlockHeight()
}

func (pool FarmPool) ExpiredHeight() (int64, error) {
	targetInteval := int64(math.MaxInt64)
	for _, r := range pool.Rules {
		inteval := r.TotalReward.Quo(r.RewardPerBlock).Int64()
		if targetInteval > inteval {
			targetInteval = inteval
		}
	}
	if int64(math.MaxInt64)-pool.StartHeight < targetInteval {
		return 0, errorsmod.Wrapf(sdkerrors.ErrInvalidHeight, "endheight overflow")
	}
	return pool.StartHeight + targetInteval, nil
}

func (pool FarmPool) CaclRewards(farmInfo FarmInfo, deltaAmt sdkmath.Int) (rewards, rewardDebt sdk.Coins) {
	for _, r := range pool.Rules {
		if farmInfo.Locked.GT(sdk.ZeroInt()) {
			pendingRewardTotal := r.RewardPerShare.MulInt(farmInfo.Locked).TruncateInt()
			pendingReward := pendingRewardTotal.Sub(farmInfo.RewardDebt.AmountOf(r.Reward))
			rewards = rewards.Add(sdk.NewCoin(r.Reward, pendingReward))
		}

		locked := farmInfo.Locked.Add(deltaAmt)
		debt := sdk.NewCoin(r.Reward, r.RewardPerShare.MulInt(locked).TruncateInt())
		rewardDebt = rewardDebt.Add(debt)
	}
	return rewards, rewardDebt
}

type RewardRules []RewardRule

func (rs RewardRules) Contains(reward sdk.Coins) bool {
	if len(rs) < len(reward) {
		return false
	}
	var allRewards sdk.Coins
	for _, r := range rs {
		allRewards = allRewards.Add(sdk.NewCoin(r.Reward, r.RemainingReward))
	}
	return reward.DenomsSubsetOf(allRewards)
}

func (rs RewardRules) UpdateWith(rewardPerBlock sdk.Coins) RewardRules {
	if rewardPerBlock == nil {
		return rs
	}
	for i := range rs {
		rewardAmt := rewardPerBlock.AmountOf(rs[i].Reward)
		if rewardAmt.IsPositive() {
			rs[i].RewardPerBlock = rewardAmt
		}
	}
	return rs
}

func (rs RewardRules) RewardsPerBlock() (coins sdk.Coins) {
	for _, r := range rs {
		coins = coins.Add(sdk.NewCoin(r.Reward, r.RewardPerBlock))
	}
	return coins
}

func (rs RewardRules) TotalReward() (total sdk.Coins) {
	for _, r := range rs {
		total = total.Add(sdk.NewCoin(r.Reward, r.TotalReward))
	}
	return total
}
