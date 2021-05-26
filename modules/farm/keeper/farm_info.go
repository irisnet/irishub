package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irismod/modules/farm/types"
)

// GetFarmer return the specified farmer
func (k Keeper) GetFarmInfo(ctx sdk.Context, poolName, address string) (info types.FarmInfo, exist bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyFarmInfo(address, poolName))
	if len(bz) == 0 {
		return info, false
	}

	k.cdc.MustUnmarshalBinaryBare(bz, &info)
	return info, true
}

func (k Keeper) IteratorFarmInfo(ctx sdk.Context, address string, fun func(farmer types.FarmInfo)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PrefixFarmInfo(address))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var farmer types.FarmInfo
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &farmer)
		fun(farmer)
	}
}

func (k Keeper) IteratorAllFarmInfo(ctx sdk.Context, fun func(farmer types.FarmInfo)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.FarmerKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var farmer types.FarmInfo
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &farmer)
		fun(farmer)
	}
}

// SetFarmer save the farmer information
func (k Keeper) SetFarmInfo(ctx sdk.Context, farmer types.FarmInfo) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&farmer)
	store.Set(types.KeyFarmInfo(farmer.Address, farmer.PoolName), bz)
}

func (k Keeper) DeleteFarmInfo(ctx sdk.Context, poolName, address string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyFarmInfo(address, poolName))
}
