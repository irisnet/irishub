package keeper

import (
	"context"
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

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

func (m msgServer) CreatePool(goCtx context.Context, msg *types.MsgCreatePool) (*types.MsgCreatePoolResponse, error) {
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
			ctx.BlockHeight(), msg.StartHeight,
		)
	}

	if maxRewardCategories := m.Keeper.MaxRewardCategories(ctx); uint32(len(msg.TotalReward)) > maxRewardCategories {
		return nil, sdkerrors.Wrapf(
			types.ErrInvalidRewardRule,
			"the max reward category num is [%d], but got [%d]",
			maxRewardCategories, len(msg.TotalReward),
		)
	}

	//check valid lp token denom
	if err := m.Keeper.validateLPToken(ctx, msg.LptDenom); err != nil {
		return nil, sdkerrors.Wrapf(
			types.ErrInvalidLPToken,
			"The lp token denom[%s] is not exist",
			msg.LptDenom,
		)
	}
	pool, err := m.Keeper.CreatePool(
		ctx,
		msg.Description,
		msg.LptDenom,
		msg.StartHeight,
		msg.RewardPerBlock.Sort(),
		msg.TotalReward.Sort(),
		msg.Editable,
		creator,
	)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreatePool,
			sdk.NewAttribute(types.AttributeValueCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeValuePoolId, pool.Id),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	})
	return &types.MsgCreatePoolResponse{}, nil
}

func (m msgServer) CreatePoolWithCommunityPool(goCtx context.Context,
	msg *types.MsgCreatePoolWithCommunityPool) (*types.MsgCreatePoolWithCommunityPoolResponse, error) {
	proposer, err := sdk.AccAddressFromBech32(msg.Proposer)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	totalReward := sdk.NewCoins(msg.Content.FundApplied...).Add(msg.Content.FundSelfBond...)
	maxRewardCategories := m.Keeper.MaxRewardCategories(ctx)
	if uint32(len(totalReward)) > maxRewardCategories {
		return nil, sdkerrors.Wrapf(
			types.ErrInvalidRewardRule,
			"the max reward category num is [%d], but got [%d]",
			maxRewardCategories, len(totalReward),
		)
	}

	//check valid lp token denom
	if err := m.Keeper.validateLPToken(ctx, msg.Content.LptDenom); err != nil {
		return nil, sdkerrors.Wrapf(
			types.ErrInvalidLPToken,
			"The lp token denom[%s] is not exist",
			msg.Content.LptDenom,
		)
	}

	//escrow FundSelfBond to EscrowCollector
	if err := m.bk.SendCoinsFromAccountToModule(ctx,
		proposer, types.EscrowCollector, msg.Content.FundSelfBond); err != nil {
		return nil, err
	}

	//escrow FundApplied to EscrowCollector
	if err := m.escrowFromFeePool(ctx, msg.Content.FundApplied); err != nil {
		return nil, err
	}

	data, err := codectypes.NewAnyWithValue(&msg.Content)
	if err != nil {
		return nil, err
	}

	msgs := []sdk.Msg{
		&govv1.MsgExecLegacyContent{
			Content:   data,
			Authority: m.gk.GetGovernanceAccount(ctx).GetAddress().String(),
		},
	}

	//create new proposal given a content
	proposal, err := m.gk.SubmitProposal(ctx, msgs, "")
	if err != nil {
		return nil, err
	}

	// adds a deposit of a specific depositor on a specific proposal
	_, err = m.gk.AddDeposit(ctx, proposal.Id, proposer, msg.InitialDeposit)
	if err != nil {
		return nil, err
	}

	// add a escrowInfo to the proposal
	m.SetEscrowInfo(ctx, types.EscrowInfo{
		Proposer:     msg.Proposer,
		FundApplied:  msg.Content.FundApplied,
		FundSelfBond: msg.Content.FundSelfBond,
		ProposalId:   proposal.Id,
	})
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreatePoolWithCommunityPool,
			sdk.NewAttribute(types.AttributeValueCreator, msg.Proposer),
			sdk.NewAttribute(types.AttributeValueProposal, fmt.Sprintf("%d", proposal.Id)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Proposer),
		),
	})
	return &types.MsgCreatePoolWithCommunityPoolResponse{}, nil
}

func (m msgServer) DestroyPool(goCtx context.Context, msg *types.MsgDestroyPool) (*types.MsgDestroyPoolResponse, error) {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	refundCoin, err := m.Keeper.DestroyPool(ctx, msg.PoolId, creator)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDestroyPool,
			sdk.NewAttribute(types.AttributeValueCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeValuePoolId, msg.PoolId),
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
	if err = m.Keeper.AdjustPool(
		ctx,
		msg.PoolId,
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
			sdk.NewAttribute(types.AttributeValuePoolId, msg.PoolId),
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
	reward, err := m.Keeper.Stake(ctx, msg.PoolId, msg.Amount, sender)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeStake,
			sdk.NewAttribute(types.AttributeValueCreator, msg.Sender),
			sdk.NewAttribute(types.AttributeValuePoolId, msg.PoolId),
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
	reward, err := m.Keeper.Unstake(ctx, msg.PoolId, msg.Amount, sender)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnstake,
			sdk.NewAttribute(types.AttributeValueCreator, msg.Sender),
			sdk.NewAttribute(types.AttributeValuePoolId, msg.PoolId),
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
	reward, err := m.Keeper.Harvest(ctx, msg.PoolId, sender)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeHarvest,
			sdk.NewAttribute(types.AttributeValueCreator, msg.Sender),
			sdk.NewAttribute(types.AttributeValuePoolId, msg.PoolId),
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
