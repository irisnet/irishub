package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/irisnet/irishub/modules/asset"
	token "github.com/irisnet/irishub/modules/asset/01-token"
	"github.com/irisnet/irishub/simapp"
)

var DefaultToken = token.FungibleToken{
	Symbol:        "btc",
	Name:          "Bitcoin Network",
	Scale:         18,
	MinUnit:       "satoshi",
	InitialSupply: sdk.NewIntWithDecimal(21, 6),
	MaxSupply:     sdk.NewIntWithDecimal(21, 6),
	Mintable:      true,
}

// test that the params can be properly set and retrieved
type TestSuite struct {
	suite.Suite

	cdc *codec.Codec
	ctx sdk.Context
	k   asset.Keeper
}

func TestQuerierSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) SetupTest() {
	app := simapp.Setup(false)

	suite.cdc = app.Codec()
	suite.ctx = app.BaseApp.NewContext(false, abci.Header{})
	suite.k = app.AssetKeeper

	err := suite.k.TokenKeeper.AddToken(suite.ctx, DefaultToken)
	suite.NoError(err)
}

func (suite *TestSuite) TestNewQuerier() {
	req := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	querier := asset.NewQuerier(suite.k)
	res, err := querier(suite.ctx, []string{"other"}, req)
	suite.Error(err)
	suite.Nil(res)

	tokenBz := codec.MustMarshalJSONIndent(suite.cdc, token.Tokens{DefaultToken})

	//test QueryToken
	bz, err1 := suite.cdc.MarshalJSON(token.QueryTokenParams{
		Symbol: DefaultToken.Symbol,
	})

	suite.NoError(err1)
	req.Data = bz
	bz, err = querier(suite.ctx, []string{token.QuerierRoute, token.QueryToken}, req)
	suite.NoError(err)
	suite.Equal(tokenBz, bz)

	//test QueryTokens
	bz, err1 = suite.cdc.MarshalJSON(token.QueryTokensParams{
		Symbol: DefaultToken.Symbol,
	})
	suite.NoError(err1)
	req.Data = bz
	bz, err = querier(suite.ctx, []string{token.QuerierRoute, token.QueryTokens}, req)
	suite.NoError(err)
	suite.NoError(err)
	suite.Equal(tokenBz, bz)

	//test QueryFees
	bz, err1 = suite.cdc.MarshalJSON(token.QueryTokenFeesParams{
		Symbol: DefaultToken.Symbol,
	})
	suite.NoError(err1)
	req.Data = bz
	bz, err = querier(suite.ctx, []string{token.QuerierRoute, token.QueryFees}, req)
	suite.NoError(err)
	var tfo token.TokenFeesOutput
	err1 = suite.cdc.UnmarshalJSON(bz, &tfo)
	suite.NoError(err)
	param := token.DefaultParams()
	mintFee := sdk.NewDecFromInt(param.IssueTokenBaseFee.Amount).Mul(param.MintTokenFeeRatio)
	suite.Equal(true, tfo.Exist)
	suite.Equal(param.IssueTokenBaseFee, tfo.IssueFee)
	suite.Equal(sdk.NewCoin(param.IssueTokenBaseFee.Denom, mintFee.TruncateInt()), tfo.MintFee)

	//test QueryParameters
	bz, err = querier(suite.ctx, []string{token.QuerierRoute, token.QueryParameters}, req)
	suite.NoError(err)
	var p token.Params
	err1 = suite.cdc.UnmarshalJSON(bz, &p)
	suite.NoError(err)
	suite.Equal(token.DefaultParams(), p)
}
