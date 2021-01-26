package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/token/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the token MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) IssueToken(goCtx context.Context, msg *types.MsgIssueToken) (*types.MsgIssueTokenResponse, error) {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// handle fee for token
	if err := m.Keeper.DeductIssueTokenFee(ctx, owner, msg.Symbol); err != nil {
		return nil, err
	}

	if err := m.Keeper.IssueToken(
		ctx, msg.Symbol, msg.Name, msg.MinUnit, msg.Scale,
		msg.InitialSupply, msg.MaxSupply, msg.Mintable, owner,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeIssueToken,
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Owner),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner),
		),
	})

	return &types.MsgIssueTokenResponse{}, nil
}

func (m msgServer) EditToken(goCtx context.Context, msg *types.MsgEditToken) (*types.MsgEditTokenResponse, error) {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := m.Keeper.EditToken(
		ctx, msg.Symbol, msg.Name,
		msg.MaxSupply, msg.Mintable, owner,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEditToken,
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner),
		),
	})

	return &types.MsgEditTokenResponse{}, nil
}

func (m msgServer) MintToken(goCtx context.Context, msg *types.MsgMintToken) (*types.MsgMintTokenResponse, error) {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	var recipient sdk.AccAddress

	if len(msg.To) != 0 {
		recipient, err = sdk.AccAddressFromBech32(msg.To)
		if err != nil {
			return nil, err
		}
	} else {
		recipient = owner
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := m.Keeper.DeductMintTokenFee(ctx, owner, msg.Symbol); err != nil {
		return nil, err
	}

	if err := m.Keeper.MintToken(ctx, msg.Symbol, msg.Amount, recipient, owner); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMintToken,
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeKeyAmount, strconv.FormatUint(msg.Amount, 10)),
			sdk.NewAttribute(types.AttributeKeyRecipient, recipient.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner),
		),
	})

	return &types.MsgMintTokenResponse{}, nil
}

func (m msgServer) BurnToken(goCtx context.Context, msg *types.MsgBurnToken) (*types.MsgBurnTokenResponse, error) {
	owner, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := m.Keeper.DeductMintTokenFee(ctx, owner, msg.Symbol); err != nil {
		return nil, err
	}

	if err := m.Keeper.BurnToken(ctx, msg.Symbol, msg.Amount, owner); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnToken,
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeKeyAmount, strconv.FormatUint(msg.Amount, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgBurnTokenResponse{}, nil
}

func (m msgServer) TransferTokenOwner(goCtx context.Context, msg *types.MsgTransferTokenOwner) (*types.MsgTransferTokenOwnerResponse, error) {
	srcOwner, err := sdk.AccAddressFromBech32(msg.SrcOwner)
	if err != nil {
		return nil, err
	}

	dstOwner, err := sdk.AccAddressFromBech32(msg.DstOwner)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := m.Keeper.TransferTokenOwner(ctx, msg.Symbol, srcOwner, dstOwner); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransferTokenOwner,
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.SrcOwner),
			sdk.NewAttribute(types.AttributeKeyDstOwner, msg.DstOwner),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.SrcOwner),
		),
	})

	return &types.MsgTransferTokenOwnerResponse{}, nil
}
