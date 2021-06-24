package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/farm/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the farm MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) CreatePool(goCtx context.Context, msg *types.MsgCreatePool) (*types.MsgCreatePoolResponse,
	error) {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	//check valid begin height
	if ctx.BlockHeight() > int64(msg.StartHeight) {
		return nil, sdkerrors.Wrapf(
			types.ErrExpiredHeight,
			"The current block height[%d] is greater than StartHeight[%d]",
			ctx.BlockHeight(),
			msg.StartHeight,
		)
	}

	if maxRewardCategories := m.Keeper.MaxRewardCategories(ctx); uint32(len(msg.TotalReward)) > maxRewardCategories {
		return nil, sdkerrors.Wrapf(
			types.ErrInvalidRewardRule,
			"the max reward category num is [%d], but got [%d]",
			maxRewardCategories,
			len(msg.TotalReward),
		)
	}

	//check pool exist
	if _, exist := m.Keeper.GetPool(ctx, msg.Name); exist {
		return nil, sdkerrors.Wrapf(types.ErrPoolExist, msg.Name)
	}

	//check valid lp token denom
	if err := m.Keeper.validateLPToken(ctx, msg.LpTokenDenom); err != nil {
		return nil, sdkerrors.Wrapf(
			types.ErrInvalidLPToken,
			"The lp token denom[%s] is not exist",
			msg.LpTokenDenom,
		)
	}
	if err = m.Keeper.CreatePool(ctx,
		msg.Name,
		msg.Description,
		msg.LpTokenDenom,
		msg.StartHeight,
		msg.RewardPerBlock.Sort(),
		msg.TotalReward.Sort(),
		msg.Editable,
		creator,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreatePool,
			sdk.NewAttribute(types.AttributeValueCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeValuePoolName, msg.Name),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	})
	return &types.MsgCreatePoolResponse{}, nil
}

func (m msgServer) DestroyPool(goCtx context.Context, msg *types.MsgDestroyPool) (*types.MsgDestroyPoolResponse, error) {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	refundCoin, err := m.Keeper.DestroyPool(ctx, msg.PoolName, creator)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDestroyPool,
			sdk.NewAttribute(types.AttributeValueCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeValuePoolName, msg.PoolName),
			sdk.NewAttribute(types.AttributeValueAmount, refundCoin.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	})
	return &types.MsgDestroyPoolResponse{}, nil
}

func (m msgServer) AdjustPool(goCtx context.Context, msg *types.MsgAdjustPool) (*types.MsgAdjustPoolResponse, error) {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err = m.Keeper.AdjustPool(ctx,
		msg.PoolName,
		msg.AdditionalReward,
		msg.RewardPerBlock,
		creator,
	); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAppendReward,
			sdk.NewAttribute(types.AttributeValueCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeValuePoolName, msg.PoolName),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	})
	return &types.MsgAdjustPoolResponse{}, nil
}

func (m msgServer) Stake(goCtx context.Context, msg *types.MsgStake) (*types.MsgStakeResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	reward, err := m.Keeper.Stake(ctx, msg.PoolName, msg.Amount, sender)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeStake,
			sdk.NewAttribute(types.AttributeValueCreator, msg.Sender),
			sdk.NewAttribute(types.AttributeValuePoolName, msg.PoolName),
			sdk.NewAttribute(types.AttributeValueAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeValueReward, reward.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})
	return &types.MsgStakeResponse{Reward: reward}, nil
}

func (m msgServer) Unstake(goCtx context.Context, msg *types.MsgUnstake) (*types.MsgUnstakeResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	reward, err := m.Keeper.Unstake(ctx, msg.PoolName, msg.Amount, sender)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnstake,
			sdk.NewAttribute(types.AttributeValueCreator, msg.Sender),
			sdk.NewAttribute(types.AttributeValuePoolName, msg.PoolName),
			sdk.NewAttribute(types.AttributeValueAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeValueReward, reward.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})
	return &types.MsgUnstakeResponse{Reward: reward}, nil
}

func (m msgServer) Harvest(goCtx context.Context, msg *types.MsgHarvest) (*types.MsgHarvestResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	reward, err := m.Keeper.Harvest(ctx, msg.PoolName, sender)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeHarvest,
			sdk.NewAttribute(types.AttributeValueCreator, msg.Sender),
			sdk.NewAttribute(types.AttributeValuePoolName, msg.PoolName),
			sdk.NewAttribute(types.AttributeValueReward, reward.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})
	return &types.MsgHarvestResponse{Reward: reward}, nil
}
