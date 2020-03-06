package types

import (
	sdk "github.com/irisnet/irishub/types"
)

const (
	ModuleName = "rand"

	// MsgRoute identifies transaction types
	MsgRoute = "rand"

	// DefaultBlockInterval is the default block interval
	DefaultBlockInterval = uint64(10)
)

var _ sdk.Msg = &MsgRequestRand{}

// MsgRequestRand represents a msg for requesting a random number
type MsgRequestRand struct {
	Consumer      sdk.AccAddress `json:"consumer"`       // request address
	BlockInterval uint64         `json:"block_interval"` // block interval after which the requested random number will be generated
	Oracle        bool           `json:"oracle"`         // oracle method
}

// NewMsgRequestRand constructs a MsgRequestRand
func NewMsgRequestRand(
	consumer sdk.AccAddress,
	blockInterval uint64,
	oracle bool,
) MsgRequestRand {
	return MsgRequestRand{
		Consumer:      consumer,
		BlockInterval: blockInterval,
		Oracle:        oracle,
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
