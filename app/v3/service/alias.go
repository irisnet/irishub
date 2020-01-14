package service

import (
	"github.com/irisnet/irishub/app/v3/service/internal/keeper"
	"github.com/irisnet/irishub/app/v3/service/internal/types"
)

// nolint
const (
	Global            = types.Global
	Local             = types.Local
	DefaultParamSpace = types.DefaultParamSpace
	DefaultCodespace  = types.DefaultCodespace
	MsgRoute          = types.MsgRoute
	MetricsSubsystem  = types.MetricsSubsystem
	QueryDefinition   = types.QueryDefinition
	QueryBinding      = types.QueryBinding
	QueryBindings     = types.QueryBindings
	QueryRequests     = types.QueryRequests
	QueryResponse     = types.QueryResponse
	QueryFees         = types.QueryFees
)

// nolint
var (
	// variables and functions aliases
	KeyTxSizeLimit             = types.KeyTxSizeLimit
	NewKeeper                  = keeper.NewKeeper
	NewQuerier                 = keeper.NewQuerier
	NewSvcBinding              = types.NewSvcBinding
	SvcBindingEqual            = types.SvcBindingEqual
	BindingTypeFromString      = types.BindingTypeFromString
	RegisterCodec              = types.RegisterCodec
	NewServiceDefinition       = types.NewServiceDefinition
	NewGenesisState            = types.NewGenesisState
	DefaultGenesisState        = types.DefaultGenesisState
	DefaultGenesisStateForTest = types.DefaultGenesisStateForTest
	ValidateGenesis            = types.ValidateGenesis
	NewSvcRequest              = types.NewSvcRequest
	ConvertRequestID           = types.ConvertRequestID
	NewSvcResponse             = types.NewSvcResponse
	NewReturnedFee             = types.NewReturnedFee
	NewIncomingFee             = types.NewIncomingFee
	NewMsgDefineService        = types.NewMsgDefineService
	NewMsgSvcBind              = types.NewMsgSvcBind
	NewMsgSvcBindingUpdate     = types.NewMsgSvcBindingUpdate
	NewMsgSvcDisable           = types.NewMsgSvcDisable
	NewMsgSvcEnable            = types.NewMsgSvcEnable
	NewMsgSvcRefundDeposit     = types.NewMsgSvcRefundDeposit
	NewMsgSvcRequest           = types.NewMsgSvcRequest
	NewMsgSvcResponse          = types.NewMsgSvcResponse
	NewMsgSvcRefundFees        = types.NewMsgSvcRefundFees
	NewMsgSvcWithdrawFees      = types.NewMsgSvcWithdrawFees
	NewMsgSvcWithdrawTax       = types.NewMsgSvcWithdrawTax
	DefaultParams              = types.DefaultParams
	PrometheusMetrics          = types.PrometheusMetrics
)

// nolint
type (
	Keeper                = keeper.Keeper
	SvcBinding            = types.SvcBinding
	Level                 = types.Level
	BindingType           = types.BindingType
	ServiceDefinition     = types.ServiceDefinition
	GenesisState          = types.GenesisState
	SvcRequest            = types.SvcRequest
	SvcResponse           = types.SvcResponse
	ReturnedFee           = types.ReturnedFee
	IncomingFee           = types.IncomingFee
	Metrics               = types.Metrics
	MsgDefineService      = types.MsgDefineService
	MsgSvcBind            = types.MsgSvcBind
	MsgSvcBindingUpdate   = types.MsgSvcBindingUpdate
	MsgSvcDisable         = types.MsgSvcDisable
	MsgSvcEnable          = types.MsgSvcEnable
	MsgSvcRefundDeposit   = types.MsgSvcRefundDeposit
	MsgSvcRequest         = types.MsgSvcRequest
	MsgSvcResponse        = types.MsgSvcResponse
	MsgSvcRefundFees      = types.MsgSvcRefundFees
	MsgSvcWithdrawFees    = types.MsgSvcWithdrawFees
	MsgSvcWithdrawTax     = types.MsgSvcWithdrawTax
	Params                = types.Params
	QueryDefinitionParams = types.QueryDefinitionParams
	QueryBindingParams    = types.QueryBindingParams
	QueryBindingsParams   = types.QueryBindingsParams
	QueryRequestsParams   = types.QueryRequestsParams
	QueryResponseParams   = types.QueryResponseParams
	QueryFeesParams       = types.QueryFeesParams
	FeesOutput            = types.FeesOutput
)
