package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/oracle/keeper"
	"github.com/irisnet/irismod/modules/oracle/types"
)

// NewHandler returns a handler for all the "oracle" type messages
func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgCreateFeed:
			m := msg.Normalize()
			res, err := msgServer.CreateFeed(sdk.WrapSDKContext(ctx), &m)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgStartFeed:
			m := msg.Normalize()
			res, err := msgServer.StartFeed(sdk.WrapSDKContext(ctx), &m)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgPauseFeed:
			m := msg.Normalize()
			res, err := msgServer.PauseFeed(sdk.WrapSDKContext(ctx), &m)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgEditFeed:
			m := msg.Normalize()
			res, err := msgServer.EditFeed(sdk.WrapSDKContext(ctx), &m)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}
