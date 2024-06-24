package keeper

import (
	"context"
	"encoding/hex"

	tmbytes "github.com/cometbft/cometbft/libs/bytes"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"mods.irisnet.org/service/types"
)

type msgServer struct {
	k Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the service MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{k: keeper}
}

// DefineService handles MsgDefineService
func (m msgServer) DefineService(
	goCtx context.Context,
	msg *types.MsgDefineService,
) (*types.MsgDefineServiceResponse, error) {
	author, err := sdk.AccAddressFromBech32(msg.Author)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.k.AddServiceDefinition(
		ctx, msg.Name, msg.Description, msg.Tags, author,
		msg.AuthorDescription, msg.Schemas,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateDefinition,
			sdk.NewAttribute(types.AttributeKeyServiceName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyAuthor, msg.Author),
		),
	})

	return &types.MsgDefineServiceResponse{}, nil
}

// BindService handles MsgBindService
func (m msgServer) BindService(
	goCtx context.Context,
	msg *types.MsgBindService,
) (*types.MsgBindServiceResponse, error) {
	if _, _, found := m.k.GetModuleServiceByServiceName(msg.ServiceName); found {
		return nil, errorsmod.Wrapf(
			types.ErrBindModuleService,
			"module service %s",
			msg.ServiceName,
		)
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
	if err := m.k.AddServiceBinding(
		ctx, msg.ServiceName, provider, msg.Deposit,
		msg.Pricing, msg.QoS, msg.Options, owner,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateBinding,
			sdk.NewAttribute(types.AttributeKeyServiceName, msg.ServiceName),
			sdk.NewAttribute(types.AttributeKeyProvider, msg.Provider),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner),
		),
	})

	return &types.MsgBindServiceResponse{}, nil
}

// UpdateServiceBinding handles MsgUpdateServiceBinding
func (m msgServer) UpdateServiceBinding(
	goCtx context.Context,
	msg *types.MsgUpdateServiceBinding,
) (*types.MsgUpdateServiceBindingResponse, error) {
	provider, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return nil, err
	}
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.k.UpdateServiceBinding(
		ctx, msg.ServiceName, provider, msg.Deposit,
		msg.Pricing, msg.QoS, msg.Options, owner,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateBinding,
			sdk.NewAttribute(types.AttributeKeyServiceName, msg.ServiceName),
			sdk.NewAttribute(types.AttributeKeyProvider, msg.Provider),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner),
		),
	})

	return &types.MsgUpdateServiceBindingResponse{}, nil
}

// SetWithdrawAddress handles MsgSetWithdrawAddress
func (m msgServer) SetWithdrawAddress(
	goCtx context.Context,
	msg *types.MsgSetWithdrawAddress,
) (*types.MsgSetWithdrawAddressResponse, error) {
	withdrawAddress, err := sdk.AccAddressFromBech32(msg.WithdrawAddress)
	if err != nil {
		return nil, err
	}
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	if m.k.blockedAddrs[msg.WithdrawAddress] {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrUnauthorized,
			"%s is a module account",
			msg.WithdrawAddress,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	m.k.SetWithdrawAddress(ctx, owner, withdrawAddress)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSetWithdrawAddress,
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner),
			sdk.NewAttribute(types.AttributeKeyWithdrawAddress, msg.WithdrawAddress),
		),
	})

	return &types.MsgSetWithdrawAddressResponse{}, nil
}

// EnableServiceBinding handles MsgEnableServiceBinding
func (m msgServer) EnableServiceBinding(
	goCtx context.Context,
	msg *types.MsgEnableServiceBinding,
) (*types.MsgEnableServiceBindingResponse, error) {
	provider, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return nil, err
	}
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.k.EnableServiceBinding(ctx, msg.ServiceName, provider, msg.Deposit, owner); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEnableBinding,
			sdk.NewAttribute(types.AttributeKeyServiceName, msg.ServiceName),
			sdk.NewAttribute(types.AttributeKeyProvider, msg.Provider),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner),
		),
	})

	return &types.MsgEnableServiceBindingResponse{}, nil
}

