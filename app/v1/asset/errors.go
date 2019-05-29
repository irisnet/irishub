//nolint
package asset

import (
	sdk "github.com/irisnet/irishub/types"
)

// Asset errors reserve 100 ~ 199.
const (
	DefaultCodespace sdk.CodespaceType = "asset"

	CodeInvalidMoniker       sdk.CodeType = 100
	CodeInvalidDetails       sdk.CodeType = 101
	CodeInvalidWebsite       sdk.CodeType = 102
	CodeUnknownGateway       sdk.CodeType = 103
	CodeGatewayAlreadyExists sdk.CodeType = 104
	CodeOwnerIsOperator      sdk.CodeType = 105
	CodeInvalidGenesis       sdk.CodeType = 106
)

// NOTE: Don't stringer this, we'll put better messages in later.
func codeToDefaultMsg(code sdk.CodeType) string {
	switch code {

	default:
		return sdk.CodeToDefaultMsg(code)
	}
}

//----------------------------------------
// Error constructors

// TODO

//----------------------------------------

func msgOrDefaultMsg(msg string, code sdk.CodeType) string {
	if msg != "" {
		return msg
	}
	return codeToDefaultMsg(code)
}

func newError(codespace sdk.CodespaceType, code sdk.CodeType, msg string) sdk.Error {
	msg = msgOrDefaultMsg(msg, code)
	return sdk.NewError(codespace, code, msg)
}

//----------------------------------------
// Error constructors

func ErrInvalidMoniker(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidMoniker, msg)
}

func ErrInvalidDetails(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDetails, msg)
}

func ErrInvalidWebsite(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidWebsite, msg)
}

func ErrUnkwownGateway(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownGateway, msg)
}

func ErrGatewayAlreadyExists(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeGatewayAlreadyExists, msg)
}

func ErrOwnerIsOperator(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeOwnerIsOperator, msg)
}

func ErrInvalidGenesis(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidGenesis, msg)
}
