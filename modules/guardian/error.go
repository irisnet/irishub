package guardian

import (
	sdk "github.com/irisnet/irishub/types"
	"fmt"
)

const (
	DefaultCodespace sdk.CodespaceType = "guardian"

	CodeProfilerExists      sdk.CodeType = 100
	CodeProfilerNotExists   sdk.CodeType = 101
	CodeInvalidProfilerName sdk.CodeType = 102
)

func ErrProfilerNotExists(codespace sdk.CodespaceType, profiler sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeProfilerNotExists, fmt.Sprintf("profiler %s is not existed", profiler))
}

func ErrProfilerExists(codespace sdk.CodespaceType, profiler sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeProfilerExists, fmt.Sprintf("profiler %s already exists", profiler))
}

func ErrInvalidProfilerName(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidProfilerName, fmt.Sprintf("invalid profiler name %s, must contain alphanumeric characters, _ and - only，length greater than 0 and less than or equal to 128", msg))
}
