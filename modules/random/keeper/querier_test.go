package keeper_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/random/keeper"
	"github.com/irisnet/irishub/modules/random/types"
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

	// init random
	random := types.NewRandom(
		types.SHA256(testTxBytes),
		testHeight,
		big.NewRat(testRandomNumerator, testRandomDenomiator).FloatString(types.RandPrec),
	)
	suite.keeper.SetRandom(suite.ctx, testReqID, random)

	storedRandom, err := suite.keeper.GetRandom(suite.ctx, testReqID)
	suite.NoError(err)

	// test queryRandom

	bz, errRes := suite.cdc.MarshalJSON(types.QueryRandomParams{ReqID: hex.EncodeToString(testReqID)})
	suite.NoError(errRes)

	req.Path = fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRandom)
	req.Data = bz

	res, err = querier(suite.ctx, []string{types.QueryRandom}, req)
	suite.NoError(err)

	var resultRandom types.Random
	errRes = suite.cdc.UnmarshalJSON(res, &resultRandom)
	suite.NoError(errRes)
	suite.Equal(storedRandom, resultRandom)

	// test queryRandomRequestQueue
	request, err := suite.keeper.RequestRandom(suite.ctx, testConsumer, testBlockInterval, false, sdk.NewCoins())

	bz, errRes = suite.cdc.MarshalJSON(types.QueryRandomRequestQueueParams{Height: int64(testBlockInterval)})
	suite.NoError(errRes)

	req.Path = fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRandomRequestQueue)
	req.Data = bz

	res, err = querier(suite.ctx, []string{types.QueryRandomRequestQueue}, req)
	suite.NoError(err)

	var resultRequests []types.Request
	errRes = suite.cdc.UnmarshalJSON(res, &resultRequests)
	suite.NoError(errRes)
	suite.Equal([]types.Request{request}, resultRequests)
}
