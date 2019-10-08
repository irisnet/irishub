package rand

import (
	"github.com/irisnet/irishub/modules/rand/internal/keeper"
	"github.com/irisnet/irishub/modules/rand/internal/types"
)

// exported types
type (
	MsgRequestRand = types.MsgRequestRand
	Rand           = types.Rand
	Request        = types.Request
	Requests       = types.Requests

	Params       = types.Params
	GenesisState = types.GenesisState

	QueryRandParams             = types.QueryRandParams
	QueryRandRequestQueueParams = types.QueryRandRequestQueueParams

	Keeper = keeper.Keeper
)

// exported variables and functions
var (
	DefaultCodespace     = types.DefaultCodespace
	DefaultParamSpace    = types.DefaultParamSpace
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
	DefaultBlockInterval = types.DefaultBlockInterval
	RandPrec             = types.RandPrec

	QueryRand             = types.QueryRand
	QueryRandRequestQueue = types.QueryRandRequestQueue

	TagReqID      = types.TagReqID
	TagRandHeight = types.TagRandHeight
	TagRand       = types.TagRand

	NewKeeper  = keeper.NewKeeper
	NewQuerier = keeper.NewQuerier
)
