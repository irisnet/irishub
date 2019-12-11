package htlc

import (
	"github.com/irisnet/irishub/modules/htlc/internal/keeper"
	"github.com/irisnet/irishub/modules/htlc/internal/types"
)

// exported types
const (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	RouterKey    = types.ModuleName
	QuerierRoute = types.QuerierRoute
)

type (
	Keeper        = keeper.Keeper
	HTLC          = types.HTLC
	HTLCSecret    = types.HTLCSecret
	HTLCHashLock  = types.HTLCHashLock
	GenesisState  = types.GenesisState
	MsgCreateHTLC = types.MsgCreateHTLC
	MsgClaimHTLC  = types.MsgClaimHTLC
	MsgRefundHTLC = types.MsgRefundHTLC
)

// exported variables and functions
var (
	// functions aliases
	NewKeeper       = keeper.NewKeeper
	NewQuerier      = keeper.NewQuerier
	RegisterCodec   = types.RegisterCodec
	NewHTLC         = types.NewHTLC
	ValidateGenesis = types.ValidateGenesis
	GetHashLock     = types.GetHashLock

	// const aliases
	DefaultCodespace = types.DefaultCodespace
	OPEN             = types.OPEN
	EXPIRED          = types.EXPIRED
	REFUNDED         = types.REFUNDED

	// variable aliases
	ModuleCdc = types.ModuleCdc
)
