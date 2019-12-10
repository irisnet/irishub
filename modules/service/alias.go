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

	SvcDef         = types.SvcDef
	SvcBinding     = types.SvcBinding
	SvcRequest     = types.SvcRequest
	SvcResponse    = types.SvcResponse
	MethodProperty = types.MethodProperty
	Level          = types.Level

	GenesisState = types.GenesisState

	Keeper = keeper.Keeper
)

// exported constants
const (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute
	RouterKey    = types.RouterKey

	DepositAccName = types.RequestAccName
	RequestAccName = types.RequestAccName
	TaxAccName     = types.TaxAccName

	DefaultCodespace = types.DefaultCodespace

	NoPrivacy        = types.NoPrivacy
	PubKeyEncryption = types.PubKeyEncryption
	Unicast          = types.Unicast
	Multicast        = types.Multicast
	Global           = types.Global
	Local            = types.Local
)

// exported variables and functions
var (
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec

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

	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis

	DefaultParams = types.DefaultParams

	ErrSvcDefExists           = types.ErrSvcDefExists
	ErrSvcBindingNotExists    = types.ErrSvcBindingNotExists
	ErrSvcBindingNotAvailable = types.ErrSvcBindingNotAvailable
	ErrMethodNotExists        = types.ErrMethodNotExists
	ErrNotProfiler            = types.ErrNotProfiler
	ErrLtServiceFee           = types.ErrLtServiceFee
	ErrRequestNotActive       = types.ErrRequestNotActive
	ErrNotMatchingProvider    = types.ErrNotMatchingProvider
	ErrNotMatchingReqChainID  = types.ErrNotMatchingReqChainID
	ErrNotTrustee             = types.ErrNotTrustee

	NewKeeper  = keeper.NewKeeper
	NewQuerier = keeper.NewQuerier
)
