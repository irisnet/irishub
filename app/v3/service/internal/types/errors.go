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

	CodeInvalidDeposit            sdk.CodeType = 108
	CodeInvalidPricing            sdk.CodeType = 108
	CodeServiceBindingExists      sdk.CodeType = 106
	CodeUnknownServiceBinding     sdk.CodeType = 107
	CodeServiceBindingUnavailable sdk.CodeType = 110
	CodeServiceBindingAvailable   sdk.CodeType = 111
	CodeIncorrectRefundTime       sdk.CodeType = 108

	CodeRequestNotActive     sdk.CodeType = 117
	CodeReturnFeeNotExists   sdk.CodeType = 118
	CodeWithdrawFeeNotExists sdk.CodeType = 119
	CodeLtServiceFee         sdk.CodeType = 120
	CodeInvalidReqId         sdk.CodeType = 121
	CodeNotMatchingProvider  sdk.CodeType = 123

	CodeInvalidRequestInput   sdk.CodeType = 127
	CodeInvalidResponseOutput sdk.CodeType = 128
	CodeInvalidResponseErr    sdk.CodeType = 129

	CodeInvalidInput   sdk.CodeType = 130
	CodeInvalidAddress sdk.CodeType = 132
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
	return sdk.NewError(codespace, CodeInvalidPricing, msg)
}

func ErrServiceBindingExists(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeServiceBindingExists, fmt.Sprintf("service binding already exists"))
}

func ErrUnknownServiceBinding(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownServiceBinding, fmt.Sprintf("service binding does not exist"))
}

func ErrServiceBindingUnavailable(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeServiceBindingUnavailable, fmt.Sprintf("service binding is unavailable"))
}

func ErrServiceBindingAvailable(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeServiceBindingAvailable, fmt.Sprintf("service binding is available"))
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

func ErrInvalidReqId(codespace sdk.CodespaceType, reqId string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidReqId, fmt.Sprintf("invalid request id [%s]", reqId))
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
	return sdk.NewError(codespace, CodeInvalidRequestInput, msg)
}

func ErrInvalidResponseOutput(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidResponseOutput, msg)
}

func ErrInvalidResponseErr(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidResponseErr, msg)
}

func ErrInvalidAddress(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAddress, msg)
}
