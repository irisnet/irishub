package keeper

import (
	"context"
	"encoding/hex"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

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

	if m.Keeper.blockedAddrs[msg.To] {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is a module account", msg.To)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	id, err := m.Keeper.CreateHTLC(
		ctx,
		sender,
		to,
		msg.ReceiverOnOtherChain,
		msg.SenderOnOtherChain,
		msg.Amount,
		hashLock,
		msg.Timestamp,
		msg.TimeLock,
		msg.Transfer,
	)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateHTLC,
			sdk.NewAttribute(types.AttributeKeyID, id.String()),
			sdk.NewAttribute(types.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyReceiver, msg.To),
			sdk.NewAttribute(types.AttributeKeyReceiverOnOtherChain, msg.ReceiverOnOtherChain),
			sdk.NewAttribute(types.AttributeKeySenderOnOtherChain, msg.SenderOnOtherChain),
			sdk.NewAttribute(types.AttributeKeyTransfer, strconv.FormatBool(msg.Transfer)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})
	return &types.MsgCreateHTLCResponse{
		Id: id.String(),
	}, nil
}

func (m msgServer) ClaimHTLC(goCtx context.Context, msg *types.MsgClaimHTLC) (*types.MsgClaimHTLCResponse, error) {
	id, err := hex.DecodeString(msg.Id)
	if err != nil {
		return nil, err
	}

	secret, err := hex.DecodeString(msg.Secret)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	hashLock, transfer, direction, err := m.Keeper.ClaimHTLC(ctx, id, secret)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeClaimHTLC,
			sdk.NewAttribute(types.AttributeKeyID, msg.Id),
			sdk.NewAttribute(types.AttributeKeyHashLock, hashLock),
			sdk.NewAttribute(types.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(types.AttributeKeySecret, msg.Secret),
			sdk.NewAttribute(types.AttributeKeyTransfer, strconv.FormatBool(transfer)),
			sdk.NewAttribute(types.AttributeKeyDirection, direction.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})
	return &types.MsgClaimHTLCResponse{}, nil
}
