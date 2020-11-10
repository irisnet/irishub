package keeper

import (
	"context"
	"encoding/hex"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irismod/modules/service/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the service MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// DefineService handles MsgDefineService
func (m msgServer) DefineService(goCtx context.Context, msg *types.MsgDefineService) (*types.MsgDefineServiceResponse, error) {
	author, err := sdk.AccAddressFromBech32(msg.Author)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.AddServiceDefinition(
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
	return &types.MsgDefineServiceResponse{}, nil
}

// BindService handles MsgBindService
func (m msgServer) BindService(goCtx context.Context, msg *types.MsgBindService) (*types.MsgBindServiceResponse, error) {
	if _, _, found := m.Keeper.GetModuleServiceByServiceName(msg.ServiceName); found {
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

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.AddServiceBinding(
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
	return &types.MsgBindServiceResponse{}, nil
}

// UpdateServiceBinding handles MsgUpdateServiceBinding
func (m msgServer) UpdateServiceBinding(goCtx context.Context, msg *types.MsgUpdateServiceBinding) (*types.MsgUpdateServiceBindingResponse, error) {
	provider, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return nil, err
	}
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.UpdateServiceBinding(
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
	return &types.MsgUpdateServiceBindingResponse{}, nil
}

// SetWithdrawAddress handles MsgSetWithdrawAddress
func (m msgServer) SetWithdrawAddress(goCtx context.Context, msg *types.MsgSetWithdrawAddress) (*types.MsgSetWithdrawAddressResponse, error) {
	withdrawAddress, err := sdk.AccAddressFromBech32(msg.WithdrawAddress)
	if err != nil {
		return nil, err
	}
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	m.Keeper.SetWithdrawAddress(ctx, owner, withdrawAddress)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner),
		),
	})
	return &types.MsgSetWithdrawAddressResponse{}, nil
}

// EnableServiceBinding handles MsgEnableServiceBinding
func (m msgServer) EnableServiceBinding(goCtx context.Context, msg *types.MsgEnableServiceBinding) (*types.MsgEnableServiceBindingResponse, error) {
	provider, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return nil, err
	}
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.EnableServiceBinding(ctx, msg.ServiceName, provider, msg.Deposit, owner); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner),
		),
	})
	return &types.MsgEnableServiceBindingResponse{}, nil
}

// DisableServiceBinding handles MsgDisableServiceBinding
func (m msgServer) DisableServiceBinding(goCtx context.Context, msg *types.MsgDisableServiceBinding) (*types.MsgDisableServiceBindingResponse, error) {
	provider, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return nil, err
	}
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.DisableServiceBinding(ctx, msg.ServiceName, provider, owner); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner),
		),
	})
	return &types.MsgDisableServiceBindingResponse{}, nil
}

// RefundServiceDeposit handles MsgRefundServiceDeposit
func (m msgServer) RefundServiceDeposit(goCtx context.Context, msg *types.MsgRefundServiceDeposit) (*types.MsgRefundServiceDepositResponse, error) {
	provider, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return nil, err
	}
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.RefundDeposit(ctx, msg.ServiceName, provider, owner); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner),
		),
	})
	return &types.MsgRefundServiceDepositResponse{}, nil
}

