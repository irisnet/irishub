package keeper_test

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/irisnet/irismod/modules/token/keeper"
	"github.com/irisnet/irismod/modules/token/types"
)

func (suite *KeeperTestSuite) TestQueryToken() {
	ctx := suite.ctx
	querier := keeper.NewQuerier(suite.keeper, suite.legacyAmino)

	params := types.QueryTokenParams{
		Denom: types.GetNativeToken().Symbol,
	}
	bz := suite.legacyAmino.MustMarshalJSON(params)
	query := abci.RequestQuery{
		Path: fmt.Sprintf("/custom/%s/%s", types.QuerierRoute, types.QueryToken),
		Data: bz,
	}

	data, err := querier(ctx, []string{types.QueryToken}, query)
	suite.Nil(err)

	data2 := codec.MustMarshalJSONIndent(suite.legacyAmino, types.GetNativeToken())
	suite.Equal(data2, data)

	//query by mint_unit
	params = types.QueryTokenParams{
		Denom: types.GetNativeToken().MinUnit,
	}

	bz = suite.legacyAmino.MustMarshalJSON(params)
	query = abci.RequestQuery{
		Path: fmt.Sprintf("/custom/%s/%s", types.QuerierRoute, types.QueryToken),
		Data: bz,
	}

	data, err = querier(ctx, []string{types.QueryToken}, query)
	suite.Nil(err)

	data2 = codec.MustMarshalJSONIndent(suite.legacyAmino, types.GetNativeToken())
	suite.Equal(data2, data)
}

func (suite *KeeperTestSuite) TestQueryTokens() {
	ctx := suite.ctx
	querier := keeper.NewQuerier(suite.keeper, suite.legacyAmino)

	params := types.QueryTokensParams{
		Owner: nil,
	}
	bz := suite.legacyAmino.MustMarshalJSON(params)
	query := abci.RequestQuery{
		Path: fmt.Sprintf("/custom/%s/%s", types.QuerierRoute, types.QueryTokens),
		Data: bz,
	}

	data, err := querier(ctx, []string{types.QueryTokens}, query)
	suite.Nil(err)

	data2 := codec.MustMarshalJSONIndent(suite.legacyAmino, []types.TokenI{types.GetNativeToken()})
	suite.Equal(data2, data)
}

func (suite *KeeperTestSuite) TestQueryFees() {
	ctx := suite.ctx
	querier := keeper.NewQuerier(suite.keeper, suite.legacyAmino)

	params := types.QueryTokenFeesParams{
		Symbol: "btc",
	}
	bz := suite.legacyAmino.MustMarshalJSON(params)
	query := abci.RequestQuery{
		Path: fmt.Sprintf("/custom/%s/%s", types.QuerierRoute, types.QueryFees),
		Data: bz,
	}

	data, err := querier(ctx, []string{types.QueryFees}, query)
	suite.Nil(err)

	var fee types.QueryFeesResponse
	suite.legacyAmino.MustUnmarshalJSON(data, &fee)
	suite.Equal(false, fee.Exist)
	suite.Equal(fmt.Sprintf("60000%s", types.GetNativeToken().MinUnit), fee.IssueFee.String())
	suite.Equal(fmt.Sprintf("6000%s", types.GetNativeToken().MinUnit), fee.MintFee.String())
}
