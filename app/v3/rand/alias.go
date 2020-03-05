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
	Params                      = types.Params
	GenesisState                = types.GenesisState
	QueryRandParams             = types.QueryRandParams
	QueryRandRequestQueueParams = types.QueryRandRequestQueueParams
	Keeper                      = keeper.Keeper
)

// exported consts
const (
	DefaultCodespace      = types.DefaultCodespace
	DefaultParamSpace     = types.DefaultParamSpace
	DefaultBlockInterval  = types.DefaultBlockInterval
	RandPrec              = types.RandPrec
	QueryRand             = types.QueryRand
	QueryRandRequestQueue = types.QueryRandRequestQueue
)

// exported variables and functions
var (
	TagReqID      = types.TagReqID
	TagRandHeight = types.TagRandHeight
	TagRand       = types.TagRand

	DefaultParams        = types.DefaultParams
	DefaultParamsForTest = types.DefaultParamsForTest
	ValidateParams       = types.ValidateParams
	RegisterCodec        = types.RegisterCodec
	NewMsgRequestRand    = types.NewMsgRequestRand
	NewRand              = types.NewRand
	NewRequest           = types.NewRequest
	MakePRNG             = types.MakePRNG
	GenerateRequestID    = types.GenerateRequestID
	CheckReqID           = types.CheckReqID
	NewKeeper            = keeper.NewKeeper
	NewQuerier           = keeper.NewQuerier
)
