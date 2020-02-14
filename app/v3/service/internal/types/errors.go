package types

import (
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

const (
	DefaultCodespace sdk.CodespaceType = "service"

	CodeInvalidServiceName       sdk.CodeType = 100
	CodeInvalidSchemas           sdk.CodeType = 101
	CodeInvalidLength            sdk.CodeType = 102
	CodeUnknownServiceDefinition sdk.CodeType = 103
	CodeServiceDefinitionExists  sdk.CodeType = 104

	CodeInvalidDeposit            sdk.CodeType = 105
	CodeInvalidPricing            sdk.CodeType = 106
	CodeServiceBindingExists      sdk.CodeType = 107
	CodeUnknownServiceBinding     sdk.CodeType = 108
	CodeServiceBindingUnavailable sdk.CodeType = 109
	CodeServiceBindingAvailable   sdk.CodeType = 110
	CodeIncorrectRefundTime       sdk.CodeType = 111

	CodeRequestNotActive     sdk.CodeType = 112
	CodeReturnFeeNotExists   sdk.CodeType = 113
	CodeWithdrawFeeNotExists sdk.CodeType = 114
	CodeLtServiceFee         sdk.CodeType = 115
	CodeInvalidReqID         sdk.CodeType = 116
	CodeNotMatchingProvider  sdk.CodeType = 117

	CodeInvalidRequestInput   sdk.CodeType = 118
	CodeInvalidResponseOutput sdk.CodeType = 119
	CodeInvalidResponseErr    sdk.CodeType = 120

	CodeInvalidInput   sdk.CodeType = 121
	CodeInvalidAddress sdk.CodeType = 122

	CodeInvalidRequest            sdk.CodeType = 123
	CodeInvalidProviders          sdk.CodeType = 124
	CodeInvalidRepeatedFreq       sdk.CodeType = 125
	CodeInvalidRepeatedTotal      sdk.CodeType = 126
	CodeInvalidRequestContextID   sdk.CodeType = 127
	CodeRequestContextNonRepeated sdk.CodeType = 128
	CodeRequestContextNotStarted  sdk.CodeType = 129
	CodeRequestContextNotPaused   sdk.CodeType = 130
	CodeModuleNameRegistered      sdk.CodeType = 131
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
	return sdk.NewError(codespace, CodeUnknownServiceBinding, fmt.Sprintf("service binding does not exist"))
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

func ErrRequestNotActive(codespace sdk.CodespaceType, requestID string) sdk.Error {
	return sdk.NewError(codespace, CodeRequestNotActive, fmt.Sprintf("request [%s] is not existed", requestID))
}

func ErrReturnFeeNotExists(codespace sdk.CodespaceType, address sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeReturnFeeNotExists, fmt.Sprintf("There is no service refund fees for [%s]", address))
}

func ErrWithdrawFeeNotExists(codespace sdk.CodespaceType, address sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeWithdrawFeeNotExists, fmt.Sprintf("There is no service withdraw fees for [%s]", address))
}

func ErrLtServiceFee(codespace sdk.CodespaceType, coins sdk.Coins) sdk.Error {
	return sdk.NewError(codespace, CodeLtServiceFee, fmt.Sprintf("service fee amount must be equal or greater than %s", coins.String()))
}

func ErrInvalidReqID(codespace sdk.CodespaceType, reqId string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidReqID, fmt.Sprintf("invalid request id [%s]", reqId))
}

func ErrNotMatchingProvider(codespace sdk.CodespaceType, provider sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeNotMatchingProvider, fmt.Sprintf("[%s] is not a matching Provider", provider.String()))
}

func ErrNotTrustee(codespace sdk.CodespaceType, trustee sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, fmt.Sprintf("[%s] is not a trustee address", trustee))
}

func ErrNotProfiler(codespace sdk.CodespaceType, profiler sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, fmt.Sprintf("[%s] is not a profiler address", profiler))
}

func ErrNoResponseFound(codespace sdk.CodespaceType, requestID string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, fmt.Sprintf("response is not existed for request %s", requestID))
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

func ErrInvalidAddress(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAddress, msg)
}

func ErrInvalidRequest(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidRequest, fmt.Sprintf("invalid request: %s", msg))
}

func ErrInvalidProviders(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidProviders, fmt.Sprintf("invalid providers: %s", msg))
}

func ErrInvalidRepeatedFreq(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidRepeatedFreq, msg)
}

func ErrInvalidRepeatedTotal(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidRepeatedTotal, msg)
}

func ErrInvalidRequestContextID(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidRequestContextID, fmt.Sprintf("invalid request context ID: %s", msg))
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

func ErrModuleNameRegistered(codespace sdk.CodespaceType, moduleName string) sdk.Error {
	return sdk.NewError(codespace, CodeModuleNameRegistered, fmt.Sprintf("module % already registered", moduleName))
}
