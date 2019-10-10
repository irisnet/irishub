package asset

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// handle all "asset" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgIssueToken:
			return handleIssueToken(ctx, k, msg)
		case MsgCreateGateway:
			return handleMsgCreateGateway(ctx, k, msg)
		case MsgEditGateway:
			return handleMsgEditGateway(ctx, k, msg)
		case MsgEditToken:
			return handleMsgEditToken(ctx, k, msg)
		case MsgTransferGatewayOwner:
			return handleMsgTransferGatewayOwner(ctx, k, msg)
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
		token = NewFungibleToken(msg.Source, msg.Gateway, msg.Symbol, msg.Name, msg.Decimal, msg.CanonicalSymbol, msg.MinUnitAlias, sdk.NewIntWithDecimal(int64(msg.InitialSupply), decimal), sdk.NewIntWithDecimal(int64(msg.MaxSupply), decimal), msg.Mintable, msg.Owner)
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
	case GATEWAY:
		// handle fee for gateway token
		if err := GatewayTokenIssueFeeHandler(ctx, k, msg.Owner, msg.Symbol); err != nil {
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

	return sdk.Result{Events: ctx.EventManager().Events()}
}

// handleMsgCreateGateway handles MsgCreateGateway
func handleMsgCreateGateway(ctx sdk.Context, k Keeper, msg MsgCreateGateway) sdk.Result {
	// handle fee
	if err := GatewayCreateFeeHandler(ctx, k, msg.Owner, msg.Moniker); err != nil {
		return err.Result()
	}

	// convert moniker to lowercase
	msg.Moniker = strings.ToLower(msg.Moniker)

	err := k.CreateGateway(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}

// handleMsgEditGateway handles MsgEditGateway
func handleMsgEditGateway(ctx sdk.Context, k Keeper, msg MsgEditGateway) sdk.Result {
	// convert moniker to lowercase
	msg.Moniker = strings.ToLower(msg.Moniker)

	err := k.EditGateway(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}

// handleMsgEditToken handles MsgEditToken
func handleMsgEditToken(ctx sdk.Context, k Keeper, msg MsgEditToken) sdk.Result {
	err := k.EditToken(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}

// handleMsgTransferGatewayOwner handles MsgTransferGatewayOwner
func handleMsgTransferGatewayOwner(ctx sdk.Context, k Keeper, msg MsgTransferGatewayOwner) sdk.Result {
	// convert moniker to lowercase
	msg.Moniker = strings.ToLower(msg.Moniker)

	err := k.TransferGatewayOwner(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}

// handleMsgTransferTokenOwner handles MsgTransferTokenOwner
func handleMsgTransferTokenOwner(ctx sdk.Context, k Keeper, msg MsgTransferTokenOwner) sdk.Result {
	err := k.TransferTokenOwner(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}

// handleMsgMintToken handles MsgMintToken
func handleMsgMintToken(ctx sdk.Context, k Keeper, msg MsgMintToken) sdk.Result {
	err := k.MintToken(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}
