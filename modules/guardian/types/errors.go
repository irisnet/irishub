package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// guardian module sentinel errors
var (
	ErrUnknownOperator       = sdkerrors.Register(ModuleName, 2, "unknown operator")
	ErrUnknownProfiler       = sdkerrors.Register(ModuleName, 3, "unknown profiler")
	ErrUnknownTrustee        = sdkerrors.Register(ModuleName, 4, "unknown trustee")
	ErrDeleteGenesisProfiler = sdkerrors.Register(ModuleName, 5, "can't delete genesis profiler")
	ErrDeleteGenesisTrustee  = sdkerrors.Register(ModuleName, 6, "can't delete genesis trustee")
)
