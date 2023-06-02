package app

import (
	"io"

	"github.com/irisnet/irishub/modules/bridgenft"
	"github.com/irisnet/irishub/modules/internft"

	convertertypes "github.com/irisnet/erc721-bridge/x/converter/types"
	"github.com/spf13/cast"
	abci "github.com/tendermint/tendermint/abci/types"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/store/streaming"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/cosmos/ibc-go/v5/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v5/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v5/modules/apps/transfer/types"
	ibcclient "github.com/cosmos/ibc-go/v5/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v5/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v5/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v5/modules/core/keeper"

	nfttransfer "github.com/bianjieai/nft-transfer"
	ibcnfttransferkeeper "github.com/bianjieai/nft-transfer/keeper"
	ibcnfttransfertypes "github.com/bianjieai/nft-transfer/types"

	coinswapkeeper "github.com/irisnet/irismod/modules/coinswap/keeper"
	coinswaptypes "github.com/irisnet/irismod/modules/coinswap/types"
	"github.com/irisnet/irismod/modules/farm"
	farmkeeper "github.com/irisnet/irismod/modules/farm/keeper"
	farmtypes "github.com/irisnet/irismod/modules/farm/types"
	htlckeeper "github.com/irisnet/irismod/modules/htlc/keeper"
	htlctypes "github.com/irisnet/irismod/modules/htlc/types"
	mtkeeper "github.com/irisnet/irismod/modules/mt/keeper"
	mttypes "github.com/irisnet/irismod/modules/mt/types"
	nftkeeper "github.com/irisnet/irismod/modules/nft/keeper"
	nfttypes "github.com/irisnet/irismod/modules/nft/types"
	oraclekeeper "github.com/irisnet/irismod/modules/oracle/keeper"
	oracletypes "github.com/irisnet/irismod/modules/oracle/types"
	randomkeeper "github.com/irisnet/irismod/modules/random/keeper"
	randomtypes "github.com/irisnet/irismod/modules/random/types"
	recordkeeper "github.com/irisnet/irismod/modules/record/keeper"
	recordtypes "github.com/irisnet/irismod/modules/record/types"
	servicekeeper "github.com/irisnet/irismod/modules/service/keeper"
	servicetypes "github.com/irisnet/irismod/modules/service/types"
	tokenkeeper "github.com/irisnet/irismod/modules/token/keeper"
	tokentypes "github.com/irisnet/irismod/modules/token/types"
	tokenv1 "github.com/irisnet/irismod/modules/token/types/v1"

	tibcmttransfer "github.com/bianjieai/tibc-go/modules/tibc/apps/mt_transfer"
	tibcmttransferkeeper "github.com/bianjieai/tibc-go/modules/tibc/apps/mt_transfer/keeper"
	tibcmttypes "github.com/bianjieai/tibc-go/modules/tibc/apps/mt_transfer/types"
	tibcnfttransfer "github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer"
	tibcnfttransferkeeper "github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/keeper"
	tibcnfttypes "github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"
	tibchost "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	tibcroutingtypes "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
	tibccli "github.com/bianjieai/tibc-go/modules/tibc/core/client/cli"
	tibckeeper "github.com/bianjieai/tibc-go/modules/tibc/core/keeper"

	"github.com/evmos/ethermint/ethereum/eip712"
	srvflags "github.com/evmos/ethermint/server/flags"
	ethermint "github.com/evmos/ethermint/types"
	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/evmos/ethermint/x/evm/vm/geth"
	feemarketkeeper "github.com/evmos/ethermint/x/feemarket/keeper"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	converterkeeper "github.com/irisnet/erc721-bridge/x/converter/keeper"
	bridgenfttransfer "github.com/irisnet/erc721-bridge/x/nft-transfer"
	bridgenfttransferkeeper "github.com/irisnet/erc721-bridge/x/nft-transfer/keeper"

	"github.com/irisnet/irishub/address"
	irishubante "github.com/irisnet/irishub/ante"
	irisappparams "github.com/irisnet/irishub/app/params"
	"github.com/irisnet/irishub/lite"
	guardiankeeper "github.com/irisnet/irishub/modules/guardian/keeper"
	guardiantypes "github.com/irisnet/irishub/modules/guardian/types"
	mintkeeper "github.com/irisnet/irishub/modules/mint/keeper"
	minttypes "github.com/irisnet/irishub/modules/mint/types"
	iristypes "github.com/irisnet/irishub/types"
)

