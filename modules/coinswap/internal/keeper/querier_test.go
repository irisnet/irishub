package keeper_test

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/irisnet/irishub/modules/coinswap/internal/keeper"
	"github.com/irisnet/irishub/modules/coinswap/internal/types"
)

func (suite *KeeperTestSuite) TestNewQuerier() {

	req := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	querier := keeper.NewQuerier(suite.app.CoinswapKeeper)

	// query with incorrect path
	res, err := querier(suite.ctx, []string{"other"}, req)
	suite.Error(err)
	suite.Nil(res)

	// query for non existent reserve pool should return an error
	req.Path = fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLiquidity)
	req.Data = suite.cdc.MustMarshalJSON("btc")
	res, err = querier(suite.ctx, []string{"liquidity"}, req)
	suite.Error(err)
	suite.Nil(res)
}
