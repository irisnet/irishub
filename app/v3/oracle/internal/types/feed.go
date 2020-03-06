package types

import (
	"bytes"
	"fmt"
	"time"

	cmn "github.com/tendermint/tendermint/libs/common"

	sdk "github.com/irisnet/irishub/types"

	service "github.com/irisnet/irishub/app/v3/service/exported"
)

type Feed struct {
	FeedName         string         `json:"feed_name"`
	Description      string         `json:"description"`
	AggregateFunc    string         `json:"aggregate_func"`
	ValueJsonPath    string         `json:"value_json_path"`
	LatestHistory    uint64         `json:"latest_history"`
	RequestContextID cmn.HexBytes   `json:"request_context_id"`
	Creator          sdk.AccAddress `json:"creator"`
}

// String implements fmt.Stringer
func (f Feed) String() string {
	return fmt.Sprintf(`Feed:
	  FeedName:                 %s
	  Description:              %s
	  AggregateFunc:            %s
	  ValueJsonPath:            %s
	  LatestHistory:            %d
	  RequestContextID:         %s
	  Creator:                  %s`,
		f.FeedName,
		f.Description,
		f.AggregateFunc,
		f.ValueJsonPath,
		f.LatestHistory,
		f.RequestContextID.String(),
		f.Creator.String(),
	)
}

type FeedValue struct {
	Data      string    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

// String implements fmt.Stringer
func (f FeedValue) String() string {
	return fmt.Sprintf(` FeedValue:
		Data:                 %s
		Timestamp:            %s`,
		f.Data,
		f.Timestamp.String(),
	)
}

type FeedValues []FeedValue

// String implements fmt.Stringer
func (fv FeedValues) String() string {
	var bf bytes.Buffer
	bf.WriteString("[")
	for _, f := range fv {
		bf.WriteString("\n")
		bf.WriteString(f.String())
		bf.WriteString("\n")
	}
	bf.WriteString("]")
	return bf.String()
}

type FeedContext struct {
	Feed              Feed                        `json:"feed"`
	ServiceName       string                      `json:"service_name"`
	Providers         []sdk.AccAddress            `json:"providers"`
	Input             string                      `json:"input"`
	Timeout           int64                       `json:"timeout"`
	ServiceFeeCap     sdk.Coins                   `json:"service_fee_cap"`
	RepeatedFrequency uint64                      `json:"repeated_frequency"`
	RepeatedTotal     int64                       `json:"repeated_total"`
	ResponseThreshold uint16                      `json:"response_threshold"`
	State             service.RequestContextState `json:"state"`
}

// String implements fmt.Stringer
func (f FeedContext) String() string {
	var bf bytes.Buffer
	for _, addr := range f.Providers {
		bf.WriteString(addr.String())
		bf.WriteString(",")
	}
	return fmt.Sprintf(` FeedContext:
	%s
	ServiceName:                %s
	Providers:                  %s
	Input:                      %s
	Timeout:                    %d
	ServiceFeeCap:              %s
	RepeatedFrequency:          %d
	RepeatedTotal:              %d
	ResponseThreshold:          %d
	State:                      %s`,
		f.Feed.String(),
		f.ServiceName,
		bf.String(),
		f.Input,
		f.Timeout,
		f.ServiceFeeCap,
		f.RepeatedFrequency,
		f.RepeatedTotal,
		f.ResponseThreshold,
		f.State.String(),
	)
}

type FeedsContext []FeedContext

// String implements fmt.Stringer
func (fc FeedsContext) String() string {
	var bf bytes.Buffer
	bf.WriteString("[")
	for _, f := range fc {
		bf.WriteString("\n")
		bf.WriteString(f.String())
		bf.WriteString("\n")
	}
	bf.WriteString("]")
	return bf.String()
}
