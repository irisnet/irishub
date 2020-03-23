package types

import "fmt"

var (
	TagFeedName  = "feed-name"
	TagCreator   = "creator"
	TagFeedValue = func(feedName string) string {
		return fmt.Sprintf("%s.%s", TagFeedName, feedName)
	}
)
