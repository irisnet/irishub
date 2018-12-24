package guardian

import (
	sdk "github.com/irisnet/irishub/types"
	"fmt"
)

const (
	DefaultCodespace sdk.CodespaceType = "guardian"

	CodeProfilerExists     sdk.CodeType = 100
	CodeProfilerNotExists  sdk.CodeType = 101
	CodeTrusteeExists      sdk.CodeType = 102
	CodeTrusteeNotExists   sdk.CodeType = 103
	CodeInvalidDescription sdk.CodeType = 104
)

func ErrInvalidOperator(codespace sdk.CodespaceType, operator sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeProfilerNotExists, fmt.Sprintf("%s is not valid operator", operator))
}

func ErrProfilerNotExists(codespace sdk.CodespaceType, profiler sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeProfilerNotExists, fmt.Sprintf("profiler %s is not existed", profiler))
}

func ErrProfilerExists(codespace sdk.CodespaceType, profiler sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeProfilerExists, fmt.Sprintf("profiler %s already exists", profiler))
}

func ErrTrusteeExists(codespace sdk.CodespaceType, trustee sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeTrusteeExists, fmt.Sprintf("trustee %s already exists", trustee))
}

func ErrTrusteeNotExists(codespace sdk.CodespaceType, profiler sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeTrusteeNotExists, fmt.Sprintf("trustee %s is not existed", profiler))
}

func ErrInvalidDescription(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDescription, "description is empty")
}
