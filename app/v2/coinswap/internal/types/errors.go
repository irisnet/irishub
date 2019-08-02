// Package types nolint
package types

import (
	sdk "github.com/irisnet/irishub/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeReservePoolNotExists      sdk.CodeType = 101
	CodeEqualDenom                sdk.CodeType = 102
	CodeInvalidDeadline           sdk.CodeType = 103
	CodeNotPositive               sdk.CodeType = 104
	CodeConstraintNotMet          sdk.CodeType = 105
	CodeLessThanMinReward         sdk.CodeType = 106
	CodeGreaterThanMaxDeposit     sdk.CodeType = 107
	CodeLessThanMinWithdrawAmount sdk.CodeType = 108
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

func ErrLessThanMinReward(msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(DefaultCodespace, CodeLessThanMinReward, msg)
	}
	return sdk.NewError(DefaultCodespace, CodeLessThanMinReward, "min liquidity is less than MinReward")
}

func ErrGreaterThanMaxDeposit(msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(DefaultCodespace, CodeGreaterThanMaxDeposit, msg)
	}
	return sdk.NewError(DefaultCodespace, CodeGreaterThanMaxDeposit, "deposited amount is greater than max deposited amount of user")
}

func ErrLessThanMinWithdrawAmount(msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(DefaultCodespace, CodeLessThanMinWithdrawAmount, msg)
	}
	return sdk.NewError(DefaultCodespace, CodeLessThanMinWithdrawAmount, "less than min withdraw amount")
}
