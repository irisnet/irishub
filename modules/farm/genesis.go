package farm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"irismod.io/farm/keeper"
	"irismod.io/farm/types"
)

// InitGenesis stores the genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	if err := types.ValidateGenesis(data); err != nil {
		panic(err.Error())
	}
	for _, pool := range data.Pools {
		for _, r := range pool.Rules {
			k.SetRewardRule(ctx, pool.Id, r)
		}
		k.SetPool(ctx, pool)
		if !k.Expired(ctx, pool) {
			k.EnqueueActivePool(ctx, pool.Id, pool.EndHeight)
		}
	}

	for _, farmInfo := range data.FarmInfos {
		_, exist := k.GetPool(ctx, farmInfo.PoolId)
		if !exist {
			panic(types.ErrPoolNotFound)
		}
		k.SetFarmInfo(ctx, farmInfo)
	}

	for _, info := range data.Escrow {
		k.SetEscrowInfo(ctx, info)
	}
	k.SetSequence(ctx, data.Sequence)
	k.SetParams(ctx, data.Params)
}

// ExportGenesis outputs the genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	var pools []types.FarmPool
	var farmInfos []types.FarmInfo
	k.IteratorAllPools(ctx, func(pool types.FarmPool) {
		pool.Rules = k.GetRewardRules(ctx, pool.Id)
		pools = append(pools, pool)
	})
	k.IteratorAllFarmInfo(ctx, func(farmInfo types.FarmInfo) {
		farmInfos = append(farmInfos, farmInfo)
	})
	return &types.GenesisState{
		Params:    k.GetParams(ctx),
		Pools:     pools,
		FarmInfos: farmInfos,
		Sequence:  k.GetSequence(ctx),
		Escrow:    k.GetAllEscrowInfo(ctx),
	}
}
