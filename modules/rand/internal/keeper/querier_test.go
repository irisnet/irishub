package keeper_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/irisnet/irishub/modules/rand"
	"github.com/irisnet/irishub/modules/rand/internal/keeper"
	"github.com/irisnet/irishub/modules/rand/internal/types"
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

	// init rand

	rand := rand.NewRand(rand.SHA256(testTxBytes), testHeight, big.NewRat(testRandNumerator, testRandDenomiator).FloatString(rand.RandPrec))
	suite.keeper.SetRand(suite.ctx, testReqID, rand)

	storedRand, err := suite.keeper.GetRand(suite.ctx, testReqID)
	suite.NoError(err)

	// test queryRand

	bz, errRes := suite.cdc.MarshalJSON(types.QueryRandParams{ReqID: hex.EncodeToString(testReqID)})
	suite.NoError(errRes)

	req.Path = fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRand)
	req.Data = bz

	res, err = querier(suite.ctx, []string{types.QueryRand}, req)
	suite.NoError(err)

	var resultRand types.Rand
	errRes = suite.cdc.UnmarshalJSON(res, &resultRand)
	suite.NoError(errRes)
	suite.Equal(storedRand, resultRand)

	// test queryRandRequestQueue

	request, err := suite.keeper.RequestRand(suite.ctx, testConsumer, testBlockInterval)

	bz, errRes = suite.cdc.MarshalJSON(types.QueryRandRequestQueueParams{Height: int64(testBlockInterval)})
	suite.NoError(errRes)

	req.Path = fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRandRequestQueue)
	req.Data = bz

	res, err = querier(suite.ctx, []string{types.QueryRandRequestQueue}, req)
	suite.NoError(err)

	var resultRequests []types.Request
	errRes = suite.cdc.UnmarshalJSON(res, &resultRequests)
	suite.NoError(errRes)
	suite.Equal([]types.Request{request}, resultRequests)
}
