package iservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = 11

	CodeInvalidIDL sdk.CodeType = 100
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
