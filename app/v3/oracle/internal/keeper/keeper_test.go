package keeper

import (
	"testing"

	"github.com/irisnet/irishub/app/v3/service/exported"

	"github.com/irisnet/irishub/app/v3/oracle/internal/types"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
)

func TestFeed(t *testing.T) {
	ctx, keeper, acc := createTestInput(t, sdk.NewInt(1000000), 2)
	msg := types.MsgCreateFeed{
		FeedName:          "ethPrice",
		ServiceName:       "GetRthPrice",
		AggregateFunc:     "avg",
		ValueJsonPath:     "high",
		LatestHistory:     5,
		Providers:         []sdk.AccAddress{acc[0].GetAddress()},
		Input:             "xxxx",
		Timeout:           10,
		ServiceFeeCap:     sdk.NewCoins(sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(100))),
		RepeatedFrequency: 1,
		RepeatedTotal:     10,
		ResponseThreshold: 1,
		Creator:           acc[0].GetAddress(),
	}

	//================test CreateFeed start================
	_, err := keeper.CreateFeed(ctx, msg)
	require.NoError(t, err)

	//check feed existed
	feed, existed := keeper.GetFeed(ctx, msg.FeedName)
	require.True(t, existed)
	require.EqualValues(t, types.Feed{
		FeedName:         msg.FeedName,
		AggregateFunc:    msg.AggregateFunc,
		ValueJsonPath:    msg.ValueJsonPath,
		LatestHistory:    msg.LatestHistory,
		RequestContextID: mockReqCtxID,
		Creator:          msg.Creator,
	}, feed)

	//check feed state
	var feeds []types.Feed
	keeper.IteratorFeedsByState(ctx, exported.PAUSED, func(feed types.Feed) {
		feeds = append(feeds, feed)
	})
	require.Len(t, feeds, 1)
	require.Equal(t, msg.FeedName, feeds[0].FeedName)
	//================test CreateFeed end================

	//================test StartFeed start================
	err = keeper.StartFeed(ctx, types.MsgStartFeed{
		FeedName: msg.FeedName,
		Creator:  acc[0].GetAddress(),
	})

	//check feed result
	result := keeper.GetFeedValues(ctx, msg.FeedName)
	require.NoError(t, err)
	require.Len(t, result, int(msg.LatestHistory))
	require.Equal(t, "250.00000000", result[0].Data)

	//check feed state
	var feeds1 []types.Feed
	keeper.IteratorFeedsByState(ctx, exported.RUNNING, func(feed types.Feed) {
		feeds1 = append(feeds1, feed)
	})
	require.Len(t, feeds1, 1)
	require.Equal(t, msg.FeedName, feeds1[0].FeedName)

	//start again, will return error
	err = keeper.StartFeed(ctx, types.MsgStartFeed{
		FeedName: msg.FeedName,
		Creator:  acc[0].GetAddress(),
	})
	require.Error(t, err)
	//================test StartFeed end================

	//================test EditFeed start================
	latestHistory := uint64(1)
	err = keeper.EditFeed(ctx, types.MsgEditFeed{
		FeedName:          msg.FeedName,
		LatestHistory:     latestHistory,
		Providers:         []sdk.AccAddress{acc[0].GetAddress()},
		ServiceFeeCap:     sdk.NewCoins(sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(100))),
		RepeatedFrequency: 1,
		RepeatedTotal:     10,
		ResponseThreshold: 1,
		Creator:           acc[0].GetAddress(),
	})
	require.NoError(t, err)

	//check feed existed
	feed, existed = keeper.GetFeed(ctx, msg.FeedName)
	require.True(t, existed)
	require.EqualValues(t, types.Feed{
		FeedName:         msg.FeedName,
		AggregateFunc:    msg.AggregateFunc,
		ValueJsonPath:    msg.ValueJsonPath,
		LatestHistory:    latestHistory,
		RequestContextID: feed.RequestContextID,
		Creator:          msg.Creator,
	}, feed)
	//================test EditFeed end================

	//================test PauseFeed start================
	err = keeper.PauseFeed(ctx, types.MsgPauseFeed{
		FeedName: msg.FeedName,
		Creator:  acc[0].GetAddress(),
	})
	require.NoError(t, err)

	reqCtx, existed := keeper.sk.GetRequestContext(ctx, feed.RequestContextID)
	require.True(t, existed)
	require.Equal(t, exported.PAUSED, reqCtx.State)

	//pause again, will return error
	err = keeper.PauseFeed(ctx, types.MsgPauseFeed{
		FeedName: msg.FeedName,
		Creator:  acc[0].GetAddress(),
	})
	require.Error(t, err)

	//Start Feed again
	err = keeper.StartFeed(ctx, types.MsgStartFeed{
		FeedName: msg.FeedName,
		Creator:  acc[0].GetAddress(),
	})

	//check feed result
	result = keeper.GetFeedValues(ctx, msg.FeedName)
	require.NoError(t, err)
	require.Len(t, result, int(latestHistory))
	require.Equal(t, "250.00000000", result[0].Data)

	//check feed state
	var feeds2 []types.Feed
	keeper.IteratorFeedsByState(ctx, exported.RUNNING, func(feed types.Feed) {
		feeds2 = append(feeds2, feed)
	})
	require.Len(t, feeds2, 1)
	require.Equal(t, msg.FeedName, feeds2[0].FeedName)
	//================test PauseFeed end================
}
