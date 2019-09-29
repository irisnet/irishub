package v2

import (
	"fmt"
	"sort"
	"strings"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v1/asset"
	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/bank"
	distr "github.com/irisnet/irishub/app/v1/distribution"
	"github.com/irisnet/irishub/app/v1/gov"
	"github.com/irisnet/irishub/app/v1/mint"
	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/app/v1/rand"
	"github.com/irisnet/irishub/app/v1/service"
	"github.com/irisnet/irishub/app/v1/slashing"
	"github.com/irisnet/irishub/app/v1/stake"
	"github.com/irisnet/irishub/app/v1/upgrade"
	"github.com/irisnet/irishub/app/v2/coinswap"
	"github.com/irisnet/irishub/app/v2/htlc"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/guardian"
	sdk "github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
)

var _ protocol.Protocol = (*ProtocolV2)(nil)

// ProtocolV2 define the protocol
type ProtocolV2 struct {
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
	assetKeeper    asset.Keeper
	randKeeper     rand.Keeper
	coinswapKeeper coinswap.Keeper
	htlcKeeper     htlc.Keeper

	router      protocol.Router      // handle any kind of message
	queryRouter protocol.QueryRouter // router for redirecting query calls

	anteHandlers         []sdk.AnteHandler        // ante handlers for fee and auth
	feeRefundHandler     sdk.FeeRefundHandler     // fee handler for fee refund
	feePreprocessHandler sdk.FeePreprocessHandler // fee handler for fee preprocessor

	// may be nil
	initChainer  sdk.InitChainer1 // initialize state with validators and state blob
	beginBlocker sdk.BeginBlocker // logic to run before any txs
	endBlocker   sdk.EndBlocker   // logic to run after all txs, and to determine valset changes
	config       *cfg.InstrumentationConfig

	metrics *Metrics
}

