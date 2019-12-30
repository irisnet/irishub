package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateHTLC = "create_htlc" // type for MsgCreateHTLC
	TypeMsgClaimHTLC  = "claim_htlc"  // type for MsgClaimHTLC
	TypeMsgRefundHTLC = "refund_htlc" // type for MsgRefundHTLC

	SecretLength                    = 32    // the length for the secret
	HashLockLength                  = 32    // the length for the hash lock
	MaxLengthForAddressOnOtherChain = 128   // maximal length for the address on other chains
	MinTimeLock                     = 50    // minimal time span for HTLC
	MaxTimeLock                     = 25480 // maximal time span for HTLC
)

var (
	_ sdk.Msg = &MsgCreateHTLC{}
	_ sdk.Msg = &MsgClaimHTLC{}
	_ sdk.Msg = &MsgRefundHTLC{}
)

// MsgCreateHTLC represents a msg for creating an HTLC
type MsgCreateHTLC struct {
	Sender               sdk.AccAddress `json:"sender" yaml:"sender"`                                   // the initiator address
	To                   sdk.AccAddress `json:"to" yaml:"to"`                                           // the destination address
	ReceiverOnOtherChain string         `json:"receiver_on_other_chain" yaml:"receiver_on_other_chain"` // the claim receiving address on the other chain
	Amount               sdk.Coins      `json:"amount" yaml:"amount"`                                   // the amount to be transferred
	HashLock             HTLCHashLock   `json:"hash_lock" yaml:"hash_lock"`                             // the hash lock generated from secret (and timestamp if provided)
	Timestamp            uint64         `json:"timestamp" yaml:"timestamp"`                             // if provided, used to generate the hash lock together with secret
	TimeLock             uint64         `json:"time_lock" yaml:"time_lock"`                             // the time span after which the HTLC will expire
}

// NewMsgCreateHTLC constructs a MsgCreateHTLC
func NewMsgCreateHTLC(
	sender sdk.AccAddress,
	to sdk.AccAddress,
	receiverOnOtherChain string,
	amount sdk.Coins,
	hashLock HTLCHashLock,
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

// Route implements Msg.
func (msg MsgCreateHTLC) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgCreateHTLC) Type() string { return TypeMsgCreateHTLC }

// ValidateBasic implements Msg.
func (msg MsgCreateHTLC) ValidateBasic() error {
	if len(msg.Sender) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "sender address missing")
	}
	if len(msg.To) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "to address missing")
	}
	if len(msg.ReceiverOnOtherChain) > MaxLengthForAddressOnOtherChain {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid length of the receiver on other chain; got: %d, max: %d", len(msg.ReceiverOnOtherChain), MaxLengthForAddressOnOtherChain)
	}
	if !msg.Amount.IsValid() || !msg.Amount.IsAllPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}
	if len(msg.HashLock) != HashLockLength {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid hash lock length; got: %d, must: %d", len(msg.HashLock), HashLockLength)
	}
	if msg.TimeLock < MinTimeLock || msg.TimeLock > MaxTimeLock {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "the time lock must be between [%d,%d]", MinTimeLock, MaxTimeLock)
	}
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgCreateHTLC) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg.
func (msg MsgCreateHTLC) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// -----------------------------------------------------------------------------

// MsgClaimHTLC represents a msg for claiming an HTLC
type MsgClaimHTLC struct {
	Sender   sdk.AccAddress `json:"sender" yaml:"sender"`       // the initiator address
	HashLock HTLCHashLock   `json:"hash_lock" yaml:"hash_lock"` // the hash lock identifying the HTLC to be claimed
	Secret   HTLCSecret     `json:"secret" yaml:"secret"`       // the secret with which to claim
}

// NewMsgClaimHTLC constructs a MsgClaimHTLC
func NewMsgClaimHTLC(
	sender sdk.AccAddress,
	hashLock HTLCHashLock,
	secret HTLCSecret,
) MsgClaimHTLC {
	return MsgClaimHTLC{
		Sender:   sender,
		HashLock: hashLock,
		Secret:   secret,
	}
}

// Route implements Msg.
func (msg MsgClaimHTLC) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgClaimHTLC) Type() string { return TypeMsgClaimHTLC }

// ValidateBasic implements Msg.
func (msg MsgClaimHTLC) ValidateBasic() error {
	if len(msg.Sender) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "sender address missing")
	}
	if len(msg.HashLock) != HashLockLength {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid hash lock length; got: %d, must: %d", len(msg.HashLock), HashLockLength)
	}
	if len(msg.Secret) != SecretLength {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid secret length; got: %d, must: %d", len(msg.Secret), SecretLength)
	}
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgClaimHTLC) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg.
func (msg MsgClaimHTLC) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// -----------------------------------------------------------------------------

// MsgRefundHTLC represents a msg for refund an HTLC
type MsgRefundHTLC struct {
	Sender   sdk.AccAddress `json:"sender" yaml:"sender"`       // the initiator address
	HashLock HTLCHashLock   `json:"hash_lock" yaml:"hash_lock"` // the hash lock identifying the HTLC to be refunded
}

// NewMsgClaimHTLC constructs a MsgClaimHTLC
func NewMsgRefundHTLC(
	sender sdk.AccAddress,
	hashLock HTLCHashLock,
) MsgRefundHTLC {
	return MsgRefundHTLC{
		Sender:   sender,
		HashLock: hashLock,
	}
}

// Route implements Msg.
func (msg MsgRefundHTLC) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgRefundHTLC) Type() string { return TypeMsgRefundHTLC }

// ValidateBasic implements Msg.
func (msg MsgRefundHTLC) ValidateBasic() error {
	if len(msg.Sender) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "sender address missing")
	}
	if len(msg.HashLock) != HashLockLength {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid hash lock length; got: %d, must: %d", len(msg.HashLock), HashLockLength)
	}
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgRefundHTLC) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg.
func (msg MsgRefundHTLC) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
