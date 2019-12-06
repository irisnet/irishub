package keeper

import (
	"fmt"
	"github.com/irisnet/irishub/config"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/coinswap/internal/types"
)

// NewQuerier creates a querier for coinswap REST endpoints
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryLiquidity:
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

	uniDenom, err := types.GetUniDenom(params.Id)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(err.Error())
	}

	tokenDenom, err := types.GetCoinMinDenomFromUniDenom(uniDenom)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(err.Error())
	}

	reservePool := k.GetReservePool(ctx, params.Id)
	// all liquidity vouchers in module account
	liquidities := k.sk.GetModuleAccount(ctx, types.ModuleName).GetCoins()

	iris := sdk.NewCoin(config.IrisAtto, reservePool.AmountOf(config.IrisAtto))
	token := sdk.NewCoin(tokenDenom, reservePool.AmountOf(tokenDenom))
	liquidity := sdk.NewCoin(uniDenom, liquidities.AmountOf(uniDenom))

	swapParams := k.GetParams(ctx)
	fee := swapParams.Fee.String()
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
