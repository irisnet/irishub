package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/coinswap/types"
)

var _ types.QueryServer = Keeper{}

// Liquidity returns the liquidity pool information of the denom
func (k Keeper) Liquidity(c context.Context, req *types.QueryLiquidityRequest) (*types.QueryLiquidityResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	tokenDenom := req.Denom
	uniDenom := types.GetUniDenomFromDenom(tokenDenom)

	ctx := sdk.UnwrapSDKContext(c)
	reservePool, err := k.GetReservePool(ctx, uniDenom)
	if err != nil {
		return nil, err
	}

	standardDenom := k.GetStandardDenom(ctx)
	standard := sdk.NewCoin(standardDenom, reservePool.AmountOf(standardDenom))
	token := sdk.NewCoin(tokenDenom, reservePool.AmountOf(tokenDenom))
	liquidity := sdk.NewCoin(uniDenom, k.bk.GetSupply(ctx).GetTotal().AmountOf(uniDenom))

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
