package rand

import (
	"github.com/irisnet/irishub/app/v3/rand/internal/keeper"
	"github.com/irisnet/irishub/app/v3/rand/internal/types"
)

// exported types
type (
	MsgRequestRand              = types.MsgRequestRand
	Rand                        = types.Rand
	Request                     = types.Request
	Requests                    = types.Requests
	GenesisState                = types.GenesisState
	QueryRandParams             = types.QueryRandParams
	QueryRandRequestQueueParams = types.QueryRandRequestQueueParams
	Keeper                      = keeper.Keeper
)

// exported consts
const (
	ModuleName            = types.ModuleName
	DefaultCodespace      = types.DefaultCodespace
	DefaultBlockInterval  = types.DefaultBlockInterval
	RandPrec              = types.RandPrec
	QueryRand             = types.QueryRand
	QueryRandRequestQueue = types.QueryRandRequestQueue
	ModuleServiceName     = types.ServiceName
)

// exported variables and functions
var (
	TagReqID            = types.TagReqID
	TagRequestContextID = types.TagRequestContextID
	TagRandHeight       = types.TagRandHeight
	TagRand             = types.TagRand

	NewKeeper         = keeper.NewKeeper
	NewQuerier        = keeper.NewQuerier
	RegisterCodec     = types.RegisterCodec
	NewMsgRequestRand = types.NewMsgRequestRand
	NewRand           = types.NewRand
	NewRequest        = types.NewRequest
	MakePRNG          = types.MakePRNG
	GenerateRequestID = types.GenerateRequestID
	CheckReqID        = types.CheckReqID
	GetSvcDefinitions = types.GetSvcDefinitions
	ValidateGenesis   = types.ValidateGenesis
)
