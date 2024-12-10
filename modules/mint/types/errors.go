package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// mint module sentinel errors
var (
	ErrInvalidMintInflation = sdkerrors.Register(ModuleName, 2, "invalid mint inflation")
	ErrInvalidMintDenom     = sdkerrors.Register(ModuleName, 3, "invalid mint denom")
)
