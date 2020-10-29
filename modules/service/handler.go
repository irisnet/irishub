package service

import (
	"encoding/hex"

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
	author, err := sdk.AccAddressFromBech32(msg.Author)
	if err != nil {
		return nil, err
	}

	if err := k.AddServiceDefinition(
		ctx, msg.Name, msg.Description, msg.Tags, author,
		msg.AuthorDescription, msg.Schemas,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Author),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgBindService(ctx sdk.Context, k keeper.Keeper, msg *types.MsgBindService) (*sdk.Result, error) {
	if _, _, found := k.GetModuleServiceByServiceName(msg.ServiceName); found {
		return nil, sdkerrors.Wrapf(types.ErrBindModuleService, "module service %s", msg.ServiceName)
	}

	provider, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return nil, err
	}
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	if err := k.AddServiceBinding(
		ctx, msg.ServiceName, provider, msg.Deposit,
		msg.Pricing, msg.QoS, msg.Options, owner,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgUpdateServiceBinding(ctx sdk.Context, k keeper.Keeper, msg *types.MsgUpdateServiceBinding) (*sdk.Result, error) {
	provider, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return nil, err
	}
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	if err := k.UpdateServiceBinding(
		ctx, msg.ServiceName, provider, msg.Deposit,
		msg.Pricing, msg.QoS, msg.Options, owner,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgSetWithdrawAddress(ctx sdk.Context, k keeper.Keeper, msg *types.MsgSetWithdrawAddress) (*sdk.Result, error) {
	withdrawAddress, err := sdk.AccAddressFromBech32(msg.WithdrawAddress)
	if err != nil {
		return nil, err
	}
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	k.SetWithdrawAddress(ctx, owner, withdrawAddress)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgDisableServiceBinding(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDisableServiceBinding) (*sdk.Result, error) {
	provider, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return nil, err
	}
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	if err := k.DisableServiceBinding(ctx, msg.ServiceName, provider, owner); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgEnableServiceBinding(ctx sdk.Context, k keeper.Keeper, msg *types.MsgEnableServiceBinding) (*sdk.Result, error) {
	provider, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return nil, err
	}
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	if err := k.EnableServiceBinding(ctx, msg.ServiceName, provider, msg.Deposit, owner); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgRefundServiceDeposit(ctx sdk.Context, k keeper.Keeper, msg *types.MsgRefundServiceDeposit) (*sdk.Result, error) {
	provider, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return nil, err
	}
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	if err := k.RefundDeposit(ctx, msg.ServiceName, provider, owner); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// handleMsgCallService handles MsgCallService
func handleMsgCallService(ctx sdk.Context, k keeper.Keeper, msg *types.MsgCallService) (*sdk.Result, error) {
	var reqContextID tmbytes.HexBytes
	var err error

	consumer, err := sdk.AccAddressFromBech32(msg.Consumer)
	if err != nil {
		return nil, err
	}

	_, moduleService, found := k.GetModuleServiceByServiceName(msg.ServiceName)
	if !found {
		pds := make([]sdk.AccAddress, len(msg.Providers))
		for i, provider := range msg.Providers {
			pd, err := sdk.AccAddressFromBech32(provider)
			if err != nil {
				return nil, err
			}
			pds[i] = pd
		}

		if reqContextID, err = k.CreateRequestContext(
			ctx, msg.ServiceName, pds, consumer,
			msg.Input, msg.ServiceFeeCap, msg.Timeout,
			msg.SuperMode, msg.Repeated, msg.RepeatedFrequency,
			msg.RepeatedTotal, types.RUNNING, 0, "",
		); err != nil {
			return nil, err
		}
	} else {
		if reqContextID, err = k.CreateRequestContext(
			ctx, msg.ServiceName, []sdk.AccAddress{moduleService.Provider}, consumer,
			msg.Input, msg.ServiceFeeCap, 1, false, false, 0, 0, types.RUNNING, 0, "",
		); err != nil {
			return nil, err
		}

		if err := k.RequestModuleService(ctx, moduleService, reqContextID, consumer, msg.Input); err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Consumer),
			sdk.NewAttribute(types.AttributeKeyRequestContextID, reqContextID.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// handleMsgRespondService handles MsgRespondService
func handleMsgRespondService(ctx sdk.Context, k keeper.Keeper, msg *types.MsgRespondService) (*sdk.Result, error) {
	provider, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return nil, err
	}

	requestId, err := hex.DecodeString(msg.RequestId)
	if err != nil {
		return nil, err
	}

	request, _, err := k.AddResponse(ctx, requestId, provider, msg.Result, msg.Output)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Provider),
			sdk.NewAttribute(types.AttributeKeyRequestContextID, request.RequestContextId),
			sdk.NewAttribute(types.AttributeKeyRequestID, msg.RequestId),
			sdk.NewAttribute(types.AttributeKeyServiceName, request.ServiceName),
			sdk.NewAttribute(types.AttributeKeyConsumer, request.Consumer),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// handleMsgPauseRequestContext handles MsgPauseRequestContext
func handleMsgPauseRequestContext(ctx sdk.Context, k keeper.Keeper, msg *types.MsgPauseRequestContext) (*sdk.Result, error) {
	consumer, err := sdk.AccAddressFromBech32(msg.Consumer)
	if err != nil {
		return nil, err
	}
	requestContextId, err := hex.DecodeString(msg.RequestContextId)
	if err != nil {
		return nil, err
	}
	if err := k.CheckAuthority(ctx, consumer, requestContextId, true); err != nil {
		return nil, err
	}

	if err := k.PauseRequestContext(ctx, requestContextId, consumer); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Consumer),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// handleMsgStartRequestContext handles MsgStartRequestContext
func handleMsgStartRequestContext(ctx sdk.Context, k keeper.Keeper, msg *types.MsgStartRequestContext) (*sdk.Result, error) {
	consumer, err := sdk.AccAddressFromBech32(msg.Consumer)
	if err != nil {
		return nil, err
	}
	requestContextId, err := hex.DecodeString(msg.RequestContextId)
	if err != nil {
		return nil, err
	}

	if err := k.CheckAuthority(ctx, consumer, requestContextId, true); err != nil {
		return nil, err
	}

	if err := k.StartRequestContext(ctx, requestContextId, consumer); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Consumer),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// handleMsgKillRequestContext handles MsgKillRequestContext
func handleMsgKillRequestContext(ctx sdk.Context, k keeper.Keeper, msg *types.MsgKillRequestContext) (*sdk.Result, error) {
	consumer, err := sdk.AccAddressFromBech32(msg.Consumer)
	if err != nil {
		return nil, err
	}
	requestContextId, err := hex.DecodeString(msg.RequestContextId)
	if err != nil {
		return nil, err
	}

	if err := k.CheckAuthority(ctx, consumer, requestContextId, true); err != nil {
		return nil, err
	}

	if err := k.KillRequestContext(ctx, requestContextId, consumer); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Consumer),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// handleMsgUpdateRequestContext handles MsgUpdateRequestContext
func handleMsgUpdateRequestContext(ctx sdk.Context, k keeper.Keeper, msg *types.MsgUpdateRequestContext) (*sdk.Result, error) {
	consumer, err := sdk.AccAddressFromBech32(msg.Consumer)
	if err != nil {
		return nil, err
	}
	requestContextId, err := hex.DecodeString(msg.RequestContextId)
	if err != nil {
		return nil, err
	}

	if err := k.CheckAuthority(ctx, consumer, requestContextId, true); err != nil {
		return nil, err
	}

	pds := make([]sdk.AccAddress, len(msg.Providers))
	for i, provider := range msg.Providers {
		pd, err := sdk.AccAddressFromBech32(provider)
		if err != nil {
			return nil, err
		}
		pds[i] = pd
	}

	if err := k.UpdateRequestContext(
		ctx, requestContextId, pds, 0, msg.ServiceFeeCap,
		msg.Timeout, msg.RepeatedFrequency, msg.RepeatedTotal, consumer,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Consumer),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// handleMsgWithdrawEarnedFees handles MsgWithdrawEarnedFees
func handleMsgWithdrawEarnedFees(ctx sdk.Context, k keeper.Keeper, msg *types.MsgWithdrawEarnedFees) (*sdk.Result, error) {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}
	provider, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return nil, err
	}

	if err := k.WithdrawEarnedFees(ctx, owner, provider); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
