package mock

import (
	"encoding/json"
	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/arbitration/params"
	"github.com/irisnet/irishub/modules/auth"
	"github.com/irisnet/irishub/modules/bank"
	distr "github.com/irisnet/irishub/modules/distribution"
	"github.com/irisnet/irishub/modules/gov"
	"github.com/irisnet/irishub/modules/gov/params"
	"github.com/irisnet/irishub/modules/guardian"
	"github.com/irisnet/irishub/modules/mint"
	"github.com/irisnet/irishub/modules/params"
	"github.com/irisnet/irishub/modules/record"
	"github.com/irisnet/irishub/modules/service"
	"github.com/irisnet/irishub/modules/service/params"
	"github.com/irisnet/irishub/modules/slashing"
	"github.com/irisnet/irishub/modules/stake"
	"github.com/irisnet/irishub/modules/upgrade/params"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/types/common"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"time"
)

var _ protocol.Protocol = (*ProtocolVersion0)(nil)

type ProtocolVersion0 struct {
	pb  *protocol.ProtocolBase
	cdc *codec.Codec

	// Manage getting and setting accounts
	AccountMapper       auth.AccountKeeper
	FeeCollectionKeeper auth.FeeCollectionKeeper
	BankKeeper          bank.Keeper
	StakeKeeper         stake.Keeper
	SlashingKeeper      slashing.Keeper
	MintKeeper          mint.Keeper
	DistrKeeper         distr.Keeper
	GovKeeper           gov.Keeper
	ParamsKeeper        params.Keeper
	ServiceKeeper       service.Keeper
	GuardianKeeper      guardian.Keeper
	RecordKeeper        record.Keeper
	// fee manager
	FeeManager auth.FeeManager

	GenesisAccounts  []auth.Account

	Router      protocol.Router      // handle any kind of message
	QueryRouter protocol.QueryRouter // Router for redirecting query calls

	AnteHandler          sdk.AnteHandler          // ante handler for fee and auth
	FeeRefundHandler     sdk.FeeRefundHandler     // fee handler for fee refund
	FeePreprocessHandler sdk.FeePreprocessHandler // fee handler for fee preprocessor

	// may be nil
	initChainer  sdk.InitChainer1 // initialize state with validators and state blob
	beginBlocker sdk.BeginBlocker // logic to run before any txs
	endBlocker   sdk.EndBlocker   // logic to run after all txs, and to determine valset changes

}

func NewProtocolVersion0(cdc *codec.Codec) *ProtocolVersion0 {
	base := protocol.ProtocolBase{
		Definition: common.ProtocolDefinition{
			uint64(0),
			"https://github.com/irisnet/irishub/releases/tag/v0.7.0",
			uint64(1),
		},
		//		engine: engine,
	}
	p0 := ProtocolVersion0{
		pb:          &base,
		cdc:         cdc,
		Router:      protocol.NewRouter(),
		QueryRouter: protocol.NewQueryRouter(),
	}
	return &p0
}

// load the configuration of this Protocol
func (p *ProtocolVersion0) Load() {
	p.configKeepers()
	p.configRouters()
	p.configFeeHandlers()
	p.configParams()
	p.configStores()
}

// verison0 don't need the init
func (p *ProtocolVersion0) Init() {

}

func (p *ProtocolVersion0) GetDefinition() common.ProtocolDefinition {
	return p.pb.GetDefinition()
}

// create all Keepers
func (p *ProtocolVersion0) configKeepers() {
	// define the AccountKeeper
	p.AccountMapper = auth.NewAccountKeeper(
		p.cdc,
		protocol.KeyAccount,   // target store
		auth.ProtoBaseAccount, // prototype
	)

	// add handlers
	p.BankKeeper = bank.NewBaseKeeper(p.AccountMapper)
	p.FeeCollectionKeeper = auth.NewFeeCollectionKeeper(
		p.cdc,
		protocol.KeyFeeCollection,
	)
	p.ParamsKeeper = params.NewKeeper(
		p.cdc,
		protocol.KeyParams, protocol.TkeyParams,
	)
	stakeKeeper := stake.NewKeeper(
		p.cdc,
		protocol.KeyStake, protocol.TkeyStake,
		p.BankKeeper, p.ParamsKeeper.Subspace(stake.DefaultParamspace),
		stake.DefaultCodespace,
	)
	p.MintKeeper = mint.NewKeeper(p.cdc, protocol.KeyMint,
		p.ParamsKeeper.Subspace(mint.DefaultParamspace),
		&stakeKeeper, p.FeeCollectionKeeper,
	)
	p.DistrKeeper = distr.NewKeeper(
		p.cdc,
		protocol.KeyDistr,
		p.ParamsKeeper.Subspace(distr.DefaultParamspace),
		p.BankKeeper, &stakeKeeper, p.FeeCollectionKeeper,
		distr.DefaultCodespace,
	)
	p.SlashingKeeper = slashing.NewKeeper(
		p.cdc,
		protocol.KeySlashing,
		&stakeKeeper, p.ParamsKeeper.Subspace(slashing.DefaultParamspace),
		slashing.DefaultCodespace,
	)

	p.GovKeeper = gov.NewKeeper(
		p.cdc,
		protocol.KeyGov,
		p.BankKeeper, &stakeKeeper,
		gov.DefaultCodespace,
	)

	p.RecordKeeper = record.NewKeeper(
		p.cdc,
		protocol.KeyRecord,
		record.DefaultCodespace,
	)
	p.ServiceKeeper = service.NewKeeper(
		p.cdc,
		protocol.KeyService,
		p.BankKeeper,
		p.GuardianKeeper,
		service.DefaultCodespace,
	)
	p.GuardianKeeper = guardian.NewKeeper(
		p.cdc,
		protocol.KeyGuardian,
		guardian.DefaultCodespace,
	)

	// register the staking hooks
	// NOTE: StakeKeeper above are passed by reference,
	// so that it can be modified like below:
	p.StakeKeeper = *stakeKeeper.SetHooks(
		NewHooks(p.DistrKeeper.Hooks(), p.SlashingKeeper.Hooks()))
	p.FeeManager = auth.NewFeeManager(p.ParamsKeeper.Subspace("Fee"))

}

