package types

import "github.com/irisnet/irishub/app/v3/service/exported"

const (
	QueryFeed      = "feed"      // QueryFeed
	QueryFeeds     = "feeds"     // QueryFeeds
	QueryFeedValue = "feedValue" // QueryFeedValue
)

// QueryFeedParams defines the params to query a feed definition
type QueryFeedParams struct {
	FeedName string
}

// QueryFeedsParams defines the params to query a feed list by state
type QueryFeedsParams struct {
	State exported.RequestContextState
}

// QueryFeedValueParams defines the params to query a feed result
type QueryFeedValueParams struct {
	FeedName string
}
