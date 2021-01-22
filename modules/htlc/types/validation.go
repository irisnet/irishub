package types

import (
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// SecretLength ength for the secret in hex string
	SecretLength = 64
	// HashLockLength for the hash lock in hex string
	HashLockLength = 64
	// MaxLengthForAddressOnOtherChain length for the address on other chains
	MaxLengthForAddressOnOtherChain = 128
	// MinTimeLock time span for HTLC
	MinTimeLock = 50
	// MaxTimeLock time span for HTLC
	MaxTimeLock = 25480
)

// ValidateReceiverOnOtherChain verifies whether the  parameters are legal
func ValidateReceiverOnOtherChain(receiverOnOtherChain string) error {
	if len(receiverOnOtherChain) > MaxLengthForAddressOnOtherChain {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "length of the receiver on other chain must be between [0,%d]", MaxLengthForAddressOnOtherChain)
	}
	return nil
}

// ValidateAmount verifies whether the  parameters are legal
func ValidateAmount(amount sdk.Coins) error {
	if !(amount.IsValid() && amount.IsAllPositive()) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "the transferred amount must be valid")
	}
	return nil
}

// ValidateHashLock verifies whether the  parameters are legal
func ValidateHashLock(hashLock string) error {
	if len(hashLock) != HashLockLength {
		return sdkerrors.Wrapf(ErrInvalidHashLock, "length of the hash lock must be %d", HashLockLength)
	}

	if _, err := hex.DecodeString(hashLock); err != nil {
		return sdkerrors.Wrapf(ErrInvalidHashLock, "hash lock must be a hex encoded string")
	}
	return nil
}

// ValidateTimeLock verifies whether the  parameters are legal
func ValidateTimeLock(timeLock uint64) error {
	if timeLock < MinTimeLock || timeLock > MaxTimeLock {
		return sdkerrors.Wrapf(ErrInvalidTimeLock, "the time lock must be between [%d,%d]", MinTimeLock, MaxTimeLock)
	}
	return nil
}

// ValidateSecret verifies whether the  parameters are legal
func ValidateSecret(secret string) error {
	if len(secret) != SecretLength {
		return sdkerrors.Wrapf(ErrInvalidSecret, "length of the secret must be %d", SecretLength)
	}

	if _, err := hex.DecodeString(secret); err != nil {
		return sdkerrors.Wrapf(ErrInvalidSecret, "secret must be a hex encoded string")
	}
	return nil
}
