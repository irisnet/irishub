package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// coinswap errors reserve 100 ~ 199.
// nolint
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeReservePoolNotExists         sdk.CodeType = 101
	CodeEqualDenom                   sdk.CodeType = 102
	CodeInvalidDeadline              sdk.CodeType = 103
	CodeNotPositive                  sdk.CodeType = 104
	CodeConstraintNotMet             sdk.CodeType = 105
	CodeIllegalDenom                 sdk.CodeType = 106
	CodeReservePoolInsufficientFunds sdk.CodeType = 107
)

// ErrReservePoolNotExists error for reserve pool not exists
func ErrReservePoolNotExists(msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(DefaultCodespace, CodeReservePoolNotExists, msg)
	}
	return sdk.NewError(DefaultCodespace, CodeReservePoolNotExists, "reserve pool not exists")
}

// ErrEqualDenom error for equal denom
func ErrEqualDenom(msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(DefaultCodespace, CodeEqualDenom, msg)
	}
	return sdk.NewError(DefaultCodespace, CodeEqualDenom, "input and output denomination are equal")
}

// ErrIllegalDenom error for illegal denom
func ErrIllegalDenom(msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(DefaultCodespace, CodeIllegalDenom, msg)
	}
	return sdk.NewError(DefaultCodespace, CodeIllegalDenom, "illegal denomination")
}

// ErrInvalidDeadline error for invalid deadline
func ErrInvalidDeadline(msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(DefaultCodespace, CodeInvalidDeadline, msg)
	}
	return sdk.NewError(DefaultCodespace, CodeInvalidDeadline, "invalid deadline")
}

// ErrNotPositive error for not positive
func ErrNotPositive(msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(DefaultCodespace, CodeNotPositive, msg)
	}
	return sdk.NewError(DefaultCodespace, CodeNotPositive, "amount is not positive")
}

// ErrConstraintNotMet error for constraint not met
func ErrConstraintNotMet(msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(DefaultCodespace, CodeConstraintNotMet, msg)
	}
	return sdk.NewError(DefaultCodespace, CodeConstraintNotMet, "constraint not met")
}

// ErrInsufficientFunds error for insufficient funds
func ErrInsufficientFunds(msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(DefaultCodespace, CodeReservePoolInsufficientFunds, msg)
	}
	return sdk.NewError(DefaultCodespace, CodeReservePoolInsufficientFunds, "constraint not met")
}
