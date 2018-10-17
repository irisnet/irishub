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

func ErrInvalidFilename(codespace sdk.CodespaceType, filename string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidFilename, fmt.Sprintf("File Name '%s' is not valid", filename))
}

func ErrInvalidDescription(codespace sdk.CodespaceType, description string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidFileDescription, fmt.Sprintf("File Desciption '%s' is not valid", description))
}

func ErrFailUploadFile(codespace sdk.CodespaceType, dataHash string) sdk.Error {
	return sdk.NewError(codespace, CodeFailUploadFile, fmt.Sprintf("File Upload failed with %s", dataHash))
}
