package htlc

import (
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/htlc/internal/types"
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
			return sdk.ErrTxDecode("invalid message parsed in HTLC module").Result()
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

	err := k.CreateHTLC(ctx, htlc, msg.HashLock)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			),
			sdk.NewEvent(
				types.EventTypeCreateHTLC,
				sdk.NewAttribute(types.AttributeValueSender, htlc.Sender.String()),
				sdk.NewAttribute(types.AttributeValueReceiver, htlc.To.String()),
				sdk.NewAttribute(types.AttributeValueReceiverOnOtherChain, htlc.ReceiverOnOtherChain),
				sdk.NewAttribute(types.AttributeValueAmount, htlc.Amount.String()),
				sdk.NewAttribute(types.AttributeValueHashLock, hex.EncodeToString(msg.HashLock)),
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
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			),
			sdk.NewEvent(
				types.EventTypeClaimHTLC,
				sdk.NewAttribute(types.AttributeValueSender, msg.Sender.String()),
				sdk.NewAttribute(types.AttributeValueReceiver, toStr),
				sdk.NewAttribute(types.AttributeValueHashLock, hex.EncodeToString(msg.HashLock)),
				sdk.NewAttribute(types.AttributeValueSecret, hex.EncodeToString(msg.Secret)),
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
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			),
			sdk.NewEvent(
				types.EventTypeRefundHTLC,
				sdk.NewAttribute(types.AttributeValueSender, sender),
				sdk.NewAttribute(types.AttributeValueHashLock, hex.EncodeToString(msg.HashLock)),
			),
		},
	)

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
