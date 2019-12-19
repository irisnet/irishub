package keeper_test

import (
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/irisnet/irishub/modules/mint/internal/keeper"
	"github.com/irisnet/irishub/modules/mint/internal/types"
)

func (suite *KeeperTestSuite) TestNewQuerier() {
	querier := keeper.NewQuerier(suite.app.MintKeeper)

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	_, err := querier(suite.ctx, []string{types.QueryParameters}, query)
	require.NoError(suite.T(), err)
}

func (suite *KeeperTestSuite) TestQueryParams() {
	querier := keeper.NewQuerier(suite.app.MintKeeper)

	var params types.Params

	res, sdkErr := querier(suite.ctx, []string{types.QueryParameters}, abci.RequestQuery{})
	require.NoError(suite.T(), sdkErr)

	err := suite.app.Codec().UnmarshalJSON(res, &params)
	require.NoError(suite.T(), err)

	require.Equal(suite.T(), suite.app.MintKeeper.GetParamSet(suite.ctx), params)
}
