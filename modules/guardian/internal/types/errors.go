package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Asset errors reserve 100 ~ 199.
const (
	DefaultCodespace sdk.CodespaceType = "guardian"

	CodeInvalidOperator       sdk.CodeType = 100
	CodeProfilerExists        sdk.CodeType = 101
	CodeProfilerNotExists     sdk.CodeType = 102
	CodeTrusteeExists         sdk.CodeType = 103
	CodeTrusteeNotExists      sdk.CodeType = 104
	CodeInvalidDescription    sdk.CodeType = 105
	CodeDeleteGenesisProfiler sdk.CodeType = 106
	CodeDeleteGenesisTrustee  sdk.CodeType = 107
	CodeInvalidGuardian       sdk.CodeType = 108
)

// ErrInvalidOperator error for invalid operator
func ErrInvalidOperator(codespace sdk.CodespaceType, operator sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidOperator, fmt.Sprintf("%s is not a valid operator", operator))
}

// ErrProfilerNotExists error for profiler not exists
func ErrProfilerNotExists(codespace sdk.CodespaceType, profiler sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeProfilerNotExists, fmt.Sprintf("profiler %s is not existed", profiler))
}

// ErrDeleteGenesisProfiler error for delete genesis profiler
func ErrDeleteGenesisProfiler(codespace sdk.CodespaceType, profiler sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeDeleteGenesisProfiler, fmt.Sprintf("can't delete profiler %s that in genesis", profiler))
}

// ErrProfilerExists error for profiler exists
func ErrProfilerExists(codespace sdk.CodespaceType, profiler sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeProfilerExists, fmt.Sprintf("profiler %s already exists", profiler))
}

// ErrTrusteeExists error for trustee exists
func ErrTrusteeExists(codespace sdk.CodespaceType, trustee sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeTrusteeExists, fmt.Sprintf("trustee %s already exists", trustee))
}

// ErrTrusteeNotExists error for trustee not exists
func ErrTrusteeNotExists(codespace sdk.CodespaceType, trustee sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeTrusteeNotExists, fmt.Sprintf("trustee %s is not existed", trustee))
}

// ErrDeleteGenesisTrustee error for delete genesis trustee
func ErrDeleteGenesisTrustee(codespace sdk.CodespaceType, trustee sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeDeleteGenesisTrustee, fmt.Sprintf("can't delete trustee %s that in genesis", trustee))
}

// ErrInvalidDescription error for invalid description
func ErrInvalidDescription(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDescription, "description is empty")
}
