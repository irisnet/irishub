package keeper_test

import (
	gocontext "context"
	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/baseapp"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/irisnet/irishub/modules/random/types"
)

func (suite *KeeperTestSuite) TestGRPCQueryRandom() {
	app, ctx := suite.app, suite.ctx
	reqID := []byte("test")
	random := types.NewRandom(reqID, 1, "test")

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.RandomKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	_, err := queryClient.Random(gocontext.Background(), &types.QueryRandomRequest{ReqId: hex.EncodeToString(reqID)})
	suite.Require().Error(err)

	app.RandomKeeper.SetRandom(ctx, reqID, random)

	randomResp, err := queryClient.Random(gocontext.Background(), &types.QueryRandomRequest{ReqId: hex.EncodeToString(reqID)})
	suite.Require().NoError(err)
	expected, _ := app.RandomKeeper.GetRandom(ctx, reqID)
	suite.Equal(expected, *randomResp.Random)
}

func (suite *KeeperTestSuite) TestGRPCRandomRequestQueue() {
	app, ctx := suite.app, suite.ctx
	_, _, addr := authtypes.KeyTestPubAddr()
	reqID := []byte("test_req_id")
	txHash := []byte("test_hash")
	request := types.NewRequest(1, addr, txHash)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.RandomKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	_, err := queryClient.RandomRequestQueue(gocontext.Background(), &types.QueryRandomRequestQueueRequest{Height: 1})
	suite.Require().NoError(err)

	app.RandomKeeper.EnqueueRandomRequest(ctx, 1, reqID, request)

	randomResp, err := queryClient.RandomRequestQueue(gocontext.Background(), &types.QueryRandomRequestQueueRequest{Height: 1})
	suite.Require().NoError(err)
	var requests = make([]types.Request, 0)

	app.RandomKeeper.IterateRandomRequestQueue(ctx, func(h int64, r types.Request) (stop bool) {
		requests = append(requests, r)
		return false
	})
	suite.Equal(requests, randomResp.Requests)
}
