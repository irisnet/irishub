package types

import (
	errorsmod "cosmossdk.io/errors"
)

// coinswap module sentinel errors
var (
	ErrReservePoolNotExists    = errorsmod.Register(ModuleName, 2, "reserve pool not exists")
	ErrEqualDenom              = errorsmod.Register(ModuleName, 3, "input and output denomination are equal")
	ErrNotContainStandardDenom = errorsmod.Register(ModuleName, 4, "must have one standard denom")
	ErrMustStandardDenom       = errorsmod.Register(ModuleName, 5, "must be standard denom")
	ErrInvalidDenom            = errorsmod.Register(ModuleName, 6, "invalid denom")
	ErrInvalidDeadline         = errorsmod.Register(ModuleName, 7, "invalid deadline")
	ErrConstraintNotMet        = errorsmod.Register(ModuleName, 8, "constraint not met")
	ErrInsufficientFunds       = errorsmod.Register(ModuleName, 9, "insufficient funds")
)
