package htlc

import (
	"github.com/irisnet/irishub/modules/htlc/internal/keeper"
	"github.com/irisnet/irishub/modules/htlc/internal/types"
)

// exported types
const (
	ModuleName        = types.ModuleName
	DefaultParamspace = types.DefaultParamspace
	StoreKey          = types.StoreKey
	QuerierRoute      = types.QuerierRoute
)

type (
	MsgCreateHTLC = types.MsgCreateHTLC
	MsgClaimHTLC = types.MsgClaimHTLC
	MsgRefundHTLC = types.MsgRefundHTLC
	HTLC = types.HTLC
	Params = types.Params
	GenesisState = types.GenesisState
	QueryHTLCParams = types.QueryHTLCParams
	Keeper = keeper.Keeper
)

// exported variables and functions
var (
	// functions aliases
	NewKeeper            = keeper.NewKeeper
	NewQuerier           = keeper.NewQuerier
	GetHashLock          = keeper.GetHashLock
	DefaultParams        = types.DefaultParams
	DefaultParamsForTest = types.DefaultParamsForTest
	ValidateParams       = types.ValidateParams
	RegisterCodec        = types.RegisterCodec
	NewMsgCreateHTLC     = types.NewMsgCreateHTLC
	NewMsgClaimHTLC      = types.NewMsgClaimHTLC
	NewMsgRefundHTLC     = types.NewMsgRefundHTLC
	NewHTLC              = types.NewHTLC
	ValidateGenesis      = types.ValidateGenesis

	// const aliases
	DefaultCodespace  = types.DefaultCodespace
	DefaultParamSpace = types.DefaultParamSpace
	SecretLength      = types.SecretLength
	OPEN              = types.OPEN
	EXPIRED           = types.EXPIRED
	REFUNDED          = types.REFUNDED
	QueryHTLC         = types.QueryHTLC

	// variable aliases
	HTLCLockedCoinsAccAddr = keeper.HTLCLockedCoinsAccAddr
	ModuleCdc              = types.ModuleCdc
)
