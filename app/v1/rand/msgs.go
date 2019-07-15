package rand

import (
	sdk "github.com/irisnet/irishub/types"
)

const (
	// MsgRoute identifies transaction types
	MsgRoute      = "rand"
	BlockNumAfter = 10 // block interval after which the requested random number will be generated
)

var _ sdk.Msg = &MsgRequestRand{}

// MsgRequestRand represents a msg for requesting a random number
type MsgRequestRand struct {
	Consumer sdk.AccAddress `json:"consumer"`
}

// NewMsgRequestRand constructs a MsgRequestRand
func NewMsgRequestRand(consumer sdk.AccAddress) MsgRequestRand {
	return MsgRequestRand{
		Consumer: consumer,
	}
}

// Implements Msg.
func (msg MsgRequestRand) Route() string { return MsgRoute }

// Implements Msg.
func (msg MsgRequestRand) Type() string { return "request_rand" }

// Implements Msg.
func (msg MsgRequestRand) ValidateBasic() sdk.Error {
	if len(msg.Consumer) == 0 {
		return ErrInvalidConsumer(DefaultCodespace, "the consumer address must be specified")
	}

	return nil
}

// Implements Msg.
func (msg MsgRequestRand) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgRequestRand) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Consumer}
}
