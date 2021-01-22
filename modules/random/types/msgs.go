package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgRequestRandom = "request_random" // type for MsgRequestRandom

	DefaultBlockInterval = uint64(10) // DefaultBlockInterval is the default block interval
)

var _ sdk.Msg = &MsgRequestRandom{}

// NewMsgRequestRandom constructs a MsgRequestRandom
func NewMsgRequestRandom(
	consumer string,
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
	msg = msg.Normalize()
	if _, err := sdk.AccAddressFromBech32(msg.Consumer); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid consumer address (%s)", err)
	}
	return ValidateServiceFeeCap(msg.ServiceFeeCap)
}

// Normalize return a string with spaces removed and lowercase
func (msg MsgRequestRandom) Normalize() MsgRequestRandom {
	return msg
}

// GetSignBytes implements Msg.
func (msg MsgRequestRandom) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg.
func (msg MsgRequestRandom) GetSigners() []sdk.AccAddress {
	consumer, err := sdk.AccAddressFromBech32(msg.Consumer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{consumer}
}
