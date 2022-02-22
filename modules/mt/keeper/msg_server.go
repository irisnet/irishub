package keeper

import (
	"context"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/irisnet/irismod/modules/mt/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the MT MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// IssueDenom issues a new denom.
func (m msgServer) IssueDenom(goCtx context.Context, msg *types.MsgIssueDenom) (*types.MsgIssueDenomResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	denom := m.Keeper.IssueDenom(ctx, strings.TrimSpace(msg.Name), sender, msg.Data)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeIssueDenom,
			sdk.NewAttribute(types.AttributeKeyDenomID, denom.Id),
			sdk.NewAttribute(types.AttributeKeyDenomName, denom.Name),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Sender),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgIssueDenomResponse{}, nil
}

// MintMT issues a new MT or mints amounts to an MT
func (m msgServer) MintMT(goCtx context.Context, msg *types.MsgMintMT) (*types.MsgMintMTResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// recipient default to the sender
	recipient := sender
	if len(strings.TrimSpace(msg.Recipient)) > 0 {
		recipient, err = sdk.AccAddressFromBech32(msg.Recipient)
		if err != nil {
			return nil, err
		}
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// only denom owner can issue/mint MTs
	if err := m.Keeper.Authorize(ctx, msg.DenomId, sender); err != nil {
		return nil, err
	}

	mtID := strings.TrimSpace(msg.Id)

	// if user inputs an MT ID, then mint amounts to the MT, else issue a new MT
	if len(mtID) > 0 {
		if !m.Keeper.HasMT(ctx, msg.DenomId, mtID) {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "mt not found (%s)", mtID)
		}
		m.Keeper.MintMT(ctx, msg.DenomId, mtID, msg.Amount, recipient)
	} else {
		mt := m.Keeper.IssueMT(ctx, msg.DenomId, msg.Amount, msg.Data, recipient)
		mtID = mt.Id
	}

	mt, err := m.Keeper.GetMT(ctx, msg.DenomId, mtID)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMintMT,
			sdk.NewAttribute(types.AttributeKeyMTID, mtID),
			sdk.NewAttribute(types.AttributeKeyDenomID, msg.DenomId),
			sdk.NewAttribute(types.AttributeKeySupply, strconv.FormatUint(mt.GetSupply(), 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgMintMTResponse{}, nil
}

// EditMT edits an MT
func (m msgServer) EditMT(goCtx context.Context, msg *types.MsgEditMT) (*types.MsgEditMTResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// only denom owner can edit MTs
	if err := m.Keeper.Authorize(ctx, msg.DenomId, sender); err != nil {
		return nil, err
	}

	if err := m.Keeper.EditMT(ctx, msg.DenomId, msg.Id, msg.Data, sender); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEditMT,
			sdk.NewAttribute(types.AttributeKeyMTID, msg.Id),
			sdk.NewAttribute(types.AttributeKeyDenomID, msg.DenomId),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Sender),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgEditMTResponse{}, nil
}

// TransferMT transfers amounts of an MT
func (m msgServer) TransferMT(goCtx context.Context, msg *types.MsgTransferMT) (*types.MsgTransferMTResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.TransferOwner(ctx, msg.DenomId, msg.Id, msg.Amount, sender, recipient); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransfer,
			sdk.NewAttribute(types.AttributeKeyMTID, msg.Id),
			sdk.NewAttribute(types.AttributeKeyDenomID, msg.DenomId),
			sdk.NewAttribute(types.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyRecipient, msg.Recipient),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgTransferMTResponse{}, nil
}

// BurnMT burns MTs of an owner
func (m msgServer) BurnMT(goCtx context.Context, msg *types.MsgBurnMT) (*types.MsgBurnMTResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.BurnMT(ctx, msg.DenomId, msg.Id, msg.Amount, sender); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnMT,
			sdk.NewAttribute(types.AttributeKeyDenomID, msg.DenomId),
			sdk.NewAttribute(types.AttributeKeyMTID, msg.Id),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Sender),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgBurnMTResponse{}, nil
}

// TransferDenom transfers denom
func (m msgServer) TransferDenom(goCtx context.Context, msg *types.MsgTransferDenom) (*types.MsgTransferDenomResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.TransferDenomOwner(ctx, msg.Id, sender, recipient); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransferDenom,
			sdk.NewAttribute(types.AttributeKeyDenomID, msg.Id),
			sdk.NewAttribute(types.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyRecipient, msg.Recipient),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgTransferDenomResponse{}, nil
}
