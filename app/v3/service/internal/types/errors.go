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
	CodeNoWithdrawAddr            sdk.CodeType = 110
	CodeServiceBindingUnavailable sdk.CodeType = 111
	CodeServiceBindingAvailable   sdk.CodeType = 112
	CodeIncorrectRefundTime       sdk.CodeType = 113

	CodeInvalidServiceFee         sdk.CodeType = 114
	CodeInvalidProviders          sdk.CodeType = 115
	CodeInvalidTimeout            sdk.CodeType = 116
	CodeInvalidRepeatedFreq       sdk.CodeType = 117
	CodeInvalidRepeatedTotal      sdk.CodeType = 118
	CodeInvalidThreshold          sdk.CodeType = 119
	CodeInvalidResponse           sdk.CodeType = 120
	CodeInvalidRequestID          sdk.CodeType = 121
	CodeUnknownRequest            sdk.CodeType = 122
	CodeUnknownResponse           sdk.CodeType = 123
	CodeUnknownRequestContext     sdk.CodeType = 124
	CodeInvalidRequestContextID   sdk.CodeType = 125
	CodeNotMatchingConsumer       sdk.CodeType = 126
	CodeRequestContextNonRepeated sdk.CodeType = 127
	CodeRequestContextNotStarted  sdk.CodeType = 128
	CodeRequestContextNotPaused   sdk.CodeType = 129
	CodeRequestContextCompleted   sdk.CodeType = 130
	CodeModuleNameRegistered      sdk.CodeType = 131
	CodeModuleNameNotRegistered   sdk.CodeType = 132
	CodeNoEarnedFees              sdk.CodeType = 133

	CodeInvalidRequestInput   sdk.CodeType = 134
	CodeInvalidResponseOutput sdk.CodeType = 135
	CodeInvalidResponseErr    sdk.CodeType = 136

	CodeInvalidAddress  sdk.CodeType = 137
	CodeInvalidProfiler sdk.CodeType = 138
	CodeInvalidTrustee  sdk.CodeType = 139
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

func ErrNoWithdrawAddr(codespace sdk.CodespaceType, provider sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeNoWithdrawAddr, fmt.Sprintf("no withdraw address for %s", provider))
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

func ErrInvalidRequestID(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidRequestID, fmt.Sprintf("invalid request ID: %s", msg))
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

func ErrNotMatchingConsumer(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeNotMatchingConsumer, "consumer does not match")
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

func ErrModuleNameRegistered(codespace sdk.CodespaceType, moduleName string) sdk.Error {
	return sdk.NewError(codespace, CodeModuleNameRegistered, fmt.Sprintf("module %s already registered", moduleName))
}

func ErrModuleNameNotRegistered(codespace sdk.CodespaceType, moduleName string) sdk.Error {
	return sdk.NewError(codespace, CodeModuleNameNotRegistered, fmt.Sprintf("module %s not registered", moduleName))
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

func ErrInvalidResponseErr(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidResponseErr, fmt.Sprintf("invalid response err: %s", msg))
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
