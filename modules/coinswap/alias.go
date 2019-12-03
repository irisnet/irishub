package coinswap

import (
	"github.com/irisnet/irishub/modules/coinswap/internal/keeper"
	"github.com/irisnet/irishub/modules/coinswap/internal/types"
)

const (
	ModuleName         = types.ModuleName
	StoreKey           = types.StoreKey
	RouterKey          = types.RouterKey
	QuerierRoute       = types.QuerierRoute
	DefaultParamSpace  = types.DefaultParamSpace
	DefaultCodespace   = types.DefaultCodespace
	QueryLiquidity     = types.QueryLiquidity
	FormatUniABSPrefix = types.FormatUniABSPrefix
)

type (
	Keeper = keeper.Keeper
	MsgSwapOrder = types.MsgSwapOrder
	MsgAddLiquidity = types.MsgAddLiquidity
	MsgRemoveLiquidity = types.MsgRemoveLiquidity
	Params = types.Params
	QueryLiquidityParams = types.QueryLiquidityParams
	Input = types.Input
	Output = types.Output
)

var (
	NewKeeper                   = keeper.NewKeeper
	NewQuerier                  = keeper.NewQuerier
	RegisterCodec               = types.RegisterCodec
	NewMsgSwapOrder             = types.NewMsgSwapOrder
	NewMsgAddLiquidity          = types.NewMsgAddLiquidity
	NewMsgRemoveLiquidity       = types.NewMsgRemoveLiquidity
	GetUniId                    = types.GetUniId
	GetCoinMinDenomFromUniDenom = types.GetCoinMinDenomFromUniDenom
	GetUniDenom                 = types.GetUniDenom
	GetUniCoinType              = types.GetUniCoinType
	CheckUniDenom               = types.CheckUniDenom
	CheckUniId                  = types.CheckUniId
	ErrInvalidDeadline          = types.ErrInvalidDeadline
	ErrNotPositive              = types.ErrNotPositive
	ErrConstraintNotMet         = types.ErrConstraintNotMet
)

// exported variables and functions
var (
	ModuleCdc = types.ModuleCdc
)
