package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// coinswap module sentinel errors
var (
	ErrReservePoolNotExists    = sdkerrors.Register(ModuleName, 2, "reserve pool not exists")
	ErrEqualDenom              = sdkerrors.Register(ModuleName, 3, "input and output denomination are equal")
	ErrNotContainStandardDenom = sdkerrors.Register(ModuleName, 4, "must have one standard denom")
	ErrMustStandardDenom       = sdkerrors.Register(ModuleName, 5, "must be standard denom")
	ErrInvalidDenom            = sdkerrors.Register(ModuleName, 6, "invalid denom")
	ErrInvalidDeadline         = sdkerrors.Register(ModuleName, 7, "invalid deadline")
	ErrConstraintNotMet        = sdkerrors.Register(ModuleName, 8, "constraint not met")
	ErrInsufficientFunds       = sdkerrors.Register(ModuleName, 9, "insufficient funds")
)
