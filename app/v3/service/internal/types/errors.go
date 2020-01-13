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
	CodeInvalidChainId           sdk.CodeType = 105

	CodeSvcBindingExists     sdk.CodeType = 106
	CodeSvcBindingNotExists  sdk.CodeType = 107
	CodeInvalidDefChainId    sdk.CodeType = 108
	CodeInvalidBindingType   sdk.CodeType = 109
	CodeInvalidLevel         sdk.CodeType = 110
	CodeInvalidPriceCount    sdk.CodeType = 111
	CodeInvalidRefundDeposit sdk.CodeType = 112
	CodeLtMinProviderDeposit sdk.CodeType = 113
	CodeInvalidDisable       sdk.CodeType = 114
	CodeInvalidEnable        sdk.CodeType = 115

	CodeMethodNotExists        sdk.CodeType = 116
	CodeRequestNotActive       sdk.CodeType = 117
	CodeReturnFeeNotExists     sdk.CodeType = 118
	CodeWithdrawFeeNotExists   sdk.CodeType = 119
	CodeLtServiceFee           sdk.CodeType = 120
	CodeInvalidReqId           sdk.CodeType = 121
	CodeSvcBindingNotAvailable sdk.CodeType = 122
	CodeNotMatchingProvider    sdk.CodeType = 123
	CodeInvalidReqChainId      sdk.CodeType = 124
	CodeInvalidBindChainId     sdk.CodeType = 125
	CodeNotMatchingReqChainID  sdk.CodeType = 126

	CodeInvalidRequestInput   sdk.CodeType = 127
	CodeInvalidResponseOutput sdk.CodeType = 128
	CodeInvalidResponseErr    sdk.CodeType = 129

	CodeIntOverflow    sdk.CodeType = 130
	CodeInvalidInput   sdk.CodeType = 131
	CodeInvalidAddress sdk.CodeType = 132
)

func ErrInvalidServiceName(codespace sdk.CodespaceType, serviceName string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidServiceName, fmt.Sprintf("invalid service name %s; only alphanumeric characters, _ and - accepted, the length ranges in (0,70]", serviceName))
}

func ErrInvalidSchemas(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidSchemas, msg)
}

func ErrServiceDefinitionExists(codespace sdk.CodespaceType, serviceName string) sdk.Error {
	return sdk.NewError(codespace, CodeServiceDefinitionExists, fmt.Sprintf("service name %s already exists", serviceName))
}

func ErrUnknownServiceDefinition(codespace sdk.CodespaceType, serviceName string) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownServiceDefinition, fmt.Sprintf("service name %s does not exist", serviceName))
}

func ErrInvalidLength(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidLength, msg)
}

func ErrInvalidChainId(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidChainId, fmt.Sprintf("chainId is empty"))
}

func ErrInvalidDefChainId(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDefChainId, fmt.Sprintf("defined chain id is empty"))
}

func ErrSvcBindingExists(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSvcBindingExists, fmt.Sprintf("service binding already exists"))
}

func ErrSvcBindingNotExists(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSvcBindingNotExists, fmt.Sprintf("service binding is not existed"))
}

func ErrInvalidBindingType(codespace sdk.CodespaceType, bindingType BindingType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidBindingType, fmt.Sprintf("invalid binding type %s", bindingType))
}

func ErrInvalidLevel(codespace sdk.CodespaceType, level Level) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidLevel, fmt.Sprintf("invalid level %v, avg_rsp_time and usable_time must be positive integer and usable_time limit to 10000", level))
}

func ErrInvalidPriceCount(codespace sdk.CodespaceType, priceCount int, methodCount int) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidPriceCount, fmt.Sprintf("invalid prices count %d, but methods count is %d", priceCount, methodCount))
}

func ErrRefundDeposit(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidRefundDeposit, fmt.Sprintf("can't refund deposit, %s", msg))
}

func ErrDisable(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDisable, fmt.Sprintf("can't disable, %s", msg))
}

func ErrEnable(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidEnable, fmt.Sprintf("can't enable, %s", msg))
}

func ErrLtMinProviderDeposit(codespace sdk.CodespaceType, coins sdk.Coins) sdk.Error {
	return sdk.NewError(codespace, CodeLtMinProviderDeposit, fmt.Sprintf("deposit amount must be equal or greater than %s", coins.String()))
}

func ErrMethodNotExists(codespace sdk.CodespaceType, methodID int16) sdk.Error {
	return sdk.NewError(codespace, CodeMethodNotExists, fmt.Sprintf("service method [%d] is not existed", methodID))
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

func ErrSvcBindingNotAvailable(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSvcBindingNotAvailable, fmt.Sprintf("service binding is unavailable"))
}

func ErrNotMatchingProvider(codespace sdk.CodespaceType, provider sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeNotMatchingProvider, fmt.Sprintf("[%s] is not a matching Provider", provider.String()))
}

func ErrInvalidReqChainId(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidReqChainId, fmt.Sprintf("request chain id is empty"))
}

func ErrInvalidBindChainId(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidBindChainId, fmt.Sprintf("bind chain id is empty"))
}

func ErrNotMatchingReqChainID(codespace sdk.CodespaceType, reqChainID string) sdk.Error {
	return sdk.NewError(codespace, CodeNotMatchingReqChainID, fmt.Sprintf("[%s] is not a matching reqChainID", reqChainID))
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

func ErrInvalidAddress(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAddress, msg)
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
