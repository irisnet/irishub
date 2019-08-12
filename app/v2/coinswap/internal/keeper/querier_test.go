package keeper

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/irisnet/irishub/app/v2/coinswap/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

func TestNewQuerier(t *testing.T) {
	ctx, keeper, _ := createTestInput(t, sdk.NewInt(100), 2)

	req := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	querier := NewQuerier(keeper)

	// query with incorrect path
	res, err := querier(ctx, []string{"other"}, req)
	require.Error(t, err)
	require.Nil(t, res)

	// query for non existent reserve pool should return an error
	req.Path = fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLiquidity)
	req.Data = keeper.cdc.MustMarshalJSON("btc")
	res, err = querier(ctx, []string{"liquidity"}, req)
	require.Error(t, err)
	require.Nil(t, res)

	// query for fee params
	fee := types.DefaultParams().Fee
	req.Path = fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryParameters, types.ParamFee)
	req.Data = []byte{}
	res, err = querier(ctx, []string{types.QueryParameters, types.ParamFee}, req)
	keeper.cdc.UnmarshalJSON(res, &fee)
	require.Nil(t, err)
	require.Equal(t, fee, types.DefaultParams().Fee)
}

// TODO: Add tests for valid liquidity queries
