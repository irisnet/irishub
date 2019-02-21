package v0

import (
	"fmt"
	"sort"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/auth"
	"github.com/irisnet/irishub/modules/bank"
	distr "github.com/irisnet/irishub/modules/distribution"
	"github.com/irisnet/irishub/modules/gov"
	"github.com/irisnet/irishub/modules/guardian"
	"github.com/irisnet/irishub/modules/mint"
	"github.com/irisnet/irishub/modules/params"
	"github.com/irisnet/irishub/modules/service"
	"github.com/irisnet/irishub/modules/slashing"
	"github.com/irisnet/irishub/modules/stake"
	"github.com/irisnet/irishub/modules/upgrade"
	sdk "github.com/irisnet/irishub/types"

	"strings"

	abci "github.com/tendermint/tendermint/abci/types"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
)

var _ protocol.Protocol = (*ProtocolV0)(nil)

type ProtocolV0 struct {
	version        uint64
	cdc            *codec.Codec
	logger         log.Logger
	invariantLevel string
	checkInvariant bool
	trackCoinFlow  bool

	// Manage getting and setting accounts
	accountMapper  auth.AccountKeeper
	feeKeeper      auth.FeeKeeper
	bankKeeper     bank.Keeper
	StakeKeeper    stake.Keeper
	slashingKeeper slashing.Keeper
	mintKeeper     mint.Keeper
	distrKeeper    distr.Keeper
	protocolKeeper sdk.ProtocolKeeper
	govKeeper      gov.Keeper
	paramsKeeper   params.Keeper
	serviceKeeper  service.Keeper
	guardianKeeper guardian.Keeper
	upgradeKeeper  upgrade.Keeper

	router      protocol.Router      // handle any kind of message
	queryRouter protocol.QueryRouter // router for redirecting query calls

	anteHandler          sdk.AnteHandler          // ante handler for fee and auth
	feeRefundHandler     sdk.FeeRefundHandler     // fee handler for fee refund
	feePreprocessHandler sdk.FeePreprocessHandler // fee handler for fee preprocessor

	// may be nil
	initChainer  sdk.InitChainer1 // initialize state with validators and state blob
	beginBlocker sdk.BeginBlocker // logic to run before any txs
	endBlocker   sdk.EndBlocker   // logic to run after all txs, and to determine valset changes
	config       *cfg.InstrumentationConfig

	metrics		*Metrics
}

func NewProtocolV0(version uint64, log log.Logger, pk sdk.ProtocolKeeper, checkInvariant bool, trackCoinFlow bool, config *cfg.InstrumentationConfig) *ProtocolV0 {
	p0 := ProtocolV0{
		version:        version,
		logger:         log,
		protocolKeeper: pk,
		invariantLevel: strings.ToLower(sdk.InvariantLevel),
		checkInvariant: checkInvariant,
		trackCoinFlow:  trackCoinFlow,
		router:         protocol.NewRouter(),
		queryRouter:    protocol.NewQueryRouter(),
		config:         config,
		metrics:        PrometheusMetrics(config),
	}
	return &p0
}

// load the configuration of this Protocol
func (p *ProtocolV0) Load() {
	p.configCodec()
	p.configKeepers()
	p.configRouters()
	p.configFeeHandlers()
	p.configParams()
}

// verison0 don't need the init
func (p *ProtocolV0) Init() {

}

// verison0 tx codec
func (p *ProtocolV0) GetCodec() *codec.Codec {
	return p.cdc
}

func (p *ProtocolV0) configCodec() {
	p.cdc = MakeCodec()
}

func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	bank.RegisterCodec(cdc)
	stake.RegisterCodec(cdc)
	distr.RegisterCodec(cdc)
	slashing.RegisterCodec(cdc)
	gov.RegisterCodec(cdc)
	upgrade.RegisterCodec(cdc)
	service.RegisterCodec(cdc)
	guardian.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

func (p *ProtocolV0) GetVersion() uint64 {
	return p.version
}

