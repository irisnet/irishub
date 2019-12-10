package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Asset errors reserve 100 ~ 199.
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidMoniker   sdk.CodeType = 100
	CodeInvalidOwner     sdk.CodeType = 101
	CodeInvalidAddress   sdk.CodeType = 102
	CodeInvalidToAddress sdk.CodeType = 103

	CodeNilAssetOwner               sdk.CodeType = 104
	CodeInvalidAssetFamily          sdk.CodeType = 105
	CodeInvalidAssetSource          sdk.CodeType = 106
	CodeInvalidAssetName            sdk.CodeType = 107
	CodeInvalidAssetSymbol          sdk.CodeType = 108
	CodeInvalidAssetCanonicalSymbol sdk.CodeType = 109
	CodeInvalidAssetMinUnitAlias    sdk.CodeType = 110
	CodeInvalidAssetInitSupply      sdk.CodeType = 111
	CodeInvalidAssetMaxSupply       sdk.CodeType = 112
	CodeInvalidAssetDecimal         sdk.CodeType = 113
	CodeAssetAlreadyExists          sdk.CodeType = 114
	CodeAssetNotExists              sdk.CodeType = 115
	CodeAssetNotMintable            sdk.CodeType = 116
)

//----------------------------------------
// Asset error constructors

func ErrNilAssetOwner(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeNilAssetOwner, msg)
}

func ErrInvalidAssetFamily(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetFamily, msg)
}

func ErrInvalidAssetSource(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetSource, msg)
}

func ErrInvalidAssetName(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetName, msg)
}

func ErrInvalidAssetCanonicalSymbol(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetCanonicalSymbol, msg)
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

//----------------------------------------
// Gateway error constructors

func ErrInvalidMoniker(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidMoniker, msg)
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
