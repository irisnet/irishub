package coinswap

import (
	"github.com/irisnet/irishub/app/v2/coinswap/internal/keeper"
	"github.com/irisnet/irishub/app/v2/coinswap/internal/types"
)

type (
	Keeper               = keeper.Keeper
	MsgSwapOrder         = types.MsgSwapOrder
	MsgAddLiquidity      = types.MsgAddLiquidity
	MsgRemoveLiquidity   = types.MsgRemoveLiquidity
	Params               = types.Params
	QueryLiquidityParams = types.QueryLiquidityParams
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
	NewQuerier          = keeper.NewQuerier
	DefaultParamSpace   = types.DefaultParamSpace

	QueryLiquidity = types.QueryLiquidity
)

const (
	DefaultCodespace = types.DefaultCodespace
	ModuleName       = types.ModuleName
)
