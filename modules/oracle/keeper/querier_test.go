package keeper_test

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irismod/service/exported"

	guardiantypes "github.com/irisnet/irishub/modules/guardian/types"
	"github.com/irisnet/irishub/modules/oracle/keeper"
	"github.com/irisnet/irishub/modules/oracle/types"
)

func (suite *KeeperTestSuite) TestNewQuerier() {
	// add profiler
	suite.app.GuardianKeeper.AddProfiler(suite.ctx, guardiantypes.NewGuardian("test", guardiantypes.Ordinary, addrs[0], addrs[0]))

	msg := &types.MsgCreateFeed{
		FeedName:          "ethPrice",
		ServiceName:       "GetEthPrice",
		AggregateFunc:     "avg",
		ValueJsonPath:     "high",
		LatestHistory:     5,
		Providers:         []sdk.AccAddress{addrs[1]},
		Input:             "xxxx",
		Timeout:           10,
		ServiceFeeCap:     sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
		RepeatedFrequency: 1,
		ResponseThreshold: 1,
		Creator:           addrs[0],
	}

	//================test CreateFeed start================
	err := suite.keeper.CreateFeed(suite.ctx, msg)
	suite.NoError(err)

	//test QueryFeed
	querier := keeper.NewQuerier(suite.keeper, suite.cdc)

	params := types.QueryFeedParams{
		FeedName: msg.FeedName,
	}
	bz := suite.cdc.MustMarshalJSON(params)
	res, err := querier(suite.ctx, []string{types.QueryFeed}, abci.RequestQuery{
		Data: bz,
	})
	suite.NoError(err)

	var feedCtx types.FeedContext
	suite.cdc.MustUnmarshalJSON(res, &feedCtx)

	suite.EqualValues(types.FeedContext{
		Feed: &types.Feed{
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
		ResponseThreshold: msg.ResponseThreshold,
		State:             exported.PAUSED,
	}, feedCtx)

	//test QueryFeeds
	params1 := types.QueryFeedsParams{
		State: "paused",
	}
	bz = suite.cdc.MustMarshalJSON(params1)
	res, err = querier(suite.ctx, []string{types.QueryFeeds}, abci.RequestQuery{
		Data: bz,
	})
	suite.NoError(err)

	var feedsCtx []types.FeedContext
	suite.cdc.MustUnmarshalJSON(res, &feedsCtx)
	suite.Len(feedsCtx, 1)
	suite.EqualValues(types.FeedContext{
		Feed: &types.Feed{
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
		ResponseThreshold: msg.ResponseThreshold,
		State:             exported.PAUSED,
	}, feedsCtx[0])

	//================test StartFeed start================
	err = suite.keeper.StartFeed(suite.ctx, &types.MsgStartFeed{
		FeedName: msg.FeedName,
		Creator:  addrs[0],
	})
	suite.NoError(err)

	//test QueryValue
	params2 := types.QueryFeedValueParams{
		FeedName: msg.FeedName,
	}
	bz = suite.cdc.MustMarshalJSON(params2)
	res, err = querier(
		suite.ctx,
		[]string{types.QueryFeedValue},
		abci.RequestQuery{Data: bz},
	)
	suite.NoError(err)
	var feedValues types.FeedValues
	suite.cdc.MustUnmarshalJSON(res, &feedValues)
	suite.Len(feedsCtx, 1)
	suite.Equal("250.00000000", feedValues[0].Data)
}
