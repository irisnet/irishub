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
		case *types.MsgAddSuper:
			return handleMsgAddSuper(ctx, k, msg)
		case *types.MsgDeleteSuper:
			return handleMsgDeleteSuper(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized bank message type: %T", msg)
		}
	}
}

// handleMsgAddSuper handles MsgAddSuper
func handleMsgAddSuper(ctx sdk.Context, k keeper.Keeper, msg *types.MsgAddSuper) (*sdk.Result, error) {
	addedBy, err := sdk.AccAddressFromBech32(msg.AddedBy)
	if err != nil {
		return nil, err
	}
	address, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, err
	}
	if super, found := k.GetSuper(ctx, addedBy); !found || super.GetAccountType() != types.Genesis {
		return nil, sdkerrors.Wrap(types.ErrUnknownOperator, msg.AddedBy)
	}
	if _, found := k.GetSuper(ctx, address); found {
		return nil, sdkerrors.Wrap(types.ErrSuperExists, msg.Address)
	}
	super := types.NewSuper(msg.Description, types.Ordinary, address, addedBy)
	k.AddSuper(ctx, super)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.AddedBy),
		),
		sdk.NewEvent(
			types.EventTypeAddSuper,
			sdk.NewAttribute(types.AttributeKeySuperAddress, msg.Address),
			sdk.NewAttribute(types.AttributeKeyAddedBy, msg.AddedBy),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// handleMsgDeleteSuper handles MsgDeleteSuper
func handleMsgDeleteSuper(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDeleteSuper) (*sdk.Result, error) {
	deletedBy, err := sdk.AccAddressFromBech32(msg.DeletedBy)
	if err != nil {
		return nil, err
	}
	address, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, err
	}

	if super, found := k.GetSuper(ctx, deletedBy); !found || super.GetAccountType() != types.Genesis {
		return nil, sdkerrors.Wrap(types.ErrUnknownOperator, msg.DeletedBy)
	}
	super, found := k.GetSuper(ctx, address)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrUnknownSuper, msg.Address)
	}
	if super.GetAccountType() == types.Genesis {
		return nil, sdkerrors.Wrap(types.ErrDeleteGenesisSuper, msg.Address)
	}

	k.DeleteSuper(ctx, address)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DeletedBy),
		),
		sdk.NewEvent(
			types.EventTypeDeleteSuper,
			sdk.NewAttribute(types.AttributeKeySuperAddress, msg.Address),
			sdk.NewAttribute(types.AttributeKeyDeletedBy, msg.DeletedBy),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
