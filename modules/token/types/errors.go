//nolint
package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// token module sentinel errors
var (
	ErrNilOwner             = sdkerrors.Register(ModuleName, 2, "the owner of the token must be specified")
	ErrInvalidName          = sdkerrors.Register(ModuleName, 3, "invalid token name")
	ErrInvalidMinUnit       = sdkerrors.Register(ModuleName, 4, "invalid token min_unit")
	ErrInvalidSymbol        = sdkerrors.Register(ModuleName, 5, "must be standard denom")
	ErrInvalidInitSupply    = sdkerrors.Register(ModuleName, 6, "invalid token initial supply")
	ErrInvalidMaxSupply     = sdkerrors.Register(ModuleName, 7, "invalid token max supply")
	ErrInvalidScale         = sdkerrors.Register(ModuleName, 8, "invalid token scale")
	ErrSymbolAlreadyExists  = sdkerrors.Register(ModuleName, 9, "symbol has existed")
	ErrMinUnitAlreadyExists = sdkerrors.Register(ModuleName, 10, "min_unit has existed")
	ErrTokenNotExists       = sdkerrors.Register(ModuleName, 11, "token does not exist")
	ErrInvalidAddress       = sdkerrors.Register(ModuleName, 12, "the owner of the token must be specified")
	ErrInvalidToAddress     = sdkerrors.Register(ModuleName, 13, "the new owner must not be same as the original owner")
	ErrInvalidOwner         = sdkerrors.Register(ModuleName, 14, "invalid token owner")
	ErrNotMintable          = sdkerrors.Register(ModuleName, 15, "the token is set to be non-mintable")
)
