package token

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HandleIssueToken handles MsgIssueToken
func HandleIssueToken(ctx sdk.Context, k Keeper, msg MsgIssueToken) (*sdk.Result, error) {
	token := NewFungibleToken(
		msg.Symbol, msg.Name, msg.Scale,
		msg.MinUnit, sdk.NewInt(int64(msg.InitialSupply)),
		sdk.NewInt(int64(msg.MaxSupply)), msg.Mintable, msg.Owner,
	)

	if err := k.IssueToken(ctx, token); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
		sdk.NewEvent(
			EventTypeIssueToken,
			sdk.NewAttribute(AttributeKeyTokenSymbol, token.GetSymbol()),
			sdk.NewAttribute(AttributeKeyTokenDenom, token.GetMinUnit()),
			sdk.NewAttribute(AttributeKeyTokenOwner, token.GetOwner().String()),
		),
	})

	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

// HandleMsgEditToken handles MsgEditToken
func HandleMsgEditToken(ctx sdk.Context, k Keeper, msg MsgEditToken) (*sdk.Result, error) {
	if err := k.EditToken(ctx, msg); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
		sdk.NewEvent(
			EventTypeEditToken,
			sdk.NewAttribute(AttributeKeyTokenSymbol, msg.Symbol),
		),
	})

	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

// HandleMsgTransferToken handles MsgTransferToken
func HandleMsgTransferToken(ctx sdk.Context, k Keeper, msg MsgTransferToken) (*sdk.Result, error) {
	if err := k.TransferToken(ctx, msg); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.SrcOwner.String()),
		),
		sdk.NewEvent(
			EventTypeTransferToken,
			sdk.NewAttribute(AttributeKeyTokenSymbol, msg.Symbol),
		),
	})

	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

// handleMsgMintToken handles MsgMintToken
func HandleMsgMintToken(ctx sdk.Context, k Keeper, msg MsgMintToken) (*sdk.Result, error) {
	mintCoin, err := k.MintToken(ctx, msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
		sdk.NewEvent(
			EventTypeMintToken,
			sdk.NewAttribute(AttributeKeyTokenSymbol, msg.Symbol),
			sdk.NewAttribute(sdk.AttributeKeyAmount, mintCoin.String()),
		),
	})

	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

// HandleMsgBurnToken handles MsgBurnToken
func HandleMsgBurnToken(ctx sdk.Context, k Keeper, msg MsgBurnToken) (*sdk.Result, error) {
	if err := k.BurnToken(ctx, msg); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
		sdk.NewEvent(
			EventTypeBurnToken,
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		),
	})

	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}
