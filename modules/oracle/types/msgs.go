package types

import (
	"regexp"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateFeed = "create_feed" // type for MsgCreateFeed
	TypeMsgStartFeed  = "start_feed"  // type for MsgStartFeed
	TypeMsgPauseFeed  = "pause_feed"  // type for MsgPauseFeed
	TypeMsgEditFeed   = "edit_feed"   // type for MsgEditFeed

	DoNotModify = "do-not-modify"
)

var (
	_ sdk.Msg = &MsgCreateFeed{}
	_ sdk.Msg = &MsgStartFeed{}
	_ sdk.Msg = &MsgPauseFeed{}
	_ sdk.Msg = &MsgEditFeed{}

	// the feed/service name only accepts alphanumeric characters, _ and -
	regPlainText = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)
)

// ______________________________________________________________________

// Route implements Msg.
func (msg MsgCreateFeed) Route() string {
	return RouterKey
}

// Type implements Msg.
func (msg MsgCreateFeed) Type() string {
	return TypeMsgCreateFeed
}

// ValidateBasic implements Msg.
func (msg MsgCreateFeed) ValidateBasic() error {
	msg = msg.Normalize()
	if err := ValidateFeedName(msg.FeedName); err != nil {
		return err
	}

	if err := ValidateDescription(msg.Description); err != nil {
		return err
	}

	if err := ValidateServiceName(msg.ServiceName); err != nil {
		return err
	}

	if err := ValidateLatestHistory(msg.LatestHistory); err != nil {
		return err
	}

	if err := ValidateTimeout(msg.Timeout, msg.RepeatedFrequency); err != nil {
		return err
	}
	if len(msg.Providers) == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "providers missing")
	}

	if err := ValidateAggregateFunc(msg.AggregateFunc); err != nil {
		return err
	}

	if err := ValidateValueJSONPath(msg.ValueJsonPath); err != nil {
		return err
	}

	if !msg.ServiceFeeCap.IsValid() {
		return sdkerrors.Wrapf(ErrInvalidServiceFeeCap, msg.ServiceFeeCap.String())
	}

	if err := ValidateCreator(msg.Creator); err != nil {
		return err
	}

	return ValidateResponseThreshold(msg.ResponseThreshold, len(msg.Providers))
}

// Normalize return a string with spaces removed and lowercase
func (msg MsgCreateFeed) Normalize() MsgCreateFeed {
	msg.FeedName = strings.TrimSpace(msg.FeedName)
	msg.ServiceName = strings.TrimSpace(msg.ServiceName)
	msg.Input = strings.TrimSpace(msg.Input)
	msg.AggregateFunc = strings.TrimSpace(msg.AggregateFunc)
	msg.ValueJsonPath = strings.TrimSpace(msg.ValueJsonPath)
	return msg
}

// GetSignBytes implements Msg.
func (msg MsgCreateFeed) GetSignBytes() []byte {
	if len(msg.Providers) == 0 {
		msg.Providers = nil
	}
	if msg.ServiceFeeCap.Empty() {
		msg.ServiceFeeCap = nil
	}
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements Msg.
func (msg MsgCreateFeed) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// _____________________________________________________________________

// Route implements Msg.
func (msg MsgStartFeed) Route() string {
	return RouterKey
}

// Type implements Msg.
func (msg MsgStartFeed) Type() string {
	return TypeMsgStartFeed
}

// ValidateBasic implements Msg.
func (msg MsgStartFeed) ValidateBasic() error {
	msg = msg.Normalize()
	if err := ValidateCreator(msg.Creator); err != nil {
		return err
	}
	return ValidateFeedName(msg.FeedName)
}

// Normalize return a string with spaces removed and lowercase
func (msg MsgStartFeed) Normalize() MsgStartFeed {
	msg.FeedName = strings.TrimSpace(msg.FeedName)
	return msg
}

// GetSignBytes implements Msg.
func (msg MsgStartFeed) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements Msg.
func (msg MsgStartFeed) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ______________________________________________________________________

// Route implements Msg.
func (msg MsgPauseFeed) Route() string {
	return RouterKey
}

// Type implements Msg.
func (msg MsgPauseFeed) Type() string {
	return TypeMsgPauseFeed
}

// ValidateBasic implements Msg.
func (msg MsgPauseFeed) ValidateBasic() error {
	msg = msg.Normalize()
	if err := ValidateCreator(msg.Creator); err != nil {
		return err
	}
	return ValidateFeedName(msg.FeedName)
}

// Normalize return a string with spaces removed and lowercase
func (msg MsgPauseFeed) Normalize() MsgPauseFeed {
	msg.FeedName = strings.TrimSpace(msg.FeedName)
	return msg
}

// GetSignBytes implements Msg.
func (msg MsgPauseFeed) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements Msg.
func (msg MsgPauseFeed) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ______________________________________________________________________

// Route implements Msg.
func (msg MsgEditFeed) Route() string {
	return RouterKey
}

// Type implements Msg.
func (msg MsgEditFeed) Type() string {
	return TypeMsgEditFeed
}

// ValidateBasic implements Msg.
func (msg MsgEditFeed) ValidateBasic() error {
	msg = msg.Normalize()
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

	if err := ValidateServiceFeeCap(msg.ServiceFeeCap); err != nil {
		return err
	}

	if msg.Timeout != 0 && msg.RepeatedFrequency != 0 {
		if err := ValidateTimeout(msg.Timeout, msg.RepeatedFrequency); err != nil {
			return err
		}
	}

	if msg.ResponseThreshold != 0 {
		if err := ValidateResponseThreshold(msg.ResponseThreshold, len(msg.Providers)); err != nil {
			return err
		}
	}
	return ValidateCreator(msg.Creator)
}

// Normalize return a string with spaces removed and lowercase
func (msg MsgEditFeed) Normalize() MsgEditFeed {
	msg.FeedName = strings.TrimSpace(msg.FeedName)
	return msg
}

// GetSignBytes implements Msg.
func (msg MsgEditFeed) GetSignBytes() []byte {
	if len(msg.Providers) == 0 {
		msg.Providers = nil
	}
	if msg.ServiceFeeCap.Empty() {
		msg.ServiceFeeCap = nil
	}
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements Msg.
func (msg MsgEditFeed) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}
