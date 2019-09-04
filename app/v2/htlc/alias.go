package htlc

import (
	"github.com/irisnet/irishub/app/v2/htlc/internal/keeper"
	"github.com/irisnet/irishub/app/v2/htlc/internal/types"
)

// exported types
type (
	MsgCreateHTLC = types.MsgCreateHTLC
	HTLC          = types.HTLC

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
	NewHTLC          = types.NewHTLC

	ValidateSecretHashLock = types.ValidateSecretHashLock

	QueryHTLC = types.QueryHTLC

	NewKeeper  = keeper.NewKeeper
	NewQuerier = keeper.NewQuerier
)
