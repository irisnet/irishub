package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/irisnet/irismod/modules/coinswap/types"
)

var _ types.QueryServer = Keeper{}

// LiquidityPool returns the liquidity pool information of the denom
func (k Keeper) LiquidityPool(c context.Context, req *types.QueryLiquidityPoolRequest) (*types.QueryLiquidityPoolResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	pool, exists := k.GetPoolByLptDenom(ctx, req.LptDenom)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrReservePoolNotExists, "liquidity pool token: %s", req.LptDenom)
	}

	balances, err := k.GetPoolBalancesByLptDenom(ctx, pool.LptDenom)
	if err != nil {
		return nil, err
	}

	supply := k.bk.GetSupply(ctx)
	standard := sdk.NewCoin(pool.StandardDenom, balances.AmountOf(pool.StandardDenom))
	token := sdk.NewCoin(pool.CounterpartyDenom, balances.AmountOf(pool.CounterpartyDenom))
	liquidity := sdk.NewCoin(pool.LptDenom, supply.GetTotal().AmountOf(pool.LptDenom))
	params := k.GetParams(ctx)
	res := types.QueryLiquidityPoolResponse{
		Pool: types.PoolInfo{
			Id:            pool.Id,
			EscrowAddress: pool.EscrowAddress,
			Standard:      standard,
			Token:         token,
			Lpt:           liquidity,
			Fee:           params.Fee.String(),
		},
	}
	return &res, nil
}

func (k Keeper) LiquidityPools(c context.Context, req *types.QueryLiquidityPoolsRequest) (*types.QueryLiquidityPoolsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	supply := k.bk.GetSupply(ctx).GetTotal()
	params := k.GetParams(ctx)

	var pools []types.PoolInfo

	store := ctx.KVStore(k.storeKey)
	nftStore := prefix.NewStore(store, []byte(types.KeyPool))
	pageRes, err := query.Paginate(nftStore, req.Pagination, func(key []byte, value []byte) error {
		var pool types.Pool
		k.cdc.MustUnmarshalBinaryBare(value, &pool)

		balances, err := k.GetPoolBalancesByLptDenom(ctx, pool.LptDenom)
		if err != nil {
			return err
		}

		pools = append(pools, types.PoolInfo{
			Id:            pool.Id,
			EscrowAddress: pool.EscrowAddress,
			Standard:      sdk.NewCoin(pool.StandardDenom, balances.AmountOf(pool.StandardDenom)),
			Token:         sdk.NewCoin(pool.CounterpartyDenom, balances.AmountOf(pool.CounterpartyDenom)),
			Lpt:           sdk.NewCoin(pool.LptDenom, supply.AmountOf(pool.LptDenom)),
			Fee:           params.Fee.String(),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &types.QueryLiquidityPoolsResponse{
		Pagination: pageRes,
		Pools:      pools,
	}, nil
}
