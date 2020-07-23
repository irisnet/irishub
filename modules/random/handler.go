package random

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irishub/modules/random/keeper"
	"github.com/irisnet/irishub/modules/random/types"
)

// NewHandler returns a handler for all rand msgs
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgRequestRandom:
			return handleMsgRequestRandom(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

// handleMsgRequestRandom handles MsgRequestRandom
func handleMsgRequestRandom(ctx sdk.Context, k keeper.Keeper, msg *types.MsgRequestRandom) (*sdk.Result, error) {
	request, err := k.RequestRandom(ctx, msg.Consumer, msg.BlockInterval, msg.Oracle, msg.ServiceFeeCap)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Consumer.String()),
			),
			sdk.NewEvent(
				types.EventTypeRequestRandom,
				sdk.NewAttribute(types.AttributeKeyRequestID, hex.EncodeToString(types.GenerateRequestID(request))),
				sdk.NewAttribute(types.AttributeKeyGenHeight, fmt.Sprintf("%d", request.Height+int64(msg.BlockInterval))),
			),
		},
	)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
