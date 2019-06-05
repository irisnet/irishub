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
	CodeInvalidOperator      sdk.CodeType = 105
	CodeInvalidOwner         sdk.CodeType = 106
	CodeNoUpdatesProvided    sdk.CodeType = 107
	CodeInvalidAddress       sdk.CodeType = 108
	CodeInvalidGenesis       sdk.CodeType = 109

	CodeNilAssetOwner          sdk.CodeType = 110
	CodeInvalidAssetFamily     sdk.CodeType = 111
	CodeInvalidAssetName       sdk.CodeType = 112
	CodeInvalidAssetSymbol     sdk.CodeType = 113
	CodeInvalidAssetInitSupply sdk.CodeType = 114
	CodeInvalidAssetMaxSupply  sdk.CodeType = 115
	CodeInvalidAssetDecimal    sdk.CodeType = 116
)

//----------------------------------------
// Asset Error constructors

func ErrNilAssetOwner(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeNilAssetOwner, fmt.Sprintf("nil asset owner"))
}

func ErrInvalidAssetFamily(codespace sdk.CodespaceType, family string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAssetFamily, fmt.Sprintf("invalid asset family %s, only accepts 00, 01", family))
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

//----------------------------------------
// Error constructors

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

func ErrInvalidOperator(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidOperator, msg)
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

func ErrInvalidGenesis(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidGenesis, msg)
}
