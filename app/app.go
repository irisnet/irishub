package app

import (
	"io"
	"os"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/irisnet/irishub/modules/asset"
	token "github.com/irisnet/irishub/modules/asset/01-token"
	"github.com/irisnet/irishub/modules/coinswap"
	"github.com/irisnet/irishub/modules/guardian"
	"github.com/irisnet/irishub/modules/htlc"
	"github.com/irisnet/irishub/modules/mint"
	"github.com/irisnet/irishub/modules/rand"
	"github.com/irisnet/irishub/modules/service"
)

const appName = "IrisApp"

var (
	// default home directories for iriscli
	DefaultCLIHome = os.ExpandEnv("$HOME/.iriscli")

	// default home directories for iris
	DefaultNodeHome = os.ExpandEnv("$HOME/.iris")

	// The module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(paramsclient.ProposalHandler, distr.ProposalHandler),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		supply.AppModuleBasic{},
		asset.AppModuleBasic{},
		guardian.AppModuleBasic{},
		htlc.AppModuleBasic{},
		coinswap.AppModuleBasic{},
		rand.AppModuleBasic{},
		service.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		mint.ModuleName:           {supply.Minter},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
		token.SubModuleName:       {supply.Minter, supply.Burner},
		htlc.ModuleName:           nil,
		coinswap.ModuleName:       {supply.Minter, supply.Burner},
		service.DepositAccName:    {supply.Burner},
		service.RequestAccName:    nil,
		service.TaxAccName:        nil,
	}
)

// MakeCodec creates the application codec. The codec is sealed before it is
// returned.
func MakeCodec() *codec.Codec {
	var cdc = codec.New()

	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	codec.RegisterEvidences(cdc)
	authvesting.RegisterCodec(cdc)

	return cdc.Seal()
}

// Verify app interface at compile time
var _ simapp.App = (*IrisApp)(nil)

// IrisApp extended ABCI application
type IrisApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tKeys map[string]*sdk.TransientStoreKey

	// keepers
	accountKeeper  auth.AccountKeeper
	bankKeeper     bank.Keeper
	supplyKeeper   supply.Keeper
	stakingKeeper  staking.Keeper
	slashingKeeper slashing.Keeper
	mintKeeper     mint.Keeper
	distrKeeper    distr.Keeper
	govKeeper      gov.Keeper
	crisisKeeper   crisis.Keeper
	paramsKeeper   params.Keeper
	evidenceKeeper *evidence.Keeper
	assetKeeper    asset.Keeper
	guardianKeeper guardian.Keeper
	htlcKeeper     htlc.Keeper
	coinswapKeeper coinswap.Keeper
	randKeeper     rand.Keeper
	serviceKeeper  service.Keeper

	// the module manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager
}

