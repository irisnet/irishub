//nolint
package bank

import (
	"fmt"
	sdk "github.com/irisnet/irishub/types"
)

// Bank errors reserve 100 ~ 199.
const (
	DefaultCodespace sdk.CodespaceType = "bank"

	CodeInvalidInput       sdk.CodeType = 101
	CodeInvalidOutput      sdk.CodeType = 102
	CodeBurnEmptyCoins     sdk.CodeType = 103
	CodeFreezeEmptyCoin   sdk.CodeType = 104
	CodeUnfreezeEmptyCoin sdk.CodeType = 105
	CodeEmptyDenom         sdk.CodeType = 106
)

// NOTE: Don't stringer this, we'll put better messages in later.
func codeToDefaultMsg(code sdk.CodeType) string {
	switch code {
	case CodeInvalidInput:
		return "invalid input coins"
	case CodeInvalidOutput:
		return "invalid output coins"
	case CodeBurnEmptyCoins:
		return "burn empty coins"
	default:
		return sdk.CodeToDefaultMsg(code)
	}
}

//----------------------------------------
// Error constructors

func ErrFreezeEmptyCoin(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeFreezeEmptyCoin, fmt.Sprintf("freeze empty coins"))
}


func ErrEmptyDenom(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyDenom, fmt.Sprintf("empty denom for token"))
}

func ErrUnfreezeEmptyCoin(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeUnfreezeEmptyCoin, fmt.Sprintf("unfreeze empty coins"))
}

func ErrNoInputs(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeInvalidInput, "")
}

func ErrInvalidOutput(codespace sdk.CodespaceType, msg string) sdk.Error {
	return newError(codespace, CodeInvalidOutput, msg)
}

func ErrNoOutputs(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeInvalidOutput, "")
}

func ErrBurnEmptyCoins(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeBurnEmptyCoins, "")
}

//----------------------------------------

func msgOrDefaultMsg(msg string, code sdk.CodeType) string {
	if msg != "" {
		return msg
	}
	return codeToDefaultMsg(code)
}

func newError(codespace sdk.CodespaceType, code sdk.CodeType, msg string) sdk.Error {
	msg = msgOrDefaultMsg(msg, code)
	return sdk.NewError(codespace, code, msg)
}
