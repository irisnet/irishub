package keeper

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/irisnet/irishub/modules/oracle/types"
)

// NewQuerier creates a querier for the oracle module
func NewQuerier(k Keeper, legacyQuerierCdc codec.JSONMarshaler) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryFeed:
			return queryFeed(ctx, req, k, legacyQuerierCdc)
		case types.QueryFeeds:
			return queryFeeds(ctx, req, k, legacyQuerierCdc)
		case types.QueryFeedValue:
			return queryFeedValue(ctx, req, k, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryFeed(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc codec.JSONMarshaler) ([]byte, error) {
	var params types.QueryFeedParams
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	feed, found := k.GetFeed(ctx, params.FeedName)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrUnknownFeedName, params.FeedName)
	}
	feedCtx := BuildFeedContext(ctx, k, feed)

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, feedCtx)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryFeeds(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc codec.JSONMarshaler) ([]byte, error) {
	var params types.QueryFeedsParams
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	state := strings.TrimSpace(params.State)
	var result types.FeedsContext
	if len(state) == 0 {
		k.IteratorFeeds(ctx, func(feed types.Feed) {
			result = append(result, BuildFeedContext(ctx, k, feed))
		})
	} else {
		state, err := types.RequestContextStateFromString(params.State)
		if err != nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, err.Error())
		}
		k.IteratorFeedsByState(ctx, state, func(feed types.Feed) {
			result = append(result, BuildFeedContext(ctx, k, feed))
		})
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, result)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryFeedValue(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc codec.JSONMarshaler) ([]byte, error) {
	var params types.QueryFeedValueParams
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	result := k.GetFeedValues(ctx, params.FeedName)
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, result)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func BuildFeedContext(ctx sdk.Context, k Keeper, feed types.Feed) (feedCtx types.FeedContext) {
	reqCtx, found := k.sk.GetRequestContext(ctx, feed.RequestContextID)
	if found {
		feedCtx.Providers = reqCtx.Providers
		feedCtx.ResponseThreshold = reqCtx.ResponseThreshold
		feedCtx.ServiceName = reqCtx.ServiceName
		feedCtx.Input = reqCtx.Input
		feedCtx.RepeatedFrequency = reqCtx.RepeatedFrequency
		feedCtx.ServiceFeeCap = reqCtx.ServiceFeeCap
		feedCtx.Timeout = reqCtx.Timeout
		feedCtx.State = reqCtx.State
	}
	feedCtx.Feed = &feed
	return feedCtx
}
