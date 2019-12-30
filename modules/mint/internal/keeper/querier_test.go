package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/irisnet/irishub/modules/mint/internal/keeper"
	"github.com/irisnet/irishub/modules/mint/internal/types"
)

func TestQuerierSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestNewQuerier() {
	querier := keeper.NewQuerier(suite.app.MintKeeper)
	res, err := querier(suite.ctx, []string{"other"}, abci.RequestQuery{})
	suite.Error(err)
	suite.Nil(res)

	// test queryParams

	res, err = querier(suite.ctx, []string{types.QueryParameters}, abci.RequestQuery{})
	suite.NoError(err)
	var params types.Params
	e := suite.cdc.UnmarshalJSON(res, &params)
	suite.NoError(e)
	suite.Equal(suite.app.MintKeeper.GetParamSet(suite.ctx), params)
}
