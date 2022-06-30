package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/irisnet/irishub/modules/guardian/types"
)

var _ types.QueryServer = Keeper{}

// Supers implements the Query/Supers gRPC method
func (k Keeper) Supers(c context.Context, req *types.QuerySupersRequest) (*types.QuerySupersResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	var supers []types.Super
	store := ctx.KVStore(k.storeKey)

	pageRes, err := query.Paginate(store, req.Pagination, func(key []byte, value []byte) error {
		var super types.Super
		k.cdc.MustUnmarshal(value, &super)
		supers = append(supers, super)
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}

	return &types.QuerySupersResponse{Supers: supers, Pagination: pageRes}, nil
}
