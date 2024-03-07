package app

import (
	"encoding/json"
	"io"

	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	tmjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/libs/log"
	tmos "github.com/cometbft/cometbft/libs/os"
	"github.com/spf13/cast"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/mempool"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	srvflags "github.com/evmos/ethermint/server/flags"

	irishubante "github.com/irisnet/irishub/v2/app/ante"
	"github.com/irisnet/irishub/v2/app/keepers"
	"github.com/irisnet/irishub/v2/app/params"
	"github.com/irisnet/irishub/v2/app/rpc"
	"github.com/irisnet/irishub/v2/client/lite"
	iristypes "github.com/irisnet/irishub/v2/types"
)

var (
	_ runtime.AppI            = (*IrisApp)(nil)
	_ servertypes.Application = (*IrisApp)(nil)
)

// IrisApp extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type IrisApp struct {
	*baseapp.BaseApp
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry
	configurator      module.Configurator

	keepers.AppKeepers

	// the module manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager
}

// NewIrisApp returns a reference to an initialized IrisApp.
func NewIrisApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	encodingConfig params.EncodingConfig,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *IrisApp {
	appCodec := encodingConfig.Marshaler
	legacyAmino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	// Setup Mempool
	baseAppOptions = append(baseAppOptions, NoOpMempoolOption())

	bApp := baseapp.NewBaseApp(
		iristypes.AppName,
		logger,
		db,
		encodingConfig.TxConfig.TxDecoder(),
		baseAppOptions...,
	)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	app := &IrisApp{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
	}

	// get skipUpgradeHeights from the app options
	skipUpgradeHeights := map[int64]bool{}
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	// Setup keepers
	app.AppKeepers = keepers.New(
		appCodec,
		bApp,
		legacyAmino,
		maccPerms,
		app.ModuleAccountAddrs(),
		app.BlockedModuleAccountAddrs(),
		skipUpgradeHeights,
		cast.ToString(appOpts.Get(flags.FlagHome)),
		cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)),
		logger,
		appOpts,
	)

	/****  Module Options ****/
	skipGenesisInvariants := false
	opt := appOpts.Get(crisis.FlagSkipGenesisInvariants)
	if opt, ok := opt.(bool); ok {
		skipGenesisInvariants = opt
	}

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(appModules(app, encodingConfig, skipGenesisInvariants)...)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	app.mm.SetOrderBeginBlockers(orderBeginBlockers()...)
	app.mm.SetOrderEndBlockers(orderEndBlockers()...)
	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: The genutils module must also occur after auth so that it can access the params from auth.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.mm.SetOrderInitGenesis(orderInitBlockers()...)
	app.mm.SetOrderExportGenesis(orderInitBlockers()...)

	app.configurator = module.NewConfigurator(
		appCodec,
		app.MsgServiceRouter(),
		app.GRPCQueryRouter(),
	)
	app.mm.RegisterInvariants(app.CrisisKeeper)
	// app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)
	// app.mm.RegisterServices(app.configurator)
	app.RegisterServices()

	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: this is not required apps that don't use the simulator for fuzz testing
	// transactions
	app.sm = module.NewSimulationManager(
		simulationModules(app, encodingConfig, skipGenesisInvariants)...)

	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(app.KvStoreKeys())
	app.MountTransientStores(app.TransientStoreKeys())
	app.MountMemoryStores(app.MemoryStoreKeys())

	maxGasWanted := cast.ToUint64(appOpts.Get(srvflags.EVMMaxTxGasWanted))
	anteHandler := irishubante.NewAnteHandler(
		irishubante.HandlerOptions{
			HandlerOptions: ante.HandlerOptions{
				AccountKeeper:   app.AccountKeeper,
				BankKeeper:      app.BankKeeper,
				FeegrantKeeper:  app.FeeGrantKeeper,
				SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
			},
			AccountKeeper:        app.AccountKeeper,
			BankKeeper:           app.BankKeeper,
			IBCKeeper:            app.IBCKeeper,
			TokenKeeper:          app.TokenKeeper,
			OracleKeeper:         app.OracleKeeper,
			GuardianKeeper:       app.GuardianKeeper,
			EvmKeeper:            app.EvmKeeper,
			FeeMarketKeeper:      app.FeeMarketKeeper,
			BypassMinFeeMsgTypes: []string{},
			MaxTxGasWanted:       maxGasWanted,
		},
	)

	app.Init()
	app.SetAnteHandler(anteHandler)
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	// Set software upgrade execution logic
	app.RegisterUpgradePlans()

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}

		// Initialize and seal the capability keeper so all persistent capabilities
		// are loaded in-memory and prevent any further modules from creating scoped
		// sub-keepers.
		// This must be done during creation of baseapp rather than in InitChain so
		// that in-memory capabilities get regenerated on app restart.
		// Note that since this reads from the store, we can only perform it when
		// `loadLatest` is set to true.
		app.CapabilityKeeper.Seal()
	}
	return app
}