func (p *ProtocolV0) ValidateTx(ctx sdk.Context, txBytes []byte, msgs []sdk.Msg) sdk.Error {

	serviceMsgNum := 0
	for _, msg := range msgs {
		if msg.Route() == service.MsgRoute {
			serviceMsgNum++
		}
	}

	if serviceMsgNum != 0 && serviceMsgNum != len(msgs) {
		return sdk.ErrServiceTxLimit("Can't mix service msgs with other types of msg in one transaction!")
	}

	if serviceMsgNum == 0 {
		subspace, found := p.paramsKeeper.GetSubspace(auth.DefaultParamSpace)
		var txSizeLimit uint64
		if found {
			subspace.Get(ctx, auth.TxSizeLimitKey, &txSizeLimit)
		} else {
			panic("The subspace " + auth.DefaultParamSpace + " cannot be found!")
		}
		if uint64(len(txBytes)) > txSizeLimit {
			return sdk.ErrExceedsTxSize(fmt.Sprintf("the tx size [%d] exceeds the limitation [%d]", len(txBytes), txSizeLimit))
		}
	}

	if serviceMsgNum == len(msgs) {
		subspace, found := p.paramsKeeper.GetSubspace(service.DefaultParamSpace)
		var serviceTxSizeLimit uint64
		if found {
			subspace.Get(ctx, service.KeyTxSizeLimit, &serviceTxSizeLimit)
		} else {
			panic("The subspace " + service.DefaultParamSpace + " cannot be found!")
		}

		if uint64(len(txBytes)) > serviceTxSizeLimit {
			return sdk.ErrExceedsTxSize(fmt.Sprintf("the tx size [%d] exceeds the limitation [%d]", len(txBytes), serviceTxSizeLimit))
		}

	}

	return nil
}

// create all Keepers
func (p *ProtocolV0) configKeepers() {
	// define the AccountKeeper
	p.accountMapper = auth.NewAccountKeeper(
		p.cdc,
		protocol.KeyAccount,   // target store
		auth.ProtoBaseAccount, // prototype
	)

	// add handlers
	p.guardianKeeper = guardian.NewKeeper(
		p.cdc,
		protocol.KeyGuardian,
		guardian.DefaultCodespace,
	)
	p.bankKeeper = bank.NewBaseKeeper(p.accountMapper)
	p.paramsKeeper = params.NewKeeper(
		p.cdc,
		protocol.KeyParams, protocol.TkeyParams,
	)
	p.feeKeeper = auth.NewFeeKeeper(
		p.cdc,
		protocol.KeyFee, p.paramsKeeper.Subspace(auth.DefaultParamSpace),
	)
	stakeKeeper := stake.NewKeeper(
		p.cdc,
		protocol.KeyStake, protocol.TkeyStake,
		p.bankKeeper, p.paramsKeeper.Subspace(stake.DefaultParamspace),
		stake.DefaultCodespace,
		stake.PrometheusMetrics(p.config),
	)
	p.mintKeeper = mint.NewKeeper(p.cdc, protocol.KeyMint,
		p.paramsKeeper.Subspace(mint.DefaultParamSpace),
		p.bankKeeper, p.feeKeeper,
	)
	p.distrKeeper = distr.NewKeeper(
		p.cdc,
		protocol.KeyDistr,
		p.paramsKeeper.Subspace(distr.DefaultParamspace),
		p.bankKeeper, &stakeKeeper, p.feeKeeper,
		distr.DefaultCodespace,  distr.PrometheusMetrics(p.config),
	)
	p.slashingKeeper = slashing.NewKeeper(
		p.cdc,
		protocol.KeySlashing,
		&stakeKeeper, p.paramsKeeper.Subspace(slashing.DefaultParamspace),
		slashing.DefaultCodespace,
		slashing.PrometheusMetrics(p.config),
	)

	p.govKeeper = gov.NewKeeper(
		protocol.KeyGov,
		p.cdc,
		p.paramsKeeper.Subspace(gov.DefaultParamSpace),
		p.paramsKeeper,
		p.protocolKeeper,
		p.bankKeeper,
		p.distrKeeper,
		p.guardianKeeper,
		&stakeKeeper,
		gov.DefaultCodespace,
		gov.PrometheusMetrics(p.config),
	)

	p.serviceKeeper = service.NewKeeper(
		p.cdc,
		protocol.KeyService,
		p.bankKeeper,
		p.guardianKeeper,
		service.DefaultCodespace,
		p.paramsKeeper.Subspace(service.DefaultParamSpace),
		service.PrometheusMetrics(p.config),
	)

	// register the staking hooks
	// NOTE: StakeKeeper above are passed by reference,
	// so that it can be modified like below:
	p.StakeKeeper = *stakeKeeper.SetHooks(
		NewHooks(p.distrKeeper.Hooks(), p.slashingKeeper.Hooks()))

	p.upgradeKeeper = upgrade.NewKeeper(p.cdc, protocol.KeyUpgrade, p.protocolKeeper, p.StakeKeeper, upgrade.PrometheusMetrics(p.config))
}

