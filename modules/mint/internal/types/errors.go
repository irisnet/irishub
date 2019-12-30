package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// mint module sentinel errors
var (
	ErrInvalidMintInflation = sdkerrors.Register(ModuleName, 1, "invalid mint inflation")
	ErrInvalidMintDenom     = sdkerrors.Register(ModuleName, 2, "invalid mint denom")
)
