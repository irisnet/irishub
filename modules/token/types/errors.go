// nolint
package types

import (
	errorsmod "cosmossdk.io/errors"
)

// token module sentinel errors
var (
	ErrInvalidName          = errorsmod.Register(ModuleName, 2, "invalid token name")
	ErrInvalidMinUnit       = errorsmod.Register(ModuleName, 3, "invalid token min unit")
	ErrInvalidSymbol        = errorsmod.Register(ModuleName, 4, "invalid standard denom")
	ErrInvalidInitSupply    = errorsmod.Register(ModuleName, 5, "invalid token initial supply")
	ErrInvalidMaxSupply     = errorsmod.Register(ModuleName, 6, "invalid token maximum supply")
	ErrInvalidScale         = errorsmod.Register(ModuleName, 7, "invalid token scale")
	ErrSymbolAlreadyExists  = errorsmod.Register(ModuleName, 8, "symbol already exists")
	ErrMinUnitAlreadyExists = errorsmod.Register(ModuleName, 9, "min unit already exists")
	ErrTokenNotExists       = errorsmod.Register(ModuleName, 10, "token does not exist")
	ErrInvalidToAddress     = errorsmod.Register(ModuleName, 11, "the new owner must not be same as the original owner")
	ErrInvalidOwner         = errorsmod.Register(ModuleName, 12, "invalid token owner")
	ErrNotMintable          = errorsmod.Register(ModuleName, 13, "token is not mintable")
	ErrNotFoundTokenAmt     = errorsmod.Register(ModuleName, 14, "burned token amount not found")
	ErrInvalidAmount        = errorsmod.Register(ModuleName, 15, "invalid amount")
	ErrInvalidBaseFee       = errorsmod.Register(ModuleName, 16, "invalid base fee")
	ErrInvalidSwap          = errorsmod.Register(ModuleName, 17, "unregistered swapable fee token")
	ErrInsufficientFee      = errorsmod.Register(ModuleName, 18, "the amount of tokens after swap is less than 1")
)
