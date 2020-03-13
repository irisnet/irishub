package simapp

import (
	"io"
	"os"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
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

const appName = "SimApp"

var (
	// default home directories for the application CLI
	DefaultCLIHome = os.ExpandEnv("$HOME/.simapp")

	// default home directories for the application daemon
	DefaultNodeHome = os.ExpandEnv("$HOME/.simapp")

	// The module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		supply.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(paramsclient.ProposalHandler, distr.ProposalHandler),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		evidence.AppModuleBasic{},
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

// MakeCodec custom tx codec
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	ModuleBasics.RegisterCodec(cdc)
	vesting.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

// SimApp extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type SimApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tkeys map[string]*sdk.TransientStoreKey

	// subspaces
	subspaces map[string]params.Subspace

	// keepers
	AccountKeeper  auth.AccountKeeper
	BankKeeper     bank.Keeper
	SupplyKeeper   supply.Keeper
	StakingKeeper  staking.Keeper
	SlashingKeeper slashing.Keeper
	MintKeeper     mint.Keeper
	DistrKeeper    distr.Keeper
	GovKeeper      gov.Keeper
	CrisisKeeper   crisis.Keeper
	ParamsKeeper   params.Keeper
	EvidenceKeeper evidence.Keeper
	AssetKeeper    asset.Keeper
	GuardianKeeper guardian.Keeper
	HTLCKeeper     htlc.Keeper
	CoinswapKeeper coinswap.Keeper
	RandKeeper     rand.Keeper
	ServiceKeeper  service.Keeper

	// the module manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager
}

// NewSimApp returns a reference to an initialized SimApp.
func NewSimApp(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp),
) *SimApp {

	cdc := MakeCodec()

	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)

	keys := sdk.NewKVStoreKeys(
		bam.MainStoreKey, auth.StoreKey, staking.StoreKey, supply.StoreKey, mint.StoreKey,
		distr.StoreKey, slashing.StoreKey, gov.StoreKey, params.StoreKey, evidence.StoreKey,
		asset.StoreKey, guardian.StoreKey, htlc.StoreKey, coinswap.StoreKey, rand.ModuleName,
		service.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(params.TStoreKey)

	app := &SimApp{
		BaseApp:        bApp,
		cdc:            cdc,
		invCheckPeriod: invCheckPeriod,
		keys:           keys,
		tkeys:          tkeys,
		subspaces:      make(map[string]params.Subspace),
	}

	// init params keeper and subspaces
	app.ParamsKeeper = params.NewKeeper(app.cdc, keys[params.StoreKey], tkeys[params.TStoreKey])
	app.subspaces[auth.ModuleName] = app.ParamsKeeper.Subspace(auth.DefaultParamspace)
	app.subspaces[bank.ModuleName] = app.ParamsKeeper.Subspace(bank.DefaultParamspace)
	app.subspaces[staking.ModuleName] = app.ParamsKeeper.Subspace(staking.DefaultParamspace)
	app.subspaces[mint.ModuleName] = app.ParamsKeeper.Subspace(mint.DefaultParamspace)
	app.subspaces[distr.ModuleName] = app.ParamsKeeper.Subspace(distr.DefaultParamspace)
	app.subspaces[slashing.ModuleName] = app.ParamsKeeper.Subspace(slashing.DefaultParamspace)
	app.subspaces[gov.ModuleName] = app.ParamsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())
	app.subspaces[crisis.ModuleName] = app.ParamsKeeper.Subspace(crisis.DefaultParamspace)
	app.subspaces[evidence.ModuleName] = app.ParamsKeeper.Subspace(evidence.DefaultParamspace)
	app.subspaces[asset.ModuleName] = app.ParamsKeeper.Subspace(asset.DefaultParamspace)
	app.subspaces[coinswap.ModuleName] = app.ParamsKeeper.Subspace(coinswap.DefaultParamspace)
	app.subspaces[service.ModuleName] = app.ParamsKeeper.Subspace(service.DefaultParamspace)

	// add keepers
	app.AccountKeeper = auth.NewAccountKeeper(
		app.cdc, keys[auth.StoreKey], app.subspaces[auth.ModuleName], auth.ProtoBaseAccount,
	)
	app.BankKeeper = bank.NewBaseKeeper(
		app.AccountKeeper, app.subspaces[bank.ModuleName], app.ModuleAccountAddrs(),
	)
	app.SupplyKeeper = supply.NewKeeper(
		app.cdc, keys[supply.StoreKey], app.AccountKeeper, app.BankKeeper, maccPerms,
	)
	stakingKeeper := staking.NewKeeper(
		app.cdc, keys[staking.StoreKey], app.SupplyKeeper, app.subspaces[staking.ModuleName],
	)
	app.MintKeeper = mint.NewKeeper(
		app.cdc, keys[mint.StoreKey], app.subspaces[mint.ModuleName],
		app.SupplyKeeper, auth.FeeCollectorName,
	)
	app.DistrKeeper = distr.NewKeeper(
		app.cdc, keys[distr.StoreKey], app.subspaces[distr.ModuleName], &stakingKeeper,
		app.SupplyKeeper, auth.FeeCollectorName, app.ModuleAccountAddrs(),
	)
	app.SlashingKeeper = slashing.NewKeeper(
		app.cdc, keys[slashing.StoreKey], &stakingKeeper,
		app.subspaces[slashing.ModuleName],
	)
	app.CrisisKeeper = crisis.NewKeeper(
		app.subspaces[crisis.ModuleName], invCheckPeriod, app.SupplyKeeper, auth.FeeCollectorName,
	)
	app.HTLCKeeper = htlc.NewKeeper(app.cdc, keys[htlc.StoreKey], app.SupplyKeeper)
	app.CoinswapKeeper = coinswap.NewKeeper(
		app.cdc, keys[coinswap.StoreKey], app.BankKeeper, app.AccountKeeper, app.SupplyKeeper, app.subspaces[coinswap.ModuleName],
	)
	// create evidence keeper with router
	evidenceKeeper := evidence.NewKeeper(
		app.cdc, keys[evidence.StoreKey], app.subspaces[evidence.ModuleName],
		stakingKeeper, app.SlashingKeeper,
	)

	app.GuardianKeeper = guardian.NewKeeper(app.cdc, keys[guardian.StoreKey])

	evidenceRouter := evidence.NewRouter()
	// TODO: Register evidence routes.
	evidenceKeeper.SetRouter(evidenceRouter)
	app.EvidenceKeeper = *evidenceKeeper

	// register the proposal types
	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(distr.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.DistrKeeper))
	app.GovKeeper = gov.NewKeeper(
		app.cdc, keys[gov.StoreKey], app.subspaces[gov.ModuleName], app.SupplyKeeper,
		&stakingKeeper, govRouter,
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.StakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)

	app.AssetKeeper = asset.NewKeeper(
		app.cdc, keys[asset.StoreKey], app.subspaces[asset.ModuleName],
		app.SupplyKeeper, auth.FeeCollectorName)

	app.RandKeeper = rand.NewKeeper(app.cdc, keys[rand.StoreKey])

	app.ServiceKeeper = service.NewKeeper(
		app.cdc, keys[service.StoreKey], app.SupplyKeeper, app.GuardianKeeper,
		app.subspaces[service.ModuleName],
	)

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.AccountKeeper),
		bank.NewAppModule(app.BankKeeper, app.AccountKeeper),
		crisis.NewAppModule(&app.CrisisKeeper),
		supply.NewAppModule(app.SupplyKeeper, app.AccountKeeper),
		gov.NewAppModule(app.GovKeeper, app.AccountKeeper, app.SupplyKeeper),
		mint.NewAppModule(app.MintKeeper),
		distr.NewAppModule(app.DistrKeeper, app.AccountKeeper, app.SupplyKeeper, app.StakingKeeper),
		slashing.NewAppModule(app.SlashingKeeper, app.AccountKeeper, app.StakingKeeper),
		staking.NewAppModule(app.StakingKeeper, app.AccountKeeper, app.SupplyKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		asset.NewAppModule(app.AssetKeeper),
		guardian.NewAppModule(app.GuardianKeeper),
		htlc.NewAppModule(app.HTLCKeeper),
		coinswap.NewAppModule(app.CoinswapKeeper),
		rand.NewAppModule(app.RandKeeper),
		service.NewAppModule(app.ServiceKeeper),
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

	// NOTE: The genutils moodule must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		distr.ModuleName, staking.ModuleName, auth.ModuleName, bank.ModuleName,
		slashing.ModuleName, gov.ModuleName, mint.ModuleName,
		asset.ModuleName, guardian.ModuleName, htlc.ModuleName,
		coinswap.ModuleName, rand.ModuleName, service.ModuleName,
		supply.ModuleName, crisis.ModuleName, genutil.ModuleName, evidence.ModuleName,
	)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: this is not required apps that don't use the simulator for fuzz testing
	// transactions
	// TODO:
	//app.sm = module.NewSimulationManager(
	//	auth.NewAppModule(app.AccountKeeper),
	//	bank.NewAppModule(app.BankKeeper, app.AccountKeeper),
	//	supply.NewAppModule(app.SupplyKeeper, app.AccountKeeper),
	//	gov.NewAppModule(app.GovKeeper, app.SupplyKeeper),
	//	mint.NewAppModule(app.MintKeeper),
	//	distr.NewAppModule(app.DistrKeeper, app.SupplyKeeper),
	//	staking.NewAppModule(app.StakingKeeper, app.AccountKeeper, app.SupplyKeeper),
	//	slashing.NewAppModule(app.SlashingKeeper, app.StakingKeeper),
	//	asset.NewAppModule(app.AssetKeeper),
	//	htlc.NewAppModule(app.HTLCKeeper),
	//	coinswap.NewAppModule(app.CoinswapKeeper),
	//	rand.NewAppModule(app.RandKeeper),
	//	service.NewAppModule(app.ServiceKeeper),
	//)

	//app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(ante.NewAnteHandler(app.AccountKeeper, app.SupplyKeeper, auth.DefaultSigVerificationGasConsumer))
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		if err := app.LoadLatestVersion(app.keys[bam.MainStoreKey]); err != nil {
			tmos.Exit(err.Error())
		}
	}

	return app
}

// application updates every begin block
func (app *SimApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// application updates every end block
func (app *SimApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// application update at chain initialization
func (app *SimApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)
	return app.mm.InitGenesis(ctx, genesisState)
}

// load a particular height
func (app *SimApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *SimApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// Codec returns SimApp's codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *SimApp) Codec() *codec.Codec {
	return app.cdc
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *SimApp) GetKey(storeKey string) *sdk.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *SimApp) GetTKey(storeKey string) *sdk.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *SimApp) GetSubspace(moduleName string) params.Subspace {
	return app.subspaces[moduleName]
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}
	return dupMaccPerms
}
