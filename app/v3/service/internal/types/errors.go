package types

import (
	"fmt"

	cmn "github.com/tendermint/tendermint/libs/common"

	sdk "github.com/irisnet/irishub/types"
)

const (
	DefaultCodespace sdk.CodespaceType = "service"

	CodeInvalidServiceName       sdk.CodeType = 100
	CodeInvalidSchemas           sdk.CodeType = 101
	CodeInvalidLength            sdk.CodeType = 102
	CodeDuplicateTags            sdk.CodeType = 103
	CodeUnknownServiceDefinition sdk.CodeType = 104
	CodeServiceDefinitionExists  sdk.CodeType = 105

	CodeInvalidDeposit            sdk.CodeType = 106
	CodeInvalidPricing            sdk.CodeType = 107
	CodeServiceBindingExists      sdk.CodeType = 108
	CodeUnknownServiceBinding     sdk.CodeType = 109
	CodeServiceBindingUnavailable sdk.CodeType = 110
	CodeServiceBindingAvailable   sdk.CodeType = 111
	CodeIncorrectRefundTime       sdk.CodeType = 112

	CodeInvalidServiceFee         sdk.CodeType = 113
	CodeInvalidProviders          sdk.CodeType = 114
	CodeInvalidTimeout            sdk.CodeType = 115
	CodeInvalidRepeatedFreq       sdk.CodeType = 116
	CodeInvalidRepeatedTotal      sdk.CodeType = 117
	CodeInvalidThreshold          sdk.CodeType = 118
	CodeInvalidResponse           sdk.CodeType = 119
	CodeInvalidRequestID          sdk.CodeType = 120
	CodeUnknownRequest            sdk.CodeType = 121
	CodeUnknownResponse           sdk.CodeType = 122
	CodeUnknownRequestContext     sdk.CodeType = 123
	CodeInvalidRequestContextID   sdk.CodeType = 124
	CodeNotAuthorized             sdk.CodeType = 125
	CodeRequestContextNonRepeated sdk.CodeType = 126
	CodeRequestContextNotStarted  sdk.CodeType = 127
	CodeRequestContextNotPaused   sdk.CodeType = 128
	CodeRequestContextCompleted   sdk.CodeType = 129
	CodeCallbackRegistered        sdk.CodeType = 130
	CodeCallbackNotRegistered     sdk.CodeType = 131
	CodeNoEarnedFees              sdk.CodeType = 132

	CodeInvalidRequestInput   sdk.CodeType = 133
	CodeInvalidResponseOutput sdk.CodeType = 134
	CodeInvalidResponseResult sdk.CodeType = 135

	CodeInvalidAddress  sdk.CodeType = 136
	CodeInvalidProfiler sdk.CodeType = 137
	CodeInvalidTrustee  sdk.CodeType = 138
)

func ErrInvalidServiceName(codespace sdk.CodespaceType, serviceName string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidServiceName, fmt.Sprintf("invalid service name %s; only alphanumeric characters, _ and - accepted, the length ranges in (0,70]", serviceName))
}

func ErrInvalidSchemas(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidSchemas, msg)
}

func ErrInvalidLength(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidLength, msg)
}

func ErrDuplicateTags(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeDuplicateTags, "there exists duplicate tags")
}

func ErrServiceDefinitionExists(codespace sdk.CodespaceType, serviceName string) sdk.Error {
	return sdk.NewError(codespace, CodeServiceDefinitionExists, fmt.Sprintf("service name %s already exists", serviceName))
}

func ErrUnknownServiceDefinition(codespace sdk.CodespaceType, serviceName string) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownServiceDefinition, fmt.Sprintf("service name %s does not exist", serviceName))
}

func ErrInvalidDeposit(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDeposit, msg)
}

func ErrInvalidPricing(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidPricing, fmt.Sprintf("invalid pricing: %s", msg))
}

func ErrServiceBindingExists(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeServiceBindingExists, "service binding already exists")
}

