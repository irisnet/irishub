package oracle

import (
	"fmt"

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
		case MsgStopFeed:
			return handleMsgStopFeed(ctx, k, msg)
		case MsgEditFeed:
			return handleMsgEditFeed(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized service message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgCreateFeed(ctx sdk.Context, k Keeper, msg MsgCreateFeed) sdk.Result {
	tags, err := k.CreateFeed(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

func handleMsgStartFeed(ctx sdk.Context, k Keeper, msg MsgStartFeed) sdk.Result {
	tags, err := k.StartFeed(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

func handleMsgStopFeed(ctx sdk.Context, k Keeper, msg MsgStopFeed) sdk.Result {
	tags, err := k.StopFeed(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

func handleMsgEditFeed(ctx sdk.Context, k Keeper, msg MsgEditFeed) sdk.Result {
	tags, err := k.EditFeed(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}
