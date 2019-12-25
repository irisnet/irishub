package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/irisnet/irishub/modules/guardian/internal/keeper"
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

	// init xxx

	// test queryProfilers

	// test queryProfilers

}
