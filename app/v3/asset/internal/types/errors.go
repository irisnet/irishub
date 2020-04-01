//nolint
package types

import (
	sdk "github.com/irisnet/irishub/types"
)

// Asset errors reserve 100 ~ 199.
const (
	DefaultCodespace sdk.CodespaceType = "asset"

	CodeNilAssetOwner            sdk.CodeType = 100
	CodeInvalidAssetName         sdk.CodeType = 101
	CodeInvalidAssetSymbol       sdk.CodeType = 102
	CodeInvalidAssetMinUnitAlias sdk.CodeType = 103
	CodeInvalidAssetInitSupply   sdk.CodeType = 104
	CodeInvalidAssetMaxSupply    sdk.CodeType = 105
	CodeInvalidAssetDecimal      sdk.CodeType = 106
	CodeAssetAlreadyExists       sdk.CodeType = 107
	CodeAssetNotExists           sdk.CodeType = 108
	CodeAssetNotMintable         sdk.CodeType = 109
	CodeInvalidOwner             sdk.CodeType = 110
	CodeInvalidAddress           sdk.CodeType = 111
	CodeInvalidToAddress         sdk.CodeType = 112

	CodeInsufficientCoins       sdk.CodeType = 113
	CodeSignersMissingInContext sdk.CodeType = 114
)

//----------------------------------------
// Asset error constructors

func ErrNilAssetOwner(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeNilAssetOwner, msg)
}

func ErrInvalidAssetName(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetName, msg)
}

func ErrInvalidAssetMinUnitAlias(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetMinUnitAlias, msg)
}

func ErrInvalidAssetSymbol(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetSymbol, msg)
}
func ErrInvalidAssetInitSupply(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetInitSupply, msg)
}

func ErrInvalidAssetMaxSupply(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetMaxSupply, msg)
}

func ErrInvalidAssetDecimal(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetDecimal, msg)
}

func ErrAssetAlreadyExists(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeAssetAlreadyExists, msg)
}

func ErrAssetNotExists(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeAssetNotExists, msg)
}

func ErrAssetNotMintable(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeAssetNotMintable, msg)
}

func ErrInvalidOwner(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidOwner, msg)
}

func ErrInvalidAddress(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAddress, msg)
}

func ErrInvalidToAddress(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidToAddress, msg)
}

//----------------------------------------
// misc

func ErrInsufficientCoins(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInsufficientCoins, msg)
}

func ErrSignersMissingInContext(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeSignersMissingInContext, msg)
}
