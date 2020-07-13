package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// oracle module sentinel errors
var (
	ErrUnknownFeedName      = sdkerrors.Register(ModuleName, 2, "unknown feed")
	ErrInvalidFeedName      = sdkerrors.Register(ModuleName, 2, "invalid feed name")
	ErrExistedFeedName      = sdkerrors.Register(ModuleName, 2, "feed already exists")
	ErrUnauthorized         = sdkerrors.Register(ModuleName, 2, "unauthorized owner")
	ErrInvalidServiceName   = sdkerrors.Register(ModuleName, 2, "invalid service name")
	ErrInvalidDescription   = sdkerrors.Register(ModuleName, 2, "invalid description")
	ErrNotRegisterFunc      = sdkerrors.Register(ModuleName, 2, "method don't register")
	ErrNotProfiler          = sdkerrors.Register(ModuleName, 2, "not a profiler address")
	ErrInvalidFeedState     = sdkerrors.Register(ModuleName, 2, "invalid state feed")
	ErrInvalidServiceFeeCap = sdkerrors.Register(ModuleName, 2, "service fee cap is invalid")
)