// NewIrisApp returns a reference to an initialized IrisApp.
func NewIrisApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp)) *IrisApp {

	cdc := MakeCodec()

	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)

	keys := sdk.NewKVStoreKeys(
		bam.MainStoreKey, auth.StoreKey, staking.StoreKey, supply.StoreKey,
		mint.StoreKey, distr.StoreKey, slashing.StoreKey, gov.StoreKey,
		params.StoreKey, evidence.StoreKey, asset.StoreKey, guardian.StoreKey,
		htlc.StoreKey, coinswap.StoreKey, rand.StoreKey, service.StoreKey,
	)
	tKeys := sdk.NewTransientStoreKeys(staking.TStoreKey, params.TStoreKey)

	app := &IrisApp{
		BaseApp:        bApp,
		cdc:            cdc,
		invCheckPeriod: invCheckPeriod,
		keys:           keys,
		tKeys:          tKeys,
	}

	// init params keeper and subspaces
	app.paramsKeeper = params.NewKeeper(app.cdc, keys[params.StoreKey], tKeys[params.TStoreKey])
	authSubspace := app.paramsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := app.paramsKeeper.Subspace(bank.DefaultParamspace)
	stakingSubspace := app.paramsKeeper.Subspace(staking.DefaultParamspace)
	mintSubspace := app.paramsKeeper.Subspace(mint.DefaultParamspace)
	distrSubspace := app.paramsKeeper.Subspace(distr.DefaultParamspace)
	slashingSubspace := app.paramsKeeper.Subspace(slashing.DefaultParamspace)
	govSubspace := app.paramsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())
	crisisSubspace := app.paramsKeeper.Subspace(crisis.DefaultParamspace)
	evidenceSubspace := app.paramsKeeper.Subspace(evidence.DefaultParamspace)
	assetSubspace := app.paramsKeeper.Subspace(asset.DefaultParamspace)
	coinswapSubspace := app.paramsKeeper.Subspace(coinswap.DefaultParamspace)
	serviceSubspace := app.paramsKeeper.Subspace(service.DefaultParamspace)

	// add keepers
	app.accountKeeper = auth.NewAccountKeeper(app.cdc, keys[auth.StoreKey], authSubspace, auth.ProtoBaseAccount)
	app.bankKeeper = bank.NewBaseKeeper(app.accountKeeper, bankSubspace, app.ModuleAccountAddrs())
	app.supplyKeeper = supply.NewKeeper(app.cdc, keys[supply.StoreKey], app.accountKeeper, app.bankKeeper, maccPerms)
	stakingKeeper := staking.NewKeeper(
		app.cdc, keys[staking.StoreKey], app.supplyKeeper, stakingSubspace,
	)
	app.mintKeeper = mint.NewKeeper(app.cdc, keys[mint.StoreKey], mintSubspace, app.supplyKeeper, auth.FeeCollectorName)
	app.distrKeeper = distr.NewKeeper(app.cdc, keys[distr.StoreKey], distrSubspace, &stakingKeeper,
		app.supplyKeeper, auth.FeeCollectorName, app.ModuleAccountAddrs())
	app.slashingKeeper = slashing.NewKeeper(
		app.cdc, keys[slashing.StoreKey], &stakingKeeper, slashingSubspace,
	)
	app.crisisKeeper = crisis.NewKeeper(crisisSubspace, invCheckPeriod, app.supplyKeeper, auth.FeeCollectorName)

	// create evidence keeper with evidence router
	app.evidenceKeeper = evidence.NewKeeper(
		app.cdc, keys[evidence.StoreKey], evidenceSubspace, stakingKeeper, app.slashingKeeper,
	)
	evidenceRouter := evidence.NewRouter()
	// TODO: Register evidence routes.
	app.evidenceKeeper.SetRouter(evidenceRouter)

	// create guardian keeper with guardian router
	app.guardianKeeper = guardian.NewKeeper(app.cdc, keys[guardian.StoreKey])

	// register the proposal types
	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper)).
		AddRoute(distr.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.distrKeeper))
	app.govKeeper = gov.NewKeeper(
		app.cdc, keys[gov.StoreKey], govSubspace,
		app.supplyKeeper, &stakingKeeper, govRouter,
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(app.distrKeeper.Hooks(), app.slashingKeeper.Hooks()),
	)

	app.assetKeeper = asset.NewKeeper(
		app.cdc, keys[asset.StoreKey], assetSubspace,
		app.supplyKeeper, auth.FeeCollectorName,
	)

	app.htlcKeeper = htlc.NewKeeper(app.cdc, keys[htlc.StoreKey], app.supplyKeeper)

	app.coinswapKeeper = coinswap.NewKeeper(
		app.cdc, keys[coinswap.StoreKey], app.bankKeeper, app.accountKeeper, app.supplyKeeper, coinswapSubspace,
	)

	app.randKeeper = rand.NewKeeper(app.cdc, keys[rand.StoreKey])

	app.serviceKeeper = service.NewKeeper(
		app.cdc, keys[service.StoreKey], app.supplyKeeper, app.guardianKeeper, serviceSubspace,
	)

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		crisis.NewAppModule(&app.crisisKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		distr.NewAppModule(app.distrKeeper, app.accountKeeper, app.supplyKeeper, app.stakingKeeper),
		gov.NewAppModule(app.govKeeper, app.accountKeeper, app.supplyKeeper),
		mint.NewAppModule(app.mintKeeper),
		slashing.NewAppModule(app.slashingKeeper, app.accountKeeper, app.stakingKeeper),
		staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper),
		evidence.NewAppModule(*app.evidenceKeeper),
		asset.NewAppModule(app.assetKeeper),
		guardian.NewAppModule(app.guardianKeeper),
		htlc.NewAppModule(app.htlcKeeper),
		coinswap.NewAppModule(app.coinswapKeeper),
		rand.NewAppModule(app.randKeeper),
		service.NewAppModule(app.serviceKeeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	app.mm.SetOrderBeginBlockers(
		mint.ModuleName,
		distr.ModuleName,
		slashing.ModuleName,
		htlc.ModuleName,
		rand.ModuleName,
	)

	app.mm.SetOrderEndBlockers(
		crisis.ModuleName,
		gov.ModuleName,
		staking.ModuleName,
		service.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		distr.ModuleName, staking.ModuleName, auth.ModuleName, bank.ModuleName,
		slashing.ModuleName, gov.ModuleName, mint.ModuleName,
		asset.ModuleName, guardian.ModuleName, htlc.ModuleName,
		coinswap.ModuleName, rand.ModuleName, service.ModuleName,
		supply.ModuleName, crisis.ModuleName, genutil.ModuleName, evidence.ModuleName,
	)

	app.mm.RegisterInvariants(&app.crisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: This is not required for apps that don't use the simulator for fuzz testing
	// transactions.
	// TODO
	//app.sm = module.NewSimulationManager(
	//	auth.NewAppModule(app.accountKeeper),
	//	bank.NewAppModule(app.bankKeeper, app.accountKeeper),
	//	supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
	//	gov.NewAppModule(app.govKeeper, app.supplyKeeper),
	//	mint.NewAppModule(app.mintKeeper),
	//	distr.NewAppModule(app.distrKeeper, app.supplyKeeper),
	//	staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper),
	//	slashing.NewAppModule(app.slashingKeeper, app.stakingKeeper),
	//	asset.NewAppModule(app.assetKeeper),
	//	htlc.NewAppModule(app.htlcKeeper),
	//	coinswap.NewAppModule(app.coinswapKeeper),
	//	rand.NewAppModule(app.randKeeper),
	//	service.NewAppModule(app.serviceKeeper),
	//)

	//app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tKeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(NewAnteHandler(app.accountKeeper, app.supplyKeeper, app.assetKeeper, auth.DefaultSigVerificationGasConsumer))
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		if err := app.LoadLatestVersion(app.keys[bam.MainStoreKey]); err != nil {
			tmos.Exit(err.Error())
		}
	}

	return app
}

// application updates every begin block
func (app *IrisApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// application updates every end block
func (app *IrisApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// application update at chain initialization
func (app *IrisApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState
	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)

	return app.mm.InitGenesis(ctx, genesisState)
}

// load a particular height
func (app *IrisApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *IrisApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// Codec returns the application's sealed codec.
func (app *IrisApp) Codec() *codec.Codec {
	return app.cdc
}

// SimulationManager implements the SimulationApp interface
func (app *IrisApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

// GetMaccPerms returns a mapping of the application's module account permissions.
func GetMaccPerms() map[string][]string {
	modAccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		modAccPerms[k] = v
	}
	return modAccPerms
}
