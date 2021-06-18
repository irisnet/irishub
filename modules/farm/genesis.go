package farm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/farm/keeper"
	"github.com/irisnet/irismod/modules/farm/types"
)

// InitGenesis stores the genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	if err := types.ValidateGenesis(data); err != nil {
		panic(err.Error())
	}
	for _, pool := range data.Pools {
		for _, r := range pool.Rules {
			k.SetRewardRule(ctx, pool.Name, r)
		}
		k.SetPool(ctx, pool)
		if !k.Expired(ctx, pool) {
			k.EnqueueActivePool(ctx, pool.Name, pool.EndHeight)
		}
	}

	for _, farmInfo := range data.FarmInfos {
		_, exist := k.GetPool(ctx, farmInfo.PoolName)
		if !exist {
			panic(types.ErrNotExistPool)
		}
		k.SetFarmInfo(ctx, farmInfo)
	}
	k.SetParams(ctx, data.Params)
}

// ExportGenesis outputs the genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	var pools []types.FarmPool
	var farmInfos []types.FarmInfo
	k.IteratorAllPools(ctx, func(pool types.FarmPool) {
		pool.Rules = k.GetRewardRules(ctx, pool.Name)
		pools = append(pools, pool)
	})
	k.IteratorAllFarmInfo(ctx, func(farmInfo types.FarmInfo) {
		farmInfos = append(farmInfos, farmInfo)
	})
	return &types.GenesisState{
		Params:    types.Params{CreatePoolFee: k.CreatePoolFee(ctx)},
		Pools:     pools,
		FarmInfos: farmInfos,
	}
}