// DisableServiceBinding handles MsgDisableServiceBinding
func (m msgServer) DisableServiceBinding(
	goCtx context.Context,
	msg *types.MsgDisableServiceBinding,
) (*types.MsgDisableServiceBindingResponse, error) {
	provider, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return nil, err
	}
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.k.DisableServiceBinding(ctx, msg.ServiceName, provider, owner); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDisableBinding,
			sdk.NewAttribute(types.AttributeKeyServiceName, msg.ServiceName),
			sdk.NewAttribute(types.AttributeKeyProvider, msg.Provider),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner),
		),
	})

	return &types.MsgDisableServiceBindingResponse{}, nil
}

// RefundServiceDeposit handles MsgRefundServiceDeposit
func (m msgServer) RefundServiceDeposit(
	goCtx context.Context,
	msg *types.MsgRefundServiceDeposit,
) (*types.MsgRefundServiceDepositResponse, error) {
	provider, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return nil, err
	}
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.k.RefundDeposit(ctx, msg.ServiceName, provider, owner); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRefundDeposit,
			sdk.NewAttribute(types.AttributeKeyServiceName, msg.ServiceName),
			sdk.NewAttribute(types.AttributeKeyProvider, msg.Provider),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner),
		),
	})

	return &types.MsgRefundServiceDepositResponse{}, nil
}

// CallService handles MsgCallService
func (m msgServer) CallService(
	goCtx context.Context,
	msg *types.MsgCallService,
) (*types.MsgCallServiceResponse, error) {
	var reqContextID tmbytes.HexBytes
	var err error

	consumer, err := sdk.AccAddressFromBech32(msg.Consumer)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	_, moduleService, found := m.k.GetModuleServiceByServiceName(msg.ServiceName)
	if !found {
		pds := make([]sdk.AccAddress, len(msg.Providers))
		for i, provider := range msg.Providers {
			pd, err := sdk.AccAddressFromBech32(provider)
			if err != nil {
				return nil, err
			}
			pds[i] = pd
		}

		if reqContextID, err = m.k.CreateRequestContext(
			ctx, msg.ServiceName, pds, consumer,
			msg.Input, msg.ServiceFeeCap, msg.Timeout,
			msg.Repeated, msg.RepeatedFrequency,
			msg.RepeatedTotal, types.RUNNING, 0, "",
		); err != nil {
			return nil, err
		}
	} else {
		if reqContextID, err = m.k.CreateRequestContext(
			ctx, msg.ServiceName, []sdk.AccAddress{moduleService.Provider}, consumer,
			msg.Input, msg.ServiceFeeCap, 1, false, 0, 0, types.RUNNING, 0, "",
		); err != nil {
			return nil, err
		}

		if err := m.k.RequestModuleService(ctx, moduleService, reqContextID, consumer, msg.Input); err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateContext,
			sdk.NewAttribute(types.AttributeKeyRequestContextID, reqContextID.String()),
			sdk.NewAttribute(types.AttributeKeyServiceName, msg.ServiceName),
			sdk.NewAttribute(types.AttributeKeyConsumer, msg.Consumer),
		),
	})

	return &types.MsgCallServiceResponse{
		RequestContextId: reqContextID.String(),
	}, nil
}

// RespondService handles MsgRespondService
func (m msgServer) RespondService(
	goCtx context.Context,
	msg *types.MsgRespondService,
) (*types.MsgRespondServiceResponse, error) {
	provider, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return nil, err
	}

	requestId, err := hex.DecodeString(msg.RequestId)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	request, _, err := m.k.AddResponse(ctx, requestId, provider, msg.Result, msg.Output)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRespondService,
			sdk.NewAttribute(types.AttributeKeyRequestContextID, request.RequestContextId),
			sdk.NewAttribute(types.AttributeKeyRequestID, msg.RequestId),
			sdk.NewAttribute(types.AttributeKeyServiceName, request.ServiceName),
			sdk.NewAttribute(types.AttributeKeyProvider, request.Provider),
			sdk.NewAttribute(types.AttributeKeyConsumer, request.Consumer),
		),
	})

	return &types.MsgRespondServiceResponse{}, nil
}

