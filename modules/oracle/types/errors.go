package types

import (
	errorsmod "cosmossdk.io/errors"
)

// oracle module sentinel errors
var (
	ErrUnknownFeedName      = errorsmod.Register(ModuleName, 2, "unknown feed")
	ErrInvalidFeedName      = errorsmod.Register(ModuleName, 3, "invalid feed name")
	ErrExistedFeedName      = errorsmod.Register(ModuleName, 4, "feed already exists")
	ErrUnauthorized         = errorsmod.Register(ModuleName, 5, "unauthorized owner")
	ErrInvalidServiceName   = errorsmod.Register(ModuleName, 6, "invalid service name")
	ErrInvalidDescription   = errorsmod.Register(ModuleName, 7, "invalid description")
	ErrNotRegisterFunc      = errorsmod.Register(ModuleName, 8, "method don't register")
	ErrInvalidFeedState     = errorsmod.Register(ModuleName, 9, "invalid state feed")
	ErrInvalidServiceFeeCap = errorsmod.Register(ModuleName, 10, "service fee cap is invalid")
)
