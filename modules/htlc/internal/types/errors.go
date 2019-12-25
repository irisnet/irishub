//nolint
package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Rand errors reserve 100 ~ 199.
const (
	DefaultCodespace sdk.CodespaceType = "htlc"

	CodeHTLCNotExists         sdk.CodeType = 100
	CodeInvalidAddress        sdk.CodeType = 101
	CodeInvalidAmount         sdk.CodeType = 102
	CodeInvalidHashLock       sdk.CodeType = 103
	CodeHashLockAlreadyExists sdk.CodeType = 104
	CodeInvalidTimeLock       sdk.CodeType = 105
	CodeInvalidSecret         sdk.CodeType = 106
	CodeStateIsNotOpen        sdk.CodeType = 107
	CodeStateIsNotExpired     sdk.CodeType = 108
)

// ErrHTLCNotExists error for HTLC not exists
func ErrHTLCNotExists(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeHTLCNotExists, msg)
}

// ErrInvalidAddress error for invalid address
func ErrInvalidAddress(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAddress, msg)
}

// ErrInvalidAmount error for invalid amount
func ErrInvalidAmount(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAmount, msg)
}

// ErrInvalidHashLock error for invalid hash lock
func ErrInvalidHashLock(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidHashLock, msg)
}

// ErrHashLockAlreadyExists error for hash lock already exists
func ErrHashLockAlreadyExists(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeHashLockAlreadyExists, msg)
}

// ErrInvalidTimeLock error for invalid lock time
func ErrInvalidTimeLock(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTimeLock, msg)
}

// ErrInvalidSecret error for invalid secret
func ErrInvalidSecret(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidSecret, msg)
}

// ErrStateIsNotOpen error for state is not open
func ErrStateIsNotOpen(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeStateIsNotOpen, msg)
}

// ErrStateIsNotExpired error for invalid state is not expired
func ErrStateIsNotExpired(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeStateIsNotExpired, msg)
}
