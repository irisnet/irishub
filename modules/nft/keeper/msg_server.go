package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"mods.irisnet.org/nft/types"
)

var _ types.MsgServer = Keeper{}

// IssueDenom issue a new denom.
func (k Keeper) IssueDenom(
	goCtx context.Context,
	msg *types.MsgIssueDenom,
) (*types.MsgIssueDenomResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.SaveDenom(ctx, msg.Id,
		msg.Name,
		msg.Schema,
		msg.Symbol,
		sender,
		msg.MintRestricted,
		msg.UpdateRestricted,
		msg.Description,
		msg.Uri,
		msg.UriHash,
		msg.Data,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeIssueDenom,
			sdk.NewAttribute(types.AttributeKeyDenomID, msg.Id),
			sdk.NewAttribute(types.AttributeKeyDenomName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Sender),
		),
	})

	return &types.MsgIssueDenomResponse{}, nil
}

func (k Keeper) MintNFT(
	goCtx context.Context,
	msg *types.MsgMintNFT,
) (*types.MsgMintNFTResponse, error) {
	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return nil, err
	}

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	denom, err := k.GetDenomInfo(ctx, msg.DenomId)
	if err != nil {
		return nil, err
	}

	if denom.MintRestricted && denom.Creator != sender.String() {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrUnauthorized,
			"%s is not allowed to mint NFT of denom %s",
			sender,
			msg.DenomId,
		)
	}

	if err := k.SaveNFT(ctx, msg.DenomId,
		msg.Id,
		msg.Name,
		msg.URI,
		msg.UriHash,
		msg.Data,
		recipient,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMintNFT,
			sdk.NewAttribute(types.AttributeKeyTokenID, msg.Id),
			sdk.NewAttribute(types.AttributeKeyDenomID, msg.DenomId),
			sdk.NewAttribute(types.AttributeKeyTokenURI, msg.URI),
			sdk.NewAttribute(types.AttributeKeyRecipient, msg.Recipient),
		),
	})

	return &types.MsgMintNFTResponse{}, nil
}

func (k Keeper) EditNFT(
	goCtx context.Context,
	msg *types.MsgEditNFT,
) (*types.MsgEditNFTResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.UpdateNFT(ctx, msg.DenomId,
		msg.Id,
		msg.Name,
		msg.URI,
		msg.UriHash,
		msg.Data,
		sender,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEditNFT,
			sdk.NewAttribute(types.AttributeKeyTokenID, msg.Id),
			sdk.NewAttribute(types.AttributeKeyDenomID, msg.DenomId),
			sdk.NewAttribute(types.AttributeKeyTokenURI, msg.URI),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Sender),
		),
	})

	return &types.MsgEditNFTResponse{}, nil
}

func (k Keeper) TransferNFT(goCtx context.Context,
	msg *types.MsgTransferNFT) (*types.MsgTransferNFTResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.TransferOwnership(ctx, msg.DenomId,
		msg.Id,
		msg.Name,
		msg.URI,
		msg.UriHash,
		msg.Data,
		sender,
		recipient,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransfer,
			sdk.NewAttribute(types.AttributeKeyTokenID, msg.Id),
			sdk.NewAttribute(types.AttributeKeyDenomID, msg.DenomId),
			sdk.NewAttribute(types.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyRecipient, msg.Recipient),
		),
	})

	return &types.MsgTransferNFTResponse{}, nil
}

func (k Keeper) BurnNFT(
	goCtx context.Context,
	msg *types.MsgBurnNFT,
) (*types.MsgBurnNFTResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.RemoveNFT(ctx, msg.DenomId, msg.Id, sender); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnNFT,
			sdk.NewAttribute(types.AttributeKeyDenomID, msg.DenomId),
			sdk.NewAttribute(types.AttributeKeyTokenID, msg.Id),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Sender),
		),
	})

	return &types.MsgBurnNFTResponse{}, nil
}

func (k Keeper) TransferDenom(
	goCtx context.Context,
	msg *types.MsgTransferDenom,
) (*types.MsgTransferDenomResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.TransferDenomOwner(ctx, msg.Id, sender, recipient); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransferDenom,
			sdk.NewAttribute(types.AttributeKeyDenomID, msg.Id),
			sdk.NewAttribute(types.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyRecipient, msg.Recipient),
		),
	})

	return &types.MsgTransferDenomResponse{}, nil
}
