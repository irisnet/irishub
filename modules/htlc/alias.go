package htlc

import (
	"github.com/irisnet/irishub/modules/htlc/internal/keeper"
	"github.com/irisnet/irishub/modules/htlc/internal/types"
)

// exported constants
const (
	ModuleName       = types.ModuleName
	StoreKey         = types.StoreKey
	RouterKey        = types.ModuleName
	QuerierRoute     = types.QuerierRoute
	DefaultCodespace = types.DefaultCodespace
	OPEN             = types.OPEN
	EXPIRED          = types.EXPIRED

	EventTypeCreateHTLC                = types.EventTypeCreateHTLC
	EventTypeClaimHTLC                 = types.EventTypeClaimHTLC
	EventTypeRefundHTLC                = types.EventTypeRefundHTLC
	EventTypeExpiredHTLC               = types.EventTypeExpiredHTLC
	AttributeValueCategory             = types.AttributeValueCategory
	AttributeValueSender               = types.AttributeValueSender
	AttributeValueReceiver             = types.AttributeValueReceiver
	AttributeValueReceiverOnOtherChain = types.AttributeValueReceiverOnOtherChain
	AttributeValueAmount               = types.AttributeValueAmount
	AttributeValueHashLock             = types.AttributeValueHashLock
	AttributeValueSecret               = types.AttributeValueSecret
)

// exported types
type (
	Keeper        = keeper.Keeper
	HTLC          = types.HTLC
	HTLCState     = types.HTLCState
	HTLCSecret    = types.HTLCSecret
	HTLCHashLock  = types.HTLCHashLock
	GenesisState  = types.GenesisState
	MsgCreateHTLC = types.MsgCreateHTLC
	MsgClaimHTLC  = types.MsgClaimHTLC
	MsgRefundHTLC = types.MsgRefundHTLC
)

// exported variables and functions
var (
	// variable aliases
	ModuleCdc = types.ModuleCdc

	// functions aliases
	NewKeeper           = keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
	RegisterCodec       = types.RegisterCodec
	NewHTLC             = types.NewHTLC
	GetHashLock         = types.GetHashLock
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
)
