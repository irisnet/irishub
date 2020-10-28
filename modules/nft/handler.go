package nft

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/nft/keeper"
	"github.com/irisnet/irismod/modules/nft/types"
)

// NewHandler routes the messages to the handlers
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgIssueDenom:
			return HandleMsgIssueDenom(ctx, msg, k)
		case *types.MsgMintNFT:
			return HandleMsgMintNFT(ctx, msg, k)
		case *types.MsgTransferNFT:
			return HandleMsgTransferNFT(ctx, msg, k)
		case *types.MsgEditNFT:
			return HandleMsgEditNFT(ctx, msg, k)
		case *types.MsgBurnNFT:
			return HandleMsgBurnNFT(ctx, msg, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized nft message type: %T", msg)
		}
	}
}

func HandleMsgIssueDenom(ctx sdk.Context, msg *types.MsgIssueDenom, k keeper.Keeper) (*sdk.Result, error) {
	id := strings.ToLower(strings.TrimSpace(msg.Id))
	name := strings.ToLower(strings.TrimSpace(msg.Name))

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	if err := k.IssueDenom(ctx, id, name, msg.Schema, sender); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeIssueDenom,
			sdk.NewAttribute(types.AttributeKeyDenom, id),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// HandleMsgTransferNFT handler for MsgTransferNFT
func HandleMsgTransferNFT(ctx sdk.Context, msg *types.MsgTransferNFT, k keeper.Keeper) (*sdk.Result, error) {
	id := strings.ToLower(strings.TrimSpace(msg.Id))
	denom := strings.ToLower(strings.TrimSpace(msg.Denom))

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return nil, err
	}

	if err := k.TransferOwner(
		ctx, denom, id,
		strings.TrimSpace(msg.Name),
		strings.TrimSpace(msg.URI),
		msg.Data,
		sender,
		recipient,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransfer,
			sdk.NewAttribute(types.AttributeKeyRecipient, msg.Recipient),
			sdk.NewAttribute(types.AttributeKeyDenom, denom),
			sdk.NewAttribute(types.AttributeKeyTokenID, id),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// HandleMsgEditNFT handler for MsgEditNFT
func HandleMsgEditNFT(ctx sdk.Context, msg *types.MsgEditNFT, k keeper.Keeper) (*sdk.Result, error) {
	id := strings.ToLower(strings.TrimSpace(msg.Id))
	denom := strings.ToLower(strings.TrimSpace(msg.Denom))

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	if err := k.EditNFT(
		ctx, denom, id,
		strings.TrimSpace(msg.Name),
		strings.TrimSpace(msg.URI),
		msg.Data,
		sender,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEditNFT,
			sdk.NewAttribute(types.AttributeKeyDenom, denom),
			sdk.NewAttribute(types.AttributeKeyTokenID, id),
			sdk.NewAttribute(types.AttributeKeyTokenURI, msg.URI),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// HandleMsgMintNFT handles MsgMintNFT
func HandleMsgMintNFT(ctx sdk.Context, msg *types.MsgMintNFT, k keeper.Keeper) (*sdk.Result, error) {
	id := strings.ToLower(strings.TrimSpace(msg.Id))
	denom := strings.ToLower(strings.TrimSpace(msg.Denom))

	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return nil, err
	}

	if err := k.MintNFT(
		ctx, denom, id,
		strings.TrimSpace(msg.Name),
		strings.TrimSpace(msg.URI),
		msg.Data,
		recipient,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMintNFT,
			sdk.NewAttribute(types.AttributeKeyRecipient, msg.Recipient),
			sdk.NewAttribute(types.AttributeKeyDenom, denom),
			sdk.NewAttribute(types.AttributeKeyTokenID, id),
			sdk.NewAttribute(types.AttributeKeyTokenURI, msg.URI),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// HandleMsgBurnNFT handles MsgBurnNFT
func HandleMsgBurnNFT(ctx sdk.Context, msg *types.MsgBurnNFT, k keeper.Keeper) (*sdk.Result, error) {
	id := strings.ToLower(strings.TrimSpace(msg.Id))
	denom := strings.ToLower(strings.TrimSpace(msg.Denom))

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	if err := k.BurnNFT(ctx, denom, id, sender); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnNFT,
			sdk.NewAttribute(types.AttributeKeyDenom, denom),
			sdk.NewAttribute(types.AttributeKeyTokenID, id),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
