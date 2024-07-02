package keeper

import (
	"context"
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"mods.irisnet.org/modules/random/types"
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

func queryRandomRequestQueueByHeight(ctx sdk.Context, height int64, k Keeper) []types.Request {
	requests := make([]types.Request, 0)

	iterator := k.IterateRandomRequestQueueByHeight(ctx, height)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var request types.Request
		k.cdc.MustUnmarshal(iterator.Value(), &request)

		requests = append(requests, request)
	}

	return requests
}

func queryAllRandomRequestsInQueue(ctx sdk.Context, k Keeper) []types.Request {
	requests := make([]types.Request, 0)

	k.IterateRandomRequestQueue(ctx, func(h int64, reqID []byte, r types.Request) (stop bool) {
		requests = append(requests, r)
		return false
	})

	return requests
}
