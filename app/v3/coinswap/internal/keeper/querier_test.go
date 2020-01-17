package keeper

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/irisnet/irishub/app/v3/coinswap/internal/types"
)

func TestNewQuerier(t *testing.T) {
	app := createTestApp(nil, 2)

	req := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	querier := NewQuerier(app.csk)

	// query with incorrect path
	res, err := querier(app.ctx, []string{"other"}, req)
	require.Error(t, err)
	require.Nil(t, res)

	// query for non existent reserve pool should return an error
	req.Path = fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLiquidity)
	req.Data = app.csk.cdc.MustMarshalJSON("btc")
	res, err = querier(app.ctx, []string{"liquidity"}, req)
	require.Error(t, err)
	require.Nil(t, res)
}
