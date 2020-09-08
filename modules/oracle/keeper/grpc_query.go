package keeper

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/oracle/types"
)

var _ types.QueryServer = Keeper{}

// Feed queries a feed by feed name
func (k Keeper) Feed(c context.Context, req *types.QueryFeedRequest) (*types.QueryFeedResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	feed, found := k.GetFeed(ctx, req.FeedName)
	if !found {
		return nil, status.Errorf(codes.NotFound, "feed %s not found", req.FeedName)
	}
	feedCtx := BuildFeedContext(ctx, k, feed)
	return &types.QueryFeedResponse{Feed: feedCtx}, nil
}

func (k Keeper) Feeds(c context.Context, req *types.QueryFeedsRequest) (*types.QueryFeedsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	state := strings.TrimSpace(req.State)
	var result types.FeedsContext
	if len(state) == 0 {
		k.IteratorFeeds(ctx, func(feed types.Feed) {
			result = append(result, BuildFeedContext(ctx, k, feed))
		})
	} else {
		state, err := types.RequestContextStateFromString(req.State)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid request state")

		}
		k.IteratorFeedsByState(ctx, state, func(feed types.Feed) {
			result = append(result, BuildFeedContext(ctx, k, feed))
		})
	}
	return &types.QueryFeedsResponse{Feeds: result}, nil
}

func (k Keeper) FeedValue(c context.Context, req *types.QueryFeedValueRequest) (*types.QueryFeedValueResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	result := k.GetFeedValues(ctx, req.FeedName)
	return &types.QueryFeedValueResponse{FeedValues: result}, nil
}
