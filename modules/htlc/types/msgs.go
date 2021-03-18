package types

import (
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
)

// NewMsgCreateHTLC creates a new MsgCreateHTLC instance
func NewMsgCreateHTLC(
	sender string,
	to string,
	receiverOnOtherChain string,
	senderOnOtherChain string,
	amount sdk.Coins,
	hashLock string,
	timestamp uint64,
	timeLock uint64,
	transfer bool,
) MsgCreateHTLC {
	return MsgCreateHTLC{
		Sender:               sender,
		To:                   to,
		ReceiverOnOtherChain: receiverOnOtherChain,
		SenderOnOtherChain:   senderOnOtherChain,
		Amount:               amount,
		HashLock:             hashLock,
		Timestamp:            timestamp,
		TimeLock:             timeLock,
		Transfer:             transfer,
	}
}

// Route implements Msg
func (msg MsgCreateHTLC) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgCreateHTLC) Type() string { return TypeMsgCreateHTLC }

// ValidateBasic implements Msg
func (msg MsgCreateHTLC) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.To); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid recipient address (%s)", err)
	}

	if err := ValidateReceiverOnOtherChain(msg.ReceiverOnOtherChain); err != nil {
		return err
	}

	if err := ValidateSenderOnOtherChain(msg.SenderOnOtherChain); err != nil {
		return err
	}

	if err := ValidateAmount(msg.Transfer, msg.Amount); err != nil {
		return err
	}

	if err := ValidateID(msg.HashLock); err != nil {
		return err
	}

	return ValidateTimeLock(msg.TimeLock)
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
	id string,
	secret string,
) MsgClaimHTLC {
	return MsgClaimHTLC{
		Sender: sender,
		Id:     id,
		Secret: secret,
	}
}

// Route implements Msg
func (msg MsgClaimHTLC) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgClaimHTLC) Type() string { return TypeMsgClaimHTLC }

// ValidateBasic implements Msg.
func (msg MsgClaimHTLC) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if err := ValidateID(msg.Id); err != nil {
		return err
	}

	if err := ValidateSecret(msg.Secret); err != nil {
		return err
	}
	return nil
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
