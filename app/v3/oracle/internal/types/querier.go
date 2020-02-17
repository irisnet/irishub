package types

const (
	QueryFeed       = "feed"       // QueryFeed
	QueryFeeds      = "feeds"      // QueryFeeds
	QueryFeedResult = "feedResult" // QueryFeedResult
)

// QueryFeedParams defines the params to query a feed definition
type QueryFeedParams struct {
	FeedName string
}

// QueryFeedsParams defines the params to query a feed list by state
type QueryFeedsParams struct {
	State string
}

// QueryFeedsResult defines the params to query a feed result
type QueryFeedsResult struct {
	FeedName string
}
