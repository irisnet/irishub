package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/coinswap/types"
)

var _ types.QueryServer = Keeper{}

// Liquidity returns the liquidity pool information of the denom
func (k Keeper) Liquidity(c context.Context, req *types.QueryLiquidityRequest) (*types.QueryLiquidityResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	pool, exists := k.GetPoolByLptDenom(ctx, req.Denom)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrReservePoolNotExists, "liquidity pool token: %s", req.Denom)
	}

	standardDenom := k.GetStandardDenom(ctx)
	reservePool, err := k.GetPoolBalancesByLptDenom(ctx, pool.LptDenom)
	if err != nil {
		return nil, err
	}

	standard := sdk.NewCoin(standardDenom, reservePool.AmountOf(standardDenom))
	token := sdk.NewCoin(pool.CounterpartyDenom, reservePool.AmountOf(pool.CounterpartyDenom))
	liquidity := sdk.NewCoin(pool.LptDenom, k.bk.GetSupply(ctx).GetTotal().AmountOf(pool.LptDenom))

	swapParams := k.GetParams(ctx)
	fee := swapParams.Fee.String()
	res := types.QueryLiquidityResponse{
		Standard:  standard,
		Token:     token,
		Liquidity: liquidity,
		Fee:       fee,
	}
	return &res, nil
}
