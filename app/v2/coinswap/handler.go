package coinswap

import (
	"fmt"
	"time"

	sdk "github.com/irisnet/irishub/types"
)

// NewHandler returns a handler for "coinswap" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSwapOrder:
			return HandleMsgSwapOrder(ctx, msg, k)

		case MsgAddLiquidity:
			return HandleMsgAddLiquidity(ctx, msg, k)

		case MsgRemoveLiquidity:
			return HandleMsgRemoveLiquidity(ctx, msg, k)

		default:
			errMsg := fmt.Sprintf("unrecognized coinswap message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSwapOrder.
func HandleMsgSwapOrder(ctx sdk.Context, msg MsgSwapOrder, k Keeper) sdk.Result {
	// check that deadline has not passed
	if ctx.BlockHeader().Time.After(time.Unix(msg.Deadline, 0)) {
		return ErrInvalidDeadline("deadline has passed for MsgSwapOrder").Result()
	}
	tags, err := k.HandleSwap(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{Tags: tags}
}

// Handle MsgAddLiquidity. If the reserve pool does not exist, it will be
// created. The first liquidity provider sets the exchange rate.
func HandleMsgAddLiquidity(ctx sdk.Context, msg MsgAddLiquidity, k Keeper) sdk.Result {
	// check that deadline has not passed
	if ctx.BlockHeader().Time.After(time.Unix(msg.Deadline, 0)) {
		return ErrInvalidDeadline("deadline has passed for MsgAddLiquidity").Result()
	}

	tags, err := k.HandleAddLiquidity(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}

// HandleMsgRemoveLiquidity handler for MsgRemoveLiquidity
func HandleMsgRemoveLiquidity(ctx sdk.Context, msg MsgRemoveLiquidity, k Keeper) sdk.Result {
	// check that deadline has not passed
	if ctx.BlockHeader().Time.After(time.Unix(msg.Deadline, 0)) {
		return ErrInvalidDeadline("deadline has passed for MsgRemoveLiquidity").Result()
	}

	tags, err := k.HandleRemoveLiquidity(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}
