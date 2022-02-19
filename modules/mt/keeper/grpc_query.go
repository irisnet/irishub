package keeper

import (
	"context"

	"github.com/irisnet/irismod/modules/mt/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Supply(context.Context, *types.QuerySupplyRequest) (*types.QuerySupplyResponse, error) {
	panic("implement me")
}

func (k Keeper) Denoms(context.Context, *types.QueryDenomsRequest) (*types.QueryDenomsResponse, error) {
	panic("implement me")
}

func (k Keeper) Denom(context.Context, *types.QueryDenomRequest) (*types.QueryDenomResponse, error) {
	panic("implement me")
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
