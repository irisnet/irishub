package asset

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/asset/internal/types"
)

// handle all "asset" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
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
			return sdk.ErrTxDecode("invalid message parse in asset module").Result()
		}

		return sdk.ErrTxDecode("invalid message parse in asset module").Result()
	}
}

// handleIssueToken handles MsgIssueToken
func handleIssueToken(ctx sdk.Context, k Keeper, msg MsgIssueToken) sdk.Result {
	var token FungibleToken
	switch msg.Family {
	case FUNGIBLE:
		decimal := int(msg.Decimal)
		token = NewFungibleToken(msg.Source, msg.Symbol, msg.Name, msg.Decimal, msg.CanonicalSymbol,
			msg.MinUnitAlias, sdk.NewIntWithDecimal(int64(msg.InitialSupply), decimal),
			sdk.NewIntWithDecimal(int64(msg.MaxSupply), decimal), msg.Mintable, msg.Owner)
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

	err := k.IssueToken(ctx, token)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
		sdk.NewEvent(
			types.EventTypeIssueToken,
			sdk.NewAttribute(types.AttributeKeyTokenID, token.GetUniqueID()),
			sdk.NewAttribute(types.AttributeKeyTokenDenom, token.GetDenom()),
			sdk.NewAttribute(types.AttributeKeyTokenSource, token.GetSource().String()),
			sdk.NewAttribute(types.AttributeKeyTokenOwner, token.GetOwner().String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// handleMsgEditToken handles MsgEditToken
func handleMsgEditToken(ctx sdk.Context, k Keeper, msg MsgEditToken) sdk.Result {
	err := k.EditToken(ctx, msg)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
		sdk.NewEvent(
			types.EventTypeEditToken,
			sdk.NewAttribute(types.AttributeKeyTokenID, msg.TokenId),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// handleMsgTransferTokenOwner handles MsgTransferTokenOwner
func handleMsgTransferTokenOwner(ctx sdk.Context, k Keeper, msg MsgTransferTokenOwner) sdk.Result {
	err := k.TransferTokenOwner(ctx, msg)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.SrcOwner.String()),
		),
		sdk.NewEvent(
			types.EventTypeTransferTokenOwner,
			sdk.NewAttribute(types.AttributeKeyTokenID, msg.TokenId),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// handleMsgMintToken handles MsgMintToken
func handleMsgMintToken(ctx sdk.Context, k Keeper, msg MsgMintToken) sdk.Result {
	err := k.MintToken(ctx, msg)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
		sdk.NewEvent(
			types.EventTypeMintToken,
			sdk.NewAttribute(types.AttributeKeyTokenID, msg.TokenId),
			sdk.NewAttribute(sdk.AttributeKeyAmount, string(msg.Amount)),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
