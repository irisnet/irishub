//nolint
package record

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = 7

	CodeInvalidFilename        sdk.CodeType = 1
	CodeInvalidFileDescription sdk.CodeType = 2
	CodeFailUploadFile         sdk.CodeType = 3
)

func ErrInvalidDataSize(codespace sdk.CodespaceType, size int64) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidFilename, fmt.Sprintf("Onchain data can't be empty and upload limit is %d bytes", size))
}

func ErrInvalidDescription(codespace sdk.CodespaceType, size int64) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidFileDescription, fmt.Sprintf("Descriprion can't be empty and upload limit is %d bytes", size))
}

func ErrInvalidDataHash(codespace sdk.CodespaceType, hash string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidFilename, fmt.Sprintf("Data hash [%s] is invalid", hash))
}
