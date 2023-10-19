package types

import (
	errormod "cosmossdk.io/errors"
)

var (
	ErrInvalidCollection = errormod.Register(ModuleName, 9, "invalid nft collection")
	ErrUnknownCollection = errormod.Register(ModuleName, 10, "unknown nft collection")
	ErrInvalidNFT        = errormod.Register(ModuleName, 11, "invalid nft")
	ErrNFTAlreadyExists  = errormod.Register(ModuleName, 12, "nft already exists")
	ErrUnknownNFT        = errormod.Register(ModuleName, 13, "unknown nft")
	ErrEmptyTokenData    = errormod.Register(ModuleName, 14, "nft data can't be empty")
	ErrUnauthorized      = errormod.Register(ModuleName, 15, "unauthorized address")
	ErrInvalidDenom      = errormod.Register(ModuleName, 16, "invalid denom")
	ErrInvalidTokenID    = errormod.Register(ModuleName, 17, "invalid nft id")
	ErrInvalidTokenURI   = errormod.Register(ModuleName, 18, "invalid nft uri")
)
