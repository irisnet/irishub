package types

import (
	sdk "github.com/irisnet/irishub/types"
	"time"
)

type Feed struct {
	FeedName              string         `json:"feed_name"`
	AggregateMethod       string         `json:"aggregate_method"`
	AggregateArgsJsonPath string         `json:"aggregate_args_json_path"`
	LatestHistory         uint64         `json:"latest_history"`
	RequestContextID      []byte         `json:"request_context_id"`
	Owner                 sdk.AccAddress `json:"owner"`
}
type Value interface{}
type FeedResult struct {
	Data      Value     `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}
type FeedResults []FeedResult
