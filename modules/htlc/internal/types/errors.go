package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// rand module sentinel errors
var (
	ErrUnknownHTLC           = sdkerrors.Register(ModuleName, 1, "unknown HTLC")
	ErrInvalidHashLock       = sdkerrors.Register(ModuleName, 2, "invalid hash lock")
	ErrInvalidSecret         = sdkerrors.Register(ModuleName, 3, "invalid secret")
	ErrHashLockAlreadyExists = sdkerrors.Register(ModuleName, 4, "hash lock already exists")
	ErrStateIsNotOpen        = sdkerrors.Register(ModuleName, 5, "the HTLC is not open")
	ErrStateIsNotExpired     = sdkerrors.Register(ModuleName, 6, "the HTLC is not expired")
)
