package types

import (
	errorsmod "cosmossdk.io/errors"
)

// HTLC module sentinel errors
var (
	ErrInvalidID                   = errorsmod.Register(ModuleName, 2, "invalid htlc id")
	ErrInvalidHashLock             = errorsmod.Register(ModuleName, 3, "invalid hash lock")
	ErrInvalidTimeLock             = errorsmod.Register(ModuleName, 4, "invalid time lock")
	ErrInvalidSecret               = errorsmod.Register(ModuleName, 5, "invalid secret")
	ErrInvalidExpirationHeight     = errorsmod.Register(ModuleName, 6, "invalid expiration height")
	ErrInvalidTimestamp            = errorsmod.Register(ModuleName, 7, "invalid timestamp")
	ErrInvalidState                = errorsmod.Register(ModuleName, 8, "invalid state")
	ErrInvalidClosedBlock          = errorsmod.Register(ModuleName, 9, "invalid closed block")
	ErrInvalidDirection            = errorsmod.Register(ModuleName, 10, "invalid direction")
	ErrHTLCExists                  = errorsmod.Register(ModuleName, 11, "htlc already exists")
	ErrUnknownHTLC                 = errorsmod.Register(ModuleName, 12, "unknown htlc")
	ErrHTLCNotOpen                 = errorsmod.Register(ModuleName, 13, "htlc not open")
	ErrAssetNotSupported           = errorsmod.Register(ModuleName, 14, "asset not found")
	ErrAssetNotActive              = errorsmod.Register(ModuleName, 15, "asset is currently inactive")
	ErrInvalidAccount              = errorsmod.Register(ModuleName, 16, "invalid account")
	ErrInvalidAmount               = errorsmod.Register(ModuleName, 17, "invalid amount")
	ErrInsufficientAmount          = errorsmod.Register(ModuleName, 18, "amount cannot cover the deputy fixed fee")
	ErrExceedsSupplyLimit          = errorsmod.Register(ModuleName, 19, "asset supply over limit")
	ErrExceedsTimeBasedSupplyLimit = errorsmod.Register(ModuleName, 20, "asset supply over limit for current time period")
	ErrInvalidCurrentSupply        = errorsmod.Register(ModuleName, 21, "supply decrease puts current asset supply below 0")
	ErrInvalidIncomingSupply       = errorsmod.Register(ModuleName, 22, "supply decrease puts incoming asset supply below 0")
	ErrInvalidOutgoingSupply       = errorsmod.Register(ModuleName, 23, "supply decrease puts outgoing asset supply below 0")
	ErrExceedsAvailableSupply      = errorsmod.Register(ModuleName, 24, "outgoing swap exceeds total available supply")
	ErrAssetSupplyNotFound         = errorsmod.Register(ModuleName, 25, "asset supply not found in store")
)
