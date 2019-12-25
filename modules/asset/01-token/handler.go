package token

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HandleIssueToken handles MsgIssueToken
func HandleIssueToken(ctx sdk.Context, k Keeper, msg MsgIssueToken) sdk.Result {
	var token FungibleToken
	switch msg.Family {
	case FUNGIBLE:
		token = NewFungibleToken(
			msg.Source, msg.Symbol, msg.Name, msg.Decimal, msg.CanonicalSymbol,
			msg.MinUnitAlias, sdk.NewIntWithDecimal(int64(msg.InitialSupply), int(msg.Decimal)),
			sdk.NewIntWithDecimal(int64(msg.MaxSupply), int(msg.Decimal)), msg.Mintable, msg.Owner,
		)
	default:
		return ErrInvalidAssetFamily(DefaultCodespace, fmt.Sprintf("invalid asset family type %s", msg.Family)).Result()
	}

	switch msg.Source {
	case NATIVE:
		// handle fee for native token
		if err := TokenIssueFeeHandler(ctx, k, msg.Owner, msg.Symbol); err != nil {
			return err.Result()
		}
		break
	default:
		break
	}

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
			sdk.NewAttribute(AttributeKeyTokenID, token.GetUniqueID()),
			sdk.NewAttribute(AttributeKeyTokenDenom, token.GetDenom()),
			sdk.NewAttribute(AttributeKeyTokenSource, token.GetSource().String()),
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
			sdk.NewAttribute(AttributeKeyTokenID, msg.TokenID),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// HandleMsgTransferTokenOwner handles MsgTransferTokenOwner
func HandleMsgTransferTokenOwner(ctx sdk.Context, k Keeper, msg MsgTransferTokenOwner) sdk.Result {
	if err := k.TransferTokenOwner(ctx, msg); err != nil {
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
			sdk.NewAttribute(AttributeKeyTokenID, msg.TokenID),
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
			sdk.NewAttribute(AttributeKeyTokenID, msg.TokenID),
			sdk.NewAttribute(sdk.AttributeKeyAmount, string(msg.Amount)),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
