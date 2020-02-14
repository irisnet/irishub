package keeper

import (
	"testing"

	"github.com/irisnet/irishub/app/v3/oracle/internal/types"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
)

func TestCreateFeed(t *testing.T) {
	ctx, keeper, acc := createTestInput(t, sdk.NewInt(1000000), 2)
	msg := types.MsgCreateFeed{
		FeedName:              "ethPrice",
		ServiceName:           "GetRthPrice",
		AggregateMethod:       "avg",
		AggregateArgsJsonPath: "high",
		LatestHistory:         1,
		Providers:             []sdk.AccAddress{acc[0].GetAddress()},
		Input:                 "xxxx",
		Timeout:               10,
		ServiceFeeCap:         sdk.NewCoins(sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(100))),
		RepeatedFrequency:     1,
		RepeatedTotal:         100,
		ResponseThreshold:     1,
		Owner:                 acc[0].GetAddress(),
	}
	err := keeper.CreateFeed(ctx, msg)
	require.NoError(t, err)

	err = keeper.StartFeed(ctx, types.MsgStartFeed{
		FeedName: msg.FeedName,
		Owner:    acc[0].GetAddress(),
	})

	result := keeper.GetFeedResults(ctx, msg.FeedName)
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, 250.0, result[0].Data)
}
