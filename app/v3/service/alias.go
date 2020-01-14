package service

import (
	"github.com/irisnet/irishub/app/v3/service/internal/keeper"
	"github.com/irisnet/irishub/app/v3/service/internal/types"
)

// nolint
const (
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
	RegisterCodec              = types.RegisterCodec
	NewGenesisState            = types.NewGenesisState
	DefaultGenesisState        = types.DefaultGenesisState
	DefaultGenesisStateForTest = types.DefaultGenesisStateForTest
	ValidateGenesis            = types.ValidateGenesis
	DefaultParams              = types.DefaultParams
	PrometheusMetrics          = types.PrometheusMetrics
	NewServiceDefinition       = types.NewServiceDefinition
	NewServiceBinding          = types.NewServiceBinding
	NewSvcRequest              = types.NewSvcRequest
	NewSvcResponse             = types.NewSvcResponse
	NewReturnedFee             = types.NewReturnedFee
	NewIncomingFee             = types.NewIncomingFee
	NewMsgDefineService        = types.NewMsgDefineService
	NewMsgBindService          = types.NewMsgBindService
	NewMsgUpdateServiceBinding = types.NewMsgUpdateServiceBinding
	NewMsgSetWithdrawAddress   = types.NewMsgSetWithdrawAddress
	NewMsgDisableService       = types.NewMsgDisableService
	NewMsgEnableService        = types.NewMsgEnableService
	NewMsgRefundServiceDeposit = types.NewMsgRefundServiceDeposit
	NewMsgSvcRequest           = types.NewMsgSvcRequest
	NewMsgSvcResponse          = types.NewMsgSvcResponse
	NewMsgSvcRefundFees        = types.NewMsgSvcRefundFees
	NewMsgSvcWithdrawFees      = types.NewMsgSvcWithdrawFees
	NewMsgSvcWithdrawTax       = types.NewMsgSvcWithdrawTax
	ConvertRequestID           = types.ConvertRequestID
)

// nolint
type (
	Keeper                  = keeper.Keeper
	GenesisState            = types.GenesisState
	Params                  = types.Params
	ServiceDefinition       = types.ServiceDefinition
	ServiceBinding          = types.ServiceBinding
	ServiceBindings         = types.ServiceBindings
	SvcRequest              = types.SvcRequest
	SvcResponse             = types.SvcResponse
	ReturnedFee             = types.ReturnedFee
	IncomingFee             = types.IncomingFee
	Metrics                 = types.Metrics
	MsgDefineService        = types.MsgDefineService
	MsgBindService          = types.MsgBindService
	MsgUpdateServiceBinding = types.MsgUpdateServiceBinding
	MsgSetWithdrawAddress   = types.MsgSetWithdrawAddress
	MsgDisableService       = types.MsgDisableService
	MsgEnableService        = types.MsgEnableService
	MsgRefundServiceDeposit = types.MsgRefundServiceDeposit
	MsgSvcRequest           = types.MsgSvcRequest
	MsgSvcResponse          = types.MsgSvcResponse
	MsgSvcRefundFees        = types.MsgSvcRefundFees
	MsgSvcWithdrawFees      = types.MsgSvcWithdrawFees
	MsgSvcWithdrawTax       = types.MsgSvcWithdrawTax
	QueryDefinitionParams   = types.QueryDefinitionParams
	QueryBindingParams      = types.QueryBindingParams
	QueryBindingsParams     = types.QueryBindingsParams
	QueryRequestsParams     = types.QueryRequestsParams
	QueryResponseParams     = types.QueryResponseParams
	QueryFeesParams         = types.QueryFeesParams
	FeesOutput              = types.FeesOutput
)
