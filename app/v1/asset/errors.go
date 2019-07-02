//nolint
package asset

import (
	sdk "github.com/irisnet/irishub/types"
)

// Asset errors reserve 100 ~ 199.
const (
	DefaultCodespace sdk.CodespaceType = "asset"

	CodeInvalidMoniker       sdk.CodeType = 100
	CodeInvalidIdentity      sdk.CodeType = 101
	CodeInvalidDetails       sdk.CodeType = 102
	CodeInvalidWebsite       sdk.CodeType = 103
	CodeUnknownGateway       sdk.CodeType = 104
	CodeGatewayAlreadyExists sdk.CodeType = 105
	CodeInvalidOwner         sdk.CodeType = 106
	CodeNoUpdatesProvided    sdk.CodeType = 107
	CodeInvalidAddress       sdk.CodeType = 108
	CodeInvalidToAddress     sdk.CodeType = 109

	CodeNilAssetOwner                 sdk.CodeType = 110
	CodeInvalidAssetFamily            sdk.CodeType = 111
	CodeInvalidAssetSource            sdk.CodeType = 112
	CodeInvalidAssetName              sdk.CodeType = 113
	CodeInvalidAssetSymbol            sdk.CodeType = 114
	CodeInvalidAssetSymbolAtSource    sdk.CodeType = 115
	CodeInvalidAssetSymbolMinAlias    sdk.CodeType = 116
	CodeInvalidAssetInitSupply        sdk.CodeType = 117
	CodeInvalidAssetMaxSupply         sdk.CodeType = 118
	CodeInvalidAssetDecimal           sdk.CodeType = 119
	CodeAssetAlreadyExists            sdk.CodeType = 120
	CodeUnauthorizedIssueGatewayAsset sdk.CodeType = 121
	CodeAssetNotExists                sdk.CodeType = 122
	CodeAssetNotMintable              sdk.CodeType = 123

	CodeInsufficientCoins sdk.CodeType = 130
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

func ErrInvalidAssetSymbolAtSource(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetSymbolAtSource, msg)
}

func ErrInvalidAssetSymbolMinAlias(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetSymbolMinAlias, msg)
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

func ErrInvalidIdentity(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidIdentity, msg)
}

func ErrInvalidDetails(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDetails, msg)
}

func ErrInvalidWebsite(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidWebsite, msg)
}

func ErrUnkwownGateway(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownGateway, msg)
}

func ErrGatewayAlreadyExists(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeGatewayAlreadyExists, msg)
}

func ErrInvalidOwner(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidOwner, msg)
}

func ErrNoUpdatesProvided(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeNoUpdatesProvided, msg)
}

func ErrInvalidAddress(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAddress, msg)
}

func ErrInvalidToAddress(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidToAddress, msg)
}

func ErrUnauthorizedIssueGatewayAsset(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeUnauthorizedIssueGatewayAsset, msg)
}

//----------------------------------------
// fee error constructor

func ErrInsufficientCoins(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInsufficientCoins, msg)
}
