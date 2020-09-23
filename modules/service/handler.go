package service

import (
	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/service/keeper"
	"github.com/irisnet/irismod/modules/service/types"
)

// NewHandler creates an sdk.Handler for all the service type messages
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgDefineService:
			return handleMsgDefineService(ctx, k, msg)

		case *types.MsgBindService:
			return handleMsgBindService(ctx, k, msg)

		case *types.MsgUpdateServiceBinding:
			return handleMsgUpdateServiceBinding(ctx, k, msg)

		case *types.MsgSetWithdrawAddress:
			return handleMsgSetWithdrawAddress(ctx, k, msg)

		case *types.MsgDisableServiceBinding:
			return handleMsgDisableServiceBinding(ctx, k, msg)

		case *types.MsgEnableServiceBinding:
			return handleMsgEnableServiceBinding(ctx, k, msg)

		case *types.MsgRefundServiceDeposit:
			return handleMsgRefundServiceDeposit(ctx, k, msg)

		case *types.MsgCallService:
			return handleMsgCallService(ctx, k, msg)

		case *types.MsgRespondService:
			return handleMsgRespondService(ctx, k, msg)

		case *types.MsgPauseRequestContext:
			return handleMsgPauseRequestContext(ctx, k, msg)

		case *types.MsgStartRequestContext:
			return handleMsgStartRequestContext(ctx, k, msg)

		case *types.MsgKillRequestContext:
			return handleMsgKillRequestContext(ctx, k, msg)

		case *types.MsgUpdateRequestContext:
			return handleMsgUpdateRequestContext(ctx, k, msg)

		case *types.MsgWithdrawEarnedFees:
			return handleMsgWithdrawEarnedFees(ctx, k, msg)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

