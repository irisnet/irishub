package rand

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for all rand msgs
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgRequestRand:
			return handleMsgRequestRand(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized rand message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleMsgRequestRand handles MsgRequestRand
func handleMsgRequestRand(ctx sdk.Context, k Keeper, msg MsgRequestRand) sdk.Result {
	request, err := k.RequestRand(ctx, msg.Consumer, msg.BlockInterval)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Consumer.String()),
			),
			sdk.NewEvent(
				EventTypeRequestRand,
				sdk.NewAttribute(AttributeKeyRequestID, hex.EncodeToString(GenerateRequestID(request))),
				sdk.NewAttribute(AttributeKeyGenHeight, fmt.Sprintf("%d", request.Height+int64(msg.BlockInterval))),
			),
		},
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}
