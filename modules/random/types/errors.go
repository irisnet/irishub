package types

import (
	errorsmod "cosmossdk.io/errors"
)

// random module sentinel errors
var (
	ErrInvalidReqID            = errorsmod.Register(ModuleName, 2, "invalid request id")
	ErrInvalidHeight           = errorsmod.Register(ModuleName, 3, "invalid height, must be greater than 0")
	ErrInvalidServiceBindings  = errorsmod.Register(ModuleName, 4, "no service bindings available")
	ErrInvalidRequestContextID = errorsmod.Register(ModuleName, 5, "invalid request context id")
	ErrInvalidServiceFeeCap    = errorsmod.Register(ModuleName, 6, "invalid service fee cap")
)
