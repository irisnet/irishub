//nolint
package asset

import (
	"fmt"
	sdk "github.com/irisnet/irishub/types"
)

// Asset errors reserve 100 ~ 199.
const (
	DefaultCodespace sdk.CodespaceType = "asset"

	CodeInvalidMoniker       sdk.CodeType = 100
	CodeInvalidDetails       sdk.CodeType = 101
	CodeInvalidWebsite       sdk.CodeType = 102
	CodeUnknownGateway       sdk.CodeType = 103
	CodeGatewayAlreadyExists sdk.CodeType = 104
	CodeInvalidOwner         sdk.CodeType = 105
	CodeNoUpdatesProvided    sdk.CodeType = 106
	CodeInvalidAddress       sdk.CodeType = 107

	CodeNilAssetOwner                 sdk.CodeType = 110
	CodeInvalidAssetFamily            sdk.CodeType = 111
	CodeInvalidAssetSource            sdk.CodeType = 112
	CodeInvalidAssetName              sdk.CodeType = 113
	CodeInvalidAssetSymbol            sdk.CodeType = 114
	CodeInvalidAssetInitSupply        sdk.CodeType = 115
	CodeInvalidAssetMaxSupply         sdk.CodeType = 116
	CodeInvalidAssetDecimal           sdk.CodeType = 117
	CodeAssetAlreadyExists            sdk.CodeType = 118
	CodeUnauthorizedIssueGatewayAsset sdk.CodeType = 119
)

//----------------------------------------
// Asset error constructors

func ErrNilAssetOwner(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeNilAssetOwner, fmt.Sprintf("nil asset owner"))
}

func ErrInvalidAssetFamily(codespace sdk.CodespaceType, family byte) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetFamily, fmt.Sprintf("invalid asset family type %d", family))
}

func ErrInvalidAssetSource(codespace sdk.CodespaceType, source byte) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetSource, fmt.Sprintf("invalid asset source type %d", source))
}

func ErrInvalidAssetName(codespace sdk.CodespaceType, name string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetName, fmt.Sprintf("invalid asset name %s, only accepts alphanumeric characters, _ and -, length between 0 and 32", name))
}

func ErrInvalidAssetSymbol(codespace sdk.CodespaceType, symbol string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetSymbol, fmt.Sprintf("invalid asset symbol %s, only accepts alphanumeric characters, _ and -, length between 3 and 6", symbol))
}

func ErrInvalidAssetInitSupply(codespace sdk.CodespaceType, initSupply uint64) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetInitSupply, fmt.Sprintf("invalid asset initial supply %s", string(initSupply)))
}

func ErrInvalidAssetMaxSupply(codespace sdk.CodespaceType, maxSupply uint64) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetMaxSupply, fmt.Sprintf("invalid asset max supply %s", string(maxSupply)))
}

func ErrInvalidAssetDecimal(codespace sdk.CodespaceType, decimal uint8) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetDecimal, fmt.Sprintf("invalid asset decimal %s, max decimal is 18", string(decimal)))
}

func ErrAssetAlreadyExists(codespace sdk.CodespaceType, symbol string) sdk.Error {
	return sdk.NewError(codespace, CodeAssetAlreadyExists, fmt.Sprintf("asset already exists:%s", symbol))
}

//----------------------------------------
// Gateway error constructors

func ErrInvalidMoniker(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidMoniker, msg)
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

func ErrUnauthorizedIssueGatewayAsset(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeUnauthorizedIssueGatewayAsset, msg)
}
