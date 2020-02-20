package keeper

import (
	"github.com/irisnet/irishub/app/v3/oracle/internal/types"
	service "github.com/irisnet/irishub/app/v3/service/exported"
	sdk "github.com/irisnet/irishub/types"
)

func (k Keeper) GetFeed(ctx sdk.Context, feedName string) (feed types.Feed, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(GetFeedKey(feedName))
	if bz == nil {
		return feed, false
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &feed)
	return feed, true
}

func (k Keeper) GetFeeds(ctx sdk.Context) (feeds []types.Feed) {
	k.IteratorFeeds(ctx, func(feed types.Feed) {
		feeds = append(feeds, feed)
	})
	return
}

func (k Keeper) IteratorFeeds(ctx sdk.Context, fn func(feed types.Feed)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStoreReversePrefixIterator(store, GetFeedPrefixKey())
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var res types.Feed
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &res)
		fn(res)
	}
}

func (k Keeper) SetFeed(ctx sdk.Context, feed types.Feed) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(feed)
	store.Set(GetFeedKey(feed.FeedName), bz)

	bz = k.cdc.MustMarshalBinaryLengthPrefixed(feed.FeedName)
	store.Set(GetReqCtxIDKey(feed.RequestContextID), bz)
}

func (k Keeper) GetFeedByReqCtxID(ctx sdk.Context, requestContextID []byte) (feed types.Feed, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(GetReqCtxIDKey(requestContextID))
	var feedName string
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &feedName)
	return k.GetFeed(ctx, feedName)
}

func (k Keeper) SetFeedValue(ctx sdk.Context,
	feedName string,
	batchCounter uint64,
	latestHistory uint64,
	value types.FeedValue) {
	store := ctx.KVStore(k.storeKey)
	counter := k.getFeedValuesCnt(ctx, feedName)
	delta := counter - int(latestHistory)
	k.deleteOldestFeedValue(ctx, feedName, delta+1)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(value)
	store.Set(GetFeedValueKey(feedName, batchCounter), bz)
}

func (k Keeper) GetFeedValues(ctx sdk.Context, feedName string) (result types.FeedValues) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStoreReversePrefixIterator(store, GetFeedValuePrefixKey(feedName))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var res types.FeedValue
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &res)
		result = append(result, res)
	}
	return
}

func (k Keeper) GetFeedByState(ctx sdk.Context, state service.RequestContextState) (feeds []types.Feed) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, GetFeedStatePrefixKey(state))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var feedName string
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &feedName)
		if feed, found := k.GetFeed(ctx, feedName); found {
			feeds = append(feeds, feed)
		}
	}
	return
}

func (k Keeper) Enqueue(ctx sdk.Context, feedName string, state service.RequestContextState) {
	var dequeueState service.RequestContextState
	if state == service.RUNNING {
		dequeueState = service.PAUSED
	} else {
		dequeueState = service.RUNNING
	}
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetFeedStateKey(feedName, dequeueState))

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(feedName)
	store.Set(GetFeedStateKey(feedName, state), bz)
}

func (k Keeper) getFeedValuesCnt(ctx sdk.Context, feedName string) (i int) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStoreReversePrefixIterator(store, GetFeedValuePrefixKey(feedName))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		i++
	}
	return
}

func (k Keeper) deleteOldestFeedValue(ctx sdk.Context, feedName string, delta int) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, GetFeedValuePrefixKey(feedName))
	defer iterator.Close()
	for i := 1; iterator.Valid() && i <= delta; iterator.Next() {
		store.Delete(iterator.Key())
		i++
	}
}
