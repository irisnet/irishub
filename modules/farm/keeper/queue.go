package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"mods.irisnet.org/farm/types"
)

func (k Keeper) Expired(ctx sdk.Context, pool types.FarmPool) bool {
	height := ctx.BlockHeader().Height
	switch {
	case height > pool.EndHeight:
		return true
	case height == pool.EndHeight:
		// When Destroy and other operations are at the same block height
		key := types.KeyActiveFarmPool(pool.EndHeight, pool.Id)
		store := ctx.KVStore(k.storeKey)
		return !store.Has(key)
	default:
		return false
	}
}

func (k Keeper) EnqueueActivePool(ctx sdk.Context, poolId string, expiredHeight int64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(
		types.KeyActiveFarmPool(expiredHeight, poolId),
		types.MustMarshalPoolId(k.cdc, poolId),
	)
}

func (k Keeper) DequeueActivePool(ctx sdk.Context, poolId string, expiredHeight int64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyActiveFarmPool(expiredHeight, poolId))
}

func (k Keeper) IteratorExpiredPool(ctx sdk.Context, height int64, fun func(pool types.FarmPool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PrefixActiveFarmPool(height))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		poolId := types.MustUnMarshalPoolId(k.cdc, iterator.Value())
		if pool, exist := k.GetPool(ctx, poolId); exist {
			fun(pool)
		}
	}
}

func (k Keeper) IteratorActivePool(ctx sdk.Context, fun func(pool types.FarmPool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ActiveFarmPoolKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		poolId := types.MustUnMarshalPoolId(k.cdc, iterator.Value())
		if pool, exist := k.GetPool(ctx, poolId); exist {
			fun(pool)
		}
	}
}

func (k Keeper) IteratorAllPools(ctx sdk.Context, fun func(pool types.FarmPool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.FarmPoolKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var pool types.FarmPool
		k.cdc.MustUnmarshal(iterator.Value(), &pool)
		fun(pool)
	}
}
