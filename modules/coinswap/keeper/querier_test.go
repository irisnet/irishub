package keeper_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/coinswap/keeper"
	"github.com/irisnet/irismod/modules/coinswap/types"
)

func TestQuerierSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestNewQuerier() {
	req := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	legacyAmino := suite.app.LegacyAmino()

	querier := keeper.NewQuerier(suite.app.CoinswapKeeper, legacyAmino)
	res, err := querier(suite.ctx, []string{"other"}, req)
	suite.Error(err)
	suite.Nil(res)

	btcAmt, _ := sdk.NewIntFromString("100")
	standardAmt, _ := sdk.NewIntFromString("10000000000000000000")
	depositCoin := sdk.NewCoin(denomBTC, btcAmt)
	minReward := sdk.NewInt(1)
	deadline := time.Now().Add(1 * time.Minute)
	msg := types.NewMsgAddLiquidity(depositCoin, standardAmt, minReward, deadline.Unix(), addrSender1.String())
	lptCoin, _ := suite.app.CoinswapKeeper.AddLiquidity(suite.ctx, msg)

	// test queryLiquidity

	bz, errRes := legacyAmino.MarshalJSON(types.QueryPoolParams{LptDenom: lptCoin.Denom})
	suite.NoError(errRes)

	req.Path = fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryPool)
	req.Data = bz

	res, err = querier(suite.ctx, []string{types.QueryPool}, req)
	suite.NoError(err)

	var respone types.QueryPoolResponse
	errRes = suite.app.LegacyAmino().UnmarshalJSON(res, &respone)
	suite.NoError(errRes)
	standard := sdk.NewCoin(denomStandard, standardAmt)
	token := sdk.NewCoin(denomBTC, btcAmt)
	liquidity := sdk.NewCoin(lptCoin.Denom, standardAmt)
	suite.Equal(standard, respone.Pool.Standard)
	suite.Equal(token, respone.Pool.Token)
	suite.Equal(liquidity, respone.Pool.Lpt)
	suite.Equal(suite.app.CoinswapKeeper.GetParams(suite.ctx).Fee.String(), respone.Pool.Fee)
}
