package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// HTLC module sentinel errors
var (
	ErrInvalidID                   = sdkerrors.Register(ModuleName, 2, "invalid htlc id")
	ErrInvalidHashLock             = sdkerrors.Register(ModuleName, 3, "invalid hash lock")
	ErrInvalidTimeLock             = sdkerrors.Register(ModuleName, 4, "invalid time lock")
	ErrInvalidSecret               = sdkerrors.Register(ModuleName, 5, "invalid secret")
	ErrInvalidExpirationHeight     = sdkerrors.Register(ModuleName, 6, "invalid expiration height")
	ErrInvalidTimestamp            = sdkerrors.Register(ModuleName, 7, "invalid timestamp")
	ErrInvalidState                = sdkerrors.Register(ModuleName, 8, "invalid state")
	ErrInvalidClosedBlock          = sdkerrors.Register(ModuleName, 9, "invalid closed block")
	ErrInvalidDirection            = sdkerrors.Register(ModuleName, 10, "invalid direction")
	ErrHTLCExists                  = sdkerrors.Register(ModuleName, 11, "htlc already exists")
	ErrUnknownHTLC                 = sdkerrors.Register(ModuleName, 12, "unknown htlc")
	ErrHTLCNotOpen                 = sdkerrors.Register(ModuleName, 13, "htlc not open")
	ErrAssetNotSupported           = sdkerrors.Register(ModuleName, 14, "asset not found")
	ErrAssetNotActive              = sdkerrors.Register(ModuleName, 15, "asset is currently inactive")
	ErrInvalidAccount              = sdkerrors.Register(ModuleName, 16, "invalid account")
	ErrInvalidAmount               = sdkerrors.Register(ModuleName, 17, "invalid amount")
	ErrInsufficientAmount          = sdkerrors.Register(ModuleName, 18, "amount cannot cover the deputy fixed fee")
	ErrExceedsSupplyLimit          = sdkerrors.Register(ModuleName, 19, "asset supply over limit")
	ErrExceedsTimeBasedSupplyLimit = sdkerrors.Register(ModuleName, 20, "asset supply over limit for current time period")
	ErrInvalidCurrentSupply        = sdkerrors.Register(ModuleName, 21, "supply decrease puts current asset supply below 0")
	ErrInvalidIncomingSupply       = sdkerrors.Register(ModuleName, 22, "supply decrease puts incoming asset supply below 0")
	ErrInvalidOutgoingSupply       = sdkerrors.Register(ModuleName, 23, "supply decrease puts outgoing asset supply below 0")
	ErrExceedsAvailableSupply      = sdkerrors.Register(ModuleName, 24, "outgoing swap exceeds total available supply")
	ErrAssetSupplyNotFound         = sdkerrors.Register(ModuleName, 25, "asset supply not found in store")
)
