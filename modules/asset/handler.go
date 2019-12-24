package asset

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for all "asset" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgIssueToken:
			return handleIssueToken(ctx, k, msg)
		case MsgEditToken:
			return handleMsgEditToken(ctx, k, msg)
		case MsgMintToken:
			return handleMsgMintToken(ctx, k, msg)
		case MsgTransferTokenOwner:
			return handleMsgTransferTokenOwner(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized asset message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleIssueToken handles MsgIssueToken
func handleIssueToken(ctx sdk.Context, k Keeper, msg MsgIssueToken) sdk.Result {
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

// handleMsgEditToken handles MsgEditToken
func handleMsgEditToken(ctx sdk.Context, k Keeper, msg MsgEditToken) sdk.Result {
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

// handleMsgTransferTokenOwner handles MsgTransferTokenOwner
func handleMsgTransferTokenOwner(ctx sdk.Context, k Keeper, msg MsgTransferTokenOwner) sdk.Result {
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
func handleMsgMintToken(ctx sdk.Context, k Keeper, msg MsgMintToken) sdk.Result {
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
