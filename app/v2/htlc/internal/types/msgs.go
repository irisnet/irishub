package types

import (
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

const (
	// type for MsgCreateHTLC
	TypeMsgCreateHTLC = "create_htlc"

	// type for MsgClaimHTLC
	TypeMsgClaimHTLC = "claim_htlc"

	// type for MsgRefundHTLC
	TypeMsgRefundHTLC = "refund_htlc"

	SecretLength                    = 32    // the length for the secret
	HashLockLength                  = 32    // the length for the hash lock
	MaxLengthForAddressOnOtherChain = 32    // maximal length in bytes for the address on other chains
	MinTimeLock                     = 50    // minimal time span for HTLC
	MaxTimeLock                     = 25480 // maximal time span for HTLC
)

var _ sdk.Msg = &MsgCreateHTLC{}
var _ sdk.Msg = &MsgClaimHTLC{}
var _ sdk.Msg = &MsgRefundHTLC{}

// MsgCreateHTLC represents a msg for creating an HTLC
type MsgCreateHTLC struct {
	Sender               sdk.AccAddress `json:"sender"`                  // the initiator address
	Receiver             sdk.AccAddress `json:"receiver"`                // the recipient address
	ReceiverOnOtherChain []byte         `json:"receiver_on_other_chain"` // the recipient address on other chain
	Amount               sdk.Coin       `json:"amount"`                  // the amount to be transferred
	HashLock             []byte         `json:"hash_lock"`               // the hash lock generated from secret (and timestamp if provided)
	Timestamp            uint64         `json:"timestamp"`               // if provided, used to generate the hash lock together with secret
	TimeLock             uint64         `json:"time_lock"`               // the time span after which the HTLC will expire
}

// NewMsgCreateHTLC constructs a MsgCreateHTLC
func NewMsgCreateHTLC(
	sender sdk.AccAddress,
	receiver sdk.AccAddress,
	receiverOnOtherChain []byte,
	amount sdk.Coin,
	hashLock []byte,
	timestamp uint64,
	timeLock uint64,
) MsgCreateHTLC {
	return MsgCreateHTLC{
		Sender:               sender,
		Receiver:             receiver,
		ReceiverOnOtherChain: receiverOnOtherChain,
		Amount:               amount,
		HashLock:             hashLock,
		Timestamp:            timestamp,
		TimeLock:             timeLock,
	}
}

// Implements Msg.
func (msg MsgCreateHTLC) Route() string { return RouterKey }

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

	if !msg.Amount.IsValid() || !msg.Amount.IsPositive() {
		return ErrInvalidAmount(DefaultCodespace, "the transferred amount must be valid")
	}

	if len(msg.HashLock) != HashLockLength {
		return ErrInvalidHashLock(DefaultCodespace, fmt.Sprintf("the hash lock must be %d bytes long", HashLockLength))
	}

	if msg.TimeLock < MinTimeLock || msg.TimeLock > MaxTimeLock {
		return ErrInvalidTimeLock(DefaultCodespace, fmt.Sprintf("the time lock must be between [%d,%d]", MinTimeLock, MaxTimeLock))
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

// -----------------------------------------------------------------------------

// MsgClaimHTLC represents a msg for claiming an HTLC
type MsgClaimHTLC struct {
	Sender   sdk.AccAddress `json:"sender"`    // the initiator address
	HashLock []byte         `json:"hash_lock"` // the hash lock identifying the HTLC to be claimed
	Secret   []byte         `json:"secret"`    // the secret with which to claim
}

// NewMsgClaimHTLC constructs a MsgClaimHTLC
func NewMsgClaimHTLC(
	sender sdk.AccAddress,
	hashLock []byte,
	secret []byte,
) MsgClaimHTLC {
	return MsgClaimHTLC{
		Sender:   sender,
		HashLock: hashLock,
		Secret:   secret,
	}
}

// Implements Msg.
func (msg MsgClaimHTLC) Route() string { return RouterKey }

// Implements Msg.
func (msg MsgClaimHTLC) Type() string { return TypeMsgClaimHTLC }

// Implements Msg.
func (msg MsgClaimHTLC) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "the sender address must be specified")
	}

	if len(msg.HashLock) != HashLockLength {
		return ErrInvalidHashLock(DefaultCodespace, fmt.Sprintf("the hash lock must be %d bytes long", HashLockLength))
	}

	if len(msg.Secret) != SecretLength {
		return ErrInvalidSecret(DefaultCodespace, fmt.Sprintf("the secret must be %d bytes long", SecretLength))
	}

	return nil
}

// Implements Msg.
func (msg MsgClaimHTLC) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgClaimHTLC) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// -----------------------------------------------------------------------------

// MsgRefundHTLC represents a msg for refund an HTLC
type MsgRefundHTLC struct {
	Sender   sdk.AccAddress `json:"sender"`    // the initiator address
	HashLock []byte         `json:"hash_lock"` // the hash lock identifying the HTLC to be refunded
}

// NewMsgClaimHTLC constructs a MsgClaimHTLC
func NewMsgRefundHTLC(
	sender sdk.AccAddress,
	hashLock []byte,
) MsgRefundHTLC {
	return MsgRefundHTLC{
		Sender:   sender,
		HashLock: hashLock,
	}
}

// Implements Msg.
func (msg MsgRefundHTLC) Route() string { return RouterKey }

// Implements Msg.
func (msg MsgRefundHTLC) Type() string { return TypeMsgRefundHTLC }

// Implements Msg.
func (msg MsgRefundHTLC) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "the sender address must be specified")
	}

	if len(msg.HashLock) != HashLockLength {
		return ErrInvalidHashLock(DefaultCodespace, fmt.Sprintf("the hash lock must be %d bytes long", HashLockLength))
	}

	return nil
}

// Implements Msg.
func (msg MsgRefundHTLC) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgRefundHTLC) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
