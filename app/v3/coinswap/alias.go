package coinswap

import (
	"github.com/irisnet/irishub/app/v3/coinswap/internal/keeper"
	"github.com/irisnet/irishub/app/v3/coinswap/internal/types"
)

type (
	Keeper               = keeper.Keeper
	MsgSwapOrder         = types.MsgSwapOrder
	MsgAddLiquidity      = types.MsgAddLiquidity
	MsgRemoveLiquidity   = types.MsgRemoveLiquidity
	Params               = types.Params
	QueryLiquidityParams = types.QueryLiquidityParams
	Input                = types.Input
	Output               = types.Output
)

const (
	DefaultCodespace       = types.DefaultCodespace
	ModuleName             = types.ModuleName
	LiquidityVoucherPrefix = types.LiquidityVoucherPrefix
	QueryLiquidity         = types.QueryLiquidity
	DefaultParamSpace      = types.DefaultParamSpace
)

var (
	RegisterCodec         = types.RegisterCodec
	NewMsgSwapOrder       = types.NewMsgSwapOrder
	NewMsgAddLiquidity    = types.NewMsgAddLiquidity
	NewMsgRemoveLiquidity = types.NewMsgRemoveLiquidity
	NewKeeper             = keeper.NewKeeper
	NewQuerier            = keeper.NewQuerier
	ErrInvalidDeadline    = types.ErrInvalidDeadline
	ErrNotPositive        = types.ErrNotPositive
	ErrConstraintNotMet   = types.ErrConstraintNotMet
	GetVoucherCoinName    = types.GetVoucherCoinName
	GetUnderlyingDenom    = types.GetUnderlyingDenom
	GetVoucherDenom       = types.GetVoucherDenom
	GetVoucherCoinType    = types.GetVoucherCoinType
	CheckVoucherDenom     = types.CheckVoucherDenom
	CheckVoucherCoinName  = types.CheckVoucherCoinName
)
