package keeper

import (
	"context"
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/htlc/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the HTLC MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// CreateHTLC creates an HTLC
func (m msgServer) CreateHTLC(goCtx context.Context, msg *types.MsgCreateHTLC) (*types.MsgCreateHTLCResponse, error) {
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

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.CreateHTLC(
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
	return &types.MsgCreateHTLCResponse{}, nil
}

func (m msgServer) ClaimHTLC(goCtx context.Context, msg *types.MsgClaimHTLC) (*types.MsgClaimHTLCResponse, error) {
	hashLock, err := hex.DecodeString(msg.HashLock)
	if err != nil {
		return nil, err
	}

	secret, err := hex.DecodeString(msg.Secret)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.ClaimHTLC(ctx, hashLock, secret); err != nil {
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
	return &types.MsgClaimHTLCResponse{}, nil
}

func (m msgServer) RefundHTLC(goCtx context.Context, msg *types.MsgRefundHTLC) (*types.MsgRefundHTLCResponse, error) {
	hashLock, err := hex.DecodeString(msg.HashLock)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.RefundHTLC(ctx, hashLock); err != nil {
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
	return &types.MsgRefundHTLCResponse{}, nil
}
