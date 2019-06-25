package asset

import (
	"fmt"
	"strings"

	sdk "github.com/irisnet/irishub/types"
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
		totalSupply := msg.InitialSupply
		decimal := int(msg.Decimal)
		token = NewFungibleToken(msg.Source, msg.Gateway, msg.Symbol, msg.Name, msg.Decimal, msg.SymbolAtSource, msg.SymbolMinAlias, sdk.NewIntWithDecimal(int64(msg.InitialSupply), decimal), sdk.NewIntWithDecimal(int64(totalSupply), decimal), sdk.NewIntWithDecimal(int64(msg.MaxSupply), decimal), msg.Mintable, msg.Owner)
	default:
		return ErrInvalidAssetFamily(DefaultCodespace, fmt.Sprintf("invalid asset family type %s", msg.Family)).Result()
	}

	switch msg.Source {
	case NATIVE:
		// handle fee for native token
		if err := TokenIssueFeeHandler(ctx, k, msg.Owner, msg.Symbol, getTokenIssueFee(ctx, k, msg.Symbol)); err != nil {
			return err.Result()
		}
		break
	case GATEWAY:
		// handle fee for gateway token
		if err := GatewayTokenIssueFeeHandler(ctx, k, msg.Owner, msg.Symbol, getGatewayTokenIssueFee(ctx, k, msg.Symbol)); err != nil {
			return err.Result()
		}
		break
	default:
		break
	}

	tags, err := k.IssueToken(ctx, token)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}

// handleMsgCreateGateway handles MsgCreateGateway
func handleMsgCreateGateway(ctx sdk.Context, k Keeper, msg MsgCreateGateway) sdk.Result {
	// handle fee
	if err := GatewayFeeHandler(ctx, k, msg.Owner, msg.Moniker, msg.Fee); err != nil {
		return err.Result()
	}

	// convert moniker to lowercase
	msg.Moniker = strings.ToLower(msg.Moniker)

	tags, err := k.CreateGateway(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}

// handleMsgEditGateway handles MsgEditGateway
func handleMsgEditGateway(ctx sdk.Context, k Keeper, msg MsgEditGateway) sdk.Result {
	// convert moniker to lowercase
	msg.Moniker = strings.ToLower(msg.Moniker)

	tags, err := k.EditGateway(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}

// handleMsgEditToken handles MsgEditToken
func handleMsgEditToken(ctx sdk.Context, k Keeper, msg MsgEditToken) sdk.Result {
	tags, err := k.EditToken(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

// handleMsgTransferGatewayOwner handles MsgTransferGatewayOwner
func handleMsgTransferGatewayOwner(ctx sdk.Context, k Keeper, msg MsgTransferGatewayOwner) sdk.Result {
	// convert moniker to lowercase
	msg.Moniker = strings.ToLower(msg.Moniker)

	tags, err := k.TransferGatewayOwner(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}
