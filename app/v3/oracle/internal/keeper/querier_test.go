package keeper

import (
	"testing"

	"github.com/irisnet/irishub/app/v3/oracle/internal/types"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestNewQuerier(t *testing.T) {
	ctx, keeper, acc := createTestInput(t, sdk.NewInt(1000000), 2)
	query := NewQuerier(keeper)

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
	err := keeper.CreateFeed(ctx, msg)
	require.NoError(t, err)

	//test QueryFeed
	params := types.QueryFeedParams{
		FeedName: msg.FeedName,
	}
	bz := keeper.cdc.MustMarshalJSON(params)
	res, err := query(ctx, []string{types.QueryFeed}, abci.RequestQuery{
		Data: bz,
	})
	require.NoError(t, err)

	var feedCtx types.FeedContext
	keeper.cdc.MustUnmarshalJSON(res, &feedCtx)

	require.EqualValues(t, types.FeedContext{
		Feed: types.Feed{
			FeedName:         msg.FeedName,
			AggregateFunc:    msg.AggregateFunc,
			ValueJsonPath:    msg.ValueJsonPath,
			LatestHistory:    msg.LatestHistory,
			RequestContextID: mockReqCtxID,
			Creator:          msg.Creator,
		},
		ServiceName:       msg.ServiceName,
		Providers:         msg.Providers,
		Input:             msg.Input,
		Timeout:           msg.Timeout,
		ServiceFeeCap:     msg.ServiceFeeCap,
		RepeatedFrequency: msg.RepeatedFrequency,
		RepeatedTotal:     msg.RepeatedTotal,
		ResponseThreshold: msg.ResponseThreshold,
		State:             types.Pause,
	}, feedCtx)

	//test QueryFeeds
	params1 := types.QueryFeedsParams{
		State: "pause",
	}
	bz = keeper.cdc.MustMarshalJSON(params1)
	res, err = query(ctx, []string{types.QueryFeeds}, abci.RequestQuery{
		Data: bz,
	})
	require.NoError(t, err)

	var feedsCtx []types.FeedContext
	keeper.cdc.MustUnmarshalJSON(res, &feedsCtx)
	require.Len(t, feedsCtx, 1)
	require.EqualValues(t, types.FeedContext{
		Feed: types.Feed{
			FeedName:         msg.FeedName,
			AggregateFunc:    msg.AggregateFunc,
			ValueJsonPath:    msg.ValueJsonPath,
			LatestHistory:    msg.LatestHistory,
			RequestContextID: mockReqCtxID,
			Creator:          msg.Creator,
		},
		ServiceName:       msg.ServiceName,
		Providers:         msg.Providers,
		Input:             msg.Input,
		Timeout:           msg.Timeout,
		ServiceFeeCap:     msg.ServiceFeeCap,
		RepeatedFrequency: msg.RepeatedFrequency,
		RepeatedTotal:     msg.RepeatedTotal,
		ResponseThreshold: msg.ResponseThreshold,
		State:             types.Pause,
	}, feedsCtx[0])

	//================test StartFeed start================
	err = keeper.StartFeed(ctx, types.MsgStartFeed{
		FeedName: msg.FeedName,
		Creator:  acc[0].GetAddress(),
	})
	require.NoError(t, err)

	//test QueryValue
	params2 := types.QueryFeedValueParams{
		FeedName: msg.FeedName,
	}
	bz = keeper.cdc.MustMarshalJSON(params2)
	res, err = query(ctx, []string{types.QueryFeedValue}, abci.RequestQuery{
		Data: bz,
	})
	require.NoError(t, err)
	var feedValues types.FeedValues
	keeper.cdc.MustUnmarshalJSON(res, &feedValues)
	require.Len(t, feedsCtx, 1)
	require.Equal(t, "250.00000000", feedValues[0].Data)
}
