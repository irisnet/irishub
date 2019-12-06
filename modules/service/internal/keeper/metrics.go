package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) InitMetrics(store sdk.KVStore) {
	iterator := k.ActiveAllRequestQueueIterator(store)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		k.metrics.ActiveRequests.Add(1)
	}
}
