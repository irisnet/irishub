package types

import (
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHTLC constructs a new HTLC instance
func NewHTLC(
	sender sdk.AccAddress,
	to sdk.AccAddress,
	receiverOnOtherChain string,
	amount sdk.Coins,
	secret tmbytes.HexBytes,
	timestamp uint64,
	expirationHeight uint64,
	state HTLCState,
) HTLC {
	return HTLC{
		Sender:               sender,
		To:                   to,
		ReceiverOnOtherChain: receiverOnOtherChain,
		Amount:               amount,
		Secret:               secret,
		Timestamp:            timestamp,
		ExpirationHeight:     expirationHeight,
		State:                state,
	}
}

// Validate validates the HTLC
func (h HTLC) Validate() error {
	if h.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "sender missing")
	}

	if len(h.To) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "recipient missing")
	}

	if len(h.ReceiverOnOtherChain) > MaxLengthForAddressOnOtherChain {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "length of the receiver on other chain must be between [0,%d]", MaxLengthForAddressOnOtherChain)
	}

	if !h.Amount.IsValid() || !h.Amount.IsAllPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "the transferred amount must be valid")
	}

	if h.State != Completed && len(h.Secret) > 0 {
		return sdkerrors.Wrapf(ErrInvalidSecret, "secret must be empty when the HTLC has not be claimed")
	}

	if h.State == Completed && len(h.Secret) != SecretLength {
		return sdkerrors.Wrapf(ErrInvalidSecret, "length of the secret must be %d in bytes", SecretLength)
	}

	return nil
}

// GetHashLock calculates the hash lock from the given secret and timestamp
func GetHashLock(secret tmbytes.HexBytes, timestamp uint64) []byte {
	if timestamp > 0 {
		return tmhash.Sum(append(secret, sdk.Uint64ToBigEndian(timestamp)...))
	}

	return tmhash.Sum(secret)
}
