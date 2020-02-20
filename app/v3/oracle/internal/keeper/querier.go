package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"
	"strings"

	"github.com/irisnet/irishub/app/v3/oracle/internal/types"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

// NewQuerier creates a querier for the oracle module
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryFeed:
			return queryFeed(ctx, req, k)
		case types.QueryFeeds:
			return queryFeeds(ctx, req, k)
		case types.QueryFeedValue:
			return queryFeedValue(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown oracle query endpoint")
		}
	}
}

func queryFeed(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryFeedParams
	if err := k.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	feed, found := k.GetFeed(ctx, params.FeedName)
	if !found {
		return nil, types.ErrUnknownFeedName(types.DefaultCodespace, params.FeedName)
	}
	feedCtx := buildFeedContext(ctx, k, feed)

	bz, err := codec.MarshalJSONIndent(k.cdc, feedCtx)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}
	return bz, nil
}

func queryFeeds(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryFeedsParams
	if err := k.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	state := strings.TrimSpace(params.State)
	var result types.FeedsContext
	if len(state) == 0 {
		k.IteratorFeeds(ctx, func(feed types.Feed) {
			result = append(result, buildFeedContext(ctx, k, feed))
		})
	} else {
		k.IteratorFeedsByState(ctx, types.StateFromString(params.State), func(feed types.Feed) {
			result = append(result, buildFeedContext(ctx, k, feed))
		})
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, result)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}
	return bz, nil
}

func queryFeedValue(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryFeedValueParams
	if err := k.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	result := k.GetFeedValues(ctx, params.FeedName)
	bz, err := codec.MarshalJSONIndent(k.cdc, result)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}
	return bz, nil
}

func buildFeedContext(ctx sdk.Context, k Keeper, feed types.Feed) (feedCtx types.FeedContext) {
	reqCtx, found := k.sk.GetRequestContext(ctx, feed.RequestContextID)
	if found {
		feedCtx.Providers = reqCtx.Providers
		feedCtx.ResponseThreshold = reqCtx.ResponseThreshold
		feedCtx.ServiceName = reqCtx.ServiceName
		feedCtx.Input = reqCtx.Input
		feedCtx.RepeatedFrequency = reqCtx.RepeatedFrequency
		feedCtx.RepeatedTotal = reqCtx.RepeatedTotal
		feedCtx.ServiceFeeCap = reqCtx.ServiceFeeCap
		feedCtx.Timeout = reqCtx.Timeout
		feedCtx.ResponseThreshold = reqCtx.ResponseThreshold
		feedCtx.State = reqCtx.State
	}
	feedCtx.Feed = feed
	return feedCtx
}
