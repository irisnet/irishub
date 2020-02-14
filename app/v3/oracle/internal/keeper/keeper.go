package keeper

import (
	"github.com/irisnet/irishub/app/v3/oracle/internal/types"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

// Keeper
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec

	sk types.ServiceKeeper
	// codespace
	codespace sdk.CodespaceType
}

// NewKeeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, codespace sdk.CodespaceType, sk types.ServiceKeeper) Keeper {
	keeper := Keeper{
		storeKey:  key,
		cdc:       cdc,
		sk:        sk,
		codespace: codespace,
	}

	return keeper
}

func (k Keeper) CreateFeed(ctx sdk.Context, msg types.MsgCreateFeed) sdk.Error {
	if _, found := k.GetFeed(ctx, msg.FeedName); found {
		return types.ErrExistedFeedName(types.DefaultCodespace, msg.FeedName)
	}
	requestContextID, err := k.sk.CreateRequestContext(ctx)
	if err != nil {
		return sdk.ErrInternal(err.Error())
	}
	feed := Feed{
		FeedName:              msg.FeedName,
		AggregateMethod:       msg.AggregateMethod,
		AggregateArgsJsonPath: msg.AggregateArgsJsonPath,
		LatestHistory:         msg.LatestHistory,
		RequestContextID:      requestContextID,
		Owner:                 msg.Owner,
	}
	k.setFeed(ctx, feed)
	return nil
}

func (k Keeper) StartFeed(ctx sdk.Context, msg types.MsgStartFeed) sdk.Error {
	feed, found := k.GetFeed(ctx, msg.FeedName)
	if !found {
		return types.ErrUnknownFeedName(types.DefaultCodespace, msg.FeedName)
	}
	if msg.Owner.Equals(feed.Owner) {
		return types.ErrUnauthorized(types.DefaultCodespace, msg.FeedName, msg.Owner)
	}
	//Can not start feed in "running" state
	//TODO
	k.sk.GetRequestContext(ctx, feed.RequestContextID)

	//TODO params ?
	if err := k.sk.StartRequestContext(ctx, feed.RequestContextID); err != nil {
		return sdk.ErrInternal(err.Error())
	}
	return nil
}

func (k Keeper) PauseFeed(ctx sdk.Context, msg types.MsgPauseFeed) sdk.Error {
	feed, found := k.GetFeed(ctx, msg.FeedName)
	if !found {
		return types.ErrUnknownFeedName(types.DefaultCodespace, msg.FeedName)
	}
	if msg.Owner.Equals(feed.Owner) {
		return types.ErrUnauthorized(types.DefaultCodespace, msg.FeedName, msg.Owner)
	}
	//Can only pause feed in "running" state
	//TODO
	k.sk.GetRequestContext(ctx, feed.RequestContextID)

	//TODO params ?
	if err := k.sk.PauseRequestContext(ctx, feed.RequestContextID); err != nil {
		return sdk.ErrInternal(err.Error())
	}
	return nil
}

func (k Keeper) KillFeed(ctx sdk.Context, msg types.MsgKillFeed) sdk.Error {
	feed, found := k.GetFeed(ctx, msg.FeedName)
	if !found {
		return types.ErrUnknownFeedName(types.DefaultCodespace, msg.FeedName)
	}
	if msg.Owner.Equals(feed.Owner) {
		return types.ErrUnauthorized(types.DefaultCodespace, msg.FeedName, msg.Owner)
	}
	//TODO params ?
	if err := k.sk.KillRequestContext(ctx, feed.RequestContextID); err != nil {
		return sdk.ErrInternal(err.Error())
	}
	k.deleteFeed(ctx, feed)
	return nil
}

func (k Keeper) EditFeed(ctx sdk.Context, msg types.MsgEditFeed) sdk.Error {
	feed, found := k.GetFeed(ctx, msg.FeedName)
	if !found {
		return types.ErrUnknownFeedName(types.DefaultCodespace, msg.FeedName)
	}
	if msg.Owner.Equals(feed.Owner) {
		return types.ErrUnauthorized(types.DefaultCodespace, msg.FeedName, msg.Owner)
	}
	if msg.LatestHistory > 1 {
		if msg.LatestHistory < feed.LatestHistory {
			count := int(feed.LatestHistory - msg.LatestHistory)
			k.deleteOldestFeedResult(ctx, feed.FeedName, count)
		}
		feed.LatestHistory = msg.LatestHistory
	}
	//TODO
	if err := k.sk.UpdateRequestContext(ctx, feed.RequestContextID); err != nil {
		return sdk.ErrInternal(err.Error())
	}
	k.setFeed(ctx, feed)
	return nil
}

func (k Keeper) HandleServiceResponse(ctx sdk.Context, requestContextID []byte, responseOutput []string) error {
	return nil
}
