package iservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = 6

	CodeInvalidIDL   sdk.CodeType = 100
	CodeSvcDefExists sdk.CodeType = 101
)

func codeToDefaultMsg(code sdk.CodeType) string {
	switch code {
	case CodeInvalidIDL:
		return "The IDL file cannot be parsed"
	default:
		return sdk.CodeToDefaultMsg(code)
	}
}

func NewError(codespace sdk.CodespaceType, code sdk.CodeType, msg string) sdk.Error {
	msg = msgOrDefaultMsg(msg, code)
	return sdk.NewError(codespace, code, msg)
}

func msgOrDefaultMsg(msg string, code sdk.CodeType) string {
	if msg != "" {
		return msg
	}
	return codeToDefaultMsg(code)
}

func ErrSvcDefExists(codespace sdk.CodespaceType, svcDefName string) sdk.Error {
	return sdk.NewError(codespace, CodeSvcDefExists, "service definition name %s already exist,must use new name", svcDefName)
}

func ErrInvalidIDL(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSvcDefExists, codeToDefaultMsg(CodeInvalidIDL))
}
