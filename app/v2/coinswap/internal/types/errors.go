// Package types nolint
package types

import (
	sdk "github.com/irisnet/irishub/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeReservePoolAlreadyExists  sdk.CodeType = 101
	CodeEqualDenom                sdk.CodeType = 102
	CodeInvalidDeadline           sdk.CodeType = 103
	CodeNotPositive               sdk.CodeType = 104
	CodeConstraintNotMet          sdk.CodeType = 105
	CodeLessThanMinReward         sdk.CodeType = 106
	CodeGreaterThanMaxDeposit     sdk.CodeType = 107
	CodeLessThanMinWithdrawAmount sdk.CodeType = 108
)

func ErrReservePoolAlreadyExists(codespace sdk.CodespaceType, msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(codespace, CodeReservePoolAlreadyExists, msg)
	}
	return sdk.NewError(codespace, CodeReservePoolAlreadyExists, "reserve pool already exists")
}

func ErrEqualDenom(codespace sdk.CodespaceType, msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(codespace, CodeEqualDenom, msg)
	}
	return sdk.NewError(codespace, CodeEqualDenom, "input and output denomination are equal")
}

func ErrInvalidDeadline(codespace sdk.CodespaceType, msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(codespace, CodeInvalidDeadline, msg)
	}
	return sdk.NewError(codespace, CodeInvalidDeadline, "invalid deadline")
}

func ErrNotPositive(codespace sdk.CodespaceType, msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(codespace, CodeNotPositive, msg)
	}
	return sdk.NewError(codespace, CodeNotPositive, "amount is not positive")
}

func ErrConstraintNotMet(codespace sdk.CodespaceType, msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(codespace, CodeConstraintNotMet, msg)
	}
	return sdk.NewError(codespace, CodeConstraintNotMet, "constraint not met")
}

func ErrLessThanMinReward(codespace sdk.CodespaceType, msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(codespace, CodeLessThanMinReward, msg)
	}
	return sdk.NewError(codespace, CodeLessThanMinReward, "min liquidity is less than MinReward")
}

func ErrGreaterThanMaxDeposit(codespace sdk.CodespaceType, msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(codespace, CodeGreaterThanMaxDeposit, msg)
	}
	return sdk.NewError(codespace, CodeGreaterThanMaxDeposit, "deposited amount is greater than max deposited amount of user")
}

func ErrLessThanMinWithdrawAmount(codespace sdk.CodespaceType, msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(codespace, CodeLessThanMinWithdrawAmount, msg)
	}
	return sdk.NewError(codespace, CodeLessThanMinWithdrawAmount, "less than min withdraw amount")
}
