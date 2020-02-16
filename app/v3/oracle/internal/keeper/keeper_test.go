package keeper

import (
	"testing"

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

	//test CreateFeed
	err := keeper.CreateFeed(ctx, msg)
	require.NoError(t, err)

	feed, existed := keeper.GetFeed(ctx, msg.FeedName)
	require.True(t, existed)
	require.EqualValues(t, types.Feed{
		FeedName:         msg.FeedName,
		AggregateFunc:    msg.AggregateFunc,
		ValueJsonPath:    msg.ValueJsonPath,
		LatestHistory:    msg.LatestHistory,
		RequestContextID: feed.RequestContextID,
		Creator:          msg.Creator,
	}, feed)

	//test StartFeed
	err = keeper.StartFeed(ctx, types.MsgStartFeed{
		FeedName: msg.FeedName,
		Creator:  acc[0].GetAddress(),
	})

	result := keeper.GetFeedResults(ctx, msg.FeedName)
	require.NoError(t, err)
	require.Len(t, result, int(msg.LatestHistory))
	require.Equal(t, "250.00000000", result[0].Data)

	//start again, will return error
	err = keeper.StartFeed(ctx, types.MsgStartFeed{
		FeedName: msg.FeedName,
		Creator:  acc[0].GetAddress(),
	})
	require.Error(t, err)

	//edit feed
	latestHistory := uint64(3)
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

	err = keeper.PauseFeed(ctx, types.MsgPauseFeed{
		FeedName: msg.FeedName,
		Creator:  acc[0].GetAddress(),
	})
	require.NoError(t, err)

	reqCtx, existed := keeper.sk.GetRequestContext(ctx, feed.RequestContextID)
	require.True(t, existed)
	require.Equal(t, types.Pause, reqCtx.State)

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

	result = keeper.GetFeedResults(ctx, msg.FeedName)
	require.NoError(t, err)
	require.Len(t, result, int(latestHistory))
	require.Equal(t, "250.00000000", result[0].Data)
}
