package types

import (
	math "math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (pool FarmPool) ExpiredHeight() (int64, error) {
	var targetInteval = int64(math.MaxInt64)
	for _, r := range pool.Rules {
		inteval := r.TotalReward.Quo(r.RewardPerBlock).Int64()
		if targetInteval > inteval {
			targetInteval = inteval
		}
	}
	if int64(math.MaxInt64)-pool.StartHeight < targetInteval {
		return 0, sdkerrors.Wrapf(sdkerrors.ErrInvalidHeight, "endheight overflow")
	}
	return pool.StartHeight + targetInteval, nil
}

func (pool FarmPool) RemainingHeight() int64 {
	if len(pool.Rules) == 0 {
		return 0
	}
	targetInteval := pool.Rules[0].RemainingReward.Quo(pool.Rules[0].RewardPerBlock).Int64()
	for _, r := range pool.Rules[1:] {
		inteval := r.RemainingReward.Quo(r.RewardPerBlock).Int64()
		if targetInteval > inteval {
			targetInteval = inteval
		}
	}
	return targetInteval
}

func (pool FarmPool) CaclRewards(farmInfo FarmInfo, deltaAmt sdk.Int) (rewards, rewardDebt sdk.Coins) {
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

func (rs RewardRules) UpdateWith(rewarPerBlockd sdk.Coins) RewardRules {
	for i := range rs {
		rewardAmt := rewarPerBlockd.AmountOf(rs[i].Reward)
		if rewardAmt.IsPositive() {
			rs[i].RewardPerBlock = rewardAmt
		}
	}
	return rs
}
