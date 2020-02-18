package oracle

import (
	"github.com/irisnet/irishub/app/v3/oracle/internal/keeper"
	"github.com/irisnet/irishub/app/v3/oracle/internal/types"
)

// nolint
type (
	Keeper = keeper.Keeper

	MsgCreateFeed = types.MsgCreateFeed
	MsgStartFeed  = types.MsgStartFeed
	MsgPauseFeed  = types.MsgPauseFeed
	MsgEditFeed   = types.MsgEditFeed

	QueryFeedParams      = types.QueryFeedParams
	QueryFeedsParams     = types.QueryFeedsParams
	QueryFeedValueParams = types.QueryFeedValueParams

	FeedContext = types.FeedContext
	FeedValue   = types.FeedValue
	FeedValues  = types.FeedValues
)

const (
	QueryFeed        = types.QueryFeed
	QueryFeeds       = types.QueryFeeds
	QueryFeedValue   = types.QueryFeedValue
	DefaultCodespace = types.DefaultCodespace
)

var (
	NewKeeper     = keeper.NewKeeper
	NewQuerier    = keeper.NewQuerier
	RegisterCodec = types.RegisterCodec
)
