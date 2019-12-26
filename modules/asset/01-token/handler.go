package token

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HandleIssueToken handles MsgIssueToken
func HandleIssueToken(ctx sdk.Context, k Keeper, msg MsgIssueToken) sdk.Result {
	token := NewFungibleToken(
		msg.Symbol, msg.Name, msg.Scale,
		msg.MinUnit, sdk.NewInt(int64(msg.InitialSupply)),
		sdk.NewInt(int64(msg.MaxSupply)), msg.Mintable, msg.Owner,
	)

	if err := k.IssueToken(ctx, token); err != nil {
		return err.Result()
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

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// HandleMsgEditToken handles MsgEditToken
func HandleMsgEditToken(ctx sdk.Context, k Keeper, msg MsgEditToken) sdk.Result {
	if err := k.EditToken(ctx, msg); err != nil {
		return err.Result()
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

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// HandleMsgTransferToken handles MsgTransferToken
func HandleMsgTransferToken(ctx sdk.Context, k Keeper, msg MsgTransferToken) sdk.Result {
	if err := k.TransferToken(ctx, msg); err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.SrcOwner.String()),
		),
		sdk.NewEvent(
			EventTypeTransferTokenOwner,
			sdk.NewAttribute(AttributeKeyTokenSymbol, msg.Symbol),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// handleMsgMintToken handles MsgMintToken
func HandleMsgMintToken(ctx sdk.Context, k Keeper, msg MsgMintToken) sdk.Result {
	if err := k.MintToken(ctx, msg); err != nil {
		return err.Result()
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
			sdk.NewAttribute(sdk.AttributeKeyAmount, string(msg.Amount)),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// HandleMsgBurnToken handles MsgMintToken
func HandleMsgBurnToken(ctx sdk.Context, k Keeper, msg MsgBurnToken) sdk.Result {
	if err := k.BurnToken(ctx, msg); err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
		sdk.NewEvent(
			EventTypeMintToken,
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
