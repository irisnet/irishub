package bank

import (
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/types/bank"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(bo BO) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case bank.MsgSend:
			return handleMsgSend(ctx, bo, msg)
		case bank.MsgIssue:
			return handleMsgIssue(ctx, bo, msg)
		default:
			errMsg := "Unrecognized bank Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSend.
func handleMsgSend(ctx sdk.Context, bo BO, msg bank.MsgSend) sdk.Result {
	// NOTE: totalIn == totalOut should already have been checked

	tags, err := bo.InputOutputCoins(ctx, msg.Inputs, msg.Outputs)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}

// Handle MsgIssue.
func handleMsgIssue(ctx sdk.Context, bo BO, msg bank.MsgIssue) sdk.Result {
	panic("not implemented yet")
}

