package rand

import (
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryRand             = "rand"
	QueryRands            = "rands"
	QueryRandRequest      = "request"
	QueryRandRequests     = "requests"
	QueryRandRequestQueue = "queue"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryRand:
			return queryRand(ctx, req, k)
		case QueryRands:
			return queryRands(ctx, req, k)
		case QueryRandRequest:
			return queryRandRequest(ctx, req, k)
		case QueryRandRequests:
			return queryRandRequests(ctx, req, k)
		case QueryRandRequestQueue:
			return queryRandRequestQueue(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown rand query endpoint")
		}
	}
}

// QueryRandParams is the query parameters for 'custom/rand/rand'
type QueryRandParams struct {
	ReqID string
}

func queryRand(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryRandParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	rand, err2 := keeper.GetRand(ctx, params.ReqID)
	if err2 != nil {
		return nil, err2
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, rand)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}

// QueryRandsParams is the query parameters for 'custom/rand/rands'
type QueryRandsParams struct {
	Consumer sdk.AccAddress
}

func queryRands(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryRandsParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	var rands []Rand

	op := func(r Rand) bool {
		if len(params.Consumer) == 0 {
			rands = append(rands, r)
		} else {
			// TODO: query the rands by the consumer
			rands = append(rands, r)
		}

		return false
	}

	keeper.IterateRands(ctx, op)

	bz, err := codec.MarshalJSONIndent(keeper.cdc, rands)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}

// QueryRandRequestParams is the query parameters for 'custom/rand/request'
type QueryRandRequestParams struct {
	ReqID string
}

func queryRandRequest(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryRandRequestParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	request, err2 := keeper.GetRandRequest(ctx, params.ReqID)
	if err2 != nil {
		return nil, err2
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, request)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}

// QueryRandRequestsParams is the query parameters for 'custom/rand/requests'
type QueryRandRequestsParams struct {
	Consumer sdk.AccAddress
}

func queryRandRequests(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryRandRequestsParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	var requests []Request

	op := func(r Request) bool {
		if len(params.Consumer) == 0 {
			requests = append(requests, r)
		} else {
			if r.Consumer.Equals(params.Consumer) {
				requests = append(requests, r)
			}
		}

		return false
	}

	keeper.IterateRandRequests(ctx, op)

	bz, err := codec.MarshalJSONIndent(keeper.cdc, requests)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}

// QueryRandRequestQueueParams is the query parameters for 'custom/rand/queue'
type QueryRandRequestQueueParams struct {
	Height int64
}

func queryRandRequestQueue(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryRandRequestQueueParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	var requests []Request

	if params.Height == 0 {
		// query all pending requests
		requests = queryAllRandRequestsInQueue(ctx, keeper)
	} else {
		// query the pending requests by the specified height
		requests = queryRandRequestQueueByHeight(ctx, params.Height, keeper)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, requests)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}

func queryRandRequestQueueByHeight(ctx sdk.Context, height int64, keeper Keeper) []Request {
	var requests = make([]Request, 0)

	iterator := keeper.IterateRandRequestQueueByHeight(ctx, height)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var reqID string
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &reqID)

		request, err := keeper.GetRandRequest(ctx, reqID)
		if err != nil {
			continue
		}

		requests = append(requests, request)
	}

	return requests
}

func queryAllRandRequestsInQueue(ctx sdk.Context, keeper Keeper) []Request {
	var requests = make([]Request, 0)

	keeper.IterateRandRequestQueue(ctx, func(r Request) (stop bool) {
		requests = append(requests, r)
		return false
	})

	return requests
}
