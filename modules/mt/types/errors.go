package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrDenomNotFound = sdkerrors.Register(ModuleName, 1, "denom not found")
	ErrMTNotFound    = sdkerrors.Register(ModuleName, 2, "mt not found")
)