var (
	_ simapp.App              = (*IrisApp)(nil)
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

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// cosmos
	FeeGrantKeeper   feegrantkeeper.Keeper
	AccountKeeper    authkeeper.AccountKeeper
	BankKeeper       bankkeeper.Keeper
	CapabilityKeeper *capabilitykeeper.Keeper
	StakingKeeper    stakingkeeper.Keeper
	SlashingKeeper   slashingkeeper.Keeper
	MintKeeper       mintkeeper.Keeper
	DistrKeeper      distrkeeper.Keeper
	GovKeeper        govkeeper.Keeper
	CrisisKeeper     crisiskeeper.Keeper
	UpgradeKeeper    upgradekeeper.Keeper
	ParamsKeeper     paramskeeper.Keeper
	EvidenceKeeper   evidencekeeper.Keeper
	AuthzKeeper      authzkeeper.Keeper

	//ibc
	IBCKeeper            *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	IBCTransferKeeper    ibctransferkeeper.Keeper
	IBCNFTTransferKeeper ibcnfttransferkeeper.Keeper

	// make scoped keepers public for test purposes
	scopedIBCKeeper         capabilitykeeper.ScopedKeeper
	scopedTransferKeeper    capabilitykeeper.ScopedKeeper
	scopedIBCMockKeeper     capabilitykeeper.ScopedKeeper
	scopedNFTTransferKeeper capabilitykeeper.ScopedKeeper
	// tibc
	scopedTIBCKeeper     capabilitykeeper.ScopedKeeper
	scopedTIBCMockKeeper capabilitykeeper.ScopedKeeper

	GuardianKeeper        guardiankeeper.Keeper
	TokenKeeper           tokenkeeper.Keeper
	RecordKeeper          recordkeeper.Keeper
	NFTKeeper             nftkeeper.Keeper
	MTKeeper              mtkeeper.Keeper
	HTLCKeeper            htlckeeper.Keeper
	CoinswapKeeper        coinswapkeeper.Keeper
	ServiceKeeper         servicekeeper.Keeper
	OracleKeeper          oraclekeeper.Keeper
	RandomKeeper          randomkeeper.Keeper
	FarmKeeper            farmkeeper.Keeper
	TIBCKeeper            *tibckeeper.Keeper
	TIBCNFTTransferKeeper tibcnfttransferkeeper.Keeper
	TIBCMTTransferKeeper  tibcmttransferkeeper.Keeper

	// Ethermint keepers
	EvmKeeper       *evmkeeper.Keeper
	FeeMarketKeeper feemarketkeeper.Keeper

	// erc721-bridge
	Erc721ConvertKeeper converterkeeper.Keeper

	// the module manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager

	transferModule       transfer.AppModule
	nfttransferModule    tibcnfttransfer.AppModule
	mttransferModule     tibcmttransfer.AppModule
	ibcnfttransferModule bridgenfttransfer.AppModule
}

