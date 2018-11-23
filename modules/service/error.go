package service

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"fmt"
)

const (
	DefaultCodespace sdk.CodespaceType = 6

	CodeInvalidIDL               sdk.CodeType = 100
	CodeSvcDefExists             sdk.CodeType = 101
	CodeSvcDefNotExists          sdk.CodeType = 102
	CodeInvalidOutputPrivacyEnum sdk.CodeType = 103
	CodeInvalidOutputCachedEnum  sdk.CodeType = 104
	CodeInvalidServiceName       sdk.CodeType = 105
	CodeInvalidChainId           sdk.CodeType = 106
	CodeInvalidAuthor            sdk.CodeType = 107
	CodeInvalidMethodName        sdk.CodeType = 108

	CodeSvcBindingExists     sdk.CodeType = 109
	CodeSvcBindingNotExists  sdk.CodeType = 110
	CodeInvalidDefChainId    sdk.CodeType = 111
	CodeInvalidBindingType   sdk.CodeType = 112
	CodeInvalidLevel         sdk.CodeType = 113
	CodeInvalidPriceCount    sdk.CodeType = 114
	CodeInvalidRefundDeposit sdk.CodeType = 115
	CodeLtMinProviderDeposit sdk.CodeType = 116
	CodeInvalidDisable       sdk.CodeType = 117
	CodeInvalidEnable        sdk.CodeType = 118

	CodeMethodNotExists      sdk.CodeType = 119
	CodeRequestNotActive     sdk.CodeType = 120
	CodeReturnFeeNotExists   sdk.CodeType = 121
	CodeWithdrawFeeNotExists sdk.CodeType = 122
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
	return sdk.NewError(codespace, CodeInvalidDefChainId, fmt.Sprintf("def-chain-id is empty"))
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

func ErrRequestNotActive(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeRequestNotActive, fmt.Sprintf("can not find request"))
}

func ErrReturnFeeNotExists(codespace sdk.CodespaceType, address sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeReturnFeeNotExists, fmt.Sprintf("There is no service refund fees for [%s]", address))
}

func ErrWithdrawFeeNotExists(codespace sdk.CodespaceType, address sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeWithdrawFeeNotExists, fmt.Sprintf("There is no service withdraw fees for [%s]", address))
}
