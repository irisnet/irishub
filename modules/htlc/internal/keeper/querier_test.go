package keeper_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/irisnet/irishub/modules/htlc/internal/keeper"
	"github.com/irisnet/irishub/modules/htlc/internal/types"
)

func TestQuerierSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestNewQuerier() {
	req := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	querier := keeper.NewQuerier(suite.app.HTLCKeeper)
	res, err := querier(suite.ctx, []string{"other"}, req)
	suite.Error(err)
	suite.Nil(res)

	// init HTLC

	htlc := types.NewHTLC(
		addrSender,
		addrTo,
		receiverOnOtherChain,
		amount,
		types.HTLCSecret{},
		timestamp,
		expireHeight,
		stateOpen,
	)
	err = suite.app.HTLCKeeper.CreateHTLC(suite.ctx, htlc, hashLock)
	suite.NoError(err)

	// test queryHTLC

	bz, errRes := suite.cdc.MarshalJSON(types.QueryHTLCParams{HashLock: hashLock})
	suite.NoError(errRes)

	req.Path = fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryHTLC)
	req.Data = bz

	res, err = querier(suite.ctx, []string{types.QueryHTLC}, req)
	suite.NoError(err)

	var resultHTLC types.HTLC
	errRes = suite.cdc.UnmarshalJSON(res, &resultHTLC)
	suite.NoError(errRes)
	suite.Equal(htlc, resultHTLC)
}
