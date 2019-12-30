package guardian

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns a handler for all "guardian" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgAddProfiler:
			return handleMsgAddProfiler(ctx, k, msg)
		case MsgAddTrustee:
			return handleMsgAddTrustee(ctx, k, msg)
		case MsgDeleteProfiler:
			return handleMsgDeleteProfiler(ctx, k, msg)
		case MsgDeleteTrustee:
			return handleMsgDeleteTrustee(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

// handleMsgAddProfiler handles MsgAddProfiler
func handleMsgAddProfiler(ctx sdk.Context, k Keeper, msg MsgAddProfiler) (*sdk.Result, error) {
	if profiler, found := k.GetProfiler(ctx, msg.AddedBy); !found || profiler.GetAccountType() != Genesis {
		return nil, sdkerrors.Wrap(ErrUnknownOperator, msg.AddedBy.String())
	}
	if _, found := k.GetProfiler(ctx, msg.Address); found {
		return nil, sdkerrors.Wrap(ErrUnknownProfiler, msg.Address.String())
	}
	profiler := NewGuardian(msg.Description, Ordinary, msg.Address, msg.AddedBy)
	k.AddProfiler(ctx, profiler)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.AddedBy.String()),
		),
		sdk.NewEvent(
			EventTypeAddProfiler,
			sdk.NewAttribute(AttributeKeyProfilerAddress, msg.Address.String()),
			sdk.NewAttribute(AttributeKeyAddedBy, msg.AddedBy.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgAddTrustee handles MsgAddTrustee
func handleMsgAddTrustee(ctx sdk.Context, k Keeper, msg MsgAddTrustee) (*sdk.Result, error) {
	if trustee, found := k.GetTrustee(ctx, msg.AddedBy); !found || trustee.GetAccountType() != Genesis {
		return nil, sdkerrors.Wrap(ErrUnknownOperator, msg.AddedBy.String())
	}
	if _, found := k.GetTrustee(ctx, msg.Address); found {
		return nil, sdkerrors.Wrap(ErrUnknownTrustee, msg.Address.String())
	}
	trustee := NewGuardian(msg.Description, Ordinary, msg.Address, msg.AddedBy)
	k.AddTrustee(ctx, trustee)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.AddedBy.String()),
		),
		sdk.NewEvent(
			EventTypeAddTrustee,
			sdk.NewAttribute(AttributeKeyTrusteeAddress, msg.Address.String()),
			sdk.NewAttribute(AttributeKeyAddedBy, msg.AddedBy.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgDeleteProfiler handles MsgDeleteProfiler
func handleMsgDeleteProfiler(ctx sdk.Context, k Keeper, msg MsgDeleteProfiler) (*sdk.Result, error) {
	if profiler, found := k.GetProfiler(ctx, msg.DeletedBy); !found || profiler.GetAccountType() != Genesis {
		return nil, sdkerrors.Wrap(ErrUnknownOperator, msg.DeletedBy.String())
	}
	profiler, found := k.GetProfiler(ctx, msg.Address)
	if !found {
		return nil, sdkerrors.Wrap(ErrUnknownProfiler, msg.Address.String())
	}
	if profiler.GetAccountType() == Genesis {
		return nil, sdkerrors.Wrap(ErrDeleteGenesisProfiler, msg.Address.String())
	}

	k.DeleteProfiler(ctx, msg.Address)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DeletedBy.String()),
		),
		sdk.NewEvent(
			EventTypeDeleteProfiler,
			sdk.NewAttribute(AttributeKeyProfilerAddress, msg.Address.String()),
			sdk.NewAttribute(AttributeKeyDeletedBy, msg.DeletedBy.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgDeleteTrustee handles MsgDeleteTrustee
func handleMsgDeleteTrustee(ctx sdk.Context, k Keeper, msg MsgDeleteTrustee) (*sdk.Result, error) {
	if trustee, found := k.GetTrustee(ctx, msg.DeletedBy); !found || trustee.GetAccountType() != Genesis {
		return nil, sdkerrors.Wrap(ErrUnknownOperator, msg.DeletedBy.String())
	}
	trustee, found := k.GetTrustee(ctx, msg.Address)
	if !found {
		return nil, sdkerrors.Wrap(ErrUnknownTrustee, msg.Address.String())
	}
	if trustee.GetAccountType() == Genesis {
		return nil, sdkerrors.Wrap(ErrDeleteGenesisTrustee, msg.Address.String())
	}

	k.DeleteTrustee(ctx, msg.Address)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DeletedBy.String()),
		),
		sdk.NewEvent(
			EventTypeDeleteTrustee,
			sdk.NewAttribute(AttributeKeyTrusteeAddress, msg.Address.String()),
			sdk.NewAttribute(AttributeKeyDeletedBy, msg.DeletedBy.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
