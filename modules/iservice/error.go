package iservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"fmt"
)

const (
	DefaultCodespace sdk.CodespaceType = 6

	CodeInvalidIDL               sdk.CodeType = 100
	CodeSvcDefExists             sdk.CodeType = 101
	CodeInvalidOutputPrivacyEnum sdk.CodeType = 102
	CodeInvalidOutputCachedEnum  sdk.CodeType = 103
	CodeInvalidServiceName       sdk.CodeType = 104
	CodeInvalidChainId           sdk.CodeType = 105
	CodeInvalidAuthor            sdk.CodeType = 106
	CodeInvalidMethodName        sdk.CodeType = 107
	CodeInvalidBroadcastEnum     sdk.CodeType = 108
	CodeMoreTags                 sdk.CodeType = 109
	CodeDuplicateTags            sdk.CodeType = 110
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
	return sdk.NewError(codespace, CodeSvcDefExists, fmt.Sprintf("service definition name %s already exist,must use new name", svcDefName))
}

func ErrInvalidIDL(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidIDL, fmt.Sprintf("The IDL content cannot be parsed, err: %s", msg))
}

func ErrInvalidOutputPrivacyEnum(codespace sdk.CodespaceType, value string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidOutputPrivacyEnum, fmt.Sprintf("invalid OutputPrivacyEnum %s", value))
}

func ErrInvalidOutputCachedEnum(codespace sdk.CodespaceType, value string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidOutputCachedEnum, fmt.Sprintf("invalid OutputCachedEnum %s", value))
}

func ErrInvalidServiceName(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidServiceName, fmt.Sprintf("service name is empty"))
}

func ErrInvalidChainId(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidChainId, fmt.Sprintf("chainId is empty"))
}

func ErrInvalidAuthor(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAuthor, fmt.Sprintf("author is empty"))
}

func ErrInvalidMethodName(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidMethodName, fmt.Sprintf("method name is empty"))
}

func ErrInvalidBroadcastEnum(codespace sdk.CodespaceType, value BroadcastEnum) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidBroadcastEnum, fmt.Sprintf("invalid BroadcastEnum %s", value))
}

func ErrMoreTags(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeMoreTags, fmt.Sprintf("tags are limited to %d", maxTagsNum))
}

func ErrDuplicateTags(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeDuplicateTags, "tags contains duplicate tag")
}
