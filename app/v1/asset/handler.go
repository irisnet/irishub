package asset

import (
	sdk "github.com/irisnet/irishub/types"
)

// handle all "asset" type messages.
// TODO
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
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

// handleMsgCreateGateway handles MsgCreateGateway
func handleMsgCreateGateway(ctx sdk.Context, k Keeper, msg MsgCreateGateway) sdk.Result {
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
	tags, err := k.EditGateway(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}
