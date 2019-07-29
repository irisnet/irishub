package coinswap

import (
	"github.com/irisnet/irishub/app/v1/coinswap/internal/keeper"
	"github.com/irisnet/irishub/app/v1/coinswap/internal/types"
)

type (
	Keeper             = keeper.Keeper
	MsgSwapOrder       = types.MsgSwapOrder
	MsgAddLiquidity    = types.MsgAddLiquidity
	MsgRemoveLiquidity = types.MsgRemoveLiquidity
	Params             = types.Params
)

var (
	NewMsgSwapOrder       = types.NewMsgSwapOrder
	NewMsgAddLiquidity    = types.NewMsgAddLiquidity
	NewMsgRemoveLiquidity = types.NewMsgRemoveLiquidity

	ErrInvalidDeadline  = types.ErrInvalidDeadline
	ErrNotPositive      = types.ErrNotPositive
	ErrConstraintNotMet = types.ErrConstraintNotMet
	RegisterCodec       = types.RegisterCodec
	NewKeeper           = keeper.NewKeeper
	DefaultParamSpace   = types.DefaultParamSpace
)

const (
	DefaultCodespace = types.DefaultCodespace
	ModuleName       = types.ModuleName
)
