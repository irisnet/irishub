package keeper

import (
	"context"
	"encoding/hex"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/random/types"
)

var _ types.QueryServer = Keeper{}

// Random implements the Query/Random gRPC method
func (k Keeper) Random(c context.Context, req *types.QueryRandomRequest) (*types.QueryRandomResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	reqID, err := hex.DecodeString(req.ReqId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request id")
	}

	ctx := sdk.UnwrapSDKContext(c)

	random, err := k.GetRandom(ctx, reqID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "random %s not found", req.ReqId)
	}

	return &types.QueryRandomResponse{Random: &random}, nil
}

// RandomRequestQueue implements the Query/RandomRequestQueue gRPC method
func (k Keeper) RandomRequestQueue(c context.Context, req *types.QueryRandomRequestQueueRequest) (*types.QueryRandomRequestQueueResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.Height < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid height")
	}

	ctx := sdk.UnwrapSDKContext(c)

	var requests []types.Request
	if req.Height == 0 {
		// query all pending requests
		requests = queryAllRandomRequestsInQueue(ctx, k)
	} else {
		// query the pending requests by the specified height
		requests = queryRandomRequestQueueByHeight(ctx, req.Height, k)
	}

	return &types.QueryRandomRequestQueueResponse{Requests: requests}, nil
}
