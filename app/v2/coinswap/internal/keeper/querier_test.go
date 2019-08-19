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
}
