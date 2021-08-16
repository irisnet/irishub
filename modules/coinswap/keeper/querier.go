package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/coinswap/types"
)

// NewQuerier creates a querier for coinswap REST endpoints
func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryLiquidity:
			return queryLiquidity(ctx, req, k, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

// queryLiquidity returns the total liquidity available for the provided denomination
// upon success or an error if the query fails.
func queryLiquidity(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryLiquidityParams

	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	pool, exists := k.GetPoolByLptDenom(ctx, params.Denom)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrReservePoolNotExists, "liquidity pool token: %s", params.Denom)
	}

	standardDenom := k.GetStandardDenom(ctx)
	reservePool, err := k.GetPoolBalancesByLptDenom(ctx, pool.LptDenom)
	if err != nil {
		return nil, err
	}

	standard := sdk.NewCoin(standardDenom, reservePool.AmountOf(standardDenom))
	token := sdk.NewCoin(pool.CounterpartyDenom, reservePool.AmountOf(pool.CounterpartyDenom))
	liquidity := sdk.NewCoin(pool.LptDenom, k.bk.GetSupply(ctx).GetTotal().AmountOf(pool.LptDenom))

	swapParams := k.GetParams(ctx)
	fee := swapParams.Fee.String()
	res := types.QueryLiquidityResponse{
		Standard:  standard,
		Token:     token,
		Liquidity: liquidity,
		Fee:       fee,
	}

	bz, errRes := codec.MarshalJSONIndent(legacyQuerierCdc, res)
	if errRes != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}