// configure all Routers
func (p *ProtocolV0) configRouters() {
	p.router.
		AddRoute("bank", bank.NewHandler(p.bankKeeper)).
		AddRoute("stake", stake.NewHandler(p.StakeKeeper)).
		AddRoute("slashing", slashing.NewHandler(p.slashingKeeper)).
		AddRoute("distr", distr.NewHandler(p.distrKeeper)).
		AddRoute("gov", gov.NewHandler(p.govKeeper)).
		AddRoute("service", service.NewHandler(p.serviceKeeper)).
		AddRoute("guardian", guardian.NewHandler(p.guardianKeeper))
	p.queryRouter.
		AddRoute("gov", gov.NewQuerier(p.govKeeper)).
		AddRoute("stake", stake.NewQuerier(p.StakeKeeper, p.cdc))
}

// configure all Stores
func (p *ProtocolV0) configFeeHandlers() {
	p.anteHandler = auth.NewAnteHandler(p.accountMapper, p.feeKeeper)
	p.feeRefundHandler = auth.NewFeeRefundHandler(p.accountMapper, p.feeKeeper)
	p.feePreprocessHandler = auth.NewFeePreprocessHandler(p.feeKeeper)
}

// configure all Stores
func (p *ProtocolV0) GetKVStoreKeyList() []*sdk.KVStoreKey {
	return []*sdk.KVStoreKey{
		protocol.KeyMain,
		protocol.KeyAccount,
		protocol.KeyStake,
		protocol.KeyMint,
		protocol.KeyDistr,
		protocol.KeySlashing,
		protocol.KeyGov,
		protocol.KeyFee,
		protocol.KeyParams,
		protocol.KeyUpgrade,
		protocol.KeyService,
		protocol.KeyGuardian}
}

// configure all Stores
func (p *ProtocolV0) configParams() {

	p.paramsKeeper.RegisterParamSet(&mint.Params{}, &slashing.Params{}, &service.Params{}, &auth.Params{}, &stake.Params{}, &distr.Params{})

}

// application updates every end block
func (p *ProtocolV0) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	// mint new tokens for this new block
	tags := mint.BeginBlocker(ctx, p.mintKeeper)

	// distribute rewards from previous block
	distr.BeginBlocker(ctx, req, p.distrKeeper)

	slashTags := slashing.BeginBlocker(ctx, req, p.slashingKeeper)

	ctx.CoinFlowTags().TagWrite()

	tags = tags.AppendTags(slashTags)
	return abci.ResponseBeginBlock{
		Tags: tags.ToKVPairs(),
	}
}

// application updates every end block
func (p *ProtocolV0) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	tags := gov.EndBlocker(ctx, p.govKeeper)
	tags = tags.AppendTags(slashing.EndBlocker(ctx, req, p.slashingKeeper))
	tags = tags.AppendTags(service.EndBlocker(ctx, p.serviceKeeper))
	tags = tags.AppendTags(upgrade.EndBlocker(ctx, p.upgradeKeeper))
	validatorUpdates := stake.EndBlocker(ctx, p.StakeKeeper)
	if p.trackCoinFlow {
		ctx.CoinFlowTags().TagWrite()
		tags = tags.AppendTags(extractCoinFlowTags(ctx))
	}
	p.assertRuntimeInvariants(ctx)

	return abci.ResponseEndBlock{
		ValidatorUpdates: validatorUpdates,
		Tags:             tags,
	}
}

func extractCoinFlowTags(ctx sdk.Context) sdk.Tags {
	var tags sdk.Tags
	for _, tag := range ctx.CoinFlowTags().GetTags() {
		ctx.Logger().Error("CoinFlowRecord","key", string(tag.Key), "value", string(tag.Value))
		tags = tags.AppendTag(string(tag.Key), tag.Value)
	}
	return tags
}

