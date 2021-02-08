//nolint
package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// token module sentinel errors
var (
	ErrInvalidName          = sdkerrors.Register(ModuleName, 2, "invalid token name")
	ErrInvalidMinUnit       = sdkerrors.Register(ModuleName, 3, "invalid token min unit")
	ErrInvalidSymbol        = sdkerrors.Register(ModuleName, 4, "invalid standard denom")
	ErrInvalidInitSupply    = sdkerrors.Register(ModuleName, 5, "invalid token initial supply")
	ErrInvalidMaxSupply     = sdkerrors.Register(ModuleName, 6, "invalid token maximum supply")
	ErrInvalidScale         = sdkerrors.Register(ModuleName, 7, "invalid token scale")
	ErrSymbolAlreadyExists  = sdkerrors.Register(ModuleName, 8, "symbol already exists")
	ErrMinUnitAlreadyExists = sdkerrors.Register(ModuleName, 9, "min unit already exists")
	ErrTokenNotExists       = sdkerrors.Register(ModuleName, 10, "token does not exist")
	ErrInvalidToAddress     = sdkerrors.Register(ModuleName, 11, "the new owner must not be same as the original owner")
	ErrInvalidOwner         = sdkerrors.Register(ModuleName, 12, "invalid token owner")
	ErrNotMintable          = sdkerrors.Register(ModuleName, 13, "token is not mintable")
	ErrNotFoundTokenAmt     = sdkerrors.Register(ModuleName, 14, "burned token amount not found")
	ErrInvalidAmount        = sdkerrors.Register(ModuleName, 15, "invalid amount")
	ErrInvalidBaseFee       = sdkerrors.Register(ModuleName, 16, "invalid base fee")
)
