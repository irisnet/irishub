package types

import (
	errorsmod "cosmossdk.io/errors"
)

// record module sentinel errors
var (
	ErrUnknownRecord = errorsmod.Register(ModuleName, 2, "unknown record")
)
