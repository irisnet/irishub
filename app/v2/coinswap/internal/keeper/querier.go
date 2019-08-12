package keeper

import (
	"fmt"
	"strings"

	"github.com/irisnet/irishub/app/v2/coinswap/internal/types"
	sdk "github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier creates a querier for coinswap REST endpoints
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryLiquidities:
			return queryLiquidity(ctx, req, k)

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

	if len(strings.TrimSpace(params.TokenId)) == 0 {
		return nil, sdk.ErrUnknownRequest("token id can not be empty")
	}

	denom, err := sdk.GetCoinMinDenom(params.TokenId)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("illegal token id", err.Error()))
	}

	reservePoolName, err := types.GetReservePoolName(sdk.IrisAtto, denom)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not retrieve reserve pool name", err.Error()))
	}
	reservePool := k.GetReservePool(ctx, reservePoolName)

	iris := sdk.NewCoin(sdk.IrisAtto, reservePool.AmountOf(sdk.IrisAtto))
	token := sdk.NewCoin(denom, reservePool.AmountOf(denom))
	liquidity := sdk.NewCoin(reservePoolName, reservePool.AmountOf(reservePoolName))

	swapParams := k.GetParams(ctx)
	fee := swapParams.Fee.DecimalString(types.MaxFeePrecision)
	res := types.QueryLiquidityResponse{
		Iris:      iris,
		Token:     token,
		Liquidity: liquidity,
		Fee:       fee,
	}

	bz, err := k.cdc.MarshalJSONIndent(res, "", " ")
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
