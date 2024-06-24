package keeper_test

import (
	gocontext "context"
	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"mods.irisnet.org/modules/random/types"
)

func (suite *KeeperTestSuite) TestGRPCQueryRandom() {
	app, ctx := suite.app, suite.ctx
	reqID := []byte("test")
	random := types.NewRandom(hex.EncodeToString(reqID), 1, "test")

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.keeper)
	queryClient := types.NewQueryClient(queryHelper)

	_, err := queryClient.Random(gocontext.Background(), &types.QueryRandomRequest{ReqId: hex.EncodeToString(reqID)})
	suite.Require().Error(err)

	suite.keeper.SetRandom(ctx, reqID, random)

	randomResp, err := queryClient.Random(gocontext.Background(), &types.QueryRandomRequest{ReqId: hex.EncodeToString(reqID)})
	suite.Require().NoError(err)
	expected, _ := suite.keeper.GetRandom(ctx, reqID)
	suite.Equal(expected, *randomResp.Random)
}

func (suite *KeeperTestSuite) TestGRPCRandomRequestQueue() {
	app, ctx := suite.app, suite.ctx
	_, _, addr := testdata.KeyTestPubAddr()
	reqID := []byte("test_req_id")
	txHash := []byte("test_hash")
	request := types.NewRequest(1, addr.String(), string(txHash), false, sdk.NewCoins(), "")

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.keeper)
	queryClient := types.NewQueryClient(queryHelper)

	_, err := queryClient.RandomRequestQueue(gocontext.Background(), &types.QueryRandomRequestQueueRequest{Height: 1})
	suite.Require().NoError(err)

	suite.keeper.EnqueueRandomRequest(ctx, 1, reqID, request)

	randomResp, err := queryClient.RandomRequestQueue(gocontext.Background(), &types.QueryRandomRequestQueueRequest{Height: 1})
	suite.Require().NoError(err)
	var requests = make([]types.Request, 0)

	suite.keeper.IterateRandomRequestQueue(ctx, func(h int64, reqID []byte, r types.Request) (stop bool) {
		requests = append(requests, r)
		return false
	})
	suite.Equal(requests, randomResp.Requests)
}
