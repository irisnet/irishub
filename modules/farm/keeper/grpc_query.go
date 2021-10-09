package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/irisnet/irismod/modules/farm/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) FarmPools(goctx context.Context, request *types.QueryFarmPoolsRequest) (*types.QueryFarmPoolsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goctx)

	var list []*types.FarmPoolEntry
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.FarmPoolKey)
	pageRes, err := query.Paginate(prefixStore, request.Pagination, func(_ []byte, value []byte) error {
		var pool types.FarmPool
		k.cdc.MustUnmarshal(value, &pool)
		var totalReward sdk.Coins
		var remainingReward sdk.Coins
		var rewardPerBlock sdk.Coins
		k.IteratorRewardRules(ctx, pool.Name, func(r types.RewardRule) {
			totalReward = totalReward.Add(sdk.NewCoin(r.Reward, r.TotalReward))
			remainingReward = remainingReward.Add(sdk.NewCoin(r.Reward, r.RemainingReward))
			rewardPerBlock = rewardPerBlock.Add(sdk.NewCoin(r.Reward, r.RewardPerBlock))
		})

		list = append(list, &types.FarmPoolEntry{
			Name:            pool.Name,
			Creator:         pool.Creator,
			Description:     pool.Description,
			StartHeight:     pool.StartHeight,
			EndHeight:       pool.EndHeight,
			Editable:        pool.Editable,
			Expired:         k.Expired(ctx, pool),
			TotalLptLocked:  pool.TotalLptLocked,
			TotalReward:     totalReward,
			RemainingReward: remainingReward,
			RewardPerBlock:  rewardPerBlock,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &types.QueryFarmPoolsResponse{
		Pools:      list,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) FarmPool(goctx context.Context,
	request *types.QueryFarmPoolRequest) (*types.QueryFarmPoolResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if len(request.Name) == 0 {
		return nil, status.Error(codes.InvalidArgument, "pool name can not be empty")
	}
	ctx := sdk.UnwrapSDKContext(goctx)

	pool, exist := k.GetPool(ctx, request.Name)
	if !exist {
		return nil, sdkerrors.Wrapf(types.ErrPoolNotFound, request.Name)
	}

	var totalReward sdk.Coins
	var remainingReward sdk.Coins
	var rewardPerBlock sdk.Coins
	k.IteratorRewardRules(ctx, pool.Name, func(r types.RewardRule) {
		totalReward = totalReward.Add(sdk.NewCoin(r.Reward, r.TotalReward))
		remainingReward = remainingReward.Add(sdk.NewCoin(r.Reward, r.RemainingReward))
		rewardPerBlock = rewardPerBlock.Add(sdk.NewCoin(r.Reward, r.RewardPerBlock))
	})

	poolEntry := &types.FarmPoolEntry{
		Name:            pool.Name,
		Creator:         pool.Creator,
		Description:     pool.Description,
		StartHeight:     pool.StartHeight,
		EndHeight:       pool.EndHeight,
		Editable:        pool.Editable,
		Expired:         k.Expired(ctx, pool),
		TotalLptLocked:  pool.TotalLptLocked,
		TotalReward:     totalReward,
		RemainingReward: remainingReward,
		RewardPerBlock:  rewardPerBlock,
	}
	return &types.QueryFarmPoolResponse{Pool: poolEntry}, nil
}

func (k Keeper) Farmer(goctx context.Context, request *types.QueryFarmerRequest) (*types.QueryFarmerResponse, error) {
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
		return nil, sdkerrors.Wrapf(types.ErrFarmerNotFound, "not found farmer: %s", request.Farmer)
	}

	for _, farmer := range farmInfos {
		pool, exist := k.GetPool(cacheCtx, farmer.PoolName)
		if !exist {
			return nil, sdkerrors.Wrapf(types.ErrPoolNotFound, farmer.PoolName)
		}

		//The farm pool has not started, no reward
		if pool.StartHeight > ctx.BlockHeight() {
			list = append(list, &types.LockedInfo{
				PoolName: farmer.PoolName,
				Locked:   sdk.NewCoin(pool.TotalLptLocked.Denom, farmer.Locked),
			})
			continue
		}

		if !k.Expired(ctx, pool) {
			pool, _, err = k.updatePool(cacheCtx, pool, sdk.ZeroInt(), false)
			if err != nil {
				return nil, err
			}
		} else {
			pool.Rules = k.GetRewardRules(ctx, pool.Name)
		}

		rewards, _ := pool.CaclRewards(farmer, sdk.ZeroInt())
		list = append(list, &types.LockedInfo{
			PoolName:      farmer.PoolName,
			Locked:        sdk.NewCoin(pool.TotalLptLocked.Denom, farmer.Locked),
			PendingReward: rewards,
		})
	}

	return &types.QueryFarmerResponse{
		List:   list,
		Height: ctx.BlockHeight(),
	}, nil
}

func (k Keeper) Params(goctx context.Context, request *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goctx)
	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}
