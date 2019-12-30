package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// TODO: use ModuleName + SubModuleName
// asset/token module sentinel errors
var (
	ErrUnknownToken          = sdkerrors.Register(SubModuleName, 1, "unknown token")
	ErrInvalidAssetMaxSupply = sdkerrors.Register(SubModuleName, 2, "invalid token max supply")
	ErrTokenAlreadyExists    = sdkerrors.Register(SubModuleName, 3, "token symbol already exists")
	ErrTokenNotExists        = sdkerrors.Register(SubModuleName, 4, "token does not exist")
	ErrInvalidOwner          = sdkerrors.Register(SubModuleName, 5, "invalid token owner")
	ErrAssetNotMintable      = sdkerrors.Register(SubModuleName, 6, "the token is non-mintable")
)
