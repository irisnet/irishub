package types

import (
	"encoding/hex"
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

	CodeInvalidRequest            sdk.CodeType = 112
	CodeInvalidServiceFee         sdk.CodeType = 113
	CodeInvalidResponse           sdk.CodeType = 114
	CodeInvalidRequestID          sdk.CodeType = 115
	CodeInvalidProviders          sdk.CodeType = 116
	CodeInvalidTimeout            sdk.CodeType = 117
	CodeInvalidRepeatedFreq       sdk.CodeType = 117
	CodeInvalidRepeatedTotal      sdk.CodeType = 118
	CodeUnknownRequestContext     sdk.CodeType = 119
	CodeInvalidRequestContextID   sdk.CodeType = 120
	CodeRequestContextNonRepeated sdk.CodeType = 121
	CodeRequestContextNotStarted  sdk.CodeType = 122
	CodeRequestContextNotPaused   sdk.CodeType = 123
	CodeModuleNameRegistered      sdk.CodeType = 124
	CodeModuleNameNotRegistered   sdk.CodeType = 125
	CodeNoEarnedFees              sdk.CodeType = 126

	CodeInvalidRequestInput   sdk.CodeType = 127
	CodeInvalidResponseOutput sdk.CodeType = 128
	CodeInvalidResponseErr    sdk.CodeType = 129

	CodeInvalidAddress  sdk.CodeType = 130
	CodeInvalidGuardian sdk.CodeType = 131
	CodeInvalidTrustee  sdk.CodeType = 132
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

func ErrInvalidRequest(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidRequest, fmt.Sprintf("invalid request: %s", msg))
}

func ErrInvalidServiceFee(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidServiceFee, msg)
}

func ErrInvalidResponse(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidResponse, fmt.Sprintf("invalid response: %s", msg))
}

func ErrInvalidRequestID(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidResponse, fmt.Sprintf("invalid request ID: %s", msg))
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

func ErrUnknownRequestContext(codespace sdk.CodespaceType, requestContextID []byte) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownRequestContext, fmt.Sprintf("unknown request context: %s", hex.EncodeToString(requestContextID)))
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

func ErrInvalidGuardian(codespace sdk.CodespaceType, address sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidGuardian, fmt.Sprintf("invalid guardian: %s", address))
}

func ErrInvalidTrustee(codespace sdk.CodespaceType, address sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTrustee, fmt.Sprintf("invalid trustee: %s", address))
}

func ErrInvalidAddress(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAddress, msg)
}
