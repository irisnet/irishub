package asset

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	token "github.com/irisnet/irishub/modules/asset/01-token"
)

// NewHandler returns a handler for all "asset" type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case token.MsgIssueToken:
			return token.HandleIssueToken(ctx, k.TokenKeeper, msg)
		case token.MsgEditToken:
			return token.HandleMsgEditToken(ctx, k.TokenKeeper, msg)
		case token.MsgMintToken:
			return token.HandleMsgMintToken(ctx, k.TokenKeeper, msg)
		case token.MsgTransferToken:
			return token.HandleMsgTransferToken(ctx, k.TokenKeeper, msg)
		case token.MsgBurnToken:
			return token.HandleMsgBurnToken(ctx, k.TokenKeeper, msg)
		//TODO NFT，IBC-Token
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}
