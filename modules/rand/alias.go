package rand

import (
	"github.com/irisnet/irishub/modules/rand/internal/keeper"
	"github.com/irisnet/irishub/modules/rand/internal/types"
)

// exported constants
const (
	ModuleName            = types.ModuleName
	StoreKey              = types.StoreKey
	QuerierRoute          = types.QuerierRoute
	RouterKey             = types.RouterKey
	DefaultCodespace      = types.DefaultCodespace
	RandPrec              = types.RandPrec

	EventTypeGenerateRand  = types.EventTypeGenerateRand
	EventTypeRequestRand   = types.EventTypeRequestRand
	AttributeKeyRequestID  = types.AttributeKeyRequestID
	AttributeKeyRand       = types.AttributeKeyRand
	AttributeKeyGenHeight  = types.AttributeKeyGenHeight
	AttributeValueCategory = types.AttributeValueCategory
)

// exported types
type (
	Keeper                      = keeper.Keeper
	MsgRequestRand              = types.MsgRequestRand
	Rand                        = types.Rand
	Request                     = types.Request
	Requests                    = types.Requests
	GenesisState                = types.GenesisState
	QueryRandParams             = types.QueryRandParams
	QueryRandRequestQueueParams = types.QueryRandRequestQueueParams
)

// exported variables and functions
var (
	// variable aliases
	ModuleCdc = types.ModuleCdc

	// functions aliases
	RegisterCodec       = types.RegisterCodec
	NewRand             = types.NewRand
	NewRequest          = types.NewRequest
	MakePRNG            = types.MakePRNG
	GenerateRequestID   = types.GenerateRequestID
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
	NewKeeper           = keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
)