// CallService handles MsgCallService
func (m msgServer) CallService(goCtx context.Context, msg *types.MsgCallService) (*types.MsgCallServiceResponse, error) {
	var reqContextID tmbytes.HexBytes
	var err error

	consumer, err := sdk.AccAddressFromBech32(msg.Consumer)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	_, moduleService, found := m.Keeper.GetModuleServiceByServiceName(msg.ServiceName)
	if !found {
		pds := make([]sdk.AccAddress, len(msg.Providers))
		for i, provider := range msg.Providers {
			pd, err := sdk.AccAddressFromBech32(provider)
			if err != nil {
				return nil, err
			}
			pds[i] = pd
		}

		if reqContextID, err = m.Keeper.CreateRequestContext(
			ctx, msg.ServiceName, pds, consumer,
			msg.Input, msg.ServiceFeeCap, msg.Timeout,
			msg.SuperMode, msg.Repeated, msg.RepeatedFrequency,
			msg.RepeatedTotal, types.RUNNING, 0, "",
		); err != nil {
			return nil, err
		}
	} else {
		if reqContextID, err = m.Keeper.CreateRequestContext(
			ctx, msg.ServiceName, []sdk.AccAddress{moduleService.Provider}, consumer,
			msg.Input, msg.ServiceFeeCap, 1, false, false, 0, 0, types.RUNNING, 0, "",
		); err != nil {
			return nil, err
		}

		if err := m.Keeper.RequestModuleService(ctx, moduleService, reqContextID, consumer, msg.Input); err != nil {
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
	return &types.MsgCallServiceResponse{
		RequestContextId: reqContextID.String(),
	}, nil
}

// RespondService handles MsgRespondService
func (m msgServer) RespondService(goCtx context.Context, msg *types.MsgRespondService) (*types.MsgRespondServiceResponse, error) {
	provider, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return nil, err
	}

	requestId, err := hex.DecodeString(msg.RequestId)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	request, _, err := m.Keeper.AddResponse(ctx, requestId, provider, msg.Result, msg.Output)
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
	return &types.MsgRespondServiceResponse{}, nil
}

// PauseRequestContext handles MsgPauseRequestContext
func (m msgServer) PauseRequestContext(goCtx context.Context, msg *types.MsgPauseRequestContext) (*types.MsgPauseRequestContextResponse, error) {
	consumer, err := sdk.AccAddressFromBech32(msg.Consumer)
	if err != nil {
		return nil, err
	}
	requestContextId, err := hex.DecodeString(msg.RequestContextId)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.CheckAuthority(ctx, consumer, requestContextId, true); err != nil {
		return nil, err
	}

	if err := m.Keeper.PauseRequestContext(ctx, requestContextId, consumer); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Consumer),
		),
	})
	return &types.MsgPauseRequestContextResponse{}, nil
}

// StartRequestContext handles MsgStartRequestContext
func (m msgServer) StartRequestContext(goCtx context.Context, msg *types.MsgStartRequestContext) (*types.MsgStartRequestContextResponse, error) {
	consumer, err := sdk.AccAddressFromBech32(msg.Consumer)
	if err != nil {
		return nil, err
	}
	requestContextId, err := hex.DecodeString(msg.RequestContextId)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.CheckAuthority(ctx, consumer, requestContextId, true); err != nil {
		return nil, err
	}

	if err := m.Keeper.StartRequestContext(ctx, requestContextId, consumer); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Consumer),
		),
	})
	return &types.MsgStartRequestContextResponse{}, nil
}

// KillRequestContext handles MsgKillRequestContext
func (m msgServer) KillRequestContext(goCtx context.Context, msg *types.MsgKillRequestContext) (*types.MsgKillRequestContextResponse, error) {
	consumer, err := sdk.AccAddressFromBech32(msg.Consumer)
	if err != nil {
		return nil, err
	}
	requestContextId, err := hex.DecodeString(msg.RequestContextId)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.CheckAuthority(ctx, consumer, requestContextId, true); err != nil {
		return nil, err
	}

	if err := m.Keeper.KillRequestContext(ctx, requestContextId, consumer); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Consumer),
		),
	})
	return &types.MsgKillRequestContextResponse{}, nil
}

// UpdateRequestContext handles MsgUpdateRequestContext
func (m msgServer) UpdateRequestContext(goCtx context.Context, msg *types.MsgUpdateRequestContext) (*types.MsgUpdateRequestContextResponse, error) {
	consumer, err := sdk.AccAddressFromBech32(msg.Consumer)
	if err != nil {
		return nil, err
	}

	requestContextId, err := hex.DecodeString(msg.RequestContextId)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.CheckAuthority(ctx, consumer, requestContextId, true); err != nil {
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

	if err := m.Keeper.UpdateRequestContext(
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
	return &types.MsgUpdateRequestContextResponse{}, nil
}

// WithdrawEarnedFees handles MsgWithdrawEarnedFees
func (m msgServer) WithdrawEarnedFees(goCtx context.Context, msg *types.MsgWithdrawEarnedFees) (*types.MsgWithdrawEarnedFeesResponse, error) {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}
	provider, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.WithdrawEarnedFees(ctx, owner, provider); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner),
		),
	})
	return &types.MsgWithdrawEarnedFeesResponse{}, nil
}
