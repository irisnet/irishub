package types

import (
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

const (
	DefaultCodespace sdk.CodespaceType = "oracle"

	CodeUnknownFeedName          sdk.CodeType = 100
	CodeEmptyFeedName            sdk.CodeType = 101
	CodeExistedFeedName          sdk.CodeType = 102
	CodeUnauthorized             sdk.CodeType = 103
	CodeEmptyServiceName         sdk.CodeType = 104
	CodeInvalidLatestHistory     sdk.CodeType = 105
	CodeEmptyProviders           sdk.CodeType = 106
	CodeInvalidServiceFeeCap     sdk.CodeType = 107
	CodeInvalidResponseThreshold sdk.CodeType = 108
	CodeInvalidAddress           sdk.CodeType = 109
	CodeEmptyAggregateFunc       sdk.CodeType = 110
	CodeEmptyValueJsonPath       sdk.CodeType = 111
	CodeUnknownRequestContextID  sdk.CodeType = 112
	CodeNotRegisterMethod        sdk.CodeType = 113
	CodeInvalidFeedState         sdk.CodeType = 114
	CodeNotProfiler              sdk.CodeType = 115
)

func ErrUnknownFeedName(codespace sdk.CodespaceType, feedName string) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownFeedName, fmt.Sprintf("feed name %s does not exist", feedName))
}

func ErrEmptyFeedName(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyFeedName, "feed name can not be empty")
}

func ErrExistedFeedName(codespace sdk.CodespaceType, feedName string) sdk.Error {
	return sdk.NewError(codespace, CodeExistedFeedName, fmt.Sprintf("feed name %s already exists", feedName))
}

func ErrUnauthorized(codespace sdk.CodespaceType, feedName string, owner sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeUnauthorized, fmt.Sprintf("feed %s does not belong to %s", feedName, owner.String()))
}

func ErrEmptyServiceName(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyServiceName, "service name can not be empty")
}

func ErrEmptyProviders(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyProviders, "provider can not be empty")
}

func ErrInvalidLatestHistory(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidLatestHistory, fmt.Sprintf("latest history is invalid, should be between 1 and %d", LatestHistory))
}

func ErrInvalidServiceFeeCap(codespace sdk.CodespaceType, fees sdk.Coins) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidServiceFeeCap, fmt.Sprintf("service fee %s is invalid", fees.String()))
}

func ErrInvalidResponseThreshold(codespace sdk.CodespaceType, limit int) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidResponseThreshold, fmt.Sprintf("response threshold should be between 1 and %d", limit))
}

func ErrInvalidAddress(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAddress, msg)
}

func ErrEmptyEmptyAggregateFunc(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyAggregateFunc, "aggregate func can not be empty")
}

func ErrEmptyValueJsonPath(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyValueJsonPath, "json path can not be empty")
}

func ErrUnknownRequestContextID(codespace sdk.CodespaceType, reqCtxID []byte) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownRequestContextID, "request context ID %s does not exist", string(reqCtxID))
}

func ErrNotRegisterMethod(codespace sdk.CodespaceType, methodName string) sdk.Error {
	return sdk.NewError(codespace, CodeNotRegisterMethod, "method %s don't register", methodName)
}

func ErrInvalidFeedState(codespace sdk.CodespaceType, feedName string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidFeedState, "feed %s may be a invalid state", feedName)
}

func ErrNotProfiler(codespace sdk.CodespaceType, profiler sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeNotProfiler, fmt.Sprintf("[%s] is not a profiler address", profiler))
}
