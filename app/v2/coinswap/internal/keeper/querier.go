package keeper

import (
	"fmt"

	"github.com/irisnet/irishub/app/v2/coinswap/internal/types"
	sdk "github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier creates a querier for coinswap REST endpoints
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryLiquidity:
			return queryLiquidity(ctx, req, k)

		case types.QueryParameters:
			return queryParameters(ctx, path[1:], req, k)

		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("%s is not a valid query request path", req.Path))
		}
	}
}

// queryLiquidity returns the total liquidity available for the provided denomination
// upon success or an error if the query fails.
func queryLiquidity(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var denom string
	err := k.cdc.UnmarshalJSON(req.Data, &denom)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	moduleName, err := k.GetModuleName(sdk.IrisAtto, denom)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not retrieve module name", err.Error()))
	}
	reservePool, found := k.GetReservePool(ctx, moduleName)
	if !found {
		return nil, sdk.ErrInternal("reserve pool does not exist")
	}

	bz, err := k.cdc.MarshalJSONIndent(reservePool.AmountOf(denom), "", " ")
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

// queryParameters returns coinswap module parameter queried for upon success
// or an error if the query fails
func queryParameters(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	switch path[0] {
	case types.ParamFee:
		bz, err := k.cdc.MarshalJSONIndent(k.GetFeeParam(ctx), "", " ")
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
		}
		return bz, nil
	case types.ParamNativeDenom:
		bz, err := k.cdc.MarshalJSONIndent(sdk.IrisAtto, "", " ")
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
		}
		return bz, nil
	default:
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("%s is not a valid query request path", req.Path))
	}
}
