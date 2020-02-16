package types

import (
	"time"

	sdk "github.com/irisnet/irishub/types"
)

type Feed struct {
	FeedName         string         `json:"feed_name"`
	AggregateFunc    string         `json:"aggregate_func"`
	ValueJsonPath    string         `json:"value_json_path"`
	LatestHistory    uint64         `json:"latest_history"`
	RequestContextID []byte         `json:"request_context_id"`
	Creator          sdk.AccAddress `json:"creator"`
}
type FeedResult struct {
	Data      string    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}
type FeedResults []FeedResult
