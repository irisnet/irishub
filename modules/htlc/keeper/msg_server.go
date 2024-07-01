package keeper

import (
	"context"
	"encoding/hex"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"mods.irisnet.org/modules/htlc/types"
)

type msgServer struct {
	k Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the HTLC MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{k: keeper}
}

// CreateHTLC creates an HTLC
func (m msgServer) CreateHTLC(
	goCtx context.Context,
	msg *types.MsgCreateHTLC,
) (*types.MsgCreateHTLCResponse, error) {
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

	if m.k.blockedAddrs[msg.To] {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "%s is a module account", msg.To)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	id, err := m.k.CreateHTLC(
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
	})
	return &types.MsgCreateHTLCResponse{
		Id: id.String(),
	}, nil
}

func (m msgServer) ClaimHTLC(
	goCtx context.Context,
	msg *types.MsgClaimHTLC,
) (*types.MsgClaimHTLCResponse, error) {
	id, err := hex.DecodeString(msg.Id)
	if err != nil {
		return nil, err
	}

	secret, err := hex.DecodeString(msg.Secret)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	hashLock, transfer, direction, err := m.k.ClaimHTLC(ctx, id, secret)
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
	})
	return &types.MsgClaimHTLCResponse{}, nil
}

func (m msgServer) UpdateParams(
	goCtx context.Context,
	msg *types.MsgUpdateParams,
) (*types.MsgUpdateParamsResponse, error) {
	if m.k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrUnauthorized,
			"invalid authority; expected %s, got %s",
			m.k.authority,
			msg.Authority,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.k.SetParams(ctx, msg.Params); err != nil {
		return nil, err
	}
	return &types.MsgUpdateParamsResponse{}, nil
}
