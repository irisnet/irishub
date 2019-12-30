package coinswap

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns a handler for all "coinswap" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgSwapOrder:
			return handleMsgSwapOrder(ctx, k, msg)
		case MsgAddLiquidity:
			return handleMsgAddLiquidity(ctx, k, msg)
		case MsgRemoveLiquidity:
			return handleMsgRemoveLiquidity(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

// handleMsgSwapOrder handles MsgSwapOrder.
func handleMsgSwapOrder(ctx sdk.Context, k Keeper, msg MsgSwapOrder) (*sdk.Result, error) {
	// check that deadline has not passed
	if ctx.BlockHeader().Time.After(time.Unix(msg.Deadline, 0)) {
		return nil, sdkerrors.Wrap(ErrInvalidDeadline, "deadline has passed for MsgSwapOrder")
	}

	if err := k.Swap(ctx, msg); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Input.Address.String()),
		),
	)

	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

// handleMsgAddLiquidity handles MsgAddLiquidity. If the reserve pool does not exist, it will be
// created. The first liquidity provider sets the exchange rate.
func handleMsgAddLiquidity(ctx sdk.Context, k Keeper, msg MsgAddLiquidity) (*sdk.Result, error) {
	// check that deadline has not passed
	if ctx.BlockHeader().Time.After(time.Unix(msg.Deadline, 0)) {
		return nil, sdkerrors.Wrap(ErrInvalidDeadline, "deadline has passed for MsgAddLiquidity")
	}

	if err := k.AddLiquidity(ctx, msg); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
	)

	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

// handleMsgRemoveLiquidity handles MsgRemoveLiquidity
func handleMsgRemoveLiquidity(ctx sdk.Context, k Keeper, msg MsgRemoveLiquidity) (*sdk.Result, error) {
	// check that deadline has not passed
	if ctx.BlockHeader().Time.After(time.Unix(msg.Deadline, 0)) {
		return nil, sdkerrors.Wrap(ErrInvalidDeadline, "deadline has passed for MsgRemoveLiquidity")
	}

	if err := k.RemoveLiquidity(ctx, msg); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
	)

	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}
