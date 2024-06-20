package mt

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/mt/keeper"
	"github.com/irisnet/irismod/mt/types"
)

// NewHandler routes the messages to the handlers
func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgIssueDenom:
			res, err := msgServer.IssueDenom(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgMintMT:
			res, err := msgServer.MintMT(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgTransferMT:
			res, err := msgServer.TransferMT(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgEditMT:
			res, err := msgServer.EditMT(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgBurnMT:
			res, err := msgServer.BurnMT(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgTransferDenom:
			res, err := msgServer.TransferDenom(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized mt message type: %T", msg)
		}
	}
}
