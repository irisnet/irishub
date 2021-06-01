package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/farm/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Pools(goctx context.Context,
	request *types.QueryPoolsRequest) (*types.QueryPoolsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goctx)

	var pools []types.FarmPool
	if len(request.Name) > 0 {
		pool, exist := k.GetPool(ctx, request.Name)
		if !exist {
			return nil, sdkerrors.Wrapf(types.ErrNotExistPool, "not found pool: %s", request.Name)
		}
		pools = append(pools, pool)
	} else {
		k.IteratorAllPools(ctx, func(pool types.FarmPool) {
			pools = append(pools, pool)
		})
	}

	var list []*types.FarmPoolEntry
	for _, pool := range pools {
		var totalReward sdk.Coins
		var remainingReward sdk.Coins
		var rewardPerBlock sdk.Coins
		k.IteratorRewardRules(ctx, pool.Name, func(r types.RewardRule) {
			totalReward = totalReward.Add(sdk.NewCoin(r.Reward, r.TotalReward))
			remainingReward = remainingReward.Add(sdk.NewCoin(r.Reward, r.RemainingReward))
			rewardPerBlock = rewardPerBlock.Add(sdk.NewCoin(r.Reward, r.RewardPerBlock))
		})

		list = append(list, &types.FarmPoolEntry{
			Name:               pool.Name,
			Creator:            pool.Creator,
			Description:        pool.Description,
			StartHeight:        pool.StartHeight,
			EndHeight:          pool.EndHeight,
			Destructible:       pool.Destructible,
			Expired:            pool.IsExpired(ctx.BlockHeight()),
			TotalLpTokenLocked: pool.TotalLpTokenLocked,
			TotalReward:        totalReward,
			RemainingReward:    remainingReward,
			RewardPerBlock:     rewardPerBlock,
		})
	}
	return &types.QueryPoolsResponse{List: list}, nil
}

func (k Keeper) Farmer(goctx context.Context,
	request *types.QueryFarmerRequest) (*types.QueryFarmerResponse, error) {
	var list []*types.LockedInfo
	var err error
	var farmInfos []types.FarmInfo

	ctx := sdk.UnwrapSDKContext(goctx)
	cacheCtx, _ := ctx.CacheContext()
	if len(request.PoolName) == 0 {
		k.IteratorFarmInfo(cacheCtx, request.Farmer, func(farmInfo types.FarmInfo) {
			farmInfos = append(farmInfos, farmInfo)
		})
	} else {
		farmInfo, existed := k.GetFarmInfo(cacheCtx, request.PoolName, request.Farmer)
		if existed {
			farmInfos = append(farmInfos, farmInfo)
		}
	}
	if len(farmInfos) == 0 {
		return nil, sdkerrors.Wrapf(types.ErrNotExistFarmer, "not found farmer: %s", request.Farmer)
	}

	for _, farmer := range farmInfos {
		pool, exist := k.GetPool(cacheCtx, farmer.PoolName)
		if !exist {
			return nil, sdkerrors.Wrapf(types.ErrNotExistPool, "not exist pool [%s]", farmer.PoolName)
		}

		//The farm pool has not started, no reward
		if pool.StartHeight > uint64(ctx.BlockHeight()) {
			list = append(list, &types.LockedInfo{
				PoolName: farmer.PoolName,
				Locked:   sdk.NewCoin(pool.TotalLpTokenLocked.Denom, farmer.Locked),
			})
			continue
		}

		if !pool.IsExpired(ctx.BlockHeight()) {
			pool, _, err = k.UpdatePool(cacheCtx, pool, sdk.ZeroInt(), false)
			if err != nil {
				return nil, err
			}
		} else {
			pool.Rules = k.GetRewardRules(ctx, pool.Name)
		}

		rewards, _ := pool.CaclRewards(farmer, sdk.ZeroInt())
		list = append(list, &types.LockedInfo{
			PoolName:      farmer.PoolName,
			Locked:        sdk.NewCoin(pool.TotalLpTokenLocked.Denom, farmer.Locked),
			PendingReward: rewards,
		})
	}

	return &types.QueryFarmerResponse{
		List:   list,
		Height: ctx.BlockHeight(),
	}, nil
}

func (k Keeper) Params(goctx context.Context,
	request *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goctx)
	return &types.QueryParamsResponse{
		Params: k.GetParams(ctx),
	}, nil
}
