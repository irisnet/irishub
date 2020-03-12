package types

import (
	"regexp"
	"strings"

	sdk "github.com/irisnet/irishub/types"
)

const (
	ModuleName = "oracle"
	MsgRoute   = ModuleName // route for oracle msg

	MaxLatestHistory    = 100
	MaxNameLen          = 70
	MaxAggregateFuncLen = 10
	MaxValueJsonPath    = 70
	MaxDescriptionLen   = 200

	TypeMsgCreateFeed = "create_feed" // type for MsgCreateFeed
	TypeMsgStartFeed  = "start_feed"  // type for MsgStartFeed
	TypeMsgPauseFeed  = "pause_feed"  // type for MsgPauseFeed
	TypeMsgEditFeed   = "edit_feed"   // type for MsgEditFeed

	DoNotModify = "do-not-modify"
)

var (
	_ sdk.Msg = MsgCreateFeed{}
	_ sdk.Msg = MsgStartFeed{}
	_ sdk.Msg = MsgPauseFeed{}
	_ sdk.Msg = MsgEditFeed{}

	// the feed/service name only accepts alphanumeric characters, _ and -
	regPlainText = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)
)

//______________________________________________________________________

// MsgCreateFeed - struct for create a feed
type MsgCreateFeed struct {
	FeedName          string           `json:"feed_name"`
	LatestHistory     uint64           `json:"latest_history"`
	Description       string           `json:"description"`
	Creator           sdk.AccAddress   `json:"creator"`
	ServiceName       string           `json:"service_name"`
	Providers         []sdk.AccAddress `json:"providers"`
	Input             string           `json:"input"`
	Timeout           int64            `json:"timeout"`
	ServiceFeeCap     sdk.Coins        `json:"service_fee_cap"`
	RepeatedFrequency uint64           `json:"repeated_frequency"`
	RepeatedTotal     int64            `json:"repeated_total"`
	AggregateFunc     string           `json:"aggregate_func"`
	ValueJsonPath     string           `json:"value_json_path"`
	ResponseThreshold uint16           `json:"response_threshold"`
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
	if err := ValidateFeedName(msg.FeedName); err != nil {
		return err
	}

	if err := ValidateDescription(msg.Description); err != nil {
		return err
	}

	if err := validateServiceName(msg.ServiceName); err != nil {
		return err
	}

	if err := ValidateLatestHistory(msg.LatestHistory); err != nil {
		return err
	}

	if err := validateTimeout(msg.Timeout, msg.RepeatedFrequency); err != nil {
		return err
	}

	if len(msg.Providers) == 0 {
		return ErrEmptyProviders(DefaultCodespace)
	}

	if err := ValidateAggregateFunc(msg.AggregateFunc); err != nil {
		return err
	}

	if err := ValidateValueJsonPath(msg.ValueJsonPath); err != nil {
		return err
	}

	if !msg.ServiceFeeCap.IsValidIrisAtto() {
		return ErrInvalidServiceFeeCap(DefaultCodespace, msg.ServiceFeeCap)
	}

	if err := ValidateCreator(msg.Creator); err != nil {
		return err
	}
	return validateResponseThreshold(msg.ResponseThreshold, len(msg.Providers))
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
	if len(msg.Creator) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "creator can not be empty")
	}
	return ValidateFeedName(msg.FeedName)
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
	if err := ValidateCreator(msg.Creator); err != nil {
		return err
	}
	return ValidateFeedName(msg.FeedName)
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
	Description       string           `json:"description"`
	LatestHistory     uint64           `json:"latest_history"`
	Providers         []sdk.AccAddress `json:"providers"`
	Timeout           int64            `json:"timeout"`
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
	if err := ValidateFeedName(msg.FeedName); err != nil {
		return err
	}

	if err := ValidateDescription(msg.Description); err != nil {
		return err
	}

	if msg.LatestHistory != 0 {
		if err := ValidateLatestHistory(msg.LatestHistory); err != nil {
			return err
		}
	}

	if msg.ServiceFeeCap != nil && !msg.ServiceFeeCap.IsValidIrisAtto() {
		return ErrInvalidServiceFeeCap(DefaultCodespace, msg.ServiceFeeCap)
	}

	if err := validateTimeout(msg.Timeout, msg.RepeatedFrequency); err != nil {
		return err
	}

	if msg.ResponseThreshold != 0 {
		if err := validateResponseThreshold(msg.ResponseThreshold, len(msg.Providers)); err != nil {
			return err
		}
	}

	return ValidateCreator(msg.Creator)
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

func ValidateFeedName(feedName string) sdk.Error {
	feedName = strings.TrimSpace(feedName)
	if len(feedName) == 0 || len(feedName) > MaxNameLen {
		return ErrInvalidFeedName(DefaultCodespace, feedName)
	}
	if !regPlainText.MatchString(feedName) {
		return ErrInvalidFeedName(DefaultCodespace, feedName)
	}
	return nil
}

func ValidateDescription(desc string) sdk.Error {
	desc = strings.TrimSpace(desc)
	if len(desc) > MaxDescriptionLen {
		return ErrInvalidDescription(DefaultCodespace, len(desc))
	}
	return nil
}

func ValidateAggregateFunc(aggregateFunc string) sdk.Error {
	aggregateFunc = strings.TrimSpace(aggregateFunc)
	if len(aggregateFunc) == 0 || len(aggregateFunc) > MaxAggregateFuncLen {
		return ErrInvalidAggregateFunc(DefaultCodespace, len(aggregateFunc))
	}
	if _, err := GetAggregateFunc(aggregateFunc); err != nil {
		return err
	}
	return nil
}

func ValidateValueJsonPath(valueJsonPath string) sdk.Error {
	valueJsonPath = strings.TrimSpace(valueJsonPath)
	if len(valueJsonPath) == 0 || len(valueJsonPath) > MaxValueJsonPath {
		return ErrInvalidValueJsonPath(DefaultCodespace, len(valueJsonPath))
	}
	return nil
}

func ValidateLatestHistory(latestHistory uint64) sdk.Error {
	if latestHistory < 1 || latestHistory > MaxLatestHistory {
		return ErrInvalidLatestHistory(DefaultCodespace)
	}
	return nil
}

func ValidateCreator(creator sdk.AccAddress) sdk.Error {
	if len(creator) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "fee creator can not be empty")
	}
	return nil
}

func validateServiceName(serviceName string) sdk.Error {
	serviceName = strings.TrimSpace(serviceName)
	if len(serviceName) == 0 || len(serviceName) > MaxNameLen {
		return ErrInvalidServiceName(DefaultCodespace, serviceName)
	}
	if !regPlainText.MatchString(serviceName) {
		return ErrInvalidServiceName(DefaultCodespace, serviceName)
	}
	return nil
}

func validateResponseThreshold(responseThreshold uint16, maxCnt int) sdk.Error {
	if (maxCnt != 0 && int(responseThreshold) > maxCnt) || responseThreshold < 1 {
		return ErrInvalidResponseThreshold(DefaultCodespace, maxCnt)
	}
	return nil
}

func validateTimeout(timeout int64, frequency uint64) sdk.Error {
	if frequency < uint64(timeout) {
		return ErrInvalidTimeout(DefaultCodespace, timeout, frequency)
	}
	return nil
}

func IsModified(target string) bool {
	target = strings.TrimSpace(target)
	return target != DoNotModify
}
