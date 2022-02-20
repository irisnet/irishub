package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/irisnet/irismod/modules/mt/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Supply(context.Context, *types.QuerySupplyRequest) (*types.QuerySupplyResponse, error) {
	panic("implement me")
}

func (k Keeper) Denoms(c context.Context, request *types.QueryDenomsRequest) (*types.QueryDenomsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	var denoms []types.Denom
	store := ctx.KVStore(k.storeKey)
	denomStore := prefix.NewStore(store, types.KeyDenomID(""))
	pageRes, err := query.Paginate(denomStore, request.Pagination, func(key []byte, value []byte) error {
		var denom types.Denom
		k.cdc.MustUnmarshal(value, &denom)
		denoms = append(denoms, denom)
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}

	return &types.QueryDenomsResponse{
		Denoms:     denoms,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) Denom(c context.Context, request *types.QueryDenomRequest) (*types.QueryDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	denom, found := k.GetDenom(ctx, request.DenomId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrInvalidDenom, "Denom not found: %s", request.DenomId)
	}

	return &types.QueryDenomResponse{Denom: &denom}, nil
}

func (k Keeper) MTSupply(context.Context, *types.QueryMTSupplyRequest) (*types.QueryMTSupplyResponse, error) {
	panic("implement me")
}

func (k Keeper) MTs(context.Context, *types.QueryMTsRequest) (*types.QueryMTsResponse, error) {
	panic("implement me")
}

func (k Keeper) MT(context.Context, *types.QueryMTRequest) (*types.QueryMTResponse, error) {
	panic("implement me")
}

func (k Keeper) Balances(context.Context, *types.QueryBalancesRequest) (*types.QueryBalancesResponse, error) {
	panic("implement me")
}
