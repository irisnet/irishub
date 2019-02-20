package bank

import (
	sdk "github.com/irisnet/irishub/types"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSend:
			return handleMsgSend(ctx, k, msg)
		case MsgIssue:
			return handleMsgIssue(ctx, k, msg)
		case MsgBurn:
			return handleMsgBurn(ctx, k, msg)
		default:
			errMsg := "Unrecognized bank Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSend.
func handleMsgSend(ctx sdk.Context, k Keeper, msg MsgSend) sdk.Result {
	// NOTE: totalIn == totalOut should already have been checked

	tags, err := k.InputOutputCoins(ctx, msg.Inputs, msg.Outputs)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}

// Handle MsgIssue.
func handleMsgIssue(ctx sdk.Context, k Keeper, msg MsgIssue) sdk.Result {
	panic("not implemented yet")
}

// Handle MsgBurn.
func handleMsgBurn(ctx sdk.Context, k Keeper, msg MsgBurn) sdk.Result {

	tags, err := k.BurnCoinsFromAddr(ctx, msg.Owner, msg.Coins)

	if err != nil {
		return err.Result()
	}
	ctx.CoinFlowTags().AppendCoinFlowTag(ctx, msg.Owner.String(), "", msg.Coins.String(), sdk.TokenTransfer)
	return sdk.Result{
		Tags: tags,
	}
}
