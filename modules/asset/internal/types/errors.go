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

// ErrNilAssetOwner error for nil asset owner
func ErrNilAssetOwner(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeNilAssetOwner, msg)
}

// ErrInvalidAssetFamily error for invalid asset family
func ErrInvalidAssetFamily(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetFamily, msg)
}

// ErrInvalidAssetSource error for invalid asset source
func ErrInvalidAssetSource(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetSource, msg)
}

// ErrInvalidAssetName error for invalid asset name
func ErrInvalidAssetName(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetName, msg)
}

// ErrInvalidAssetCanonicalSymbol error for invalid asset canonical symbol
func ErrInvalidAssetCanonicalSymbol(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetCanonicalSymbol, msg)
}

// ErrInvalidAssetMinUnitAlias error for invalid asset min unit alias
func ErrInvalidAssetMinUnitAlias(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetMinUnitAlias, msg)
}

// ErrInvalidAssetSymbol error for invalid asset symbol
func ErrInvalidAssetSymbol(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetSymbol, msg)
}

// ErrInvalidAssetInitSupply error for invalid asset init supply
func ErrInvalidAssetInitSupply(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetInitSupply, msg)
}

// ErrInvalidAssetMaxSupply error for invalid asset max supply
func ErrInvalidAssetMaxSupply(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetMaxSupply, msg)
}

// ErrInvalidAssetDecimal error for invalid asset decimal
func ErrInvalidAssetDecimal(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetDecimal, msg)
}

// ErrAssetAlreadyExists error for invalid asset already exists
func ErrAssetAlreadyExists(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeAssetAlreadyExists, msg)
}

// ErrAssetNotExists error for asset not exists
func ErrAssetNotExists(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeAssetNotExists, msg)
}

// ErrAssetNotMintable error for asset not mintable
func ErrAssetNotMintable(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeAssetNotMintable, msg)
}

//----------------------------------------
// Gateway error constructors

// ErrInvalidMoniker error for invalid moniker
func ErrInvalidMoniker(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidMoniker, msg)
}

// ErrInvalidOwner error for invalid owner
func ErrInvalidOwner(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidOwner, msg)
}

// ErrInvalidAddress error for invalid address
func ErrInvalidAddress(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAddress, msg)
}

// ErrInvalidToAddress error for invalid to address
func ErrInvalidToAddress(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidToAddress, msg)
}
