package keeper

import (
	"context"
	"time"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"mods.irisnet.org/modules/coinswap/types"
)

type msgServer struct {
	k Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the coinswap MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{k: keeper}
}

func (m msgServer) AddLiquidity(
	goCtx context.Context,
	msg *types.MsgAddLiquidity,
) (*types.MsgAddLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// check that deadline has not passed
	if ctx.BlockHeader().Time.After(time.Unix(msg.Deadline, 0)) {
		return nil, errorsmod.Wrap(
			types.ErrInvalidDeadline,
			"deadline has passed for MsgAddLiquidity",
		)
	}

	mintToken, err := m.k.AddLiquidity(ctx, msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgAddLiquidityResponse{
		MintToken: &mintToken,
	}, nil
}

func (m msgServer) AddUnilateralLiquidity(
	goCtx context.Context,
	msg *types.MsgAddUnilateralLiquidity,
) (*types.MsgAddUnilateralLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// check that deadline has not passed
	if ctx.BlockHeader().Time.After(time.Unix(msg.Deadline, 0)) {
		return nil, errorsmod.Wrap(
			types.ErrInvalidDeadline,
			"deadline has passed for MsgAddUnilateralLiquidity",
		)
	}

	mintToken, err := m.k.AddUnilateralLiquidity(ctx, msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgAddUnilateralLiquidityResponse{
		MintToken: &mintToken,
	}, nil
}

func (m msgServer) RemoveLiquidity(
	goCtx context.Context,
	msg *types.MsgRemoveLiquidity,
) (*types.MsgRemoveLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// check that deadline has not passed
	if ctx.BlockHeader().Time.After(time.Unix(msg.Deadline, 0)) {
		return nil, errorsmod.Wrap(
			types.ErrInvalidDeadline,
			"deadline has passed for MsgRemoveLiquidity",
		)
	}
	withdrawCoins, err := m.k.RemoveLiquidity(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgRemoveLiquidityResponse{
		WithdrawCoins: withdrawCoins,
	}, nil
}

func (m msgServer) RemoveUnilateralLiquidity(
	goCtx context.Context,
	msg *types.MsgRemoveUnilateralLiquidity,
) (*types.MsgRemoveUnilateralLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// check that deadline has not passed
	if ctx.BlockHeader().Time.After(time.Unix(msg.Deadline, 0)) {
		return nil, errorsmod.Wrap(
			types.ErrInvalidDeadline,
			"deadline has passed for MsgRemoveLiquidity",
		)
	}
	withdrawCoins, err := m.k.RemoveUnilateralLiquidity(ctx, msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgRemoveUnilateralLiquidityResponse{
		WithdrawCoins: withdrawCoins,
	}, nil
}

func (m msgServer) SwapCoin(
	goCtx context.Context,
	msg *types.MsgSwapOrder,
) (*types.MsgSwapCoinResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// check that deadline has not passed
	if ctx.BlockHeader().Time.After(time.Unix(msg.Deadline, 0)) {
		return nil, errorsmod.Wrap(types.ErrInvalidDeadline, "deadline has passed for MsgSwapOrder")
	}

	if m.k.blockedAddrs[msg.Output.Address] {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrUnauthorized,
			"%s is not allowed to receive external funds",
			msg.Output.Address,
		)
	}

	if err := m.k.Swap(ctx, msg); err != nil {
		return nil, err
	}
	return &types.MsgSwapCoinResponse{}, nil
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
