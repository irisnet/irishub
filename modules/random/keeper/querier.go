package keeper

import (
	"encoding/hex"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irishub/modules/random/types"
)

// NewQuerier creates a new random Querier instance
func NewQuerier(k Keeper, legacyQuerierCdc codec.JSONMarshaler) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryRandom:
			return queryRandom(ctx, req, k, legacyQuerierCdc)
		case types.QueryRandomRequestQueue:
			return queryRandomRequestQueue(ctx, req, k, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryRandom(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc codec.JSONMarshaler) ([]byte, error) {
	var params types.QueryRandomParams
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	reqID, err := hex.DecodeString(params.ReqID)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidReqID, params.ReqID)
	}

	random, err2 := k.GetRandom(ctx, reqID)
	if err2 != nil {
		return nil, err2
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, random)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryRandomRequestQueue(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc codec.JSONMarshaler) ([]byte, error) {
	var params types.QueryRandomRequestQueueParams
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	if params.Height < 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidHeight, string(rune(params.Height)))
	}

	var requests []types.Request
	if params.Height == 0 {
		// query all pending requests
		requests = queryAllRandomRequestsInQueue(ctx, k)
	} else {
		// query the pending requests by the specified height
		requests = queryRandomRequestQueueByHeight(ctx, params.Height, k)
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, requests)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryRandomRequestQueueByHeight(ctx sdk.Context, height int64, k Keeper) []types.Request {
	var requests = make([]types.Request, 0)

	iterator := k.IterateRandomRequestQueueByHeight(ctx, height)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var request types.Request
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &request)

		requests = append(requests, request)
	}

	return requests
}

func queryAllRandomRequestsInQueue(ctx sdk.Context, k Keeper) []types.Request {
	var requests = make([]types.Request, 0)

	k.IterateRandomRequestQueue(ctx, func(h int64, r types.Request) (stop bool) {
		requests = append(requests, r)
		return false
	})

	return requests
}
