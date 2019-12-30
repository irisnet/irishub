package htlc

import (
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns a handler for all htlc messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgCreateHTLC:
			return handleMsgCreateHTLC(ctx, k, msg)
		case MsgClaimHTLC:
			return handleMsgClaimHTLC(ctx, k, msg)
		case MsgRefundHTLC:
			return handleMsgRefundHTLC(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

// handleMsgCreateHTLC handles MsgCreateHTLC
func handleMsgCreateHTLC(ctx sdk.Context, k Keeper, msg MsgCreateHTLC) (*sdk.Result, error) {
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
		return nil, err
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

	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

// handleMsgClaimHTLC handles MsgClaimHTLC
func handleMsgClaimHTLC(ctx sdk.Context, k Keeper, msg MsgClaimHTLC) (*sdk.Result, error) {
	senderStr, toStr, err := k.ClaimHTLC(ctx, msg.HashLock, msg.Secret)
	if err != nil {
		return nil, err
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
				sdk.NewAttribute(AttributeValueSender, senderStr),
				sdk.NewAttribute(AttributeValueReceiver, toStr),
				sdk.NewAttribute(AttributeValueHashLock, hex.EncodeToString(msg.HashLock)),
				sdk.NewAttribute(AttributeValueSecret, hex.EncodeToString(msg.Secret)),
			),
		},
	)

	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

// handleMsgRefundHTLC handles MsgRefundHTLC
func handleMsgRefundHTLC(ctx sdk.Context, k Keeper, msg MsgRefundHTLC) (*sdk.Result, error) {
	sender, err := k.RefundHTLC(ctx, msg.HashLock)
	if err != nil {
		return nil, err
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

	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}
