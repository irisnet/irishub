package guardian

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irishub/modules/guardian/keeper"
	"github.com/irisnet/irishub/modules/guardian/types"
)

// NewHandler returns a handler for all "guardian" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgAddProfiler:
			return handleMsgAddProfiler(ctx, k, msg)
		case *types.MsgAddTrustee:
			return handleMsgAddTrustee(ctx, k, msg)
		case *types.MsgDeleteProfiler:
			return handleMsgDeleteProfiler(ctx, k, msg)
		case *types.MsgDeleteTrustee:
			return handleMsgDeleteTrustee(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized bank message type: %T", msg)
		}
	}
}

// handleMsgAddProfiler handles MsgAddProfiler
func handleMsgAddProfiler(ctx sdk.Context, k keeper.Keeper, msg *types.MsgAddProfiler) (*sdk.Result, error) {
	if profiler, found := k.GetProfiler(ctx, msg.AddGuardian.AddedBy); !found || profiler.GetAccountType() != types.Genesis {
		return nil, sdkerrors.Wrap(types.ErrUnknownOperator, msg.AddGuardian.AddedBy.String())
	}
	if _, found := k.GetProfiler(ctx, msg.AddGuardian.Address); found {
		return nil, sdkerrors.Wrap(types.ErrProfilerExists, msg.AddGuardian.Address.String())
	}
	profiler := types.NewGuardian(msg.AddGuardian.Description, types.Ordinary, msg.AddGuardian.Address, msg.AddGuardian.AddedBy)
	k.AddProfiler(ctx, profiler)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.AddGuardian.AddedBy.String()),
		),
		sdk.NewEvent(
			types.EventTypeAddProfiler,
			sdk.NewAttribute(types.AttributeKeyProfilerAddress, msg.AddGuardian.Address.String()),
			sdk.NewAttribute(types.AttributeKeyAddedBy, msg.AddGuardian.AddedBy.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// handleMsgAddTrustee handles MsgAddTrustee
func handleMsgAddTrustee(ctx sdk.Context, k keeper.Keeper, msg *types.MsgAddTrustee) (*sdk.Result, error) {
	if trustee, found := k.GetTrustee(ctx, msg.AddGuardian.AddedBy); !found || trustee.GetAccountType() != types.Genesis {
		return nil, sdkerrors.Wrap(types.ErrUnknownOperator, msg.AddGuardian.AddedBy.String())
	}
	if _, found := k.GetTrustee(ctx, msg.AddGuardian.Address); found {
		return nil, sdkerrors.Wrap(types.ErrTrusteeExists, msg.AddGuardian.Address.String())
	}
	trustee := types.NewGuardian(msg.AddGuardian.Description, types.Ordinary, msg.AddGuardian.Address, msg.AddGuardian.AddedBy)
	k.AddTrustee(ctx, trustee)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.AddGuardian.AddedBy.String()),
		),
		sdk.NewEvent(
			types.EventTypeAddTrustee,
			sdk.NewAttribute(types.AttributeKeyTrusteeAddress, msg.AddGuardian.Address.String()),
			sdk.NewAttribute(types.AttributeKeyAddedBy, msg.AddGuardian.AddedBy.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// handleMsgDeleteProfiler handles MsgDeleteProfiler
func handleMsgDeleteProfiler(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDeleteProfiler) (*sdk.Result, error) {
	if profiler, found := k.GetProfiler(ctx, msg.DeleteGuardian.DeletedBy); !found || profiler.GetAccountType() != types.Genesis {
		return nil, sdkerrors.Wrap(types.ErrUnknownOperator, msg.DeleteGuardian.DeletedBy.String())
	}
	profiler, found := k.GetProfiler(ctx, msg.DeleteGuardian.Address)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrUnknownProfiler, msg.DeleteGuardian.Address.String())
	}
	if profiler.GetAccountType() == types.Genesis {
		return nil, sdkerrors.Wrap(types.ErrDeleteGenesisProfiler, msg.DeleteGuardian.Address.String())
	}

	k.DeleteProfiler(ctx, msg.DeleteGuardian.Address)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DeleteGuardian.DeletedBy.String()),
		),
		sdk.NewEvent(
			types.EventTypeDeleteProfiler,
			sdk.NewAttribute(types.AttributeKeyProfilerAddress, msg.DeleteGuardian.Address.String()),
			sdk.NewAttribute(types.AttributeKeyDeletedBy, msg.DeleteGuardian.DeletedBy.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// handleMsgDeleteTrustee handles MsgDeleteTrustee
func handleMsgDeleteTrustee(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDeleteTrustee) (*sdk.Result, error) {
	if trustee, found := k.GetTrustee(ctx, msg.DeleteGuardian.DeletedBy); !found || trustee.GetAccountType() != types.Genesis {
		return nil, sdkerrors.Wrap(types.ErrUnknownOperator, msg.DeleteGuardian.DeletedBy.String())
	}
	trustee, found := k.GetTrustee(ctx, msg.DeleteGuardian.Address)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrUnknownTrustee, msg.DeleteGuardian.Address.String())
	}
	if trustee.GetAccountType() == types.Genesis {
		return nil, sdkerrors.Wrap(types.ErrDeleteGenesisTrustee, msg.DeleteGuardian.Address.String())
	}

	k.DeleteTrustee(ctx, msg.DeleteGuardian.Address)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DeleteGuardian.DeletedBy.String()),
		),
		sdk.NewEvent(
			types.EventTypeDeleteTrustee,
			sdk.NewAttribute(types.AttributeKeyTrusteeAddress, msg.DeleteGuardian.Address.String()),
			sdk.NewAttribute(types.AttributeKeyDeletedBy, msg.DeleteGuardian.DeletedBy.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
