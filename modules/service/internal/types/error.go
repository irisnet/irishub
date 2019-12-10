package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidIDL               sdk.CodeType = 100
	CodeInvalidLength            sdk.CodeType = 101
	CodeSvcDefExists             sdk.CodeType = 102
	CodeSvcDefNotExists          sdk.CodeType = 103
	CodeInvalidOutputPrivacyEnum sdk.CodeType = 104
	CodeInvalidOutputCachedEnum  sdk.CodeType = 105
	CodeInvalidServiceName       sdk.CodeType = 106
	CodeInvalidChainId           sdk.CodeType = 107
	CodeInvalidAuthor            sdk.CodeType = 108
	CodeInvalidMethodName        sdk.CodeType = 109

	CodeSvcBindingExists     sdk.CodeType = 110
	CodeSvcBindingNotExists  sdk.CodeType = 111
	CodeInvalidDefChainId    sdk.CodeType = 112
	CodeInvalidBindingType   sdk.CodeType = 113
	CodeInvalidLevel         sdk.CodeType = 114
	CodeInvalidPriceCount    sdk.CodeType = 115
	CodeInvalidRefundDeposit sdk.CodeType = 116
	CodeLtMinProviderDeposit sdk.CodeType = 117
	CodeInvalidDisable       sdk.CodeType = 118
	CodeInvalidEnable        sdk.CodeType = 119

	CodeMethodNotExists        sdk.CodeType = 120
	CodeRequestNotActive       sdk.CodeType = 121
	CodeReturnFeeNotExists     sdk.CodeType = 122
	CodeWithdrawFeeNotExists   sdk.CodeType = 123
	CodeLtServiceFee           sdk.CodeType = 124
	CodeInvalidReqId           sdk.CodeType = 125
	CodeSvcBindingNotAvailable sdk.CodeType = 126
	CodeNotMatchingProvider    sdk.CodeType = 127
	CodeInvalidReqChainId      sdk.CodeType = 128
	CodeInvalidBindChainId     sdk.CodeType = 129
	CodeNotMatchingReqChainID  sdk.CodeType = 130

	CodeIntOverflow  sdk.CodeType = 131
	CodeInvalidInput sdk.CodeType = 132
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

func ErrSvcDefExists(codespace sdk.CodespaceType, defChainId, svcDefName string) sdk.Error {
	return sdk.NewError(codespace, CodeSvcDefExists, fmt.Sprintf("service definition name %s already exists in %s", svcDefName, defChainId))
}

func ErrSvcDefNotExists(codespace sdk.CodespaceType, defChainId, svcDefName string) sdk.Error {
	return sdk.NewError(codespace, CodeSvcDefNotExists, fmt.Sprintf("service definition name %s is not existed in %s", svcDefName, defChainId))
}

func ErrInvalidIDL(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidIDL, fmt.Sprintf("The IDL content cannot be parsed, %s", msg))
}

func ErrInvalidLength(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidLength, msg)
}

func ErrInvalidOutputPrivacyEnum(codespace sdk.CodespaceType, value string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidOutputPrivacyEnum, fmt.Sprintf("invalid OutputPrivacyEnum %s", value))
}

func ErrInvalidOutputCachedEnum(codespace sdk.CodespaceType, value string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidOutputCachedEnum, fmt.Sprintf("invalid OutputCachedEnum %s", value))
}

func ErrInvalidServiceName(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidServiceName, fmt.Sprintf("invalid service name %s, must contain alphanumeric characters, _ and - only，length greater than 0 and less than or equal to 128", msg))
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
