package types

import (
	"strings"

	sdk "github.com/irisnet/irishub/types"
)

const (
	MsgRoute = "oracle" // route for oracle msg

	MaxHistory = 10

	TypeMsgCreateFeed = "create_feed" // type for MsgCreateFeed
	TypeMsgStartFeed  = "start_feed"  // type for MsgStartFeed
	TypeMsgStopFeed   = "stop_feed"   // type for MsgStopFeed
	TypeMsgEditFeed   = "edit_feed"   // type for MsgEditFeed
)

var (
	_ sdk.Msg = MsgCreateFeed{}
	_ sdk.Msg = MsgStartFeed{}
	_ sdk.Msg = MsgStopFeed{}
	_ sdk.Msg = MsgEditFeed{}
)

//______________________________________________________________________

// MsgCreateFeed - struct for create a feed
type MsgCreateFeed struct {
	FeedKey       string           `json:"feed_key"`
	ServiceName   string           `json:"service_name"`
	ResHandler    string           `json:"res_handler"`
	MaxHistory    uint64           `json:"max_history"`
	Providers     []sdk.AccAddress `json:"providers"`
	Input         string           `json:"input"`
	MaxServiceFee sdk.Coins        `json:"max_service_fee"`
	Frequency     uint64           `json:"frequency"`
	MaxCount      int64            `json:"max_count"`
	ResThreshold  uint16           `json:"res_threshold"`
	Sender        sdk.AccAddress   `json:"sender"`
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
	feedKey := strings.TrimSpace(msg.FeedKey)
	if len(feedKey) == 0 {
		return ErrEmptyFeedKey(DefaultCodespace)
	}

	serviceName := strings.TrimSpace(msg.ServiceName)
	if len(serviceName) == 0 {
		return ErrEmptyServiceName(DefaultCodespace)
	}

	if err := validateMaxHistory(msg.MaxHistory); err != nil {
		return err
	}

	if len(msg.Providers) == 0 {
		return ErrEmptyProviders(DefaultCodespace)
	}

	if len(msg.Input) > 0 {
		//TODO
	}

	if !msg.MaxServiceFee.IsValidIrisAtto() {
		return ErrInvalidMaxServiceFee(DefaultCodespace, msg.MaxServiceFee)
	}

	if len(msg.Sender) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "sender can not be empty")
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
	return []sdk.AccAddress{msg.Sender}
}

//______________________________________________________________________

// MsgStartFeed - struct for start a feed
type MsgStartFeed struct {
	FeedKey string         `json:"feed_key"`
	Sender  sdk.AccAddress `json:"sender"`
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
	feedKey := strings.TrimSpace(msg.FeedKey)
	if len(feedKey) == 0 {
		return ErrEmptyFeedKey(DefaultCodespace)
	}

	if len(msg.Sender) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "sender can not be empty")
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
	return []sdk.AccAddress{msg.Sender}
}

//______________________________________________________________________

// MsgStopFeed - struct for stop a started feed
type MsgStopFeed struct {
	FeedKey string         `json:"feed_key"`
	Sender  sdk.AccAddress `json:"sender"`
}

// Route implements Msg.
func (msg MsgStopFeed) Route() string {
	return MsgRoute
}

// Type implements Msg.
func (msg MsgStopFeed) Type() string {
	return TypeMsgStopFeed
}

// ValidateBasic implements Msg.
func (msg MsgStopFeed) ValidateBasic() sdk.Error {
	feedKey := strings.TrimSpace(msg.FeedKey)
	if len(feedKey) == 0 {
		return ErrEmptyFeedKey(DefaultCodespace)
	}

	if len(msg.Sender) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "sender can not be empty")
	}
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgStopFeed) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg.
func (msg MsgStopFeed) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

//______________________________________________________________________

// MsgEditFeed - struct for edit a existed feed
type MsgEditFeed struct {
	FeedKey       string           `json:"feed_key"`
	ResHandler    string           `json:"res_handler"`
	MaxHistory    uint64           `json:"max_history"`
	Providers     []sdk.AccAddress `json:"providers"`
	MaxServiceFee sdk.Coins        `json:"max_service_fee"`
	Frequency     uint64           `json:"frequency"`
	MaxCount      int64            `json:"max_count"`
	ResThreshold  uint16           `json:"res_threshold"`
	Sender        sdk.AccAddress   `json:"sender"`
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
	feedKey := strings.TrimSpace(msg.FeedKey)
	if len(feedKey) == 0 {
		return ErrEmptyFeedKey(DefaultCodespace)
	}

	if len(msg.Sender) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "sender can not be empty")
	}

	if !msg.MaxServiceFee.IsValidIrisAtto() {
		return ErrInvalidMaxServiceFee(DefaultCodespace, msg.MaxServiceFee)
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
	return []sdk.AccAddress{msg.Sender}
}
