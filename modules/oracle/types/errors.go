package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// oracle module sentinel errors
var (
	ErrUnknownFeedName      = sdkerrors.Register(ModuleName, 2, "unknown feed")
	ErrInvalidFeedName      = sdkerrors.Register(ModuleName, 3, "invalid feed name")
	ErrExistedFeedName      = sdkerrors.Register(ModuleName, 4, "feed already exists")
	ErrUnauthorized         = sdkerrors.Register(ModuleName, 5, "unauthorized owner")
	ErrInvalidServiceName   = sdkerrors.Register(ModuleName, 6, "invalid service name")
	ErrInvalidDescription   = sdkerrors.Register(ModuleName, 7, "invalid description")
	ErrNotRegisterFunc      = sdkerrors.Register(ModuleName, 8, "method don't register")
	ErrNotProfiler          = sdkerrors.Register(ModuleName, 9, "not a profiler address")
	ErrInvalidFeedState     = sdkerrors.Register(ModuleName, 10, "invalid state feed")
	ErrInvalidServiceFeeCap = sdkerrors.Register(ModuleName, 11, "service fee cap is invalid")
)
