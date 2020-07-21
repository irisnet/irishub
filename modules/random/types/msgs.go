package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgRequestRandom = "request_rand" // type for MsgRequestRandom

	DefaultBlockInterval = uint64(10) // DefaultBlockInterval is the default block interval
)

var _ sdk.Msg = &MsgRequestRandom{}

// NewMsgRequestRandom constructs a MsgRequestRandom
func NewMsgRequestRandom(
	consumer sdk.AccAddress,
	blockInterval uint64,
	oracle bool,
	serviceFeeCap sdk.Coins,
) *MsgRequestRandom {
	return &MsgRequestRandom{
		Consumer:      consumer,
		BlockInterval: blockInterval,
		Oracle:        oracle,
		ServiceFeeCap: serviceFeeCap,
	}
}

// Route implements Msg.
func (msg MsgRequestRandom) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgRequestRandom) Type() string { return TypeMsgRequestRandom }

// ValidateBasic implements Msg.
func (msg MsgRequestRandom) ValidateBasic() error {
	if len(msg.Consumer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "the consumer address must be specified")
	}
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgRequestRandom) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg.
func (msg MsgRequestRandom) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Consumer}
}
