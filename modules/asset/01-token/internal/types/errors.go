package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Asset errors reserve 100 ~ 199.
const (
	DefaultCodespace sdk.CodespaceType = SubModuleName

	CodeInvalidOwner     sdk.CodeType = 100
	CodeInvalidAddress   sdk.CodeType = 101
	CodeInvalidToAddress sdk.CodeType = 102

	CodeNilAssetOwner          sdk.CodeType = 103
	CodeInvalidAssetName       sdk.CodeType = 104
	CodeInvalidAssetSymbol     sdk.CodeType = 105
	CodeInvalidAssetMinUnit    sdk.CodeType = 106
	CodeInvalidAssetInitSupply sdk.CodeType = 107
	CodeInvalidAssetMaxSupply  sdk.CodeType = 108
	CodeInvalidAssetScale      sdk.CodeType = 109
	CodeTokenAlreadyExists     sdk.CodeType = 110
	CodeTokenNotExists         sdk.CodeType = 111
	CodeTokenNotMintable       sdk.CodeType = 112
	CodeInvalidMintAmount      sdk.CodeType = 113
)

//----------------------------------------
// Asset error constructors

// ErrNilAssetOwner error for nil asset owner
func ErrNilAssetOwner(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeNilAssetOwner, msg)
}

// ErrInvalidAssetName error for invalid asset name
func ErrInvalidAssetName(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetName, msg)
}

// ErrInvalidAssetMinUnit error for invalid asset min unit alias
func ErrInvalidAssetMinUnit(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetMinUnit, msg)
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

// ErrInvalidAssetScale error for invalid asset decimal
func ErrInvalidAssetScale(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetScale, msg)
}

// ErrAssetAlreadyExists error for invalid asset already exists
func ErrAssetAlreadyExists(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenAlreadyExists, msg)
}

// ErrAssetNotExists error for asset not exists
func ErrAssetNotExists(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenNotExists, msg)
}

// ErrAssetNotMintable error for asset not mintable
func ErrAssetNotMintable(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenNotMintable, msg)
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

// ErrInvalidMintAmount error for invalid to amount
func ErrInvalidMintAmount(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidMintAmount, msg)
}