// NewProtocolV2 protocol v2 constructor
func NewProtocolV2(version uint64, log log.Logger, pk sdk.ProtocolKeeper, checkInvariant bool, trackCoinFlow bool, config *cfg.InstrumentationConfig) *ProtocolV2 {
	p0 := ProtocolV2{
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

// Load load the configuration of this Protocol
func (p *ProtocolV2) Load() {
	p.configCodec()
	p.configKeepers()
	p.configRouters()
	p.configFeeHandlers()
	p.configParams()
}

// Init initializes the configuration of this Protocol
func (p *ProtocolV2) Init(ctx sdk.Context) {
	// initialize coinswap params
	p.coinswapKeeper.Init(ctx)
}

// GetCodec get codec
func (p *ProtocolV2) GetCodec() *codec.Codec {
	return p.cdc
}

// InitMetrics init prometheus metrics
func (p *ProtocolV2) InitMetrics(store sdk.MultiStore) {
	p.StakeKeeper.InitMetrics(store.GetKVStore(protocol.KeyStake))
	p.serviceKeeper.InitMetrics(store.GetKVStore(protocol.KeyService))
}

func (p *ProtocolV2) configCodec() {
	p.cdc = MakeCodec()
}

// MakeCodec register codec
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	params.RegisterCodec(cdc) // only used by querier
	mint.RegisterCodec(cdc)   // only used by querier
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
	asset.RegisterCodec(cdc)
	rand.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	coinswap.RegisterCodec(cdc)
	htlc.RegisterCodec(cdc)
	return cdc
}

// GetVersion get protocol version
func (p *ProtocolV2) GetVersion() uint64 {
	return p.version
}

// ValidateTx validate txs
func (p *ProtocolV2) ValidateTx(ctx sdk.Context, txBytes []byte, msgs []sdk.Msg) sdk.Error {

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
func (p *ProtocolV2) configKeepers() {
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
	p.bankKeeper = bank.NewBaseKeeper(
		p.cdc,
		p.accountMapper,
	)
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
		distr.DefaultCodespace, distr.PrometheusMetrics(p.config),
	)
	p.slashingKeeper = slashing.NewKeeper(
		p.cdc,
		protocol.KeySlashing,
		&stakeKeeper, p.paramsKeeper.Subspace(slashing.DefaultParamspace),
		slashing.DefaultCodespace,
		slashing.PrometheusMetrics(p.config),
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

	p.assetKeeper = asset.NewKeeper(p.cdc, protocol.KeyAsset, p.bankKeeper, asset.DefaultCodespace, p.paramsKeeper.Subspace(asset.DefaultParamSpace))

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
		p.assetKeeper,
	)

	p.randKeeper = rand.NewKeeper(p.cdc, protocol.KeyRand, rand.DefaultCodespace)
	p.coinswapKeeper = coinswap.NewKeeper(p.cdc, protocol.KeySwap, p.bankKeeper, p.accountMapper, p.paramsKeeper.Subspace(coinswap.DefaultParamSpace))
	p.htlcKeeper = htlc.NewKeeper(p.cdc, protocol.KeyHtlc, p.bankKeeper, htlc.DefaultCodespace)
}

// configure all Routers
func (p *ProtocolV2) configRouters() {
	p.router.
		AddRoute(protocol.BankRoute, bank.NewHandler(p.bankKeeper)).
		AddRoute(protocol.StakeRoute, stake.NewHandler(p.StakeKeeper)).
		AddRoute(protocol.SlashingRoute, slashing.NewHandler(p.slashingKeeper)).
		AddRoute(protocol.DistrRoute, distr.NewHandler(p.distrKeeper)).
		AddRoute(protocol.GovRoute, gov.NewHandler(p.govKeeper)).
		AddRoute(protocol.ServiceRoute, service.NewHandler(p.serviceKeeper)).
		AddRoute(protocol.GuardianRoute, guardian.NewHandler(p.guardianKeeper)).
		AddRoute(protocol.AssetRoute, asset.NewHandler(p.assetKeeper)).
		AddRoute(protocol.RandRoute, rand.NewHandler(p.randKeeper)).
		AddRoute(protocol.SwapRoute, coinswap.NewHandler(p.coinswapKeeper)).
		AddRoute(protocol.HtlcRoute, htlc.NewHandler(p.htlcKeeper))

	p.queryRouter.
		AddRoute(protocol.AccountRoute, bank.NewQuerier(p.bankKeeper, p.cdc)).
		AddRoute(protocol.GovRoute, gov.NewQuerier(p.govKeeper)).
		AddRoute(protocol.StakeRoute, stake.NewQuerier(p.StakeKeeper, p.cdc)).
		AddRoute(protocol.DistrRoute, distr.NewQuerier(p.distrKeeper)).
		AddRoute(protocol.GuardianRoute, guardian.NewQuerier(p.guardianKeeper)).
		AddRoute(protocol.ServiceRoute, service.NewQuerier(p.serviceKeeper)).
		AddRoute(protocol.ParamsRoute, params.NewQuerier(p.paramsKeeper)).
		AddRoute(protocol.AssetRoute, asset.NewQuerier(p.assetKeeper)).
		AddRoute(protocol.RandRoute, rand.NewQuerier(p.randKeeper)).
		AddRoute(protocol.SwapRoute, coinswap.NewQuerier(p.coinswapKeeper)).
		AddRoute(protocol.HtlcRoute, htlc.NewQuerier(p.htlcKeeper))
}

// configure all FeeHandlers
func (p *ProtocolV2) configFeeHandlers() {
	authAnteHandler := auth.NewAnteHandler(p.accountMapper, p.feeKeeper)
	assetAnteHandler := asset.NewAnteHandler(p.assetKeeper)
	bankAnteHandler := bank.NewAnteHandler(p.accountMapper)

	p.anteHandlers = []sdk.AnteHandler{authAnteHandler, bankAnteHandler, assetAnteHandler}
	p.feeRefundHandler = auth.NewFeeRefundHandler(p.accountMapper, p.feeKeeper)
	p.feePreprocessHandler = auth.NewFeePreprocessHandler(p.feeKeeper)
}

// GetKVStoreKeyList get KVStore Key List
func (p *ProtocolV2) GetKVStoreKeyList() []*sdk.KVStoreKey {
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
		protocol.KeyGuardian,
		protocol.KeyAsset,
		protocol.KeyRand,
		protocol.KeySwap,
		protocol.KeyHtlc,
	}
}

