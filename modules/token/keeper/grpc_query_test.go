package keeper_test

import (
	gocontext "context"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/token/types"
)

func (suite *KeeperTestSuite) TestGRPCQueryToken() {
	app, ctx := suite.app, suite.ctx
	_, _, addr := testdata.KeyTestPubAddr()
	token := types.NewToken("btc", "Bitcoin Token", "satoshi", 18, 21000000, 22000000, true, addr)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.TokenKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	_ = suite.app.TokenKeeper.AddToken(ctx, token)

	// Query token
	tokenResp1, err := queryClient.Token(gocontext.Background(), &types.QueryTokenRequest{Denom: "btc"})
	suite.Require().NoError(err)
	suite.Require().NotNil(tokenResp1)

	tokenResp2, err := queryClient.Token(gocontext.Background(), &types.QueryTokenRequest{Denom: "satoshi"})
	suite.Require().NoError(err)
	suite.Require().NotNil(tokenResp2)

	// Query tokens
	tokensResp1, err := queryClient.Tokens(gocontext.Background(), &types.QueryTokensRequest{})
	suite.Require().NoError(err)
	suite.Require().NotNil(tokensResp1)
	suite.Len(tokensResp1.Tokens, 2)
}

func (suite *KeeperTestSuite) TestGRPCQueryFees() {
	app, ctx := suite.app, suite.ctx

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.TokenKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	_, err := queryClient.Fees(gocontext.Background(), &types.QueryFeesRequest{Symbol: "test"})
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) TestGRPCQueryParams() {
	app, ctx := suite.app, suite.ctx

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.TokenKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	paramsResp, err := queryClient.Params(gocontext.Background(), &types.QueryParamsRequest{})
	params := app.TokenKeeper.GetParamSet(ctx)
	suite.Require().NoError(err)
	suite.Equal(params, paramsResp.Params)
}

func (suite *KeeperTestSuite) TestGRPCQueryTotalBurn() {
	app, ctx := suite.app, suite.ctx

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.TokenKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	_, _, addr := testdata.KeyTestPubAddr()
	token := types.NewToken("btc", "Bitcoin Token", "satoshi", 18, 21000000, 22000000, true, addr)
	err := suite.app.TokenKeeper.AddToken(ctx, token)
	suite.Require().NoError(err)

	buinCoin := sdk.NewInt64Coin("satoshi", 1000000000000000000)
	app.TokenKeeper.AddBurnCoin(ctx, buinCoin)

	resp, err := queryClient.TotalBurn(gocontext.Background(), &types.QueryTotalBurnRequest{})
	suite.Require().NoError(err)
	suite.Len(resp.BurnedCoins, 1)
	suite.EqualValues(buinCoin, resp.BurnedCoins[0])
}
