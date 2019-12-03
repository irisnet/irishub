package rand

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/rand/internal/types"
)

// NewHandler handles all "rand" messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgRequestRand:
			return handleMsgRequestRand(ctx, k, msg)
		default:
			return sdk.ErrTxDecode("invalid message parsed in rand module").Result()
		}

		return sdk.ErrTxDecode("invalid message parsed in rand module").Result()
	}
}

// handleMsgRequestRand handles MsgRequestRand
func handleMsgRequestRand(ctx sdk.Context, k Keeper, msg MsgRequestRand) sdk.Result {
	request, err := k.RequestRand(ctx, msg.Consumer, msg.BlockInterval)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Consumer.String()),
		),
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRequestRand,
			sdk.NewAttribute(types.AttributeKeyRequestID, GenerateRequestID(request)),
			sdk.NewAttribute(types.AttributeKeyGenHeight, fmt.Sprintf("%d", request.Height+int64(msg.BlockInterval))),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}
