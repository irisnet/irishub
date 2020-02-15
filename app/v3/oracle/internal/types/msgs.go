package types

import (
	"strings"

	sdk "github.com/irisnet/irishub/types"
)

const (
	ModuleName = "oracle"
	MsgRoute   = ModuleName // route for oracle msg

	LatestHistory = 100

	TypeMsgCreateFeed = "create_feed" // type for MsgCreateFeed
	TypeMsgStartFeed  = "start_feed"  // type for MsgStartFeed
	TypeMsgPauseFeed  = "pause_feed"  // type for MsgPauseFeed
	TypeMsgEditFeed   = "edit_feed"   // type for MsgEditFeed
)

var (
	_ sdk.Msg = MsgCreateFeed{}
	_ sdk.Msg = MsgStartFeed{}
	_ sdk.Msg = MsgPauseFeed{}
	_ sdk.Msg = MsgEditFeed{}
)

//______________________________________________________________________

// MsgCreateFeed - struct for create a feed
type MsgCreateFeed struct {
	FeedName          string           `json:"feed_name"`
	AggregateFunc     string           `json:"aggregate_func"`
	ValueJsonPath     string           `json:"value_json_path"`
	LatestHistory     uint64           `json:"latest_history"`
	Description       string           `json:"description"`
	ServiceName       string           `json:"service_name"`
	Providers         []sdk.AccAddress `json:"providers"`
	Input             string           `json:"input"`
	Timeout           int64            `json:"timeout"`
	ServiceFeeCap     sdk.Coins        `json:"service_fee_cap"`
	RepeatedFrequency uint64           `json:"repeated_frequency"`
	RepeatedTotal     int64            `json:"repeated_total"`
	ResponseThreshold uint16           `json:"response_threshold"`
	Creator           sdk.AccAddress   `json:"creator"`
}

// Route implements Msg.
func (msg MsgCreateFeed) Route() string {
	return MsgRoute
}

// Type implements Msg.
func (msg MsgCreateFeed) Type() string {
	return TypeMsgCreateFeed
}

// ValidateBasic implements Msg.
func (msg MsgCreateFeed) ValidateBasic() sdk.Error {
	feedKey := strings.TrimSpace(msg.ServiceName)
	if len(feedKey) == 0 {
		return ErrEmptyFeedName(DefaultCodespace)
	}

	serviceName := strings.TrimSpace(msg.ServiceName)
	if len(serviceName) == 0 {
		return ErrEmptyServiceName(DefaultCodespace)
	}

	if msg.LatestHistory < 1 || msg.LatestHistory > LatestHistory {
		return ErrInvalidLatestHistory(DefaultCodespace)
	}

	if len(msg.Providers) == 0 {
		return ErrEmptyProviders(DefaultCodespace)
	}

	aggregateFunc := strings.TrimSpace(msg.AggregateFunc)
	if len(aggregateFunc) == 0 {
		return ErrEmptyEmptyAggregateFunc(DefaultCodespace)
	}

	valueJsonPath := strings.TrimSpace(msg.ValueJsonPath)
	if len(valueJsonPath) == 0 {
		return ErrEmptyValueJsonPath(DefaultCodespace)
	}

	if !msg.ServiceFeeCap.IsValidIrisAtto() {
		return ErrInvalidServiceFeeCap(DefaultCodespace, msg.ServiceFeeCap)
	}

	if int(msg.ResponseThreshold) > len(msg.Providers) || msg.ResponseThreshold < 1 {
		return ErrInvalidResponseThreshold(DefaultCodespace, len(msg.Providers))
	}

	if len(msg.Creator) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "fee creator can not be empty")
	}
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgCreateFeed) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg.
func (msg MsgCreateFeed) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

//______________________________________________________________________

// MsgStartFeed - struct for start a feed
type MsgStartFeed struct {
	FeedName string         `json:"feed_name"`
	Creator  sdk.AccAddress `json:"creator"`
}

// Route implements Msg.
func (msg MsgStartFeed) Route() string {
	return MsgRoute
}

// Type implements Msg.
func (msg MsgStartFeed) Type() string {
	return TypeMsgStartFeed
}

// ValidateBasic implements Msg.
func (msg MsgStartFeed) ValidateBasic() sdk.Error {
	feedName := strings.TrimSpace(msg.FeedName)
	if len(feedName) == 0 {
		return ErrEmptyFeedName(DefaultCodespace)
	}

	if len(msg.Creator) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "fee creator can not be empty")
	}
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgStartFeed) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg.
func (msg MsgStartFeed) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

//______________________________________________________________________

// MsgPauseFeed - struct for stop a started feed
type MsgPauseFeed struct {
	FeedName string         `json:"feed_name"`
	Creator  sdk.AccAddress `json:"creator"`
}

// Route implements Msg.
func (msg MsgPauseFeed) Route() string {
	return MsgRoute
}

// Type implements Msg.
func (msg MsgPauseFeed) Type() string {
	return TypeMsgPauseFeed
}

// ValidateBasic implements Msg.
func (msg MsgPauseFeed) ValidateBasic() sdk.Error {
	feedKey := strings.TrimSpace(msg.FeedName)
	if len(feedKey) == 0 {
		return ErrEmptyFeedName(DefaultCodespace)
	}

	if len(msg.Creator) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "fee creator can not be empty")
	}
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgPauseFeed) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg.
func (msg MsgPauseFeed) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

//______________________________________________________________________

// MsgEditFeed - struct for edit a existed feed
type MsgEditFeed struct {
	FeedName          string           `json:"feed_name"`
	LatestHistory     uint64           `json:"latest_history"`
	Providers         []sdk.AccAddress `json:"providers"`
	ServiceFeeCap     sdk.Coins        `json:"service_fee_cap"`
	RepeatedFrequency uint64           `json:"repeated_frequency"`
	RepeatedTotal     int64            `json:"repeated_total"`
	ResponseThreshold uint16           `json:"response_threshold"`
	Creator           sdk.AccAddress   `json:"creator"`
}

// Route implements Msg.
func (msg MsgEditFeed) Route() string {
	return MsgRoute
}

// Type implements Msg.
func (msg MsgEditFeed) Type() string {
	return TypeMsgEditFeed
}

// ValidateBasic implements Msg.
func (msg MsgEditFeed) ValidateBasic() sdk.Error {
	feedName := strings.TrimSpace(msg.FeedName)
	if len(feedName) == 0 {
		return ErrEmptyFeedName(DefaultCodespace)
	}

	if len(msg.Creator) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "fee creator can not be empty")
	}

	if !msg.ServiceFeeCap.IsValidIrisAtto() {
		return ErrInvalidServiceFeeCap(DefaultCodespace, msg.ServiceFeeCap)
	}

	if int(msg.ResponseThreshold) > len(msg.Providers) || msg.ResponseThreshold < 1 {
		return ErrInvalidResponseThreshold(DefaultCodespace, len(msg.Providers))
	}
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgEditFeed) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg.
func (msg MsgEditFeed) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}
