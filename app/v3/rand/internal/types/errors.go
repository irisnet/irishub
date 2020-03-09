//nolint
package types

import (
	sdk "github.com/irisnet/irishub/types"
)

// Rand errors reserve 100 ~ 199.
const (
	DefaultCodespace sdk.CodespaceType = "rand"

	CodeInvalidConsumer        sdk.CodeType = 100
	CodeInvalidReqID           sdk.CodeType = 101
	CodeInvalidHeight          sdk.CodeType = 102
	CodeInvalidServiceBindings sdk.CodeType = 103
	CodeInvalidServiceFee      sdk.CodeType = 104
	CodeInsufficientBalance    sdk.CodeType = 105
)

//----------------------------------------
// Rand error constructors

func ErrInvalidConsumer(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidConsumer, msg)
}

func ErrInvalidReqID(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidReqID, msg)
}

func ErrInvalidHeight(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidHeight, msg)
}

func ErrInvalidServiceBindings(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidServiceBindings, msg)
}

func ErrInvalidServiceFee(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidServiceFee, msg)
}

func ErrInsufficientBalance(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInsufficientBalance, msg)
}
