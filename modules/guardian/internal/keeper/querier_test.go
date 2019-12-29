package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/irisnet/irishub/modules/guardian/internal/keeper"
	"github.com/irisnet/irishub/modules/guardian/internal/types"
)

func TestQuerierSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestNewQuerier() {
	req := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	querier := keeper.NewQuerier(suite.keeper)
	res, err := querier(suite.ctx, []string{"other"}, req)
	suite.Error(err)
	suite.Nil(res)

	// test queryProfilers

	res, err = querier(suite.ctx, []string{types.QueryProfilers}, abci.RequestQuery{})
	suite.NoError(err)
	var guardianProfilers []types.Guardian
	e := suite.cdc.UnmarshalJSON(res, &guardianProfilers)
	suite.NoError(e)

	for i, val := range guardianProfilers {
		equal := val.Equal(types.DefaultGenesisState().Profilers[i])
		suite.True(equal)
	}

	// test queryTrustees

	res, err = querier(suite.ctx, []string{types.QueryTrustees}, abci.RequestQuery{})
	suite.NoError(err)
	var guardianTrustees []types.Guardian
	e = suite.cdc.UnmarshalJSON(res, &guardianTrustees)
	suite.NoError(e)

	for i, val := range guardianTrustees {
		equal := val.Equal(types.DefaultGenesisState().Trustees[i])
		suite.True(equal)
	}
}