// configure all Routers
func (p *ProtocolVersion0) configRouters() {
	p.Router.
		AddRoute("bank", bank.NewHandler(p.BankKeeper)).
		AddRoute("stake", stake.NewHandler(p.StakeKeeper)).
		AddRoute("slashing", slashing.NewHandler(p.SlashingKeeper)).
		AddRoute("distr", distr.NewHandler(p.DistrKeeper)).
		AddRoute("gov", gov.NewHandler(p.GovKeeper)).
		AddRoute("record", record.NewHandler(p.RecordKeeper)).
		AddRoute("service", service.NewHandler(p.ServiceKeeper)).
		AddRoute("guardian", guardian.NewHandler(p.GuardianKeeper))
	p.QueryRouter.
		AddRoute("gov", gov.NewQuerier(p.GovKeeper))
}

// configure all Stores
func (p *ProtocolVersion0) configFeeHandlers() {

	p.AnteHandler = auth.NewAnteHandler(p.AccountMapper, p.FeeCollectionKeeper)
	p.FeeRefundHandler = auth.NewFeeRefundHandler(p.AccountMapper, p.FeeCollectionKeeper, p.FeeManager)
	p.FeePreprocessHandler = auth.NewFeePreprocessHandler(p.FeeManager)
}

// configure all Stores
func (p *ProtocolVersion0) configStores() {

}

// configure all Stores
func (p *ProtocolVersion0) configParams() {
	params.SetParamReadWriter(p.ParamsKeeper.Subspace(params.SignalParamspace).WithTypeTable(
		params.NewTypeTable(
			upgradeparams.CurrentUpgradeProposalIdParameter.GetStoreKey(), uint64((0)),
			upgradeparams.ProposalAcceptHeightParameter.GetStoreKey(), int64(0),
			upgradeparams.SwitchPeriodParameter.GetStoreKey(), int64(0),
		)),
		&upgradeparams.CurrentUpgradeProposalIdParameter,
		&upgradeparams.ProposalAcceptHeightParameter,
		&upgradeparams.SwitchPeriodParameter)

	params.SetParamReadWriter(p.ParamsKeeper.Subspace(params.GovParamspace).WithTypeTable(
		params.NewTypeTable(
			govparams.DepositProcedureParameter.GetStoreKey(), govparams.DepositProcedure{},
			govparams.VotingProcedureParameter.GetStoreKey(), govparams.VotingProcedure{},
			govparams.TallyingProcedureParameter.GetStoreKey(), govparams.TallyingProcedure{},
			serviceparams.MaxRequestTimeoutParameter.GetStoreKey(), int64(0),
			serviceparams.MinDepositMultipleParameter.GetStoreKey(), int64(0),
			arbitrationparams.ComplaintRetrospectParameter.GetStoreKey(), time.Duration(0),
			arbitrationparams.ArbitrationTimelimitParameter.GetStoreKey(), time.Duration(0),
		)),
		&govparams.DepositProcedureParameter,
		&govparams.VotingProcedureParameter,
		&govparams.TallyingProcedureParameter,
		&serviceparams.MaxRequestTimeoutParameter,
		&serviceparams.MinDepositMultipleParameter,
		&arbitrationparams.ComplaintRetrospectParameter,
		&arbitrationparams.ArbitrationTimelimitParameter)

	params.RegisterGovParamMapping(
		&govparams.DepositProcedureParameter,
		&govparams.VotingProcedureParameter,
		&govparams.TallyingProcedureParameter,
		&serviceparams.MaxRequestTimeoutParameter,
		&serviceparams.MinDepositMultipleParameter)
}

