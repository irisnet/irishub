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
		case MsgClaimHTLC:
			return handleMsgClaimHTLC(ctx, k, msg)
		case MsgRefundHTLC:
			return handleMsgRefundHTLC(ctx, k, msg)
		default:
			return sdk.ErrTxDecode("invalid message parsed in HTLC module").Result()
		}
	}
}

// handleMsgCreateHTLC handles MsgCreateHTLC
func handleMsgCreateHTLC(ctx sdk.Context, k Keeper, msg MsgCreateHTLC) sdk.Result {
	secret := make([]byte, 0)
	expireHeight := msg.TimeLock + uint64(ctx.BlockHeight())
	state := OPEN

	htlc := NewHTLC(
		msg.Sender,
		msg.To,
		msg.ReceiverOnOtherChain,
		msg.Amount,
		secret,
		msg.Timestamp,
		expireHeight,
		state,
	)

	tags, err := k.CreateHTLC(ctx, htlc, msg.HashLock)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}

// handleMsgClaimHTLC handles MsgClaimHTLC
func handleMsgClaimHTLC(ctx sdk.Context, k Keeper, msg MsgClaimHTLC) sdk.Result {
	tags, err := k.ClaimHTLC(ctx, msg.HashLock, msg.Secret)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}

// handleMsgRefundHTLC handles MsgRefundHTLC
func handleMsgRefundHTLC(ctx sdk.Context, k Keeper, msg MsgRefundHTLC) sdk.Result {
	tags, err := k.RefundHTLC(ctx, msg.HashLock)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}
