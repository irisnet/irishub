package keeper

import (
	gogotypes "github.com/gogo/protobuf/types"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irismod/service/exported"
	servicetypes "github.com/irismod/service/types"

	"github.com/irisnet/irishub/modules/oracle/types"
)

//GetFeed return the feed by feedName
func (k Keeper) GetFeed(ctx sdk.Context, feedName string) (feed types.Feed, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetFeedKey(feedName))
	if bz == nil {
		return feed, false
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &feed)
	return feed, true
}

//GetFeedByReqCtxID return feed by requestContextID
func (k Keeper) GetFeedByReqCtxID(ctx sdk.Context, requestContextID tmbytes.HexBytes) (feed types.Feed, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetReqCtxIDKey(requestContextID))
	var feedName gogotypes.StringValue
	k.cdc.MustUnmarshalBinaryBare(bz, &feedName)
	return k.GetFeed(ctx, feedName.Value)
}

//IteratorFeeds will foreach all feeds
func (k Keeper) IteratorFeeds(ctx sdk.Context, fn func(feed types.Feed)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStoreReversePrefixIterator(store, types.GetFeedPrefixKey())
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var res types.Feed
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &res)
		fn(res)
	}
}

//IteratorFeedsByState will foreach all feeds by state
func (k Keeper) IteratorFeedsByState(
	ctx sdk.Context,
	state servicetypes.RequestContextState,
	fn func(feed types.Feed),
) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetFeedStatePrefixKey(state))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var feedName gogotypes.StringValue
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &feedName)
		if feed, found := k.GetFeed(ctx, feedName.Value); found {
			fn(feed)
		}
	}
}

//SetFeed will save a feed to store
func (k Keeper) SetFeed(ctx sdk.Context, feed types.Feed) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&feed)
	store.Set(types.GetFeedKey(feed.FeedName), bz)

	bz = k.cdc.MustMarshalBinaryBare(&gogotypes.StringValue{Value: feed.FeedName})
	store.Set(types.GetReqCtxIDKey(feed.RequestContextID), bz)
}

//SetFeedValue will save a feed result to store
func (k Keeper) SetFeedValue(
	ctx sdk.Context,
	feedName string,
	batchCounter uint64,
	latestHistory uint64,
	value types.FeedValue,
) {
	store := ctx.KVStore(k.storeKey)
	counter := k.getFeedValuesCnt(ctx, feedName)
	delta := counter - int(latestHistory)
	k.deleteOldestFeedValue(ctx, feedName, delta+1)
	bz := k.cdc.MustMarshalBinaryBare(&value)
	store.Set(types.GetFeedValueKey(feedName, batchCounter), bz)
}

//GetFeedValues return all feed values by feedName
func (k Keeper) GetFeedValues(ctx sdk.Context, feedName string) (result types.FeedValues) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStoreReversePrefixIterator(store, types.GetFeedValuePrefixKey(feedName))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var res types.FeedValue
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &res)
		result = append(result, res)
	}
	return
}

//Enqueue will put feedName to a 'state' queue
func (k Keeper) Enqueue(ctx sdk.Context, feedName string, state servicetypes.RequestContextState) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&gogotypes.StringValue{Value: feedName})
	store.Set(types.GetFeedStateKey(feedName, state), bz)
}

//Dequeue will remove from the 'state' queue
func (k Keeper) Dequeue(ctx sdk.Context, feedName string, state servicetypes.RequestContextState) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetFeedStateKey(feedName, state))
}

//dequeueAndEnqueue will move feedName  from the 'dequeueState' queue to a 'enqueueState' queue
func (k Keeper) dequeueAndEnqueue(
	ctx sdk.Context,
	feedName string,
	dequeueState servicetypes.RequestContextState,
	enqueueState servicetypes.RequestContextState,
) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetFeedStateKey(feedName, dequeueState))

	bz := k.cdc.MustMarshalBinaryBare(&gogotypes.StringValue{Value: feedName})
	store.Set(types.GetFeedStateKey(feedName, enqueueState), bz)
}

func (k Keeper) getFeedValuesCnt(ctx sdk.Context, feedName string) (i int) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStoreReversePrefixIterator(store, types.GetFeedValuePrefixKey(feedName))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		i++
	}
	return
}

func (k Keeper) deleteOldestFeedValue(ctx sdk.Context, feedName string, delta int) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetFeedValuePrefixKey(feedName))
	defer iterator.Close()
	for i := 1; iterator.Valid() && i <= delta; iterator.Next() {
		store.Delete(iterator.Key())
		i++
	}
}

func (k Keeper) ResetFeedEntryState(ctx sdk.Context) error {
	k.IteratorFeedsByState(
		ctx,
		exported.RUNNING,
		func(feed types.Feed) {
			k.dequeueAndEnqueue(
				ctx,
				feed.FeedName,
				exported.RUNNING,
				exported.PAUSED,
			)
		},
	)
	return nil
}
