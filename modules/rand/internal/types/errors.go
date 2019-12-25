//nolint
package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Rand errors reserve 100 ~ 199.
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidConsumer sdk.CodeType = 100
	CodeInvalidReqID    sdk.CodeType = 101
	CodeInvalidHeight   sdk.CodeType = 102
)

// ErrInvalidConsumer error for invalid consumer
func ErrInvalidConsumer(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidConsumer, msg)
}

// ErrInvalidReqID error for invalid request id
func ErrInvalidReqID(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidReqID, msg)
}

// ErrInvalidHeight error for invalid height
func ErrInvalidHeight(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidHeight, msg)
}
