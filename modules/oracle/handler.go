package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irishub/modules/oracle/keeper"
	"github.com/irisnet/irishub/modules/oracle/types"
)

// NewHandler returns a handler for all the "oracle" type messages
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.MsgCreateFeed:
			return handleMsgCreateFeed(ctx, k, msg)
		case *types.MsgStartFeed:
			return handleMsgStartFeed(ctx, k, msg)
		case *types.MsgPauseFeed:
			return handleMsgPauseFeed(ctx, k, msg)
		case *types.MsgEditFeed:
			return handleMsgEditFeed(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

// handleMsgCreateFeed handles MsgCreateFeed
func handleMsgCreateFeed(ctx sdk.Context, k keeper.Keeper, msg *types.MsgCreateFeed) (*sdk.Result, error) {
	err := k.CreateFeed(ctx, msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// handleMsgStartFeed handles MsgStartFeed
func handleMsgStartFeed(ctx sdk.Context, k keeper.Keeper, msg *types.MsgStartFeed) (*sdk.Result, error) {
	if err := k.StartFeed(ctx, msg); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// handleMsgPauseFeed handles MsgPauseFeed
func handleMsgPauseFeed(ctx sdk.Context, k keeper.Keeper, msg *types.MsgPauseFeed) (*sdk.Result, error) {
	if err := k.PauseFeed(ctx, msg); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// handleMsgEditFeed handles MsgEditFeed
func handleMsgEditFeed(ctx sdk.Context, k keeper.Keeper, msg *types.MsgEditFeed) (*sdk.Result, error) {
	if err := k.EditFeed(ctx, msg); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
