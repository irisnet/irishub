package types

import (
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

const (
	DefaultCodespace sdk.CodespaceType = "oracle"

	CodeUnknownFeedKey       sdk.CodeType = 100
	CodeEmptyFeedKey         sdk.CodeType = 101
	CodeExistedFeedKey       sdk.CodeType = 102
	CodeEmptyServiceName     sdk.CodeType = 103
	CodeInvalidMaxHistory    sdk.CodeType = 104
	CodeEmptyProviders       sdk.CodeType = 105
	CodeInvalidMaxServiceFee sdk.CodeType = 106
	CodeInvalidAddress       sdk.CodeType = 107
)

func ErrUnknownFeedKey(codespace sdk.CodespaceType, feedKey string) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownFeedKey, fmt.Sprintf("feedKey %s does not exist", feedKey))
}

func ErrEmptyFeedKey(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyFeedKey, "feedKey can not be empty")
}

func ErrExistedFeedKey(codespace sdk.CodespaceType, feedKey string) sdk.Error {
	return sdk.NewError(codespace, CodeExistedFeedKey, fmt.Sprintf("feedKey %s already exists", feedKey))
}

func ErrEmptyServiceName(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyServiceName, "service name can not be empty")
}

func ErrEmptyProviders(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyProviders, "provider can not be empty")
}

func ErrInvalidMaxHistory(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidMaxHistory, fmt.Sprintf("max History is invalid, should be between 1 and %d", MaxHistory))
}

func ErrInvalidMaxServiceFee(codespace sdk.CodespaceType, fees sdk.Coins) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidMaxServiceFee, fmt.Sprintf("max service fee %s is invalid", fees.String()))
}

func ErrInvalidAddress(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAddress, msg)
}
