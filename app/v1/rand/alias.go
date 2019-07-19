package rand

import (
	"github.com/irisnet/irishub/app/v1/rand/internal/keeper"
	"github.com/irisnet/irishub/app/v1/rand/internal/types"
)

// exported types
type (
	MsgRequestRand = types.MsgRequestRand
	Rand           = types.Rand
	Request        = types.Request

	Params       = types.Params
	GenesisState = types.GenesisState

	QueryRandParams             = types.QueryRandParams
	QueryRandsParams            = types.QueryRandsParams
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
	CheckReqID           = types.CheckReqID
	DefaultBlockInterval = types.DefaultBlockInterval

	QueryRand             = types.QueryRand
	QueryRands            = types.QueryRands
	QueryRandRequestQueue = types.QueryRandRequestQueue

	NewKeeper  = keeper.NewKeeper
	NewQuerier = keeper.NewQuerier
)
