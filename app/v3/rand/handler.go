package rand

import (
	sdk "github.com/irisnet/irishub/types"
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
	}
}

// handleMsgRequestRand handles MsgRequestRand
func handleMsgRequestRand(ctx sdk.Context, k Keeper, msg MsgRequestRand) sdk.Result {
	tags, err := k.RequestRand(ctx, msg.Consumer, msg.BlockInterval)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}
