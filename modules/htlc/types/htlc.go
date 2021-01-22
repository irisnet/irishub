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
		Sender:               sender.String(),
		To:                   to.String(),
		ReceiverOnOtherChain: receiverOnOtherChain,
		Amount:               amount,
		Secret:               secret.String(),
		Timestamp:            timestamp,
		ExpirationHeight:     expirationHeight,
		State:                state,
	}
}

// Validate validates the HTLC
func (h HTLC) Validate() error {
	_, err := sdk.AccAddressFromBech32(h.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if _, err := sdk.AccAddressFromBech32(h.To); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid recipient address (%s)", err)
	}

	if err := ValidateReceiverOnOtherChain(h.ReceiverOnOtherChain); err != nil {
		return err
	}

	if err := ValidateAmount(h.Amount); err != nil {
		return err
	}

	if h.State != Completed && len(h.Secret) > 0 {
		return sdkerrors.Wrapf(ErrInvalidSecret, "secret must be empty when the HTLC has not be claimed")
	}

	if h.State == Completed && len(h.Secret) != SecretLength {
		return sdkerrors.Wrapf(ErrInvalidSecret, "length of the secret must be %d", SecretLength)
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
