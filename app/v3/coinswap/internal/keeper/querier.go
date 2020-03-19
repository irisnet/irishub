package keeper

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/irisnet/irishub/app/v3/coinswap/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

// NewQuerier creates a querier for coinswap REST endpoints
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryLiquidity:
			return queryLiquidity(ctx, req, k)

		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("%s is not a valid coinswap query request path", req.Path))
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

	voucherDenom, err := types.GetVoucherDenom(params.VoucherCoinName)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(err.Error())
	}

	tokenDenom, err := types.GetUnderlyingDenom(voucherDenom)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(err.Error())
	}

	pool, existed := k.GetPool(ctx, params.VoucherCoinName)
	if !existed {
		return nil, types.ErrReservePoolNotExists(fmt.Sprintf("liquidity pool for %s not found", params.VoucherCoinName))
	}

	iris := sdk.NewCoin(sdk.IrisAtto, pool.BalanceOf(sdk.IrisAtto))
	token := sdk.NewCoin(tokenDenom, pool.BalanceOf(tokenDenom))
	liquidity := sdk.NewCoin(voucherDenom, pool.BalanceOf(voucherDenom))

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
