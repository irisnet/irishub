package asset

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	token "github.com/irisnet/irishub/modules/asset/01-token"
)

// handle all "asset" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case token.MsgIssueToken:
			return token.HandleIssueToken(ctx, k.TokenKeeper, msg)
		case token.MsgEditToken:
			return token.HandleMsgEditToken(ctx, k.TokenKeeper, msg)
		case token.MsgMintToken:
			return token.HandleMsgMintToken(ctx, k.TokenKeeper, msg)
		case token.MsgTransferTokenOwner:
			return token.HandleMsgTransferTokenOwner(ctx, k.TokenKeeper, msg)
		//TODO NFT，IBC-Token
		default:
			return sdk.ErrTxDecode("invalid message parse in asset module").Result()
		}
	}
}
