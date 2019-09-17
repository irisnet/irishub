// nolint
package distribution

import (
	"github.com/irisnet/irishub/modules/distribution/keeper"
	"github.com/irisnet/irishub/modules/distribution/tags"
	"github.com/irisnet/irishub/modules/distribution/types"
)

type (
	Keeper = keeper.Keeper
	Hooks  = keeper.Hooks
	Params = types.Params

	DelegatorWithdrawInfo = types.DelegatorWithdrawInfo
	DelegationDistInfo    = types.DelegationDistInfo
	ValidatorDistInfo     = types.ValidatorDistInfo
	TotalAccum            = types.TotalAccum
	FeePool               = types.FeePool
	Rewards               = keeper.Rewards
	CommunityTax          = keeper.CommunityTax

	MsgWithdrawDelegatorRewardsAll = types.MsgWithdrawDelegatorRewardsAll
	MsgWithdrawDelegatorReward     = types.MsgWithdrawDelegatorReward
	MsgWithdrawValidatorRewardsAll = types.MsgWithdrawValidatorRewardsAll

	GenesisState = types.GenesisState

	// expected keepers
	StakeKeeper         = types.StakeKeeper
	BankKeeper          = types.BankKeeper
	FeeCollectionKeeper = types.FeeKeeper
)

var (
	NewKeeper = keeper.NewKeeper

	PrometheusMetrics = keeper.PrometheusMetrics
	NopMetrics        = keeper.NopMetrics

	GetValidatorDistInfoKey     = keeper.GetValidatorDistInfoKey
	GetDelegationDistInfoKey    = keeper.GetDelegationDistInfoKey
	GetDelegationDistInfosKey   = keeper.GetDelegationDistInfosKey
	GetDelegatorWithdrawAddrKey = keeper.GetDelegatorWithdrawAddrKey
	FeePoolKey                  = keeper.FeePoolKey
	ValidatorDistInfoKey        = keeper.ValidatorDistInfoKey
	DelegationDistInfoKey       = keeper.DelegationDistInfoKey
	DelegatorWithdrawInfoKey    = keeper.DelegatorWithdrawInfoKey
	ProposerKey                 = keeper.ProposerKey
	DefaultParamspace           = keeper.DefaultParamspace

	InitialFeePool = types.InitialFeePool

	NewGenesisState              = types.NewGenesisState
	DefaultGenesisState          = types.DefaultGenesisState
	DefaultGenesisWithValidators = types.DefaultGenesisWithValidators

	RegisterCodec = types.RegisterCodec

	NewMsgWithdrawDelegatorRewardsAll = types.NewMsgWithdrawDelegatorRewardsAll
	NewMsgWithdrawDelegatorReward     = types.NewMsgWithdrawDelegatorReward
	NewMsgWithdrawValidatorRewardsAll = types.NewMsgWithdrawValidatorRewardsAll

	NewQuerier                       = keeper.NewQuerier
	NewQueryDelegatorParams          = keeper.NewQueryDelegatorParams
	NewQueryDelegationDistInfoParams = keeper.NewQueryDelegationDistInfoParams
	NewQueryValidatorDistInfoParams  = keeper.NewQueryValidatorDistInfoParams
	NewQueryRewardsParams            = keeper.NewQueryRewardsParams

	NewTotalAccum = types.NewTotalAccum
)

const (
	DefaultCodespace = types.DefaultCodespace
	CodeInvalidInput = types.CodeInvalidInput
)

const (
	QueryWithdrawAddr          = keeper.QueryWithdrawAddr
	QueryDelegationDistInfo    = keeper.QueryDelegationDistInfo
	QueryAllDelegationDistInfo = keeper.QueryAllDelegationDistInfo
	QueryValidatorDistInfo     = keeper.QueryValidatorDistInfo
	QueryRewards               = keeper.QueryRewards
	QueryCommunityTax          = keeper.QueryCommunityTax
)

var (
	ErrNilDelegatorAddr = types.ErrNilDelegatorAddr
	ErrNilWithdrawAddr  = types.ErrNilWithdrawAddr
	ErrNilValidatorAddr = types.ErrNilValidatorAddr

	ActionModifyWithdrawAddress       = tags.ActionModifyWithdrawAddress
	ActionWithdrawDelegatorRewardsAll = tags.ActionWithdrawDelegatorRewardsAll
	ActionWithdrawDelegatorReward     = tags.ActionWithdrawDelegatorReward
	ActionWithdrawValidatorRewardsAll = tags.ActionWithdrawValidatorRewardsAll

	TagAction    = tags.Action
	TagValidator = tags.Validator
	TagDelegator = tags.Delegator
)
