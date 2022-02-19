package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidCollection = sdkerrors.Register(ModuleName, 2, "invalid mt collection")
	ErrUnknownCollection = sdkerrors.Register(ModuleName, 3, "unknown mt collection")
	ErrInvalidMT        = sdkerrors.Register(ModuleName, 4, "invalid mt")
	ErrMTAlreadyExists  = sdkerrors.Register(ModuleName, 5, "mt already exists")
	ErrUnknownMT        = sdkerrors.Register(ModuleName, 6, "unknown mt")
	ErrEmptyTokenData    = sdkerrors.Register(ModuleName, 7, "mt data can't be empty")
	ErrUnauthorized      = sdkerrors.Register(ModuleName, 8, "unauthorized address")
	ErrInvalidDenom      = sdkerrors.Register(ModuleName, 9, "invalid denom")
	ErrInvalidTokenID    = sdkerrors.Register(ModuleName, 10, "invalid mt id")
	ErrInvalidTokenURI   = sdkerrors.Register(ModuleName, 11, "invalid mt uri")
)
