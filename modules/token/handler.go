package token

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/token/keeper"
	"github.com/irisnet/irismod/modules/token/types"
)

// NewHandler handle all "token" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgIssueToken:
			m := msg.Normalize()
			res, err := msgServer.IssueToken(sdk.WrapSDKContext(ctx), &m)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgEditToken:
			m := msg.Normalize()
			res, err := msgServer.EditToken(sdk.WrapSDKContext(ctx), &m)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgMintToken:
			m := msg.Normalize()
			res, err := msgServer.MintToken(sdk.WrapSDKContext(ctx), &m)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgBurnToken:
			m := msg.Normalize()
			res, err := msgServer.BurnToken(sdk.WrapSDKContext(ctx), &m)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgTransferTokenOwner:
			m := msg.Normalize()
			res, err := msgServer.TransferTokenOwner(sdk.WrapSDKContext(ctx), &m)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized nft message type: %T", msg)
		}
	}
}
