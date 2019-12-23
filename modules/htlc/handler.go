package htlc

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler handles all htlc messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgCreateHTLC:
			return handleMsgCreateHTLC(ctx, k, msg)
		case MsgClaimHTLC:
			return handleMsgClaimHTLC(ctx, k, msg)
		case MsgRefundHTLC:
			return handleMsgRefundHTLC(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized HTLC message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleMsgCreateHTLC handles MsgCreateHTLC
func handleMsgCreateHTLC(ctx sdk.Context, k Keeper, msg MsgCreateHTLC) sdk.Result {
	secret := HTLCSecret{}
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

	if err := k.CreateHTLC(ctx, htlc, msg.HashLock); err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			),
			sdk.NewEvent(
				EventTypeCreateHTLC,
				sdk.NewAttribute(AttributeValueSender, htlc.Sender.String()),
				sdk.NewAttribute(AttributeValueReceiver, htlc.To.String()),
				sdk.NewAttribute(AttributeValueReceiverOnOtherChain, htlc.ReceiverOnOtherChain),
				sdk.NewAttribute(AttributeValueAmount, htlc.Amount.String()),
				sdk.NewAttribute(AttributeValueHashLock, hex.EncodeToString(msg.HashLock)),
			),
		},
	)

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// handleMsgClaimHTLC handles MsgClaimHTLC
func handleMsgClaimHTLC(ctx sdk.Context, k Keeper, msg MsgClaimHTLC) sdk.Result {
	toStr, err := k.ClaimHTLC(ctx, msg.HashLock, msg.Secret)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			),
			sdk.NewEvent(
				EventTypeClaimHTLC,
				sdk.NewAttribute(AttributeValueSender, msg.Sender.String()),
				sdk.NewAttribute(AttributeValueReceiver, toStr),
				sdk.NewAttribute(AttributeValueHashLock, hex.EncodeToString(msg.HashLock)),
				sdk.NewAttribute(AttributeValueSecret, hex.EncodeToString(msg.Secret)),
			),
		},
	)

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// handleMsgRefundHTLC handles MsgRefundHTLC
func handleMsgRefundHTLC(ctx sdk.Context, k Keeper, msg MsgRefundHTLC) sdk.Result {
	sender, err := k.RefundHTLC(ctx, msg.HashLock)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			),
			sdk.NewEvent(
				EventTypeRefundHTLC,
				sdk.NewAttribute(AttributeValueSender, sender),
				sdk.NewAttribute(AttributeValueHashLock, hex.EncodeToString(msg.HashLock)),
			),
		},
	)

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
