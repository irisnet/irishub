package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

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
	k.IterateSupers(
		ctx,
		func(super types.Super) bool {
			supers = append(supers, super)
			return false
		},
	)

	return &types.QuerySupersResponse{Supers: supers}, nil
}
