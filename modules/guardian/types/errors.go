package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// guardian module sentinel errors
var (
	ErrUnknownOperator    = sdkerrors.Register(ModuleName, 2, "unknown operator")
	ErrUnknownSuper       = sdkerrors.Register(ModuleName, 3, "unknown super")
	ErrSuperExists        = sdkerrors.Register(ModuleName, 4, "super already exists")
	ErrDeleteGenesisSuper = sdkerrors.Register(ModuleName, 5, "can't delete genesis super")
)
