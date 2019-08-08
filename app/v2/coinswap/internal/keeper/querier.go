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
	var params types.QueryLiquidityParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	denom, err := sdk.GetCoinMinDenom(params.TokenId)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("illegal token id", err.Error()))
	}

	reservePoolName, err := k.GetReservePoolName(sdk.IrisAtto, denom)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not retrieve reserve pool name", err.Error()))
	}
	reservePool, _ := k.GetReservePool(ctx, reservePoolName)
	// clean reserve pool to remove non-pool coins
	cleanedReservePool, err := k.CleanReservePool(reservePool, reservePoolName)
	bz, err := k.cdc.MarshalJSONIndent(cleanedReservePool, "", " ")
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
