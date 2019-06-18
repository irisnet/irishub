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
		return ErrInvalidAssetFamily(DefaultCodespace, fmt.Sprintf("invalid token family type %s", msg.Family)).Result()
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