// custom logic for iris initialization
// just 0 version need Initchainer
func (p *ProtocolV0) InitChainer(ctx sdk.Context, DeliverTx sdk.DeliverTx, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	var genesisFileState GenesisFileState
	p.cdc.MustUnmarshalJSON(stateJSON, &genesisFileState)

	genesisState := convertToGenesisState(genesisFileState)
	// sort by account number to maintain consistency
	sort.Slice(genesisState.Accounts, func(i, j int) bool {
		return genesisState.Accounts[i].AccountNumber < genesisState.Accounts[j].AccountNumber
	})

	// load the accounts
	for _, gacc := range genesisState.Accounts {
		acc := gacc.ToAccount()
		acc.AccountNumber = p.accountMapper.GetNextAccountNumber(ctx)
		p.accountMapper.SetGenesisAccount(ctx, acc)
	}

	//upgrade.InitGenesis(ctx, p.upgradeKeeper, p.Router(), genesisState.UpgradeData)

	// load the initial stake information
	validators, err := stake.InitGenesis(ctx, p.StakeKeeper, genesisState.StakeData)
	if err != nil {
		panic(err)
	}
	gov.InitGenesis(ctx, p.govKeeper, genesisState.GovData)
	auth.InitGenesis(ctx, p.feeKeeper, p.accountMapper, genesisState.AuthData)
	slashing.InitGenesis(ctx, p.slashingKeeper, genesisState.SlashingData, genesisState.StakeData)
	mint.InitGenesis(ctx, p.mintKeeper, genesisState.MintData)
	distr.InitGenesis(ctx, p.distrKeeper, genesisState.DistrData)
	service.InitGenesis(ctx, p.serviceKeeper, genesisState.ServiceData)
	guardian.InitGenesis(ctx, p.guardianKeeper, genesisState.GuardianData)
	upgrade.InitGenesis(ctx, p.upgradeKeeper, genesisState.UpgradeData)

	// load the address to pubkey map
	err = IrisValidateGenesisState(genesisState)
	if err != nil {
		panic(err) // TODO find a way to do this w/o panics
	}

	if len(genesisState.GenTxs) > 0 {
		for _, genTx := range genesisState.GenTxs {
			var tx auth.StdTx
			err = p.cdc.UnmarshalJSON(genTx, &tx)
			if err != nil {
				panic(err)
			}
			bz := p.cdc.MustMarshalBinaryLengthPrefixed(tx)
			res := DeliverTx(bz)
			if !res.IsOK() {
				panic(res.Log)
			}
		}

		validators = p.StakeKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	}

	// sanity check
	if len(req.Validators) > 0 {
		if len(req.Validators) != len(validators) {
			panic(fmt.Errorf("len(RequestInitChain.Validators) != len(validators) (%d != %d)",
				len(req.Validators), len(validators)))
		}
		sort.Sort(abci.ValidatorUpdates(req.Validators))
		sort.Sort(abci.ValidatorUpdates(validators))
		for i, val := range validators {
			if !val.Equal(req.Validators[i]) {
				panic(fmt.Errorf("validators[%d] != req.Validators[%d] ", i, i))
			}
		}
	}
	return abci.ResponseInitChain{
		Validators: validators,
	}
}

func (p *ProtocolV0) GetRouter() protocol.Router {
	return p.router
}
func (p *ProtocolV0) GetQueryRouter() protocol.QueryRouter {
	return p.queryRouter
}
func (p *ProtocolV0) GetAnteHandler() sdk.AnteHandler {
	return p.anteHandler
}
func (p *ProtocolV0) GetFeeRefundHandler() sdk.FeeRefundHandler {
	return p.feeRefundHandler
}
func (p *ProtocolV0) GetFeePreprocessHandler() sdk.FeePreprocessHandler {
	return p.feePreprocessHandler
}
func (p *ProtocolV0) GetInitChainer() sdk.InitChainer1 {
	return p.InitChainer
}
func (p *ProtocolV0) GetBeginBlocker() sdk.BeginBlocker {
	return p.BeginBlocker
}
func (p *ProtocolV0) GetEndBlocker() sdk.EndBlocker {
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
