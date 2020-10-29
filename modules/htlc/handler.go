package htlc

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/htlc/keeper"
	"github.com/irisnet/irismod/modules/htlc/types"
)

// NewHandler creates an sdk.Handler for all the HTLC type messages
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgCreateHTLC:
			return handleMsgCreateHTLC(ctx, k, msg)

		case *types.MsgClaimHTLC:
			return handleMsgClaimHTLC(ctx, k, msg)

		case *types.MsgRefundHTLC:
			return handleMsgRefundHTLC(ctx, k, msg)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

// handleMsgCreateHTLC handles MsgCreateHTLC
func handleMsgCreateHTLC(ctx sdk.Context, k keeper.Keeper, msg *types.MsgCreateHTLC) (*sdk.Result, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	to, err := sdk.AccAddressFromBech32(msg.To)
	if err != nil {
		return nil, err
	}

	hashLock, err := hex.DecodeString(msg.HashLock)
	if err != nil {
		return nil, err
	}

	if err := k.CreateHTLC(
		ctx,
		sender,
		to,
		msg.ReceiverOnOtherChain,
		msg.Amount,
		hashLock,
		msg.Timestamp,
		msg.TimeLock,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateHTLC,
			sdk.NewAttribute(types.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyReceiver, msg.To),
			sdk.NewAttribute(types.AttributeKeyReceiverOnOtherChain, msg.ReceiverOnOtherChain),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyHashLock, msg.HashLock),
			sdk.NewAttribute(types.AttributeKeyTimeLock, fmt.Sprintf("%d", msg.TimeLock)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// handleMsgClaimHTLC handles MsgClaimHTLC
func handleMsgClaimHTLC(ctx sdk.Context, k keeper.Keeper, msg *types.MsgClaimHTLC) (*sdk.Result, error) {
	hashLock, err := hex.DecodeString(msg.HashLock)
	if err != nil {
		return nil, err
	}

	secret, err := hex.DecodeString(msg.Secret)
	if err != nil {
		return nil, err
	}

	if err := k.ClaimHTLC(ctx, hashLock, secret); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeClaimHTLC,
			sdk.NewAttribute(types.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyHashLock, msg.HashLock),
			sdk.NewAttribute(types.AttributeKeySecret, msg.Secret),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// handleMsgRefundHTLC handles MsgRefundHTLC
func handleMsgRefundHTLC(ctx sdk.Context, k keeper.Keeper, msg *types.MsgRefundHTLC) (*sdk.Result, error) {
	hashLock, err := hex.DecodeString(msg.HashLock)
	if err != nil {
		return nil, err
	}

	if err := k.RefundHTLC(ctx, hashLock); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRefundHTLC,
			sdk.NewAttribute(types.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyHashLock, msg.HashLock),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
