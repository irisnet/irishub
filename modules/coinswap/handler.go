package coinswap

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for all "coinswap" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgSwapOrder:
			return handleMsgSwapOrder(ctx, k, msg)
		case MsgAddLiquidity:
			return handleMsgAddLiquidity(ctx, k, msg)
		case MsgRemoveLiquidity:
			return handleMsgRemoveLiquidity(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized coinswap message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleMsgSwapOrder handles MsgSwapOrder.
func handleMsgSwapOrder(ctx sdk.Context, k Keeper, msg MsgSwapOrder) sdk.Result {
	// check that deadline has not passed
	if ctx.BlockHeader().Time.After(time.Unix(msg.Deadline, 0)) {
		return ErrInvalidDeadline("deadline has passed for MsgSwapOrder").Result()
	}

	if err := k.Swap(ctx, msg); err != nil {
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

// handleMsgAddLiquidity handles MsgAddLiquidity. If the reserve pool does not exist, it will be
// created. The first liquidity provider sets the exchange rate.
func handleMsgAddLiquidity(ctx sdk.Context, k Keeper, msg MsgAddLiquidity) sdk.Result {
	// check that deadline has not passed
	if ctx.BlockHeader().Time.After(time.Unix(msg.Deadline, 0)) {
		return ErrInvalidDeadline("deadline has passed for MsgAddLiquidity").Result()
	}

	if err := k.AddLiquidity(ctx, msg); err != nil {
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

// handleMsgRemoveLiquidity handles MsgRemoveLiquidity
func handleMsgRemoveLiquidity(ctx sdk.Context, k Keeper, msg MsgRemoveLiquidity) sdk.Result {
	// check that deadline has not passed
	if ctx.BlockHeader().Time.After(time.Unix(msg.Deadline, 0)) {
		return ErrInvalidDeadline("deadline has passed for MsgRemoveLiquidity").Result()
	}

	if err := k.RemoveLiquidity(ctx, msg); err != nil {
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
