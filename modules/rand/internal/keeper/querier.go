package keeper

import (
	"encoding/hex"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irishub/modules/rand/internal/types"
)

// NewQuerier creates a new rand Querier instance
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryRand:
			return queryRand(ctx, req, k)
		case types.QueryRandRequestQueue:
			return queryRandRequestQueue(ctx, req, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryRand(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryRandParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	reqID, err := hex.DecodeString(params.ReqID)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidReqID, params.ReqID)
	}

	rand, err2 := keeper.GetRand(ctx, reqID)
	if err2 != nil {
		return nil, err2
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, rand)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryRandRequestQueue(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryRandRequestQueueParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	if params.Height < 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidHeight, string(params.Height))
	}

	var requests []types.Request

	if params.Height == 0 {
		// query all pending requests
		requests = queryAllRandRequestsInQueue(ctx, keeper)
	} else {
		// query the pending requests by the specified height
		requests = queryRandRequestQueueByHeight(ctx, params.Height, keeper)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, requests)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryRandRequestQueueByHeight(ctx sdk.Context, height int64, keeper Keeper) []types.Request {
	var requests = make([]types.Request, 0)

	iterator := keeper.IterateRandRequestQueueByHeight(ctx, height)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var request types.Request
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &request)

		requests = append(requests, request)
	}

	return requests
}

func queryAllRandRequestsInQueue(ctx sdk.Context, keeper Keeper) []types.Request {
	var requests = make([]types.Request, 0)

	keeper.IterateRandRequestQueue(ctx, func(h int64, r types.Request) (stop bool) {
		requests = append(requests, r)
		return false
	})

	return requests
}
