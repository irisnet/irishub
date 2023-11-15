package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irishub/v2/modules/guardian/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the guardian MsgServer interface for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) AddSuper(goCtx context.Context, msg *types.MsgAddSuper) (*types.MsgAddSuperResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	addedBy, err := sdk.AccAddressFromBech32(msg.AddedBy)
	if err != nil {
		return nil, err
	}
	address, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, err
	}
	if super, found := m.Keeper.GetSuper(ctx, addedBy); !found || super.GetAccountType() != types.Genesis {
		return nil, sdkerrors.Wrap(types.ErrUnknownOperator, msg.AddedBy)
	}
	if _, found := m.Keeper.GetSuper(ctx, address); found {
		return nil, sdkerrors.Wrap(types.ErrSuperExists, msg.Address)
	}
	super := types.NewSuper(msg.Description, types.Ordinary, address, addedBy)
	m.Keeper.AddSuper(ctx, super)

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

	return &types.MsgAddSuperResponse{}, nil
}

func (m msgServer) DeleteSuper(goCtx context.Context, msg *types.MsgDeleteSuper) (*types.MsgDeleteSuperResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	deletedBy, err := sdk.AccAddressFromBech32(msg.DeletedBy)
	if err != nil {
		return nil, err
	}
	address, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, err
	}

	if super, found := m.Keeper.GetSuper(ctx, deletedBy); !found || super.GetAccountType() != types.Genesis {
		return nil, sdkerrors.Wrap(types.ErrUnknownOperator, msg.DeletedBy)
	}
	super, found := m.Keeper.GetSuper(ctx, address)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrUnknownSuper, msg.Address)
	}
	if super.GetAccountType() == types.Genesis {
		return nil, sdkerrors.Wrap(types.ErrDeleteGenesisSuper, msg.Address)
	}

	m.Keeper.DeleteSuper(ctx, address)

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

	return &types.MsgDeleteSuperResponse{}, nil
}
