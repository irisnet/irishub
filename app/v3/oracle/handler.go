package oracle

import (
	"fmt"

	"github.com/irisnet/irishub/app/v3/oracle/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

// NewHandler returns a handler for all the "oracle" type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgCreateFeed:
			return handleMsgCreateFeed(ctx, k, msg)
		case MsgStartFeed:
			return handleMsgStartFeed(ctx, k, msg)
		case MsgPauseFeed:
			return handleMsgPauseFeed(ctx, k, msg)
		case MsgEditFeed:
			return handleMsgEditFeed(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized oracle message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleMsgCreateFeed handles MsgCreateFeed
func handleMsgCreateFeed(ctx sdk.Context, k Keeper, msg MsgCreateFeed) sdk.Result {
	if err := k.CreateFeed(ctx, msg); err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: sdk.NewTags(
			types.TagFeedName, []byte(msg.FeedName),
			types.TagCreator, []byte(msg.Creator.String()),
		),
	}
}

// handleMsgStartFeed handles MsgStartFeed
func handleMsgStartFeed(ctx sdk.Context, k Keeper, msg MsgStartFeed) sdk.Result {
	if err := k.StartFeed(ctx, msg); err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: sdk.NewTags(
			types.TagFeedName, []byte(msg.FeedName),
			types.TagCreator, []byte(msg.Creator.String()),
		),
	}
}

// handleMsgPauseFeed handles MsgPauseFeed
func handleMsgPauseFeed(ctx sdk.Context, k Keeper, msg MsgPauseFeed) sdk.Result {
	if err := k.PauseFeed(ctx, msg); err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: sdk.NewTags(
			types.TagFeedName, []byte(msg.FeedName),
			types.TagCreator, []byte(msg.Creator.String()),
		),
	}
}

// handleMsgEditFeed handles MsgEditFeed
func handleMsgEditFeed(ctx sdk.Context, k Keeper, msg MsgEditFeed) sdk.Result {
	if err := k.EditFeed(ctx, msg); err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: sdk.NewTags(
			types.TagFeedName, []byte(msg.FeedName),
			types.TagCreator, []byte(msg.Creator.String()),
		),
	}
}
