package htlc

import (
	sdk "github.com/irisnet/irishub/types"
)

// NewHandler handles all "htlc" messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgCreateHTLC:
			return handleMsgCreateHTLC(ctx, k, msg)
		default:
			return sdk.ErrTxDecode("invalid message parsed in HTLC module").Result()
		}

		return sdk.ErrTxDecode("invalid message parsed in HTLC module").Result()
	}
}

// handleMsgCreateHTLC handles MsgCreateHTLC
func handleMsgCreateHTLC(ctx sdk.Context, k Keeper, msg MsgCreateHTLC) sdk.Result {
	secret := make([]byte, 32)
	expireHeight := msg.TimeLock + uint64(ctx.BlockHeight())
	state := uint8(0)

	htlc := NewHTLC(msg.Sender, msg.Receiver, msg.ReceiverOnOtherChain, msg.OutAmount, msg.InAmount, secret, msg.Timestamp, expireHeight, state)

	tags, err := k.CreateHTLC(ctx, htlc)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}
