package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// random module sentinel errors
var (
	ErrInvalidReqID            = sdkerrors.Register(ModuleName, 2, "invalid request id")
	ErrInvalidHeight           = sdkerrors.Register(ModuleName, 3, "invalid height, must be greater than 0")
	ErrInvalidServiceBindings  = sdkerrors.Register(ModuleName, 4, "no service bindings available")
	ErrInvalidRequestContextID = sdkerrors.Register(ModuleName, 5, "invalid request context id")
	ErrInvalidServiceFeeCap    = sdkerrors.Register(ModuleName, 6, "invalid service fee cap")
)
