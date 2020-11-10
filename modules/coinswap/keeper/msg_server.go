package keeper

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/irisnet/irismod/modules/coinswap/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) AddLiquidity(goCtx context.Context, msg *types.MsgAddLiquidity) (*types.MsgAddLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// check that deadline has not passed
	if ctx.BlockHeader().Time.After(time.Unix(msg.Deadline, 0)) {
		return nil, sdkerrors.Wrap(types.ErrInvalidDeadline, "deadline has passed for MsgAddLiquidity")
	}

	mintToken, err := m.Keeper.AddLiquidity(ctx, msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	)

	return &types.MsgAddLiquidityResponse{
		MintToken: &mintToken,
	}, nil
}

func (m msgServer) RemoveLiquidity(goCtx context.Context, msg *types.MsgRemoveLiquidity) (*types.MsgRemoveLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// check that deadline has not passed
	if ctx.BlockHeader().Time.After(time.Unix(msg.Deadline, 0)) {
		return nil, sdkerrors.Wrap(types.ErrInvalidDeadline, "deadline has passed for MsgRemoveLiquidity")
	}
	withdrawCoins, err := m.Keeper.RemoveLiquidity(ctx, msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	)

	var coins = make([]*sdk.Coin, 0, withdrawCoins.Len())
	for _, coin := range withdrawCoins {
		coins = append(coins, &coin)
	}
	return &types.MsgRemoveLiquidityResponse{
		WithdrawCoins: coins,
	}, nil
}

func (m msgServer) SwapCoin(goCtx context.Context, msg *types.MsgSwapOrder) (*types.MsgSwapCoinResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// check that deadline has not passed
	if ctx.BlockHeader().Time.After(time.Unix(msg.Deadline, 0)) {
		return nil, sdkerrors.Wrap(types.ErrInvalidDeadline, "deadline has passed for MsgSwapOrder")
	}

	if err := m.Keeper.Swap(ctx, msg); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Input.Address),
		),
	)
	return &types.MsgSwapCoinResponse{}, nil
}
