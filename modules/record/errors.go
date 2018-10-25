//nolint
package record

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = 3

	CodeInvalidDataSize        sdk.CodeType = 1
	CodeInvalidFileDescription sdk.CodeType = 2
	CodeInvalidDataHash        sdk.CodeType = 3
)

func ErrInvalidDataSize(codespace sdk.CodespaceType, limit int64) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDataSize, fmt.Sprintf("Onchain data can't be empty and upload limit is %d bytes", limit))
}

func ErrInvalidDescription(codespace sdk.CodespaceType, limit int64) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidFileDescription, fmt.Sprintf("Descriprion can't be empty and upload limit is %d bytes", limit))
}

func ErrInvalidDataHash(codespace sdk.CodespaceType, hash string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDataHash, fmt.Sprintf("Data hash [%s] is invalid", hash))
}
