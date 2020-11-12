package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/irisnet/irishub/modules/guardian/keeper"
	"github.com/irisnet/irishub/modules/guardian/types"
)

func TestQuerierSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestNewQuerier() {
	req := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	querier := keeper.NewQuerier(suite.keeper, suite.cdc)
	res, err := querier(suite.ctx, []string{"other"}, req)
	suite.Error(err)
	suite.Nil(res)

	// test querySupers
	res, err = querier(suite.ctx, []string{types.QuerySupers}, abci.RequestQuery{})
	suite.NoError(err)
	var supers []types.Super
	e := suite.cdc.UnmarshalJSON(res, &supers)
	suite.NoError(e)

	for i, val := range supers {
		equal := val.Equal(types.DefaultGenesisState().Supers[i])
		suite.True(equal)
	}
}
