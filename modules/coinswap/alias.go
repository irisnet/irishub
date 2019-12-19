package coinswap

import (
	"github.com/irisnet/irishub/modules/coinswap/internal/keeper"
	"github.com/irisnet/irishub/modules/coinswap/internal/types"
)

// exported constants
const (
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
	RouterKey         = types.RouterKey
	QuerierRoute      = types.QuerierRoute
	DefaultParamspace = types.DefaultParamspace
	DefaultCodespace  = types.DefaultCodespace

	EventTypeSwap            = types.EventTypeSwap
	EventTypeAddLiquidity    = types.EventTypeAddLiquidity
	EventTypeRemoveLiquidity = types.EventTypeRemoveLiquidity
	AttributeValueCategory   = types.AttributeValueCategory
	AttributeValueAmount     = types.AttributeValueAmount
	AttributeValueSender     = types.AttributeValueSender
	AttributeValueRecipient  = types.AttributeValueRecipient
	AttributeValueIsBuyOrder = types.AttributeValueIsBuyOrder
	AttributeValueTokenPair  = types.AttributeValueTokenPair
)

// exported types
type (
	Keeper               = keeper.Keeper
	MsgSwapOrder         = types.MsgSwapOrder
	MsgAddLiquidity      = types.MsgAddLiquidity
	MsgRemoveLiquidity   = types.MsgRemoveLiquidity
	Params               = types.Params
	QueryLiquidityParams = types.QueryLiquidityParams
	Input                = types.Input
	Output               = types.Output
	GenesisState         = types.GenesisState
)

// exported variables and functions
var (
	// functions aliases
	ModuleCdc = types.ModuleCdc

	// variable aliases
	NewKeeper           = keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
	RegisterCodec       = types.RegisterCodec
	ErrInvalidDeadline  = types.ErrInvalidDeadline
	DefaultParams       = types.DefaultParams
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
)
