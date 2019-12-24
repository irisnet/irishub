package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Asset errors reserve 100 ~ 199.
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidIDL               sdk.CodeType = 100
	CodeInvalidLength            sdk.CodeType = 101
	CodeSvcDefExists             sdk.CodeType = 102
	CodeSvcDefNotExists          sdk.CodeType = 103
	CodeInvalidOutputPrivacyEnum sdk.CodeType = 104
	CodeInvalidOutputCachedEnum  sdk.CodeType = 105
	CodeInvalidServiceName       sdk.CodeType = 106
	CodeInvalidChainID           sdk.CodeType = 107
	CodeInvalidAuthor            sdk.CodeType = 108
	CodeInvalidMethodName        sdk.CodeType = 109

	CodeSvcBindingExists     sdk.CodeType = 110
	CodeSvcBindingNotExists  sdk.CodeType = 111
	CodeInvalidDefChainID    sdk.CodeType = 112
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
	CodeInvalidReqID           sdk.CodeType = 125
	CodeSvcBindingNotAvailable sdk.CodeType = 126
	CodeNotMatchingProvider    sdk.CodeType = 127
	CodeInvalidReqChainID      sdk.CodeType = 128
	CodeInvalidBindChainID     sdk.CodeType = 129
	CodeNotMatchingReqChainID  sdk.CodeType = 130

	CodeIntOverflow  sdk.CodeType = 131
	CodeInvalidInput sdk.CodeType = 132
)

// ErrSvcDefExists error for service definition already exists
func ErrSvcDefExists(codespace sdk.CodespaceType, defChainID, svcDefName string) sdk.Error {
	return sdk.NewError(codespace, CodeSvcDefExists, fmt.Sprintf("service definition name %s already exists in %s", svcDefName, defChainID))
}

// ErrSvcDefNotExists error for service definition not exists
func ErrSvcDefNotExists(codespace sdk.CodespaceType, defChainID, svcDefName string) sdk.Error {
	return sdk.NewError(codespace, CodeSvcDefNotExists, fmt.Sprintf("service definition name %s is not existed in %s", svcDefName, defChainID))
}

// ErrInvalidIDL error for invalid IDL
func ErrInvalidIDL(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidIDL, fmt.Sprintf("The IDL content cannot be parsed, %s", msg))
}

// ErrInvalidLength error for invalid length
func ErrInvalidLength(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidLength, msg)
}

// ErrInvalidOutputPrivacyEnum error for invalid output privacy enum
func ErrInvalidOutputPrivacyEnum(codespace sdk.CodespaceType, value string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidOutputPrivacyEnum, fmt.Sprintf("invalid OutputPrivacyEnum %s", value))
}

// ErrInvalidOutputCachedEnum error for invalid output cached enum
func ErrInvalidOutputCachedEnum(codespace sdk.CodespaceType, value string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidOutputCachedEnum, fmt.Sprintf("invalid OutputCachedEnum %s", value))
}

// ErrInvalidServiceName error for invalid service name
func ErrInvalidServiceName(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidServiceName, fmt.Sprintf("invalid service name %s, must contain alphanumeric characters, _ and - only，length greater than 0 and less than or equal to 128", msg))
}

// ErrInvalidChainID error for invalid chain id
func ErrInvalidChainID(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidChainID, fmt.Sprintf("chainID is empty"))
}

// ErrInvalidAuthor error for invalid author
func ErrInvalidAuthor(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAuthor, fmt.Sprintf("author is empty"))
}

// ErrInvalidMethodName error for invalid method name
func ErrInvalidMethodName(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidMethodName, fmt.Sprintf("method name is empty"))
}

// ErrInvalidDefChainID error for invalid defined chain id
func ErrInvalidDefChainID(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDefChainID, fmt.Sprintf("defined chain id is empty"))
}

// ErrSvcBindingExists error for service binding already exists
func ErrSvcBindingExists(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSvcBindingExists, fmt.Sprintf("service binding already exists"))
}

// ErrSvcBindingNotExists error for service binding not exists
func ErrSvcBindingNotExists(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSvcBindingNotExists, fmt.Sprintf("service binding is not existed"))
}

// ErrInvalidBindingType error for invalid binding type
func ErrInvalidBindingType(codespace sdk.CodespaceType, bindingType BindingType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidBindingType, fmt.Sprintf("invalid binding type %s", bindingType))
}

