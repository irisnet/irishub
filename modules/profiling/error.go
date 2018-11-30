package profiling

import (
	sdk "github.com/irisnet/irishub/types"
	"fmt"
)

const (
	DefaultCodespace sdk.CodespaceType = 25

	CodeProfilerExists    sdk.CodeType = 100
	CodeProfilerNotExists sdk.CodeType = 101
)

func ErrProfilerNotExists(codespace sdk.CodespaceType, profiler sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeProfilerNotExists, fmt.Sprintf("profiler %s is not existed", profiler))
}
func ErrProfilerExists(codespace sdk.CodespaceType, profiler sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeProfilerExists, fmt.Sprintf("profiler %s already exists", profiler))
}
