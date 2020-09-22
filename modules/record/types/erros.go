package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// record module sentinel errors
var (
	ErrUnknownRecord = sdkerrors.Register(ModuleName, 2, "unknown record")
)
