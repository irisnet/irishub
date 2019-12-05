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

	GenesisState = types.GenesisState

	QueryRandParams             = types.QueryRandParams
	QueryRandRequestQueueParams = types.QueryRandRequestQueueParams

	Keeper = keeper.Keeper
)

// exported variables and functions
var (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	TStoreKey    = types.TStoreKey
	QuerierRoute = types.QuerierRoute
	RouterKey    = types.RouterKey
	ModuleCdc    = types.ModuleCdc

	DefaultCodespace = types.DefaultCodespace
	RegisterCodec    = types.RegisterCodec

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

	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis

	NewKeeper  = keeper.NewKeeper
	NewQuerier = keeper.NewQuerier
)
