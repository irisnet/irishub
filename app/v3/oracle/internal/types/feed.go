package types

import (
	"time"

	"github.com/irisnet/irishub/app/v3/service/exported"

	sdk "github.com/irisnet/irishub/types"
)

type Feed struct {
	FeedName         string         `json:"feed_name"`
	Description      string         `json:"description"`
	AggregateFunc    string         `json:"aggregate_func"`
	ValueJsonPath    string         `json:"value_json_path"`
	LatestHistory    uint64         `json:"latest_history"`
	RequestContextID []byte         `json:"request_context_id"`
	Creator          sdk.AccAddress `json:"creator"`
}
type FeedValue struct {
	Data      string    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}
type FeedValues []FeedValue

type FeedContext struct {
	Feed              Feed                         `json:"feed"`
	ServiceName       string                       `json:"service_name"`
	Providers         []sdk.AccAddress             `json:"providers"`
	Input             string                       `json:"input"`
	Timeout           int64                        `json:"timeout"`
	ServiceFeeCap     sdk.Coins                    `json:"service_fee_cap"`
	RepeatedFrequency uint64                       `json:"repeated_frequency"`
	RepeatedTotal     int64                        `json:"repeated_total"`
	ResponseThreshold uint16                       `json:"response_threshold"`
	State             exported.RequestContextState `json:"state"`
}
