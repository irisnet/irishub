package keeper_test

import (
	gocontext "context"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"

	"mods.irisnet.org/modules/oracle/keeper"
	"mods.irisnet.org/modules/oracle/types"
)

func (suite *KeeperTestSuite) TestGRPCQueryFeed() {
	app, ctx := suite.app, suite.ctx
	_, _, addr := testdata.KeyTestPubAddr()

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.keeper)
	queryClient := types.NewQueryClient(queryHelper)

	// Query feed
	_, err := queryClient.Feed(gocontext.Background(), &types.QueryFeedRequest{})
	suite.Error(err)

	// Query feeds
	_, err = queryClient.Feeds(gocontext.Background(), &types.QueryFeedsRequest{})
	suite.NoError(err)

	// Add feed
	feedName := "test"
	feed := types.Feed{FeedName: feedName, Creator: addr.String()}
	suite.keeper.SetFeed(ctx, feed)

	// Query feed
	feedResp, err := queryClient.Feed(gocontext.Background(), &types.QueryFeedRequest{FeedName: feedName})
	suite.NoError(err)
	expectedFeed, _ := suite.keeper.GetFeed(ctx, feedName)
	expectedFeedCtx := keeper.BuildFeedContext(ctx, suite.keeper, expectedFeed)

	suite.Equal(expectedFeedCtx, feedResp.Feed)

	// Query feeds
	feedsResp, err := queryClient.Feeds(gocontext.Background(), &types.QueryFeedsRequest{})
	suite.NoError(err)
	suite.Len(feedsResp.Feeds, 1)
	suite.Equal([]types.FeedContext{expectedFeedCtx}, feedsResp.Feeds)
}

func (suite *KeeperTestSuite) TestGRPCQueryFeedValue() {
	app, ctx := suite.app, suite.ctx

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.keeper)
	queryClient := types.NewQueryClient(queryHelper)

	// Query feed
	_, err := queryClient.FeedValue(gocontext.Background(), &types.QueryFeedValueRequest{})
	suite.NoError(err)

	// Add feed value
	feedName := "test"
	feedValue := types.FeedValue{Data: "test", Timestamp: time.Now()}
	suite.keeper.SetFeedValue(ctx, feedName, 1, 10, feedValue)

	// Query feed
	valueResp, err := queryClient.FeedValue(gocontext.Background(), &types.QueryFeedValueRequest{FeedName: feedName})
	suite.NoError(err)
	expectedValues := suite.keeper.GetFeedValues(ctx, feedName)
	suite.Equal([]types.FeedValue(expectedValues), valueResp.FeedValues)
}
