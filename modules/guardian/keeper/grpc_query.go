package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/guardian/types"
)

var _ types.QueryServer = Keeper{}

// Profilers implements the Query/Profilers gRPC method
func (k Keeper) Profilers(c context.Context, req *types.QueryProfilersRequest) (*types.QueryProfilersResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	var profilers []types.Guardian
	k.IterateProfilers(
		ctx,
		func(profiler types.Guardian) bool {
			profilers = append(profilers, profiler)
			return false
		},
	)

	return &types.QueryProfilersResponse{Profilers: profilers}, nil
}

// Trustees implements the Query/Trustees gRPC method
func (k Keeper) Trustees(c context.Context, req *types.QueryTrusteesRequest) (*types.QueryTrusteesResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	var trustees []types.Guardian
	k.IterateTrustees(
		ctx,
		func(trustee types.Guardian) bool {
			trustees = append(trustees, trustee)
			return false
		},
	)

	return &types.QueryTrusteesResponse{Trustees: trustees}, nil
}
