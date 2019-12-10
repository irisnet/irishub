package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) InitMetrics(ctx sdk.Context) {
	iterator := k.ActiveAllRequestQueueIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		k.metrics.ActiveRequests.Add(1)
	}
}
