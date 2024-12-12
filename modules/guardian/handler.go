package guardian

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorstypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irishub/v4/modules/guardian/keeper"
	"github.com/irisnet/irishub/v4/modules/guardian/types"
)

// NewHandler returns a handler for all "guardian" type messages.
func NewHandler(k keeper.Keeper) func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgAddSuper:
			res, err := msgServer.AddSuper(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgDeleteSuper:
			res, err := msgServer.DeleteSuper(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, sdkerrors.Wrapf(errorstypes.ErrUnknownRequest, "unrecognized bank message type: %T", msg)
		}
	}
}
