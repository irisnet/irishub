package types

import (
	errorsmod "cosmossdk.io/errors"
)

// farm module sentinel errors
var (
	ErrExpiredHeight      = errorsmod.Register(ModuleName, 2, "expired block height")
	ErrInvalidLPToken     = errorsmod.Register(ModuleName, 3, "invalid lp token denom")
	ErrNotMatch           = errorsmod.Register(ModuleName, 4, "the data does not match")
	ErrPoolExpired        = errorsmod.Register(ModuleName, 5, "the farm pool has expired")
	ErrPoolNotStart       = errorsmod.Register(ModuleName, 6, "the farm pool don't start")
	ErrPoolExist          = errorsmod.Register(ModuleName, 7, "the farm pool already exist")
	ErrPoolNotFound       = errorsmod.Register(ModuleName, 8, "the farm pool does not exist")
	ErrInvalidOperate     = errorsmod.Register(ModuleName, 9, "invalid operate")
	ErrFarmerNotFound     = errorsmod.Register(ModuleName, 10, "the farmer does not exist")
	ErrInvalidPoolId      = errorsmod.Register(ModuleName, 11, "invalid pool id")
	ErrInvalidDescription = errorsmod.Register(ModuleName, 12, "invalid pool description, length must be less than or equal to 280")
	ErrInvalidAppend      = errorsmod.Register(ModuleName, 13, "cannot add new token as a reward")
	ErrInvalidRewardRule  = errorsmod.Register(ModuleName, 14, "invalid reward rule")
	ErrAllEmpty           = errorsmod.Register(ModuleName, 15, "shouldn't all be empty")
	ErrBadDistribution    = errorsmod.Register(ModuleName, 16, "community pool does not have sufficient coins to distribute")
	ErrInvalidRefund      = errorsmod.Register(ModuleName, 17, "invalid refund amount")
	ErrInvalidProposal    = errorsmod.Register(ModuleName, 18, "invalid proposal")
	ErrEscrowInfoNotFound = errorsmod.Register(ModuleName, 19, "the escrow information does not exist")
)
