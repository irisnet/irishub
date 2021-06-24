package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/farm/types"
)

func (k Keeper) Expired(ctx sdk.Context, pool types.FarmPool) bool {
	height := ctx.BlockHeader().Height
	switch {
	case height > pool.EndHeight:
		return true
	case height == pool.EndHeight:
		//When Destroy and other operations are at the same block height
		key := types.KeyActiveFarmPool(pool.EndHeight, pool.Name)
		store := ctx.KVStore(k.storeKey)
		return !store.Has(key)
	default:
		return false
	}
}

func (k Keeper) EnqueueActivePool(ctx sdk.Context, poolName string, expiredHeight int64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(
		types.KeyActiveFarmPool(expiredHeight, poolName),
		types.MustMarshalPoolName(k.cdc, poolName),
	)
}

func (k Keeper) DequeueActivePool(ctx sdk.Context, poolName string, expiredHeight int64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyActiveFarmPool(expiredHeight, poolName))
}

func (k Keeper) IteratorExpiredPool(ctx sdk.Context, height int64, fun func(pool types.FarmPool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store,
		types.PrefixActiveFarmPool(height))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		poolName := types.MustUnMarshalPoolName(k.cdc, iterator.Value())
		if pool, exist := k.GetPool(ctx, poolName); exist {
			fun(pool)
		}
	}
}

func (k Keeper) IteratorActivePool(ctx sdk.Context, fun func(pool types.FarmPool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ActiveFarmPoolKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		poolName := types.MustUnMarshalPoolName(k.cdc, iterator.Value())
		if pool, exist := k.GetPool(ctx, poolName); exist {
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
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &pool)
		fun(pool)
	}
}
