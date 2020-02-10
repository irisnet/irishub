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
	NewRequest                 = types.NewRequest
	NewResponse                = types.NewResponse
	NewEarnedFees              = types.NewEarnedFees
	NewMsgDefineService        = types.NewMsgDefineService
	NewMsgBindService          = types.NewMsgBindService
	NewMsgUpdateServiceBinding = types.NewMsgUpdateServiceBinding
	NewMsgSetWithdrawAddress   = types.NewMsgSetWithdrawAddress
	NewMsgDisableService       = types.NewMsgDisableService
	NewMsgEnableService        = types.NewMsgEnableService
	NewMsgRefundServiceDeposit = types.NewMsgRefundServiceDeposit
	NewMsgRequestService       = types.NewMsgRequestService
	NewMsgRespondService       = types.NewMsgRespondService
	NewMsgStopRepeated         = types.NewMsgStopRepeated
	NewMsgUpdateRequestContext = types.NewMsgUpdateRequestContext
	NewMsgWithdrawEarnedFees   = types.NewMsgWithdrawEarnedFees
	NewMsgWithdrawTax          = types.NewMsgWithdrawTax
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
	Request                 = types.Request
	Response                = types.Response
	ReturnedFee             = types.EarnedFees
	Metrics                 = types.Metrics
	MsgDefineService        = types.MsgDefineService
	MsgBindService          = types.MsgBindService
	MsgUpdateServiceBinding = types.MsgUpdateServiceBinding
	MsgSetWithdrawAddress   = types.MsgSetWithdrawAddress
	MsgDisableService       = types.MsgDisableService
	MsgEnableService        = types.MsgEnableService
	MsgRefundServiceDeposit = types.MsgRefundServiceDeposit
	MsgRequestService       = types.MsgRequestService
	MsgRespondService       = types.MsgRespondService
	MsgStopRepeated         = types.MsgStopRepeated
	MsgUpdateRequestContext = types.MsgUpdateRequestContext
	MsgWithdrawEarnedFees   = types.MsgWithdrawEarnedFees
	MsgWithdrawTax          = types.MsgWithdrawTax
	QueryDefinitionParams   = types.QueryDefinitionParams
	QueryBindingParams      = types.QueryBindingParams
	QueryBindingsParams     = types.QueryBindingsParams
	QueryRequestsParams     = types.QueryRequestsParams
	QueryResponseParams     = types.QueryResponseParams
	QueryFeesParams         = types.QueryFeesParams
	FeesOutput              = types.FeesOutput
)