func handleMsgDefineService(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDefineService) (*sdk.Result, error) {
	if err := k.AddServiceDefinition(
		ctx, msg.Name, msg.Description, msg.Tags, msg.Author,
		msg.AuthorDescription, msg.Schemas,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Author.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgBindService(ctx sdk.Context, k keeper.Keeper, msg *types.MsgBindService) (*sdk.Result, error) {
	if _, _, found := k.GetModuleServiceByServiceName(msg.ServiceName); found {
		return nil, sdkerrors.Wrapf(types.ErrBindModuleService, "module service %s", msg.ServiceName)
	}

	if err := k.AddServiceBinding(
		ctx, msg.ServiceName, msg.Provider, msg.Deposit,
		msg.Pricing, msg.QoS, msg.Options, msg.Owner,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgUpdateServiceBinding(ctx sdk.Context, k keeper.Keeper, msg *types.MsgUpdateServiceBinding) (*sdk.Result, error) {
	if err := k.UpdateServiceBinding(
		ctx, msg.ServiceName, msg.Provider, msg.Deposit,
		msg.Pricing, msg.QoS, msg.Options, msg.Owner,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgSetWithdrawAddress(ctx sdk.Context, k keeper.Keeper, msg *types.MsgSetWithdrawAddress) (*sdk.Result, error) {
	k.SetWithdrawAddress(ctx, msg.Owner, msg.WithdrawAddress)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgDisableServiceBinding(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDisableServiceBinding) (*sdk.Result, error) {
	if err := k.DisableServiceBinding(ctx, msg.ServiceName, msg.Provider, msg.Owner); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgEnableServiceBinding(ctx sdk.Context, k keeper.Keeper, msg *types.MsgEnableServiceBinding) (*sdk.Result, error) {
	if err := k.EnableServiceBinding(ctx, msg.ServiceName, msg.Provider, msg.Deposit, msg.Owner); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgRefundServiceDeposit(ctx sdk.Context, k keeper.Keeper, msg *types.MsgRefundServiceDeposit) (*sdk.Result, error) {
	if err := k.RefundDeposit(ctx, msg.ServiceName, msg.Provider, msg.Owner); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// handleMsgCallService handles MsgCallService
func handleMsgCallService(ctx sdk.Context, k keeper.Keeper, msg *types.MsgCallService) (*sdk.Result, error) {
	var reqContextID tmbytes.HexBytes
	var err error

	_, moduleService, found := k.GetModuleServiceByServiceName(msg.ServiceName)
	if !found {
		if reqContextID, err = k.CreateRequestContext(
			ctx, msg.ServiceName, msg.Providers, msg.Consumer,
			msg.Input, msg.ServiceFeeCap, msg.Timeout,
			msg.SuperMode, msg.Repeated, msg.RepeatedFrequency,
			msg.RepeatedTotal, types.RUNNING, 0, "",
		); err != nil {
			return nil, err
		}
	} else {
		if reqContextID, err = k.CreateRequestContext(
			ctx, msg.ServiceName, []sdk.AccAddress{moduleService.Provider}, msg.Consumer,
			msg.Input, msg.ServiceFeeCap, 1, false, false, 0, 0, types.RUNNING, 0, "",
		); err != nil {
			return nil, err
		}

		if err := k.RequestModuleService(ctx, moduleService, reqContextID, msg.Consumer, msg.Input); err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Consumer.String()),
			sdk.NewAttribute(types.AttributeKeyRequestContextID, reqContextID.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// handleMsgRespondService handles MsgRespondService
func handleMsgRespondService(ctx sdk.Context, k keeper.Keeper, msg *types.MsgRespondService) (*sdk.Result, error) {
	request, _, err := k.AddResponse(ctx, msg.RequestId, msg.Provider, msg.Result, msg.Output)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Provider.String()),
			sdk.NewAttribute(types.AttributeKeyRequestContextID, request.RequestContextId.String()),
			sdk.NewAttribute(types.AttributeKeyRequestID, msg.RequestId.String()),
			sdk.NewAttribute(types.AttributeKeyConsumer, request.Consumer.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// handleMsgPauseRequestContext handles MsgPauseRequestContext
func handleMsgPauseRequestContext(ctx sdk.Context, k keeper.Keeper, msg *types.MsgPauseRequestContext) (*sdk.Result, error) {
	if err := k.CheckAuthority(ctx, msg.Consumer, msg.RequestContextId, true); err != nil {
		return nil, err
	}

	if err := k.PauseRequestContext(ctx, msg.RequestContextId, msg.Consumer); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Consumer.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// handleMsgStartRequestContext handles MsgStartRequestContext
func handleMsgStartRequestContext(ctx sdk.Context, k keeper.Keeper, msg *types.MsgStartRequestContext) (*sdk.Result, error) {
	if err := k.CheckAuthority(ctx, msg.Consumer, msg.RequestContextId, true); err != nil {
		return nil, err
	}

	if err := k.StartRequestContext(ctx, msg.RequestContextId, msg.Consumer); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Consumer.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// handleMsgKillRequestContext handles MsgKillRequestContext
func handleMsgKillRequestContext(ctx sdk.Context, k keeper.Keeper, msg *types.MsgKillRequestContext) (*sdk.Result, error) {
	if err := k.CheckAuthority(ctx, msg.Consumer, msg.RequestContextId, true); err != nil {
		return nil, err
	}

	if err := k.KillRequestContext(ctx, msg.RequestContextId, msg.Consumer); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Consumer.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// handleMsgUpdateRequestContext handles MsgUpdateRequestContext
func handleMsgUpdateRequestContext(ctx sdk.Context, k keeper.Keeper, msg *types.MsgUpdateRequestContext) (*sdk.Result, error) {
	if err := k.CheckAuthority(ctx, msg.Consumer, msg.RequestContextId, true); err != nil {
		return nil, err
	}

	if err := k.UpdateRequestContext(
		ctx, msg.RequestContextId, msg.Providers, 0, msg.ServiceFeeCap,
		msg.Timeout, msg.RepeatedFrequency, msg.RepeatedTotal, msg.Consumer,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Consumer.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// handleMsgWithdrawEarnedFees handles MsgWithdrawEarnedFees
func handleMsgWithdrawEarnedFees(ctx sdk.Context, k keeper.Keeper, msg *types.MsgWithdrawEarnedFees) (*sdk.Result, error) {
	if err := k.WithdrawEarnedFees(ctx, msg.Owner, msg.Provider); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
