package types

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	service "github.com/irismod/service/exported"
)

type FeedValues []FeedValue

// String implements fmt.Stringer
func (fv FeedValues) String() string {
	if len(fv) == 0 {
		return "[]"
	}

	var str string
	for _, f := range fv {
		str += f.String() + "\n"
	}
	return str
}

// FeedContext represents a struct that keeps track of the feed runtime state
type FeedContext struct {
	Feed              Feed                        `json:"feed"`
	ServiceName       string                      `json:"service_name"`
	Providers         []sdk.AccAddress            `json:"providers"`
	Input             string                      `json:"input"`
	Timeout           int64                       `json:"timeout"`
	ServiceFeeCap     sdk.Coins                   `json:"service_fee_cap"`
	RepeatedFrequency uint64                      `json:"repeated_frequency"`
	ResponseThreshold uint32                      `json:"response_threshold"`
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
	ResponseThreshold:          %d
	State:                      %s`,
		f.Feed.String(),
		f.ServiceName,
		bf.String(),
		f.Input,
		f.Timeout,
		f.ServiceFeeCap,
		f.RepeatedFrequency,
		f.ResponseThreshold,
		f.State.String(),
	)
}

type FeedsContext []FeedContext

// String implements fmt.Stringer
func (fc FeedsContext) String() string {
	if len(fc) == 0 {
		return "[]"
	}

	var str string
	for _, f := range fc {
		str += f.String() + "\n"
	}
	return str
}
