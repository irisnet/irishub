package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// TypeMsgCreateHTLC is the type for MsgCreateHTLC
	TypeMsgCreateHTLC = "create_htlc"

	// TypeMsgClaimHTLC is the type for MsgClaimHTLC
	TypeMsgClaimHTLC = "claim_htlc"

	// TypeMsgRefundHTLC is the type for MsgRefundHTLC
	TypeMsgRefundHTLC = "refund_htlc"
)

var (
	_ sdk.Msg = &MsgCreateHTLC{}
	_ sdk.Msg = &MsgClaimHTLC{}
	_ sdk.Msg = &MsgRefundHTLC{}
)

// NewMsgCreateHTLC creates a new MsgCreateHTLC instance
func NewMsgCreateHTLC(
	sender string,
	to string,
	receiverOnOtherChain string,
	amount sdk.Coins,
	hashLock string,
	timestamp uint64,
	timeLock uint64,
) MsgCreateHTLC {
	return MsgCreateHTLC{
		Sender:               sender,
		To:                   to,
		ReceiverOnOtherChain: receiverOnOtherChain,
		Amount:               amount,
		HashLock:             hashLock,
		Timestamp:            timestamp,
		TimeLock:             timeLock,
	}
}

// Route implements Msg
func (msg MsgCreateHTLC) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgCreateHTLC) Type() string { return TypeMsgCreateHTLC }

// ValidateBasic implements Msg
func (msg MsgCreateHTLC) ValidateBasic() error {
	msg = msg.Normalize()
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.To); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid recipient address (%s)", err)
	}

	if err := ValidateReceiverOnOtherChain(msg.ReceiverOnOtherChain); err != nil {
		return err
	}

	if err := ValidateAmount(msg.Amount); err != nil {
		return err
	}

	if err := ValidateHashLock(msg.HashLock); err != nil {
		return err
	}

	if err := ValidateTimeLock(msg.TimeLock); err != nil {
		return err
	}
	return nil
}

// Normalize return a string with spaces removed and lowercase
func (msg MsgCreateHTLC) Normalize() MsgCreateHTLC {
	msg.HashLock = strings.TrimSpace(msg.HashLock)
	msg.ReceiverOnOtherChain = strings.TrimSpace(msg.ReceiverOnOtherChain)
	return msg
}

// GetSignBytes implements Msg
func (msg MsgCreateHTLC) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgCreateHTLC) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// -----------------------------------------------------------------------------

// NewMsgClaimHTLC constructs a new MsgClaimHTLC instance
func NewMsgClaimHTLC(
	sender string,
	hashLock string,
	secret string,
) MsgClaimHTLC {
	return MsgClaimHTLC{
		Sender:   sender,
		HashLock: hashLock,
		Secret:   secret,
	}
}

// Route implements Msg
func (msg MsgClaimHTLC) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgClaimHTLC) Type() string { return TypeMsgClaimHTLC }

// ValidateBasic implements Msg.
func (msg MsgClaimHTLC) ValidateBasic() error {
	msg = msg.Normalize()
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if err := ValidateHashLock(msg.HashLock); err != nil {
		return err
	}

	if err := ValidateSecret(msg.Secret); err != nil {
		return err
	}
	return nil
}

// Normalize return a string with spaces removed and lowercase
func (msg MsgClaimHTLC) Normalize() MsgClaimHTLC {
	msg.HashLock = strings.TrimSpace(msg.HashLock)
	msg.Secret = strings.TrimSpace(msg.Secret)
	return msg
}

// GetSignBytes implements Msg
func (msg MsgClaimHTLC) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgClaimHTLC) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// -----------------------------------------------------------------------------

// NewMsgRefundHTLC constructs a new MsgRefundHTLC instance
func NewMsgRefundHTLC(
	sender string,
	hashLock string,
) MsgRefundHTLC {
	return MsgRefundHTLC{
		Sender:   sender,
		HashLock: hashLock,
	}
}

// Route implements Msg
func (msg MsgRefundHTLC) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgRefundHTLC) Type() string { return TypeMsgRefundHTLC }

// ValidateBasic implements Msg
func (msg MsgRefundHTLC) ValidateBasic() error {
	msg.Normalize()
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if err := ValidateHashLock(msg.HashLock); err != nil {
		return err
	}
	return nil
}

// Normalize return a string with spaces removed and lowercase
func (msg MsgRefundHTLC) Normalize() MsgRefundHTLC {
	msg.HashLock = strings.TrimSpace(msg.HashLock)
	return msg
}

// GetSignBytes implements Msg
func (msg MsgRefundHTLC) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgRefundHTLC) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}
