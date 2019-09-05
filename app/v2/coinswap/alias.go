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
	Input                = types.Input
	Output               = types.Output
)

var (
	DefaultParamSpace = types.DefaultParamSpace
	QueryLiquidity    = types.QueryLiquidity

	RegisterCodec = types.RegisterCodec

	NewMsgSwapOrder       = types.NewMsgSwapOrder
	NewMsgAddLiquidity    = types.NewMsgAddLiquidity
	NewMsgRemoveLiquidity = types.NewMsgRemoveLiquidity
	NewKeeper             = keeper.NewKeeper
	NewQuerier            = keeper.NewQuerier

	ErrInvalidDeadline  = types.ErrInvalidDeadline
	ErrNotPositive      = types.ErrNotPositive
	ErrConstraintNotMet = types.ErrConstraintNotMet

	GetUniId                    = types.GetUniId
	GetCoinMinDenomFromUniDenom = types.GetCoinMinDenomFromUniDenom
	GetUniDenom                 = types.GetUniDenom
	GetUniCoinType              = types.GetUniCoinType
	CheckUniDenom               = types.CheckUniDenom
	CheckUniId                  = types.CheckUniId
)

const (
	DefaultCodespace   = types.DefaultCodespace
	ModuleName         = types.ModuleName
	FormatUniABSPrefix = types.FormatUniABSPrefix
)