// PauseRequestContext handles MsgPauseRequestContext
func (m msgServer) PauseRequestContext(
	goCtx context.Context,
	msg *types.MsgPauseRequestContext,
) (*types.MsgPauseRequestContextResponse, error) {
	consumer, err := sdk.AccAddressFromBech32(msg.Consumer)
	if err != nil {
		return nil, err
	}
	requestContextId, err := hex.DecodeString(msg.RequestContextId)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.k.CheckAuthority(ctx, consumer, requestContextId, true); err != nil {
		return nil, err
	}

	if err := m.k.PauseRequestContext(ctx, requestContextId, consumer); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypePauseContext,
			sdk.NewAttribute(types.AttributeKeyRequestContextID, msg.RequestContextId),
			sdk.NewAttribute(types.AttributeKeyConsumer, msg.Consumer),
		),
	})

	return &types.MsgPauseRequestContextResponse{}, nil
}

// StartRequestContext handles MsgStartRequestContext
func (m msgServer) StartRequestContext(
	goCtx context.Context,
	msg *types.MsgStartRequestContext,
) (*types.MsgStartRequestContextResponse, error) {
	consumer, err := sdk.AccAddressFromBech32(msg.Consumer)
	if err != nil {
		return nil, err
	}
	requestContextId, err := hex.DecodeString(msg.RequestContextId)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.k.CheckAuthority(ctx, consumer, requestContextId, true); err != nil {
		return nil, err
	}

	if err := m.k.StartRequestContext(ctx, requestContextId, consumer); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeStartContext,
			sdk.NewAttribute(types.AttributeKeyRequestContextID, msg.RequestContextId),
			sdk.NewAttribute(types.AttributeKeyConsumer, msg.Consumer),
		),
	})

	return &types.MsgStartRequestContextResponse{}, nil
}

// KillRequestContext handles MsgKillRequestContext
func (m msgServer) KillRequestContext(
	goCtx context.Context,
	msg *types.MsgKillRequestContext,
) (*types.MsgKillRequestContextResponse, error) {
	consumer, err := sdk.AccAddressFromBech32(msg.Consumer)
	if err != nil {
		return nil, err
	}
	requestContextId, err := hex.DecodeString(msg.RequestContextId)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.k.CheckAuthority(ctx, consumer, requestContextId, true); err != nil {
		return nil, err
	}

	if err := m.k.KillRequestContext(ctx, requestContextId, consumer); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeKillContext,
			sdk.NewAttribute(types.AttributeKeyRequestContextID, msg.RequestContextId),
			sdk.NewAttribute(types.AttributeKeyConsumer, msg.Consumer),
		),
	})

	return &types.MsgKillRequestContextResponse{}, nil
}

// UpdateRequestContext handles MsgUpdateRequestContext
func (m msgServer) UpdateRequestContext(
	goCtx context.Context,
	msg *types.MsgUpdateRequestContext,
) (*types.MsgUpdateRequestContextResponse, error) {
	consumer, err := sdk.AccAddressFromBech32(msg.Consumer)
	if err != nil {
		return nil, err
	}

	requestContextId, err := hex.DecodeString(msg.RequestContextId)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.k.CheckAuthority(ctx, consumer, requestContextId, true); err != nil {
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

	if err := m.k.UpdateRequestContext(
		ctx, requestContextId, pds, 0, msg.ServiceFeeCap,
		msg.Timeout, msg.RepeatedFrequency, msg.RepeatedTotal, consumer,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateContext,
			sdk.NewAttribute(types.AttributeKeyRequestContextID, msg.RequestContextId),
			sdk.NewAttribute(types.AttributeKeyConsumer, msg.Consumer),
		),
	})

	return &types.MsgUpdateRequestContextResponse{}, nil
}

// WithdrawEarnedFees handles MsgWithdrawEarnedFees
func (m msgServer) WithdrawEarnedFees(
	goCtx context.Context,
	msg *types.MsgWithdrawEarnedFees,
) (*types.MsgWithdrawEarnedFeesResponse, error) {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}
	provider, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.k.WithdrawEarnedFees(ctx, owner, provider); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdrawEarnedFees,
			sdk.NewAttribute(types.AttributeKeyProvider, msg.Provider),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner),
		),
	})

	return &types.MsgWithdrawEarnedFeesResponse{}, nil
}

func (m msgServer) UpdateParams(
	goCtx context.Context,
	msg *types.MsgUpdateParams,
) (*types.MsgUpdateParamsResponse, error) {
	if m.k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrUnauthorized,
			"invalid authority; expected %s, got %s",
			m.k.authority,
			msg.Authority,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.k.SetParams(ctx, msg.Params); err != nil {
		return nil, err
	}
	return &types.MsgUpdateParamsResponse{}, nil
}
