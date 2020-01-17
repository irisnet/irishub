package htlc

import (
	"github.com/irisnet/irishub/app/v2/htlc/internal/keeper"
	"github.com/irisnet/irishub/app/v2/htlc/internal/types"
)

// exported types
type (
	MsgCreateHTLC   = types.MsgCreateHTLC
	MsgClaimHTLC    = types.MsgClaimHTLC
	MsgRefundHTLC   = types.MsgRefundHTLC
	HTLC            = types.HTLC
	HTLCState       = types.HTLCState
	GenesisState    = types.GenesisState
	QueryHTLCParams = types.QueryHTLCParams
	Keeper          = keeper.Keeper
)

const (
	DefaultCodespace = types.DefaultCodespace
	SecretLength     = types.SecretLength
	OPEN             = types.OPEN
	EXPIRED          = types.EXPIRED
	REFUNDED         = types.REFUNDED
	QueryHTLC        = types.QueryHTLC
)

// exported variables and functions
var (
	TagHashLock = types.TagHashLock

	NewKeeper        = keeper.NewKeeper
	GetHashLock      = keeper.GetHashLock
	NewQuerier       = keeper.NewQuerier
	RegisterCodec    = types.RegisterCodec
	NewMsgCreateHTLC = types.NewMsgCreateHTLC
	NewMsgClaimHTLC  = types.NewMsgClaimHTLC
	NewMsgRefundHTLC = types.NewMsgRefundHTLC
	NewHTLC          = types.NewHTLC
)
