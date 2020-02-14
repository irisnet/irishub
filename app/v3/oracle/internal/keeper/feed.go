package keeper

import (
	"time"

	sdk "github.com/irisnet/irishub/types"
)

type Feed struct {
	FeedName              string         `json:"feed_name"`
	AggregateMethod       string         `json:"aggregate_method"`
	AggregateArgsJsonPath string         `json:"aggregate_args_json_path"`
	LatestHistory         uint64         `json:"latest_history"`
	RequestContextID      []byte         `json:"request_context_id"`
	Owner                 sdk.AccAddress `json:"owner"`
}
type Value interface{}
type FeedResult struct {
	Data      Value     `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}
type FeedResults []FeedResult

func (k Keeper) GetFeed(ctx sdk.Context, feedName string) (feed Feed) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(GetFeedKey(feedName))
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, feed)
	return
}

func (k Keeper) hasFeed(ctx sdk.Context, feedName string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(GetFeedKey(feedName))
}

func (k Keeper) setFeed(ctx sdk.Context, feed Feed) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(feed)
	store.Set(GetFeedKey(feed.FeedName), bz)
}
