package upgrade

import (
	sdk "github.com/irisnet/irishub/types"
)

const (
	DefaultCodespace sdk.CodespaceType = "upgrade"

	CodeInvalidMsgType     sdk.CodeType = 100
	CodeUnSupportedMsgType sdk.CodeType = 101
	CodeUnknownRequest     sdk.CodeType = sdk.CodeUnknownRequest
	CodeNotCurrentProposal sdk.CodeType = 102
	CodeNotValidator       sdk.CodeType = 103
	CodeDoubleSwitch       sdk.CodeType = 104
)

func codeToDefaultMsg(code sdk.CodeType) string {
	switch code {
	case CodeInvalidMsgType:
		return "Invalid msg type"
	case CodeUnSupportedMsgType:
		return "Current version software doesn't support the msg type"
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