func ErrUnknownServiceBinding(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownServiceBinding, "service binding does not exist")
}

func ErrServiceBindingUnavailable(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeServiceBindingUnavailable, "service binding is unavailable")
}

func ErrServiceBindingAvailable(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeServiceBindingAvailable, "service binding is available")
}

func ErrIncorrectRefundTime(codespace sdk.CodespaceType, refundableTime string) sdk.Error {
	return sdk.NewError(codespace, CodeIncorrectRefundTime, fmt.Sprintf("can not refund before %s", refundableTime))
}

func ErrInvalidServiceFee(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidServiceFee, msg)
}

func ErrInvalidResponse(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidResponse, fmt.Sprintf("invalid response: %s", msg))
}

func ErrInvalidRequestID(codespace sdk.CodespaceType, requestID cmn.HexBytes) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidRequestID, fmt.Sprintf("invalid request ID: %s", requestID.String()))
}

func ErrInvalidProviders(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidProviders, fmt.Sprintf("invalid providers: %s", msg))
}

func ErrInvalidTimeout(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTimeout, msg)
}

func ErrInvalidRepeatedFreq(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidRepeatedFreq, msg)
}

func ErrInvalidRepeatedTotal(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidRepeatedTotal, msg)
}

func ErrInvalidThreshold(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidThreshold, msg)
}

func ErrUnknownRequest(codespace sdk.CodespaceType, requestID cmn.HexBytes) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownRequest, fmt.Sprintf("unknown request: %s", requestID.String()))
}

func ErrUnknownResponse(codespace sdk.CodespaceType, requestID cmn.HexBytes) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownResponse, fmt.Sprintf("unknown response: %s", requestID.String()))
}

func ErrUnknownRequestContext(codespace sdk.CodespaceType, requestContextID cmn.HexBytes) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownRequestContext, fmt.Sprintf("unknown request context: %s", requestContextID.String()))
}

func ErrInvalidRequestContextID(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidRequestContextID, fmt.Sprintf("invalid request context ID: %s", msg))
}

func ErrNotAuthorized(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeNotAuthorized, msg)
}

func ErrRequestContextNonRepeated(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeRequestContextNonRepeated, "request context is non repeated")
}

func ErrRequestContextNotStarted(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeRequestContextNotStarted, "request context not started")
}

func ErrRequestContextNotPaused(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeRequestContextNotPaused, "request context not paused")
}

func ErrRequestContextCompleted(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeRequestContextCompleted, "request context completed")
}

func ErrCallbackRegistered(codespace sdk.CodespaceType, cbType, moduleName string) sdk.Error {
	return sdk.NewError(codespace, CodeCallbackRegistered, fmt.Sprintf("%s already registered for module %s", cbType, moduleName))
}

func ErrCallbackNotRegistered(codespace sdk.CodespaceType, cbType, moduleName string) sdk.Error {
	return sdk.NewError(codespace, CodeCallbackNotRegistered, fmt.Sprintf("%s not registered for module %s", cbType, moduleName))
}

func ErrNoEarnedFees(codespace sdk.CodespaceType, provider sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeNoEarnedFees, fmt.Sprintf("no earned fees for %s", provider))
}

func ErrInvalidRequestInput(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidRequestInput, fmt.Sprintf("invalid request input: %s", msg))
}

func ErrInvalidResponseOutput(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidResponseOutput, fmt.Sprintf("invalid response output: %s", msg))
}

func ErrInvalidResponseResult(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidResponseResult, fmt.Sprintf("invalid response result: %s", msg))
}

func ErrInvalidProfiler(codespace sdk.CodespaceType, address sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidProfiler, fmt.Sprintf("invalid profiler: %s", address))
}

func ErrInvalidTrustee(codespace sdk.CodespaceType, address sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTrustee, fmt.Sprintf("invalid trustee: %s", address))
}

func ErrInvalidAddress(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAddress, msg)
}
