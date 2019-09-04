package types

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

const (
	// MsgRoute identifies transaction types
	MsgRoute = "htlc"

	// type for MsgCreateHTLC
	TypeMsgCreateHTLC = "create_htlc"

	SecretLength                    = 32    // the length for secret
	MaxLengthForAddressOnOtherChain = 32    // maximal length in bytes for the address on other chains
	DecimalNumForInAmount           = 8     // the default decimal number for InAmount
	MinTimeLock                     = 50    // minimal time span for HTLC
	MaxTimeLock                     = 25480 // maximal time span for HTLC
)

var _ sdk.Msg = &MsgCreateHTLC{}

// MsgCreateHTLC represents a msg for creating a HTLC
type MsgCreateHTLC struct {
	Sender               sdk.AccAddress `json:"sender"`                  // the initiator address
	Receiver             sdk.AccAddress `json:"receiver"`                // the recipient address
	ReceiverOnOtherChain []byte         `json:"receiver_on_other_chain"` // the recipient address on other chain
	OutAmount            sdk.Coin       `json:"out_amount"`              // the amount to be transferred
	InAmount             uint64         `json:"in_amount"`               // expected amount to be received from another HTLC
	SecretHashLock       string         `json:"secret_hash_lock"`        // the hash lock generated from secret and timestamp
	Timestamp            uint64         `json:"timestamp"`               // the time used to generate the hash lock together with secret
	TimeLock             uint64         `json:"time_lock"`               // the time span after which the HTLC will expire
}

// NewMsgCreateHTLC constructs a MsgCreateHTLC
func NewMsgCreateHTLC(sender sdk.AccAddress, receiver sdk.AccAddress, receiverOnOtherChain []byte, outAmount sdk.Coin, inAmount uint64, secretHashLock string, timestamp uint64, timeLock uint64) MsgCreateHTLC {
	return MsgCreateHTLC{
		Sender:               sender,
		Receiver:             receiver,
		ReceiverOnOtherChain: receiverOnOtherChain,
		OutAmount:            outAmount,
		InAmount:             inAmount,
		SecretHashLock:       secretHashLock,
		Timestamp:            timestamp,
		TimeLock:             timeLock,
	}
}

// Implements Msg.
func (msg MsgCreateHTLC) Route() string { return MsgRoute }

// Implements Msg.
func (msg MsgCreateHTLC) Type() string { return TypeMsgCreateHTLC }

// Implements Msg.
func (msg MsgCreateHTLC) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "the sender address must be specified")
	}

	if len(msg.Receiver) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "the receiver address must be specified")
	}

	if len(msg.ReceiverOnOtherChain) > MaxLengthForAddressOnOtherChain {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("the length of the receiver on other chain must be between [0,%d]", MaxLengthForAddressOnOtherChain))
	}

	if !msg.OutAmount.IsPositive() {
		return ErrInvalidAmount(DefaultCodespace, "the transferred amount must be positive")
	}

	if err := ValidateSecretHashLock(msg.SecretHashLock); err != nil {
		return ErrInvalidSecretHashLock(DefaultCodespace, err.Error())
	}

	if msg.TimeLock < MinTimeLock || msg.TimeLock > MaxTimeLock {
		return ErrInvalidSecretHashLock(DefaultCodespace, fmt.Sprintf("the time lock must be between [%d,%d]", MinTimeLock, MaxTimeLock))
	}

	return nil
}

// ValidateSecretHashLock validates the secret hash lock
func ValidateSecretHashLock(secretHashLock string) sdk.Error {
	secretHash, err := hex.DecodeString(secretHashLock)
	if err != nil {
		return ErrInvalidSecretHashLock(DefaultCodespace, fmt.Sprintf("invalid secret hash lock: %s", err.Error()))
	}

	if len(secretHash) != 32 {
		return ErrInvalidSecretHashLock(DefaultCodespace, fmt.Sprintf("invalid secret hash lock: %s", secretHashLock))
	}

	return nil
}

// Implements Msg.
func (msg MsgCreateHTLC) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgCreateHTLC) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
