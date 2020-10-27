package types

import (
	"encoding/hex"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"

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

	SecretLength                    = 64    // length for the secret in bytes
	HashLockLength                  = 64    // length for the hash lock in bytes
	MaxLengthForAddressOnOtherChain = 128   // maximum length for the address on other chains
	MinTimeLock                     = 50    // minimum time span for HTLC
	MaxTimeLock                     = 25480 // maximum time span for HTLC
)

var (
	_ sdk.Msg = &MsgCreateHTLC{}
	_ sdk.Msg = &MsgClaimHTLC{}
	_ sdk.Msg = &MsgRefundHTLC{}
)

// NewMsgCreateHTLC creates a new MsgCreateHTLC instance
func NewMsgCreateHTLC(
	sender sdk.AccAddress,
	to sdk.AccAddress,
	receiverOnOtherChain string,
	amount sdk.Coins,
	hashLock tmbytes.HexBytes,
	timestamp uint64,
	timeLock uint64,
) MsgCreateHTLC {
	return MsgCreateHTLC{
		Sender:               sender.String(),
		To:                   to.String(),
		ReceiverOnOtherChain: receiverOnOtherChain,
		Amount:               amount,
		HashLock:             hashLock.String(),
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
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if len(msg.To) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "recipient missing")
	}

	if len(msg.ReceiverOnOtherChain) > MaxLengthForAddressOnOtherChain {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "length of the receiver on other chain must be between [0,%d]", MaxLengthForAddressOnOtherChain)
	}

	if !msg.Amount.IsValid() || !msg.Amount.IsAllPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "the transferred amount must be valid")
	}

	_, err = hex.DecodeString(msg.HashLock)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidHashLock, "hash lock must be a hex encoded string")
	}
	if len(msg.HashLock) != HashLockLength {
		return sdkerrors.Wrapf(ErrInvalidHashLock, "length of the hash lock must be %d in bytes", HashLockLength)
	}

	if msg.TimeLock < MinTimeLock || msg.TimeLock > MaxTimeLock {
		return sdkerrors.Wrapf(ErrInvalidTimeLock, "the time lock must be between [%d,%d]", MinTimeLock, MaxTimeLock)
	}

	return nil
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
	sender sdk.AccAddress,
	hashLock tmbytes.HexBytes,
	secret tmbytes.HexBytes,
) MsgClaimHTLC {
	return MsgClaimHTLC{
		Sender:   sender.String(),
		HashLock: hashLock.String(),
		Secret:   secret.String(),
	}
}

// Route implements Msg
func (msg MsgClaimHTLC) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgClaimHTLC) Type() string { return TypeMsgClaimHTLC }

// ValidateBasic implements Msg.
func (msg MsgClaimHTLC) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	_, err = hex.DecodeString(msg.HashLock)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidHashLock, "hash lock must be a hex encoded string")
	}
	if len(msg.HashLock) != HashLockLength {
		return sdkerrors.Wrapf(ErrInvalidHashLock, "length of the hash lock must be %d in bytes", HashLockLength)
	}

	_, err = hex.DecodeString(msg.Secret)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidSecret, "secret must be a hex encoded string")
	}
	if len(msg.Secret) != SecretLength {
		return sdkerrors.Wrapf(ErrInvalidSecret, "length of the secret must be %d in bytes", SecretLength)
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

// -----------------------------------------------------------------------------

// NewMsgRefundHTLC constructs a new MsgRefundHTLC instance
func NewMsgRefundHTLC(
	sender sdk.AccAddress,
	hashLock tmbytes.HexBytes,
) MsgRefundHTLC {
	return MsgRefundHTLC{
		Sender:   sender.String(),
		HashLock: hashLock.String(),
	}
}

// Route implements Msg
func (msg MsgRefundHTLC) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgRefundHTLC) Type() string { return TypeMsgRefundHTLC }

// ValidateBasic implements Msg
func (msg MsgRefundHTLC) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	_, err = hex.DecodeString(msg.HashLock)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidHashLock, "hash lock must be a hex encoded string")
	}

	if len(msg.HashLock) != HashLockLength {
		return sdkerrors.Wrapf(ErrInvalidHashLock, "length of the hash lock must be %d in bytes", HashLockLength)
	}

	return nil
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