// Name returns the name of the App
func (app *IrisApp) Name() string { return app.BaseApp.Name() }

// BeginBlocker application updates every begin block
func (app *IrisApp) BeginBlocker(
	ctx sdk.Context,
	req abci.RequestBeginBlock,
) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *IrisApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer application update at chain initialization
func (app *IrisApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState iristypes.GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// LoadHeight loads a particular height
func (app *IrisApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *IrisApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// BlockedModuleAccountAddrs returns all the app's blocked module account
// addresses.
func (app *IrisApp) BlockedModuleAccountAddrs() map[string]bool {
	modAccAddrs := app.ModuleAccountAddrs()

	// remove module accounts that are ALLOWED to received funds
	delete(modAccAddrs, authtypes.NewModuleAddress(govtypes.ModuleName).String())

	return modAccAddrs
}

// LegacyAmino returns SimApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *IrisApp) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

// AppCodec returns IrisApp's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *IrisApp) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns IrisApp's InterfaceRegistry
func (app *IrisApp) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// SimulationManager implements the SimulationApp interface
func (app *IrisApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

// RegisterAPIRoutes registers all application module routes with the provided API server.
func (app *IrisApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register node gRPC service for grpc-gateway.
	nodeservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register grpc-gateway routes for all modules.
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if apiConfig.Swagger {
		lite.RegisterSwaggerAPI(clientCtx, apiSvr.Router)
	}
}

// RegisterServices implements the Application.RegisterTxService method.
func (app *IrisApp) RegisterServices() {
	for _, mod := range app.mm.Modules {
		m, ok := mod.(module.AppModule)
		if !ok {
			panic("unable to cast mod into AppModule")
		}

		if m.Name() == authtypes.ModuleName {
			rpc.OverrideAuthServices(app.configurator,
				app.GetKey(authtypes.StoreKey),
				app.AccountKeeper,
				app.GetSubspace(authtypes.ModuleName))
		} else {
			if mod, ok := mod.(module.HasServices); ok {
				mod.RegisterServices(app.configurator)
			}
		}
	}
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *IrisApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(
		app.BaseApp.GRPCQueryRouter(),
		clientCtx,
		app.BaseApp.Simulate,
		app.interfaceRegistry,
	)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *IrisApp) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(
		clientCtx,
		app.BaseApp.GRPCQueryRouter(),
		app.interfaceRegistry,
		app.Query,
	)
}

// RegisterNodeService registers the node service.
//
// It takes a client context as a parameter and does not return anything.
func (app *IrisApp) RegisterNodeService(clientCtx client.Context) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter())
}

// DefaultGenesis returns a default genesis from the registered AppModuleBasic's.
func (app *IrisApp) DefaultGenesis() map[string]json.RawMessage {
	return ModuleBasics.DefaultGenesis(app.AppCodec())
}

// Init initializes the IrisApp.
func (app *IrisApp) Init() {
	iristypes.Init(app.legacyAmino, app.interfaceRegistry)
}

// NoOpMempoolOption returns a function that sets up a no-op mempool for the given BaseApp.
//
// The function takes a pointer to a BaseApp as a parameter and returns nothing.
func NoOpMempoolOption() func(*baseapp.BaseApp) {
	return func(app *baseapp.BaseApp) {
		memPool := mempool.NoOpMempool{}
		app.SetMempool(memPool)
		handler := baseapp.NewDefaultProposalHandler(memPool, app)
		app.SetPrepareProposal(handler.PrepareProposalHandler())
		app.SetProcessProposal(handler.ProcessProposalHandler())
	}
}
