package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irishub/modules/coinswap/internal/types"
)

// NewQuerier creates a querier for coinswap REST endpoints
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryLiquidity:
			return queryLiquidity(ctx, req, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

// queryLiquidity returns the total liquidity available for the provided denomination
// upon success or an error if the query fails.
func queryLiquidity(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryLiquidityParams
	standardDenom := k.GetParams(ctx).StandardDenom
	if err := k.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	if err := types.CheckUniDenom(params.ID); err != nil {
		return nil, err
	}

	uniDenom := params.ID

	tokenDenom, err := types.GetCoinDenomFromUniDenom(uniDenom)
	if err != nil {
		return nil, err
	}

	reservePool := k.GetReservePool(ctx, params.ID)
	// all liquidity vouchers in module account
	liquidities := k.sk.GetSupply(ctx).GetTotal()

	standard := sdk.NewCoin(standardDenom, reservePool.AmountOf(standardDenom))
	token := sdk.NewCoin(tokenDenom, reservePool.AmountOf(tokenDenom))
	liquidity := sdk.NewCoin(uniDenom, liquidities.AmountOf(uniDenom))

	swapParams := k.GetParams(ctx)
	fee := swapParams.Fee.String()
	res := types.QueryLiquidityResponse{
		Standard:  standard,
		Token:     token,
		Liquidity: liquidity,
		Fee:       fee,
	}

	bz, errRes := k.cdc.MarshalJSONIndent(res, "", " ")
	if errRes != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}