// ErrInvalidLevel error for invalid level
func ErrInvalidLevel(codespace sdk.CodespaceType, level Level) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidLevel, fmt.Sprintf("invalid level %v, avg_rsp_time and usable_time must be positive integer and usable_time limit to 10000", level))
}

// ErrInvalidPriceCount error for invalid price count
func ErrInvalidPriceCount(codespace sdk.CodespaceType, priceCount int, methodCount int) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidPriceCount, fmt.Sprintf("invalid prices count %d, but methods count is %d", priceCount, methodCount))
}

// ErrRefundDeposit error for can't refund deposit
func ErrRefundDeposit(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidRefundDeposit, fmt.Sprintf("can't refund deposit, %s", msg))
}

// ErrDisable error for can't disable
func ErrDisable(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDisable, fmt.Sprintf("can't disable, %s", msg))
}

// ErrEnable error for can't enable
func ErrEnable(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidEnable, fmt.Sprintf("can't enable, %s", msg))
}

// ErrLtMinProviderDeposit error for insufficient provider deposit
func ErrLtMinProviderDeposit(codespace sdk.CodespaceType, coins sdk.Coins) sdk.Error {
	return sdk.NewError(codespace, CodeLtMinProviderDeposit, fmt.Sprintf("deposit amount must be equal or greater than %s", coins.String()))
}

// ErrMethodNotExists error for method not exists
func ErrMethodNotExists(codespace sdk.CodespaceType, methodID int16) sdk.Error {
	return sdk.NewError(codespace, CodeMethodNotExists, fmt.Sprintf("service method [%d] is not existed", methodID))
}

// ErrRequestNotActive error for request not active
func ErrRequestNotActive(codespace sdk.CodespaceType, requestID string) sdk.Error {
	return sdk.NewError(codespace, CodeRequestNotActive, fmt.Sprintf("request [%s] is not existed", requestID))
}

// ErrReturnFeeNotExists error for return fee not exists
func ErrReturnFeeNotExists(codespace sdk.CodespaceType, address sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeReturnFeeNotExists, fmt.Sprintf("There is no service refund fees for [%s]", address))
}

// ErrWithdrawFeeNotExists error for withdraw fee not exists
func ErrWithdrawFeeNotExists(codespace sdk.CodespaceType, address sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeWithdrawFeeNotExists, fmt.Sprintf("There is no service withdraw fees for [%s]", address))
}

// ErrLtServiceFee error for insufficient service fee
func ErrLtServiceFee(codespace sdk.CodespaceType, coins sdk.Coins) sdk.Error {
	return sdk.NewError(codespace, CodeLtServiceFee, fmt.Sprintf("service fee amount must be equal or greater than %s", coins.String()))
}

// ErrInvalidReqID error for invalid request id
func ErrInvalidReqID(codespace sdk.CodespaceType, reqID string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidReqID, fmt.Sprintf("invalid request id [%s]", reqID))
}

// ErrSvcBindingNotAvailable error for service binding not available
func ErrSvcBindingNotAvailable(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSvcBindingNotAvailable, fmt.Sprintf("service binding is unavailable"))
}

// ErrNotMatchingProvider error for not matching provider
func ErrNotMatchingProvider(codespace sdk.CodespaceType, provider sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeNotMatchingProvider, fmt.Sprintf("[%s] is not a matching Provider", provider.String()))
}

// ErrInvalidReqChainID error for invalid request chain id
func ErrInvalidReqChainID(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidReqChainID, fmt.Sprintf("request chain id is empty"))
}

// ErrInvalidBindChainID error for invalid bind chain id
func ErrInvalidBindChainID(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidBindChainID, fmt.Sprintf("bind chain id is empty"))
}

// ErrNotMatchingReqChainID error for not matching request chain id
func ErrNotMatchingReqChainID(codespace sdk.CodespaceType, reqChainID string) sdk.Error {
	return sdk.NewError(codespace, CodeNotMatchingReqChainID, fmt.Sprintf("[%s] is not a matching reqChainID", reqChainID))
}

// ErrNotTrustee error for not trustee
func ErrNotTrustee(codespace sdk.CodespaceType, trustee sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, fmt.Sprintf("[%s] is not a trustee address", trustee))
}

// ErrNotProfiler error for not profiler
func ErrNotProfiler(codespace sdk.CodespaceType, profiler sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, fmt.Sprintf("[%s] is not a profiler address", profiler))
}

// ErrNoResponseFound error for no response found
func ErrNoResponseFound(codespace sdk.CodespaceType, requestID string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, fmt.Sprintf("response is not existed for request %s", requestID))
}