// NewIrisApp returns a reference to an initialized IrisApp.
func NewIrisApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	encodingConfig irisappparams.EncodingConfig,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *IrisApp {
	appCodec := encodingConfig.Marshaler
	legacyAmino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	eip712.SetConfig(eip712.Config{
		InterfaceRegistry: interfaceRegistry,
		Amino:             legacyAmino,
		ChainIDBuilder:    iristypes.BuildEthChainID,
	})

	bApp := baseapp.NewBaseApp(
		iristypes.AppName,
		logger,
		db,
		encodingConfig.TxConfig.TxDecoder(),
		baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey,
		banktypes.StoreKey,
		stakingtypes.StoreKey,
		minttypes.StoreKey,
		distrtypes.StoreKey,
		slashingtypes.StoreKey,
		govtypes.StoreKey,
		paramstypes.StoreKey,
		ibchost.StoreKey,
		upgradetypes.StoreKey,
		evidencetypes.StoreKey,
		ibctransfertypes.StoreKey,
		capabilitytypes.StoreKey,
		guardiantypes.StoreKey,
		tokentypes.StoreKey,
		nfttypes.StoreKey,
		htlctypes.StoreKey,
		recordtypes.StoreKey,
		coinswaptypes.StoreKey,
		servicetypes.StoreKey,
		oracletypes.StoreKey,
		randomtypes.StoreKey,
		farmtypes.StoreKey,
		feegrant.StoreKey,
		tibchost.StoreKey,
		tibcnfttypes.StoreKey,
		tibcmttypes.StoreKey,
		mttypes.StoreKey,
		authzkeeper.StoreKey,
		// ethermint keys
		evmtypes.StoreKey,
		feemarkettypes.StoreKey,
		ibcnfttransfertypes.StoreKey,
		// erc721-bridge
		convertertypes.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(
		paramstypes.TStoreKey,
		evmtypes.TransientKey,
		feemarkettypes.TransientKey,
	)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	// configure state listening capabilities using AppOptions
	// we are doing nothing with the returned streamingServices and waitGroup in this case
	if _, _, err := streaming.LoadStreamingServices(bApp, appOpts, appCodec, keys); err != nil {
		tmos.Exit(err.Error())
	}

	app := &IrisApp{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	app.ParamsKeeper = initParamsKeeper(
		appCodec,
		legacyAmino,
		keys[paramstypes.StoreKey],
		tkeys[paramstypes.TStoreKey],
	)
	// set the BaseApp's parameter store
	bApp.SetParamStore(
		app.ParamsKeeper.Subspace(baseapp.Paramspace).
			WithKeyTable(paramstypes.ConsensusParamsKeyTable()),
	)

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(
		appCodec,
		keys[capabilitytypes.StoreKey],
		memKeys[capabilitytypes.MemStoreKey],
	)
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	scopedNFTTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibcnfttransfertypes.ModuleName)

	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		keys[authtypes.StoreKey],
		app.GetSubspace(authtypes.ModuleName),
		ethermint.ProtoAccount,
		maccPerms,
		address.Bech32PrefixAccAddr,
	)

	app.FeeGrantKeeper = feegrantkeeper.NewKeeper(
		appCodec,
		keys[feegrant.StoreKey],
		app.AccountKeeper,
	)

	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		keys[banktypes.StoreKey],
		app.AccountKeeper,
		app.GetSubspace(banktypes.ModuleName),
		app.BlockedModuleAccountAddrs(),
	)

	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec,
		keys[stakingtypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		app.GetSubspace(stakingtypes.ModuleName),
	)

	app.MintKeeper = mintkeeper.NewKeeper(
		appCodec,
		keys[minttypes.StoreKey],
		app.GetSubspace(minttypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.FeeCollectorName,
	)

	app.DistrKeeper = distrkeeper.NewKeeper(
		appCodec,
		keys[distrtypes.StoreKey],
		app.GetSubspace(distrtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		&stakingKeeper,
		authtypes.FeeCollectorName,
	)

	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		keys[slashingtypes.StoreKey],
		&stakingKeeper,
		app.GetSubspace(slashingtypes.ModuleName),
	)

	app.CrisisKeeper = crisiskeeper.NewKeeper(
		app.GetSubspace(crisistypes.ModuleName),
		invCheckPeriod,
		app.BankKeeper,
		authtypes.FeeCollectorName,
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.StakingKeeper = *stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)

	// set the governance module account as the authority for conducting upgrades
	// UpgradeKeeper must be created before IBCKeeper
	app.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		keys[upgradetypes.StoreKey],
		appCodec,
		homePath,
		app.BaseApp,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.AuthzKeeper = authzkeeper.NewKeeper(
		keys[authzkeeper.StoreKey],
		appCodec,
		app.MsgServiceRouter(),
		app.AccountKeeper,
	)

	scopedTIBCKeeper := app.CapabilityKeeper.ScopeToModule(tibchost.ModuleName)
	// UpgradeKeeper must be created before IBCKeeper
	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		keys[ibchost.StoreKey],
		app.GetSubspace(ibchost.ModuleName),
		app.StakingKeeper,
		app.UpgradeKeeper,
		scopedIBCKeeper,
	)

	// register the proposal types
	app.TIBCKeeper = tibckeeper.NewKeeper(
		appCodec,
		keys[tibchost.StoreKey],
		app.GetSubspace(tibchost.ModuleName), app.StakingKeeper,
	)

	app.NFTKeeper = nftkeeper.NewKeeper(
		appCodec,
		keys[nfttypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
	)

	app.TIBCNFTTransferKeeper = tibcnfttransferkeeper.NewKeeper(
		appCodec,
		keys[tibcnfttypes.StoreKey],
		app.GetSubspace(tibcnfttypes.ModuleName),
		app.AccountKeeper,
		nftkeeper.NewLegacyKeeper(app.NFTKeeper),
		app.TIBCKeeper.PacketKeeper,
		app.TIBCKeeper.ClientKeeper,
	)

	app.MTKeeper = mtkeeper.NewKeeper(
		appCodec, keys[mttypes.StoreKey],
	)

	app.TIBCMTTransferKeeper = tibcmttransferkeeper.NewKeeper(
		appCodec,
		keys[tibcnfttypes.StoreKey],
		app.GetSubspace(tibcnfttypes.ModuleName),
		app.AccountKeeper, app.MTKeeper,
		app.TIBCKeeper.PacketKeeper,
		app.TIBCKeeper.ClientKeeper,
	)

	app.IBCTransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		keys[ibctransfertypes.StoreKey],
		app.GetSubspace(ibctransfertypes.ModuleName),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		scopedTransferKeeper,
	)
	app.transferModule = transfer.NewAppModule(app.IBCTransferKeeper)
	transferIBCModule := transfer.NewIBCModule(app.IBCTransferKeeper)

	app.nfttransferModule = tibcnfttransfer.NewAppModule(app.TIBCNFTTransferKeeper)
	app.mttransferModule = tibcmttransfer.NewAppModule(app.TIBCMTTransferKeeper)

	tibcRouter := tibcroutingtypes.NewRouter()
	tibcRouter.AddRoute(tibcnfttypes.ModuleName, app.nfttransferModule).
		AddRoute(tibcmttypes.ModuleName, app.mttransferModule)
	app.TIBCKeeper.SetRouter(tibcRouter)

	// create evidence keeper with router
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec,
		keys[evidencetypes.StoreKey],
		&app.StakingKeeper,
		app.SlashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.EvidenceKeeper = *evidenceKeeper

	app.GuardianKeeper = guardiankeeper.NewKeeper(
		appCodec,
		keys[guardiantypes.StoreKey],
	)

	app.TokenKeeper = tokenkeeper.NewKeeper(
		appCodec,
		keys[tokentypes.StoreKey],
		app.GetSubspace(tokentypes.ModuleName),
		app.BankKeeper,
		app.ModuleAccountAddrs(),
		authtypes.FeeCollectorName,
	).WithSwapRegistry(tokenv1.SwapRegistry{
		iristypes.NativeToken.MinUnit: tokenv1.SwapParams{
			MinUnit: iristypes.EvmToken.MinUnit,
			Ratio:   sdk.OneDec(),
		},
		iristypes.EvmToken.MinUnit: tokenv1.SwapParams{
			MinUnit: iristypes.NativeToken.MinUnit,
			Ratio:   sdk.OneDec(),
		},
	})

	app.RecordKeeper = recordkeeper.NewKeeper(
		appCodec,
		keys[recordtypes.StoreKey],
	)

	app.HTLCKeeper = htlckeeper.NewKeeper(
		appCodec, keys[htlctypes.StoreKey],
		app.GetSubspace(htlctypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.ModuleAccountAddrs(),
	)

	app.CoinswapKeeper = coinswapkeeper.NewKeeper(
		appCodec,
		keys[coinswaptypes.StoreKey],
		app.GetSubspace(coinswaptypes.ModuleName),
		app.BankKeeper,
		app.AccountKeeper,
		app.ModuleAccountAddrs(),
		authtypes.FeeCollectorName,
	)

	app.ServiceKeeper = servicekeeper.NewKeeper(
		appCodec,
		keys[servicetypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		app.GetSubspace(servicetypes.ModuleName),
		app.ModuleAccountAddrs(),
		servicetypes.FeeCollectorName,
	)

	app.OracleKeeper = oraclekeeper.NewKeeper(
		appCodec,
		keys[oracletypes.StoreKey],
		app.GetSubspace(oracletypes.ModuleName),
		app.ServiceKeeper,
	)

	app.RandomKeeper = randomkeeper.NewKeeper(
		appCodec,
		keys[randomtypes.StoreKey],
		app.BankKeeper,
		app.ServiceKeeper,
	)

	app.FarmKeeper = farmkeeper.NewKeeper(appCodec,
		keys[farmtypes.StoreKey],
		app.BankKeeper,
		app.AccountKeeper,
		app.DistrKeeper,
		&app.GovKeeper,
		app.CoinswapKeeper.ValidatePool,
		app.GetSubspace(farmtypes.ModuleName),
		authtypes.FeeCollectorName,
		distrtypes.ModuleName,
	)

	govRouter := govv1beta1.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(distrtypes.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.DistrKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(app.IBCKeeper.ClientKeeper)).
		AddRoute(tibchost.RouterKey, tibccli.NewProposalHandler(app.TIBCKeeper)).
		AddRoute(farmtypes.RouterKey, farm.NewCommunityPoolCreateFarmProposalHandler(app.FarmKeeper))

	govConfig := govtypes.DefaultConfig()

	app.GovKeeper = govkeeper.NewKeeper(
		appCodec,
		keys[govtypes.StoreKey],
		app.GetSubspace(govtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		govRouter,
		app.MsgServiceRouter(),
		govConfig,
	)

	govHooks := govtypes.NewMultiGovHooks(farmkeeper.NewGovHook(app.FarmKeeper))
	app.GovKeeper.SetHooks(govHooks)

	tracer := cast.ToString(appOpts.Get(srvflags.EVMTracer))

	// Create Ethermint keepers
	app.FeeMarketKeeper = feemarketkeeper.NewKeeper(
		appCodec,
		app.GetSubspace(feemarkettypes.ModuleName),
		keys[feemarkettypes.StoreKey],
		tkeys[feemarkettypes.TransientKey],
	)

	app.EvmKeeper = evmkeeper.NewKeeper(
		appCodec,
		keys[evmtypes.StoreKey],
		tkeys[evmtypes.TransientKey],
		app.GetSubspace(evmtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		&stakingKeeper,
		app.FeeMarketKeeper,
		nil,
		geth.NewEVM,
		tracer,
	)

	// Create Convert Keepers
	app.Erc721ConvertKeeper = converterkeeper.NewKeeper(
		appCodec,
		keys[convertertypes.StoreKey],
		app.AccountKeeper,
		app.EvmKeeper,
		bridgenft.NewBridgeNftKeeper(appCodec, app.NFTKeeper, app.AccountKeeper),
	)

	erc721Keeper := app.Erc721ConvertKeeper.ERC721Keeper()
	internftKeeper := internft.NewInterNftKeeper(appCodec, app.NFTKeeper, app.AccountKeeper)

	nftRouter := ibcnfttransfertypes.NewRouter()
	nftRouter.AddRoute(ibcnfttransfertypes.NativePortID, internftKeeper).
		AddRoute(ibcnfttransfertypes.ERC721PortID, erc721Keeper)
	nftRouter.Seal()

	app.IBCNFTTransferKeeper = ibcnfttransferkeeper.NewKeeper(
		appCodec,
		keys[ibcnfttransfertypes.StoreKey],
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		nftRouter,
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper, scopedNFTTransferKeeper,
	)

	bridenfttransferKeeper := bridgenfttransferkeeper.NewKeeper(
		app.IBCNFTTransferKeeper,
		erc721Keeper,
	)
	app.ibcnfttransferModule = bridgenfttransfer.NewAppModule(
		nfttransfer.NewAppModule(app.IBCNFTTransferKeeper),
		bridenfttransferKeeper,
	)
	nfttransferIBCModule := bridgenfttransfer.NewIBCModule(
		nfttransfer.NewIBCModule(app.IBCNFTTransferKeeper),
		bridenfttransferKeeper,
	)

	// routerModule := router.NewAppModule(app.RouterKeeper, transferIBCModule)
	// create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferIBCModule).
		AddRoute(ibcnfttransfertypes.ModuleName, nfttransferIBCModule)
	app.IBCKeeper.SetRouter(ibcRouter)

	/****  Module Options ****/
	var skipGenesisInvariants = false
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

	app.configurator = module.NewConfigurator(
		appCodec,
		app.MsgServiceRouter(),
		app.GRPCQueryRouter(),
	)
	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)
	app.mm.RegisterServices(app.configurator)

	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: this is not required apps that don't use the simulator for fuzz testing
	// transactions
	app.sm = module.NewSimulationManager(
		simulationModules(app, encodingConfig, skipGenesisInvariants)...)

	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

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
	app.scopedTIBCKeeper = scopedTIBCKeeper
	app.scopedIBCKeeper = scopedIBCKeeper
	app.scopedTransferKeeper = scopedTransferKeeper
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
	//
	// TODO: Blocked on updating to v0.46.x
	// delete(modAccAddrs, authtypes.NewModuleAddress(grouptypes.ModuleName).String())
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

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *IrisApp) GetKey(storeKey string) *storetypes.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *IrisApp) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *IrisApp) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *IrisApp) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
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

	// Register grpc-gateway routes for all modules.
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if apiConfig.Swagger {
		lite.RegisterSwaggerAPI(clientCtx, apiSvr.Router)
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

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(
	appCodec codec.BinaryCodec,
	legacyAmino *codec.LegacyAmino,
	key, tkey storetypes.StoreKey,
) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govv1.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(tokentypes.ModuleName)
	paramsKeeper.Subspace(recordtypes.ModuleName)
	paramsKeeper.Subspace(htlctypes.ModuleName)
	paramsKeeper.Subspace(coinswaptypes.ModuleName)
	paramsKeeper.Subspace(servicetypes.ModuleName)
	paramsKeeper.Subspace(ibchost.ModuleName)
	paramsKeeper.Subspace(farmtypes.ModuleName)
	paramsKeeper.Subspace(tibchost.ModuleName)

	// ethermint subspaces
	paramsKeeper.Subspace(evmtypes.ModuleName)
	paramsKeeper.Subspace(feemarkettypes.ModuleName)

	return paramsKeeper
}
