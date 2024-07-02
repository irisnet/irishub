package types

import (
	"testing"

	"github.com/cometbft/cometbft/crypto/tmhash"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var (
	emptyAddr = ""
	addr1     = sdk.AccAddress(tmhash.SumTruncated([]byte("addr1"))).String()
	addr2     = sdk.AccAddress(tmhash.SumTruncated([]byte("addr2"))).String()
)

func TestMsgCreateFeed_ValidateBasic(t *testing.T) {
	tests := []struct {
		testCase string
		MsgCreateFeed
		expectPass bool
	}{{
		"basic good",
		MsgCreateFeed{
			FeedName:          "feedEthPrice",
			AggregateFunc:     "avg",
			ValueJsonPath:     "data.price",
			LatestHistory:     10,
			Description:       "feed eth price",
			ServiceName:       "GetEthPrice",
			Providers:         []string{addr1, addr2},
			Input:             "eth",
			Timeout:           5,
			ServiceFeeCap:     sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
			RepeatedFrequency: 5,
			ResponseThreshold: 1,
			Creator:           addr1,
		},
		true,
	}, {
		"wrong FeedName,invalid char",
		MsgCreateFeed{
			FeedName:          "$feedEthPrice",
			AggregateFunc:     "avg",
			ValueJsonPath:     "data.price",
			LatestHistory:     10,
			Description:       "feed eth price",
			ServiceName:       "GetEthPrice",
			Providers:         []string{addr1, addr2},
			Input:             "eth",
			Timeout:           5,
			ServiceFeeCap:     sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
			RepeatedFrequency: 5,
			ResponseThreshold: 1,
			Creator:           addr1,
		},
		false,
	}, {
		"wrong AggregateFunc",
		MsgCreateFeed{
			FeedName:          "feedEthPrice",
			AggregateFunc:     "",
			ValueJsonPath:     "data.price",
			LatestHistory:     10,
			Description:       "feed eth price",
			ServiceName:       "GetEthPrice",
			Providers:         []string{addr1, addr2},
			Input:             "eth",
			Timeout:           5,
			ServiceFeeCap:     sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
			RepeatedFrequency: 5,
			ResponseThreshold: 1,
			Creator:           addr1,
		},
		false,
	}, {
		"wrong MaxLatestHistory(=0)",
		MsgCreateFeed{
			FeedName:          "feedEthPrice",
			AggregateFunc:     "avg",
			ValueJsonPath:     "data.price",
			LatestHistory:     0,
			Description:       "feed eth price",
			ServiceName:       "GetEthPrice",
			Providers:         []string{addr1, addr2},
			Input:             "eth",
			Timeout:           5,
			ServiceFeeCap:     sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
			RepeatedFrequency: 5,
			ResponseThreshold: 1,
			Creator:           addr1,
		},
		false,
	}, {
		"wrong MaxLatestHistory(>100)",
		MsgCreateFeed{
			FeedName:          "feedEthPrice",
			AggregateFunc:     "avg",
			ValueJsonPath:     "data.price",
			LatestHistory:     101,
			Description:       "feed eth price",
			ServiceName:       "GetEthPrice",
			Providers:         []string{addr1, addr2},
			Input:             "eth",
			Timeout:           5,
			ServiceFeeCap:     sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
			RepeatedFrequency: 5,
			ResponseThreshold: 1,
			Creator:           addr1,
		},
		false,
	}, {
		"wrong Description(>200)",
		MsgCreateFeed{
			FeedName:          "feedEthPrice",
			AggregateFunc:     "avg",
			ValueJsonPath:     "data.price",
			LatestHistory:     10,
			Description:       "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			ServiceName:       "GetEthPrice",
			Providers:         []string{addr1, addr2},
			Input:             "eth",
			Timeout:           5,
			ServiceFeeCap:     sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
			RepeatedFrequency: 5,
			ResponseThreshold: 1,
			Creator:           addr1,
		},
		false,
	}, {
		"wrong ServiceName(>70)",
		MsgCreateFeed{
			FeedName:          "feedEthPrice",
			AggregateFunc:     "avg",
			ValueJsonPath:     "data.price",
			LatestHistory:     10,
			Description:       "feed eth price",
			ServiceName:       "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			Providers:         []string{addr1, addr2},
			Input:             "eth",
			Timeout:           5,
			ServiceFeeCap:     sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
			RepeatedFrequency: 5,
			ResponseThreshold: 1,
			Creator:           addr1,
		},
		false,
	}, {
		"wrong ServiceName,invalid char",
		MsgCreateFeed{
			FeedName:          "feedEthPrice",
			AggregateFunc:     "avg",
			ValueJsonPath:     "data.price",
			LatestHistory:     10,
			Description:       "feed eth price",
			ServiceName:       "$GetEthPrice",
			Providers:         []string{addr1, addr2},
			Input:             "eth",
			Timeout:           5,
			ServiceFeeCap:     sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
			RepeatedFrequency: 5,
			ResponseThreshold: 1,
			Creator:           addr1,
		},
		false,
	}, {
		"empty Providers",
		MsgCreateFeed{
			FeedName:          "feedEthPrice",
			AggregateFunc:     "avg",
			ValueJsonPath:     "data.price",
			LatestHistory:     10,
			Description:       "feed eth price",
			ServiceName:       "GetEthPrice",
			Providers:         []string{},
			Input:             "eth",
			Timeout:           5,
			ServiceFeeCap:     sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
			RepeatedFrequency: 5,
			ResponseThreshold: 1,
			Creator:           addr1,
		},
		false,
	}, {
		"invalid ResponseThreshold",
		MsgCreateFeed{
			FeedName:          "feedEthPrice",
			AggregateFunc:     "avg",
			ValueJsonPath:     "data.price",
			LatestHistory:     10,
			Description:       "feed eth price",
			ServiceName:       "GetEthPrice",
			Providers:         []string{addr1, addr2},
			Input:             "eth",
			Timeout:           5,
			ServiceFeeCap:     sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
			RepeatedFrequency: 5,
			ResponseThreshold: 3,
			Creator:           addr1,
		},
		false,
	}, {
		"invalid ResponseThreshold",
		MsgCreateFeed{
			FeedName:          "feedEthPrice",
			AggregateFunc:     "avg",
			ValueJsonPath:     "data.price",
			LatestHistory:     10,
			Description:       "feed eth price",
			ServiceName:       "GetEthPrice",
			Providers:         []string{addr1, addr2},
			Input:             "eth",
			Timeout:           5,
			ServiceFeeCap:     sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
			RepeatedFrequency: 5,
			ResponseThreshold: 1,
			Creator:           emptyAddr,
		},
		false,
	}}

	for _, tc := range tests {
		if tc.expectPass {
			require.Nil(t, tc.MsgCreateFeed.ValidateBasic(), "test: %v", tc.testCase)
		} else {
			require.NotNil(t, tc.MsgCreateFeed.ValidateBasic(), "test: %v", tc.testCase)
		}
	}
}

func TestMsgStartFeed_ValidateBasic(t *testing.T) {
	tests := []struct {
		testCase string
		MsgStartFeed
		expectPass bool
	}{{
		"basic good",
		MsgStartFeed{
			FeedName: "feedEthPrice",
			Creator:  addr1,
		},
		true,
	}, {
		"wrong FeedName,invalid char",
		MsgStartFeed{
			FeedName: "$feedEthPrice",
			Creator:  addr1,
		},
		false,
	}, {
		"empty Creator",
		MsgStartFeed{
			FeedName: "feedEthPrice",
			Creator:  emptyAddr,
		},
		false,
	}}

	for _, tc := range tests {
		if tc.expectPass {
			require.Nil(t, tc.MsgStartFeed.ValidateBasic(), "test: %v", tc.testCase)
		} else {
			require.NotNil(t, tc.MsgStartFeed.ValidateBasic(), "test: %v", tc.testCase)
		}
	}
}

func TestMsgPauseFeed_ValidateBasic(t *testing.T) {
	tests := []struct {
		testCase string
		MsgPauseFeed
		expectPass bool
	}{{
		"basic good",
		MsgPauseFeed{
			FeedName: "feedEthPrice",
			Creator:  addr1,
		},
		true,
	}, {
		"wrong FeedName,invalid char",
		MsgPauseFeed{
			FeedName: "$feedEthPrice",
			Creator:  addr1,
		},
		false,
	}, {
		"empty Creator",
		MsgPauseFeed{
			FeedName: "feedEthPrice",
			Creator:  emptyAddr,
		},
		false,
	}}

	for _, tc := range tests {
		if tc.expectPass {
			require.Nil(t, tc.MsgPauseFeed.ValidateBasic(), "test: %v", tc.testCase)
		} else {
			require.NotNil(t, tc.MsgPauseFeed.ValidateBasic(), "test: %v", tc.testCase)
		}
	}
}

func TestMsgEditFeed_ValidateBasic(t *testing.T) {
	tests := []struct {
		testCase string
		MsgEditFeed
		expectPass bool
	}{
		{
			"basic good",
			MsgEditFeed{
				FeedName:          "feedEthPrice",
				LatestHistory:     10,
				Providers:         []string{addr1, addr2},
				ServiceFeeCap:     sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
				RepeatedFrequency: 5,
				ResponseThreshold: 1,
				Creator:           addr1,
			},
			true,
		}, {
			"wrong FeedName, invalid char",
			MsgEditFeed{
				FeedName:          "$feedEthPrice",
				LatestHistory:     10,
				Providers:         []string{addr1, addr2},
				ServiceFeeCap:     sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
				RepeatedFrequency: 5,
				ResponseThreshold: 1,
				Creator:           addr1,
			},
			false,
		}, {
			"wrong MaxLatestHistory(>100)",
			MsgEditFeed{
				FeedName:          "feedEthPrice",
				LatestHistory:     101,
				Providers:         []string{addr1, addr2},
				ServiceFeeCap:     sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
				RepeatedFrequency: 5,
				ResponseThreshold: 1,
				Creator:           addr1,
			},
			false,
		}, {
			"invalid ResponseThreshold",
			MsgEditFeed{
				FeedName:          "feedEthPrice",
				LatestHistory:     10,
				Providers:         []string{addr1, addr2},
				ServiceFeeCap:     sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
				RepeatedFrequency: 5,
				ResponseThreshold: 3,
				Creator:           addr1,
			},
			false,
		}, {
			"empty Creator",
			MsgEditFeed{
				FeedName:          "feedEthPrice",
				LatestHistory:     10,
				Providers:         []string{addr1, addr2},
				ServiceFeeCap:     sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
				RepeatedFrequency: 5,
				ResponseThreshold: 1,
				Creator:           emptyAddr,
			},
			false,
		},
	}

	for _, tc := range tests {
		if tc.expectPass {
			require.Nil(t, tc.MsgEditFeed.ValidateBasic(), "test: %v", tc.testCase)
		} else {
			require.NotNil(t, tc.MsgEditFeed.ValidateBasic(), "test: %v", tc.testCase)
		}
	}
}
