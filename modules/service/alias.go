package service

import (
	"github.com/irisnet/irishub/modules/service/internal/keeper"
	"github.com/irisnet/irishub/modules/service/internal/types"
)

// exported types
type (
	MsgSvcDef           = types.MsgSvcDef
	MsgSvcBind          = types.MsgSvcBind
	MsgSvcBindingUpdate = types.MsgSvcBindingUpdate
	MsgSvcDisable       = types.MsgSvcDisable
	MsgSvcEnable        = types.MsgSvcEnable
	MsgSvcRefundDeposit = types.MsgSvcRefundDeposit
	MsgSvcRequest       = types.MsgSvcRequest
	MsgSvcResponse      = types.MsgSvcResponse
	MsgSvcRefundFees    = types.MsgSvcRefundFees
	MsgSvcWithdrawFees  = types.MsgSvcWithdrawFees
	MsgSvcWithdrawTax   = types.MsgSvcWithdrawTax

	SvcDef           = types.SvcDef
	SvcBinding       = types.SvcBinding
	SvcRequest       = types.SvcRequest
	SvcResponse      = types.SvcResponse
	MethodProperty   = types.MethodProperty
	Level            = types.Level
	NoPrivacy        = types.NoPrivacy
	PubKeyEncryption = types.PubKeyEncryption
	Unicast          = types.Unicast
	Multicast        = types.Multicast
	Global           = types.Global
	Local            = types.Local

	GenesisState = types.GenesisState

	Keeper = keeper.Keeper
)

// exported variables and functions
var (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	TStoreKey    = types.TStoreKey
	QuerierRoute = types.QuerierRoute
	RouterKey    = types.RouterKey
	ModuleCdc    = types.ModuleCdc

	DefaultCodespace = types.DefaultCodespace
	RegisterCodec    = types.RegisterCodec

	NewMsgSvcDef           = types.NewMsgSvcDef
	NewMsgSvcBind          = types.NewMsgSvcBind
	NewMsgSvcBindingUpdate = types.NewMsgSvcBindingUpdate
	NewMsgSvcDisable       = types.NewMsgSvcDisable
	NewMsgSvcEnable        = types.NewMsgSvcEnable
	NewMsgSvcRefundDeposit = types.NewMsgSvcRefundDeposit
	NewMsgSvcRequest       = types.NewMsgSvcRequest
	NewMsgSvcResponse      = types.NewMsgSvcResponse
	NewMsgSvcRefundFees    = types.NewMsgSvcRefundFees
	NewMsgSvcWithdrawFees  = types.NewMsgSvcWithdrawFees
	NewMsgSvcWithdrawTax   = types.NewMsgSvcWithdrawTax

	QueryDefinition = types.QueryDefinition
	QueryBinding    = types.QueryBinding
	QueryBindings   = types.QueryBindings
	QueryRequests   = types.QueryRequests
	QueryResponse   = types.QueryResponse
	QueryFees       = types.QueryFees

	NewGenesisState      = types.NewGenesisState
	DefaultParams        = types.DefaultParam
	DefaultParamsForTest = types.DefaultParamsForTest

	NewKeeper  = keeper.NewKeeper
	NewQuerier = keeper.NewQuerier
)
