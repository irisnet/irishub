package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// rand module sentinel errors
var (
	ErrInvalidReqID  = sdkerrors.Register(ModuleName, 2, "invalid request id")
	ErrInvalidHeight = sdkerrors.Register(ModuleName, 3, "invalid height, must be greater than 0")
)
