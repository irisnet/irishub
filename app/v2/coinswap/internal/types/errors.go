// Package types nolint
package types

import (
	sdk "github.com/irisnet/irishub/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeReservePoolNotExists         sdk.CodeType = 101
	CodeEqualDenom                   sdk.CodeType = 102
	CodeInvalidDeadline              sdk.CodeType = 103
	CodeNotPositive                  sdk.CodeType = 104
	CodeConstraintNotMet             sdk.CodeType = 105
	CodeIllegalDenom                 sdk.CodeType = 106
	CodeIllegalUniId                 sdk.CodeType = 107
	CodeReservePoolInsufficientFunds sdk.CodeType = 108
)

func ErrReservePoolNotExists(msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(DefaultCodespace, CodeReservePoolNotExists, msg)
	}
	return sdk.NewError(DefaultCodespace, CodeReservePoolNotExists, "reserve pool not exists")
}

func ErrEqualDenom(msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(DefaultCodespace, CodeEqualDenom, msg)
	}
	return sdk.NewError(DefaultCodespace, CodeEqualDenom, "input and output denomination are equal")
}

func ErrIllegalUniId(msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(DefaultCodespace, CodeIllegalUniId, msg)
	}
	return sdk.NewError(DefaultCodespace, CodeIllegalUniId, "illegal liquidity id")
}

func ErrIllegalDenom(msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(DefaultCodespace, CodeIllegalDenom, msg)
	}
	return sdk.NewError(DefaultCodespace, CodeIllegalDenom, "illegal denomination")
}

func ErrInvalidDeadline(msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(DefaultCodespace, CodeInvalidDeadline, msg)
	}
	return sdk.NewError(DefaultCodespace, CodeInvalidDeadline, "invalid deadline")
}

func ErrNotPositive(msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(DefaultCodespace, CodeNotPositive, msg)
	}
	return sdk.NewError(DefaultCodespace, CodeNotPositive, "amount is not positive")
}

func ErrConstraintNotMet(msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(DefaultCodespace, CodeConstraintNotMet, msg)
	}
	return sdk.NewError(DefaultCodespace, CodeConstraintNotMet, "constraint not met")
}

func ErrInsufficientFunds(msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(DefaultCodespace, CodeReservePoolInsufficientFunds, msg)
	}
	return sdk.NewError(DefaultCodespace, CodeReservePoolInsufficientFunds, "constraint not met")
}
