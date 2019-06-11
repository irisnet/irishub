package asset

import (
	"strings"

	sdk "github.com/irisnet/irishub/types"
)

// handle all "asset" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgIssueAsset:
			return handleIssueAsset(ctx, k, msg)
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

// handleIssueAsset handles MsgIssueAsset
func handleIssueAsset(ctx sdk.Context, k Keeper, msg MsgIssueAsset) sdk.Result {
	var asset Asset
	switch msg.Family {
	case FUNGIBLE:
		asset = NewFungibleToken(msg.Source, msg.Gateway, msg.Symbol, msg.Name, msg.Decimal, msg.SymbolMinAlias, msg.InitialSupply, msg.MaxSupply, msg.Mintable, msg.Owner)
	default:
		return ErrInvalidAssetFamily(k.codespace, msg.Family).Result()
	}
	tags, err := k.IssueAsset(ctx, asset)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}

// handleMsgCreateGateway handles MsgCreateGateway
func handleMsgCreateGateway(ctx sdk.Context, k Keeper, msg MsgCreateGateway) sdk.Result {
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
