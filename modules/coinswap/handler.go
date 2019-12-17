package coinswap

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "coinswap" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgSwapOrder:
			return HandleMsgSwapOrder(ctx, k, msg)
		case MsgAddLiquidity:
			return HandleMsgAddLiquidity(ctx, k, msg)
		case MsgRemoveLiquidity:
			return HandleMsgRemoveLiquidity(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized coinswap message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSwapOrder.
func HandleMsgSwapOrder(ctx sdk.Context, k Keeper, msg MsgSwapOrder) sdk.Result {
	// check that deadline has not passed
	if ctx.BlockHeader().Time.After(time.Unix(msg.Deadline, 0)) {
		return ErrInvalidDeadline("deadline has passed for MsgSwapOrder").Result()
	}

	err := k.Swap(ctx, msg)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Input.Address.String()),
		),
	)

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// Handle MsgAddLiquidity. If the reserve pool does not exist, it will be
// created. The first liquidity provider sets the exchange rate.
func HandleMsgAddLiquidity(ctx sdk.Context, k Keeper, msg MsgAddLiquidity) sdk.Result {
	// check that deadline has not passed
	if ctx.BlockHeader().Time.After(time.Unix(msg.Deadline, 0)) {
		return ErrInvalidDeadline("deadline has passed for MsgAddLiquidity").Result()
	}

	err := k.AddLiquidity(ctx, msg)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
	)

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// HandleMsgRemoveLiquidity handler for MsgRemoveLiquidity
func HandleMsgRemoveLiquidity(ctx sdk.Context, k Keeper, msg MsgRemoveLiquidity) sdk.Result {
	// check that deadline has not passed
	if ctx.BlockHeader().Time.After(time.Unix(msg.Deadline, 0)) {
		return ErrInvalidDeadline("deadline has passed for MsgRemoveLiquidity").Result()
	}

	err := k.RemoveLiquidity(ctx, msg)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
	)

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
