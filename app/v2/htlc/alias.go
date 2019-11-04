package htlc

import (
	"github.com/irisnet/irishub/app/v2/htlc/internal/keeper"
	"github.com/irisnet/irishub/app/v2/htlc/internal/types"
)

// exported types
type (
	MsgCreateHTLC = types.MsgCreateHTLC
	MsgClaimHTLC  = types.MsgClaimHTLC
	MsgRefundHTLC = types.MsgRefundHTLC

	HTLC = types.HTLC

	Params       = types.Params
	GenesisState = types.GenesisState

	QueryHTLCParams = types.QueryHTLCParams

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

	NewMsgCreateHTLC = types.NewMsgCreateHTLC
	NewMsgClaimHTLC  = types.NewMsgClaimHTLC
	NewMsgRefundHTLC = types.NewMsgRefundHTLC
	NewHTLC          = types.NewHTLC
	GetHashLock      = keeper.GetHashLock

	SecretLength = types.SecretLength

	OPEN     = types.OPEN
	EXPIRED  = types.EXPIRED
	REFUNDED = types.REFUNDED

	QueryHTLC = types.QueryHTLC

	TagHashLock = types.TagHashLock

	NewKeeper  = keeper.NewKeeper
	NewQuerier = keeper.NewQuerier
)
