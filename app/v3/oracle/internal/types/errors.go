package types

import (
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

const (
	DefaultCodespace sdk.CodespaceType = "oracle"

	CodeUnknownFeedName            sdk.CodeType = 100
	CodeEmptyFeedName              sdk.CodeType = 101
	CodeExistedFeedName            sdk.CodeType = 102
	CodeEmptyServiceName           sdk.CodeType = 103
	CodeInvalidLatestHistory       sdk.CodeType = 104
	CodeEmptyProviders             sdk.CodeType = 105
	CodeInvalidServiceFeeCap       sdk.CodeType = 106
	CodeInvalidResponseThreshold   sdk.CodeType = 107
	CodeInvalidAddress             sdk.CodeType = 108
	CodeEmptyAggregateArgsJsonPath sdk.CodeType = 109
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

func ErrEmptyAggregateArgsJsonPath(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyAggregateArgsJsonPath, "aggregate args json path can not be empty")
}