// configure all Params
func (p *ProtocolV2) configParams() {
	p.paramsKeeper.RegisterParamSet(&mint.Params{}, &slashing.Params{}, &service.Params{}, &auth.Params{}, &stake.Params{}, &distr.Params{}, &asset.Params{}, &gov.GovParams{}, &coinswap.Params{})
}

// BeginBlocker application updates every begin block
func (p *ProtocolV2) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	// mint new tokens for this new block
	tags := mint.BeginBlocker(ctx, p.mintKeeper)

	// distribute rewards from previous block
	distr.BeginBlocker(ctx, req, p.distrKeeper)

	slashTags := slashing.BeginBlocker(ctx, req, p.slashingKeeper)

	// handle pending random number requests
	randTags := rand.BeginBlocker(ctx, req, p.randKeeper)

	// handle HTLC expiration queue
	htlcTags := htlc.BeginBlocker(ctx, p.htlcKeeper)

	ctx.CoinFlowTags().TagWrite()

	tags = tags.AppendTags(slashTags).AppendTags(randTags).AppendTags(htlcTags)
	return abci.ResponseBeginBlock{
		Tags: tags.ToKVPairs(),
	}
}

// EndBlocker application updates every end block
func (p *ProtocolV2) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	tags := gov.EndBlocker(ctx, p.govKeeper)
	tags = tags.AppendTags(slashing.EndBlocker(ctx, req, p.slashingKeeper))
	tags = tags.AppendTags(service.EndBlocker(ctx, p.serviceKeeper))
	tags = tags.AppendTags(upgrade.EndBlocker(ctx, p.upgradeKeeper))
	validatorUpdates := stake.EndBlocker(ctx, p.StakeKeeper)
	if p.trackCoinFlow {
		ctx.CoinFlowTags().TagWrite()
		tags = tags.AppendTags(ctx.CoinFlowTags().GetTags())
	}
	p.assertRuntimeInvariants(ctx)

	return abci.ResponseEndBlock{
		ValidatorUpdates: validatorUpdates,
		Tags:             tags,
	}
}

// InitChainer Custom logic for iris initialization
// Only v0 needs InitChainer
func (p *ProtocolV2) InitChainer(ctx sdk.Context, DeliverTx sdk.DeliverTx, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	var genesisFileState GenesisFileState
	p.cdc.MustUnmarshalJSON(stateJSON, &genesisFileState)

	genesisState := convertToGenesisState(genesisFileState)
	// sort by account number to maintain consistency
	sort.Slice(genesisState.Accounts, func(i, j int) bool {
		return genesisState.Accounts[i].AccountNumber < genesisState.Accounts[j].AccountNumber
	})

	// init system accounts
	p.bankKeeper.AddCoins(ctx, auth.BurnedCoinsAccAddr, sdk.Coins{})
	p.bankKeeper.AddCoins(ctx, auth.CommunityTaxCoinsAccAddr, sdk.Coins{})

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
	asset.InitGenesis(ctx, p.assetKeeper, genesisState.AssetData)
	rand.InitGenesis(ctx, p.randKeeper, genesisState.RandData)
	coinswap.InitGenesis(ctx, p.coinswapKeeper, genesisState.SwapData)
	htlc.InitGenesis(ctx, p.htlcKeeper, genesisState.HtlcData)

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

func (p *ProtocolV2) GetRouter() protocol.Router {
	return p.router
}
func (p *ProtocolV2) GetQueryRouter() protocol.QueryRouter {
	return p.queryRouter
}
func (p *ProtocolV2) GetAnteHandlers() []sdk.AnteHandler {
	return p.anteHandlers
}
func (p *ProtocolV2) GetFeeRefundHandler() sdk.FeeRefundHandler {
	return p.feeRefundHandler
}
func (p *ProtocolV2) GetFeePreprocessHandler() sdk.FeePreprocessHandler {
	return p.feePreprocessHandler
}
func (p *ProtocolV2) GetInitChainer() sdk.InitChainer1 {
	return p.InitChainer
}
func (p *ProtocolV2) GetBeginBlocker() sdk.BeginBlocker {
	return p.BeginBlocker
}
func (p *ProtocolV2) GetEndBlocker() sdk.EndBlocker {
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
