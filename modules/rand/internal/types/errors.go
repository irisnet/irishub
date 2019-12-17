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
