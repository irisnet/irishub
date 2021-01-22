package service

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/service/keeper"
	"github.com/irisnet/irismod/modules/service/types"
)

// NewHandler creates an sdk.Handler for all the service type messages
func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgDefineService:
			m := msg.Normalize()
			res, err := msgServer.DefineService(sdk.WrapSDKContext(ctx), &m)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgBindService:
			m := msg.Normalize()
			res, err := msgServer.BindService(sdk.WrapSDKContext(ctx), &m)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgUpdateServiceBinding:
			m := msg.Normalize()
			res, err := msgServer.UpdateServiceBinding(sdk.WrapSDKContext(ctx), &m)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgSetWithdrawAddress:
			m := msg.Normalize()
			res, err := msgServer.SetWithdrawAddress(sdk.WrapSDKContext(ctx), &m)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgDisableServiceBinding:
			m := msg.Normalize()
			res, err := msgServer.DisableServiceBinding(sdk.WrapSDKContext(ctx), &m)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgEnableServiceBinding:
			m := msg.Normalize()
			res, err := msgServer.EnableServiceBinding(sdk.WrapSDKContext(ctx), &m)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgRefundServiceDeposit:
			m := msg.Normalize()
			res, err := msgServer.RefundServiceDeposit(sdk.WrapSDKContext(ctx), &m)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgCallService:
			m := msg.Normalize()
			res, err := msgServer.CallService(sdk.WrapSDKContext(ctx), &m)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgRespondService:
			m := msg.Normalize()
			res, err := msgServer.RespondService(sdk.WrapSDKContext(ctx), &m)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgPauseRequestContext:
			m := msg.Normalize()
			res, err := msgServer.PauseRequestContext(sdk.WrapSDKContext(ctx), &m)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgStartRequestContext:
			m := msg.Normalize()
			res, err := msgServer.StartRequestContext(sdk.WrapSDKContext(ctx), &m)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgKillRequestContext:
			m := msg.Normalize()
			res, err := msgServer.KillRequestContext(sdk.WrapSDKContext(ctx), &m)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgUpdateRequestContext:
			m := msg.Normalize()
			res, err := msgServer.UpdateRequestContext(sdk.WrapSDKContext(ctx), &m)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgWithdrawEarnedFees:
			m := msg.Normalize()
			res, err := msgServer.WithdrawEarnedFees(sdk.WrapSDKContext(ctx), &m)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}
