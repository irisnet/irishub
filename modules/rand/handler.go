package rand

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns a handler for all rand msgs
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgRequestRand:
			return handleMsgRequestRand(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

// handleMsgRequestRand handles MsgRequestRand
func handleMsgRequestRand(ctx sdk.Context, k Keeper, msg MsgRequestRand) (*sdk.Result, error) {
	request, err := k.RequestRand(ctx, msg.Consumer, msg.BlockInterval)
	if err != nil {
		return nil, err
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

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