// application updates every end block
func (p *ProtocolVersion0) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	tags := slashing.BeginBlocker(ctx, req, p.SlashingKeeper)

	// distribute rewards from previous block
	distr.BeginBlocker(ctx, req, p.DistrKeeper)

	// mint new tokens for this new block
	mint.BeginBlocker(ctx, p.MintKeeper)

	return abci.ResponseBeginBlock{
		Tags: tags.ToKVPairs(),
	}
}

// application updates every end block
func (p *ProtocolVersion0) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	tags := gov.EndBlocker(ctx, p.GovKeeper)
	validatorUpdates := stake.EndBlocker(ctx, p.StakeKeeper)
	tags = tags.AppendTags(service.EndBlocker(ctx, p.ServiceKeeper))
	return abci.ResponseEndBlock{
		ValidatorUpdates: validatorUpdates,
		Tags:             tags,
	}
}

func (p *ProtocolVersion0) ExportAppStateAndValidators(ctx sdk.Context) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	return
}

// custom logic for iris initialization
// just 0 version need Initchainer
func (p *ProtocolVersion0) InitChainer(ctx sdk.Context, DeliverTx sdk.DeliverTx, req abci.RequestInitChain) abci.ResponseInitChain {
	// Load the genesis accounts
	for _, genacc := range p.GenesisAccounts {
		acc := p.AccountMapper.NewAccountWithAddress(ctx, genacc.GetAddress())
		acc.SetCoins(genacc.GetCoins())
		p.AccountMapper.SetAccount(ctx, acc)
	}

	feeTokenGensisConfig := auth.FeeGenesisStateConfig{
		FeeTokenNative:    IrisCt.MinUnit.Denom,
		GasPriceThreshold: 0, // for mock test
	}

	auth.InitGenesis(ctx, p.FeeCollectionKeeper, auth.DefaultGenesisState(), p.FeeManager, feeTokenGensisConfig)

	return abci.ResponseInitChain{}
}

func (p *ProtocolVersion0) GetRouter() protocol.Router {
	return p.Router
}
func (p *ProtocolVersion0) GetQueryRouter() protocol.QueryRouter {
	return p.QueryRouter
}
func (p *ProtocolVersion0) GetAnteHandler() sdk.AnteHandler {
	return p.AnteHandler
}
func (p *ProtocolVersion0) GetFeeRefundHandler() sdk.FeeRefundHandler {
	return p.FeeRefundHandler
}
func (p *ProtocolVersion0) GetFeePreprocessHandler() sdk.FeePreprocessHandler {
	return p.FeePreprocessHandler
}
func (p *ProtocolVersion0) GetInitChainer() sdk.InitChainer1 {
	return p.InitChainer
}
func (p *ProtocolVersion0) GetBeginBlocker() sdk.BeginBlocker {
	return p.BeginBlocker
}
func (p *ProtocolVersion0) GetEndBlocker() sdk.EndBlocker {
	return p.EndBlocker
}

// Combined Staking Hooks
type Hooks struct {
	dh distr.Hooks
	sh slashing.Hooks
}

func NewHooks(dh distr.Hooks, sh slashing.Hooks) Hooks {
	return Hooks{dh, sh}
}

var _ sdk.StakingHooks = Hooks{}

func (h Hooks) OnValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) {
	h.dh.OnValidatorCreated(ctx, valAddr)
	h.sh.OnValidatorCreated(ctx, valAddr)
}
func (h Hooks) OnValidatorModified(ctx sdk.Context, valAddr sdk.ValAddress) {
	h.dh.OnValidatorModified(ctx, valAddr)
	h.sh.OnValidatorModified(ctx, valAddr)
}

func (h Hooks) OnValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.dh.OnValidatorRemoved(ctx, consAddr, valAddr)
	h.sh.OnValidatorRemoved(ctx, consAddr, valAddr)
}

func (h Hooks) OnValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.dh.OnValidatorBonded(ctx, consAddr, valAddr)
	h.sh.OnValidatorBonded(ctx, consAddr, valAddr)
}

func (h Hooks) OnValidatorPowerDidChange(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.dh.OnValidatorPowerDidChange(ctx, consAddr, valAddr)
	h.sh.OnValidatorPowerDidChange(ctx, consAddr, valAddr)
}

func (h Hooks) OnValidatorBeginUnbonding(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.dh.OnValidatorBeginUnbonding(ctx, consAddr, valAddr)
	h.sh.OnValidatorBeginUnbonding(ctx, consAddr, valAddr)
}

func (h Hooks) OnDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.dh.OnDelegationCreated(ctx, delAddr, valAddr)
	h.sh.OnDelegationCreated(ctx, delAddr, valAddr)
}

func (h Hooks) OnDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.dh.OnDelegationSharesModified(ctx, delAddr, valAddr)
	h.sh.OnDelegationSharesModified(ctx, delAddr, valAddr)
}

func (h Hooks) OnDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.dh.OnDelegationRemoved(ctx, delAddr, valAddr)
	h.sh.OnDelegationRemoved(ctx, delAddr, valAddr)
}
