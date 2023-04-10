package simapp

import (
	"io"
	"os"
	"path/filepath"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
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
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/cosmos/ibc-go/v5/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v5/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v5/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v5/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v5/modules/core/02-client"
	ibcclientclient "github.com/cosmos/ibc-go/v5/modules/core/02-client/client"
	ibcclienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v5/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v5/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v5/modules/core/keeper"

	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/irisnet/irismod/modules/coinswap"
	coinswapkeeper "github.com/irisnet/irismod/modules/coinswap/keeper"
	coinswaptypes "github.com/irisnet/irismod/modules/coinswap/types"
	"github.com/irisnet/irismod/modules/htlc"
	htlckeeper "github.com/irisnet/irismod/modules/htlc/keeper"
	htlctypes "github.com/irisnet/irismod/modules/htlc/types"
	"github.com/irisnet/irismod/modules/mt"
	mtkeeper "github.com/irisnet/irismod/modules/mt/keeper"
	mttypes "github.com/irisnet/irismod/modules/mt/types"
	nftkeeper "github.com/irisnet/irismod/modules/nft/keeper"
	nftmodule "github.com/irisnet/irismod/modules/nft/module"
	nfttypes "github.com/irisnet/irismod/modules/nft/types"
	"github.com/irisnet/irismod/modules/oracle"
	oraclekeeper "github.com/irisnet/irismod/modules/oracle/keeper"
	oracletypes "github.com/irisnet/irismod/modules/oracle/types"
	"github.com/irisnet/irismod/modules/random"
	randomkeeper "github.com/irisnet/irismod/modules/random/keeper"
	randomtypes "github.com/irisnet/irismod/modules/random/types"
	"github.com/irisnet/irismod/modules/record"
	recordkeeper "github.com/irisnet/irismod/modules/record/keeper"
	recordtypes "github.com/irisnet/irismod/modules/record/types"
	"github.com/irisnet/irismod/modules/service"
	servicekeeper "github.com/irisnet/irismod/modules/service/keeper"
	servicetypes "github.com/irisnet/irismod/modules/service/types"

	tokenkeeper "github.com/irisnet/irismod/modules/token/keeper"
	tokentypes "github.com/irisnet/irismod/modules/token/types"
	tokenv1 "github.com/irisnet/irismod/modules/token/types/v1"

	"github.com/irisnet/irishub/address"
	"github.com/irisnet/irishub/lite"
	"github.com/irisnet/irishub/modules/guardian"
	guardiankeeper "github.com/irisnet/irishub/modules/guardian/keeper"
	guardiantypes "github.com/irisnet/irishub/modules/guardian/types"
	"github.com/irisnet/irishub/modules/mint"
	mintkeeper "github.com/irisnet/irishub/modules/mint/keeper"
	minttypes "github.com/irisnet/irishub/modules/mint/types"

	"github.com/irisnet/irismod/modules/farm"
	farmkeeper "github.com/irisnet/irismod/modules/farm/keeper"
	farmtypes "github.com/irisnet/irismod/modules/farm/types"
	"github.com/irisnet/irismod/modules/token"

	tibcmttransfer "github.com/bianjieai/tibc-go/modules/tibc/apps/mt_transfer"
	tibcmttransferkeeper "github.com/bianjieai/tibc-go/modules/tibc/apps/mt_transfer/keeper"
	tibcmttypes "github.com/bianjieai/tibc-go/modules/tibc/apps/mt_transfer/types"
	tibcnfttransfer "github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer"
	tibcnfttransferkeeper "github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/keeper"
	tibcnfttypes "github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"
	tibc "github.com/bianjieai/tibc-go/modules/tibc/core"
	tibchost "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	tibcroutingtypes "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
	tibccli "github.com/bianjieai/tibc-go/modules/tibc/core/client/cli"
	tibckeeper "github.com/bianjieai/tibc-go/modules/tibc/core/keeper"
)

const appName = "SimApp"

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string

	// Denominations can be 3 ~ 128 characters long and support letters, followed by either
	// a letter, a number, ('-'), or a separator ('/').
	// overwite sdk reDnmString
	reDnmString = `[a-zA-Z][a-zA-Z0-9/-]{2,127}`

	legacyProposalHandlers = []govclient.ProposalHandler{
		paramsclient.ProposalHandler,
		distrclient.ProposalHandler,
		upgradeclient.LegacyProposalHandler,
		upgradeclient.LegacyCancelProposalHandler,
		ibcclientclient.UpdateClientProposalHandler,
		ibcclientclient.UpgradeProposalHandler,
	}

	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		//authzmodule.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(
			append(legacyProposalHandlers, tibccli.GovHandlers...),
		),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		ibc.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		transfer.AppModuleBasic{},
		vesting.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},

		guardian.AppModuleBasic{},
		token.AppModuleBasic{},
		record.AppModuleBasic{},
		nftmodule.AppModuleBasic{},
		htlc.AppModuleBasic{},
		coinswap.AppModuleBasic{},
		service.AppModuleBasic{},
		oracle.AppModuleBasic{},
		random.AppModuleBasic{},
		farm.AppModuleBasic{},
		tibc.AppModuleBasic{},
		tibcnfttransfer.AppModuleBasic{},
		tibcmttransfer.AppModuleBasic{},
		mt.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		tokentypes.ModuleName:          {authtypes.Minter, authtypes.Burner},
		htlctypes.ModuleName:           {authtypes.Minter, authtypes.Burner},
		coinswaptypes.ModuleName:       {authtypes.Minter, authtypes.Burner},
		servicetypes.DepositAccName:    {authtypes.Burner},
		servicetypes.RequestAccName:    nil,
		servicetypes.FeeCollectorName:  {authtypes.Burner},
		farmtypes.ModuleName:           {authtypes.Burner},
		farmtypes.RewardCollector:      nil,
		farmtypes.EscrowCollector:      nil,
		tibcnfttypes.ModuleName:        nil,
		tibcmttypes.ModuleName:         nil,
		nfttypes.ModuleName:            nil,
	}

	nativeToken tokenv1.Token
)

var (
	_ simapp.App              = (*SimApp)(nil)
	_ servertypes.Application = (*SimApp)(nil)
)

// SimApp extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type SimApp struct {
	*baseapp.BaseApp
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// keepers
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
	IBCKeeper        *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	EvidenceKeeper   evidencekeeper.Keeper
	TransferKeeper   ibctransferkeeper.Keeper

	// make scoped keepers public for test purposes
	scopedIBCKeeper      capabilitykeeper.ScopedKeeper
	scopedTransferKeeper capabilitykeeper.ScopedKeeper
	scopedIBCMockKeeper  capabilitykeeper.ScopedKeeper
	// tibc
	scopedTIBCKeeper     capabilitykeeper.ScopedKeeper
	scopedTIBCMockKeeper capabilitykeeper.ScopedKeeper

	GuardianKeeper    guardiankeeper.Keeper
	TokenKeeper       tokenkeeper.Keeper
	RecordKeeper      recordkeeper.Keeper
	NFTKeeper         nftkeeper.Keeper
	MTKeeper          mtkeeper.Keeper
	HTLCKeeper        htlckeeper.Keeper
	CoinswapKeeper    coinswapkeeper.Keeper
	ServiceKeeper     servicekeeper.Keeper
	OracleKeeper      oraclekeeper.Keeper
	RandomKeeper      randomkeeper.Keeper
	FarmKeeper        farmkeeper.Keeper
	TIBCKeeper        *tibckeeper.Keeper
	NFTTransferKeeper tibcnfttransferkeeper.Keeper
	MTTransferKeeper  tibcmttransferkeeper.Keeper

	// the module manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager
}

func init() {
	// set bech32 prefix
	address.ConfigureBech32Prefix()

	// set coin denom regexs
	sdk.SetCoinDenomRegex(DefaultCoinDenomRegex)

	nativeToken = tokenv1.Token{
		Symbol:        "iris",
		Name:          "Irishub staking token",
		Scale:         6,
		MinUnit:       "uiris",
		InitialSupply: 2000000000,
		MaxSupply:     10000000000,
		Mintable:      true,
		Owner:         sdk.AccAddress(crypto.AddressHash([]byte(tokentypes.ModuleName))).String(),
	}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, ".iris")
	owner, err := sdk.AccAddressFromBech32(nativeToken.Owner)
	if err != nil {
		panic(err)
	}

	tokenv1.SetNativeToken(
		nativeToken.Symbol,
		nativeToken.Name,
		nativeToken.MinUnit,
		nativeToken.Scale,
		nativeToken.InitialSupply,
		nativeToken.MaxSupply,
		nativeToken.Mintable,
		owner,
	)
}

// DefaultCoinDenomRegex returns the default regex string
func DefaultCoinDenomRegex() string {
	return reDnmString
}

// NewSimApp returns a reference to an initialized IrisApp.
func NewSimApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	encodingConfig simappparams.EncodingConfig,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *SimApp {

	// TODO: Remove cdc in favor of appCodec once all modules are migrated.
	appCodec := encodingConfig.Codec
	legacyAmino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp(appName, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		minttypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, ibchost.StoreKey, upgradetypes.StoreKey,
		evidencetypes.StoreKey, ibctransfertypes.StoreKey, capabilitytypes.StoreKey,
		guardiantypes.StoreKey, tokentypes.StoreKey, nfttypes.StoreKey, htlctypes.StoreKey, recordtypes.StoreKey,
		coinswaptypes.StoreKey, servicetypes.StoreKey, oracletypes.StoreKey, randomtypes.StoreKey,
		farmtypes.StoreKey, feegrant.StoreKey, tibchost.StoreKey, tibcnfttypes.StoreKey, tibcmttypes.StoreKey, mttypes.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	// configure state listening capabilities using AppOptions
	// we are doing nothing with the returned streamingServices and waitGroup in this case
	if _, _, err := streaming.LoadStreamingServices(bApp, appOpts, appCodec, keys); err != nil {
		tmos.Exit(err.Error())
	}

	app := &SimApp{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	app.ParamsKeeper = initParamsKeeper(appCodec, legacyAmino, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])
	// set the BaseApp's parameter store
	bApp.SetParamStore(
		app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable()),
	)

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)

	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		keys[authtypes.StoreKey],
		app.GetSubspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount,
		maccPerms,
		sdk.Bech32MainPrefix,
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
		appCodec, keys[stakingtypes.StoreKey], app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName),
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

	app.NFTTransferKeeper = tibcnfttransferkeeper.NewKeeper(
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

	app.MTTransferKeeper = tibcmttransferkeeper.NewKeeper(
		appCodec,
		keys[tibcnfttypes.StoreKey],
		app.GetSubspace(tibcnfttypes.ModuleName),
		app.AccountKeeper, app.MTKeeper,
		app.TIBCKeeper.PacketKeeper,
		app.TIBCKeeper.ClientKeeper,
	)

	app.TransferKeeper = ibctransferkeeper.NewKeeper(
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
	transferModule := transfer.NewAppModule(app.TransferKeeper)
	transferIBCModule := transfer.NewIBCModule(app.TransferKeeper)

	// routerModule := router.NewAppModule(app.RouterKeeper, transferIBCModule)
	// create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferIBCModule)
	app.IBCKeeper.SetRouter(ibcRouter)

	nfttransferModule := tibcnfttransfer.NewAppModule(app.NFTTransferKeeper)
	mttransferModule := tibcmttransfer.NewAppModule(app.MTTransferKeeper)

	tibcRouter := tibcroutingtypes.NewRouter()
	tibcRouter.AddRoute(tibcnfttypes.ModuleName, nfttransferModule).
		AddRoute(tibcmttypes.ModuleName, mttransferModule)
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
	)

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

	/****  Module Options ****/
	var skipGenesisInvariants = false
	opt := appOpts.Get(crisis.FlagSkipGenesisInvariants)
	if opt, ok := opt.(bool); ok {
		skipGenesisInvariants = opt
	}

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		crisis.NewAppModule(&app.CrisisKeeper, skipGenesisInvariants),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		mint.NewAppModule(appCodec, app.MintKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		ibc.NewAppModule(app.IBCKeeper), tibc.NewAppModule(app.TIBCKeeper),
		params.NewAppModule(app.ParamsKeeper),
		transferModule,
		nfttransferModule,
		mttransferModule,
		guardian.NewAppModule(appCodec, app.GuardianKeeper),
		token.NewAppModule(appCodec, app.TokenKeeper, app.AccountKeeper, app.BankKeeper),
		record.NewAppModule(appCodec, app.RecordKeeper, app.AccountKeeper, app.BankKeeper),
		nftmodule.NewAppModule(appCodec, app.NFTKeeper, app.AccountKeeper, app.BankKeeper),
		mt.NewAppModule(appCodec, app.MTKeeper, app.AccountKeeper, app.BankKeeper),
		htlc.NewAppModule(appCodec, app.HTLCKeeper, app.AccountKeeper, app.BankKeeper),
		coinswap.NewAppModule(appCodec, app.CoinswapKeeper, app.AccountKeeper, app.BankKeeper),
		service.NewAppModule(appCodec, app.ServiceKeeper, app.AccountKeeper, app.BankKeeper),
		oracle.NewAppModule(appCodec, app.OracleKeeper, app.AccountKeeper, app.BankKeeper),
		random.NewAppModule(appCodec, app.RandomKeeper, app.AccountKeeper, app.BankKeeper),
		farm.NewAppModule(appCodec, app.FarmKeeper, app.AccountKeeper, app.BankKeeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	app.mm.SetOrderBeginBlockers(
		//sdk module
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibchost.ModuleName,
		vestingtypes.ModuleName,

		//self module
		tokentypes.ModuleName,
		tibchost.ModuleName,
		ibctransfertypes.ModuleName,
		nfttypes.ModuleName,
		htlctypes.ModuleName,
		recordtypes.ModuleName,
		coinswaptypes.ModuleName,
		servicetypes.ModuleName,
		oracletypes.ModuleName,
		randomtypes.ModuleName,
		farmtypes.ModuleName,
		mttypes.ModuleName,
		tibcnfttypes.ModuleName,
		tibcmttypes.ModuleName,
		guardiantypes.ModuleName,
	)
	app.mm.SetOrderEndBlockers(
		//sdk module
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibchost.ModuleName,
		vestingtypes.ModuleName,

		//self module
		tokentypes.ModuleName,
		tibchost.ModuleName,
		ibctransfertypes.ModuleName,
		nfttypes.ModuleName,
		htlctypes.ModuleName,
		recordtypes.ModuleName,
		coinswaptypes.ModuleName,
		servicetypes.ModuleName,
		oracletypes.ModuleName,
		randomtypes.ModuleName,
		farmtypes.ModuleName,
		mttypes.ModuleName,
		tibcnfttypes.ModuleName,
		tibcmttypes.ModuleName,
		guardiantypes.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.mm.SetOrderInitGenesis(
		//sdk module
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibchost.ModuleName,
		vestingtypes.ModuleName,

		//self module
		tokentypes.ModuleName,
		tibchost.ModuleName,
		ibctransfertypes.ModuleName,
		nfttypes.ModuleName,
		htlctypes.ModuleName,
		recordtypes.ModuleName,
		coinswaptypes.ModuleName,
		servicetypes.ModuleName,
		oracletypes.ModuleName,
		randomtypes.ModuleName,
		farmtypes.ModuleName,
		mttypes.ModuleName,
		tibcnfttypes.ModuleName,
		tibcmttypes.ModuleName,
		guardiantypes.ModuleName,
	)

	cfg := module.NewConfigurator(appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)
	app.mm.RegisterServices(cfg)

	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: this is not required apps that don't use the simulator for fuzz testing
	// transactions
	app.sm = module.NewSimulationManager(
		auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		mint.NewAppModule(appCodec, app.MintKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		params.NewAppModule(app.ParamsKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		transferModule,
		nfttransferModule,
		mttransferModule,
		guardian.NewAppModule(appCodec, app.GuardianKeeper),
		token.NewAppModule(appCodec, app.TokenKeeper, app.AccountKeeper, app.BankKeeper),
		record.NewAppModule(appCodec, app.RecordKeeper, app.AccountKeeper, app.BankKeeper),
		nftmodule.NewAppModule(appCodec, app.NFTKeeper, app.AccountKeeper, app.BankKeeper),
		htlc.NewAppModule(appCodec, app.HTLCKeeper, app.AccountKeeper, app.BankKeeper),
		coinswap.NewAppModule(appCodec, app.CoinswapKeeper, app.AccountKeeper, app.BankKeeper),
		service.NewAppModule(appCodec, app.ServiceKeeper, app.AccountKeeper, app.BankKeeper),
		oracle.NewAppModule(appCodec, app.OracleKeeper, app.AccountKeeper, app.BankKeeper),
		random.NewAppModule(appCodec, app.RandomKeeper, app.AccountKeeper, app.BankKeeper),
		farm.NewAppModule(appCodec, app.FarmKeeper, app.AccountKeeper, app.BankKeeper),
		tibc.NewAppModule(app.TIBCKeeper),
	)

	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	//app.SetAnteHandler(anteHandler)
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

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

// MakeCodecs constructs the *std.Codec and *codec.LegacyAmino instances used by
// irisapp. It is useful for tests and clients who do not want to construct the
// full irisapp
func MakeCodecs() (codec.Codec, *codec.LegacyAmino) {
	config := MakeTestEncodingConfig()
	return config.Codec, config.Amino
}

// Name returns the name of the App
func (app *SimApp) Name() string { return app.BaseApp.Name() }

// BeginBlocker application updates every begin block
func (app *SimApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *SimApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer application update at chain initialization
func (app *SimApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}

	// add system service at InitChainer, overwrite if it exists
	var serviceGenState servicetypes.GenesisState
	app.appCodec.MustUnmarshalJSON(genesisState[servicetypes.ModuleName], &serviceGenState)
	serviceGenState.Definitions = append(serviceGenState.Definitions, servicetypes.GenOraclePriceSvcDefinition())
	serviceGenState.Bindings = append(serviceGenState.Bindings, servicetypes.GenOraclePriceSvcBinding(nativeToken.MinUnit))
	serviceGenState.Definitions = append(serviceGenState.Definitions, randomtypes.GetSvcDefinition())
	genesisState[servicetypes.ModuleName] = app.appCodec.MustMarshalJSON(&serviceGenState)

	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// LoadHeight loads a particular height
func (app *SimApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *SimApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// BlockedModuleAccountAddrs returns all the app's blocked module account
// addresses.
func (app *SimApp) BlockedModuleAccountAddrs() map[string]bool {
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
func (app *SimApp) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

// AppCodec returns IrisApp's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *SimApp) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns IrisApp's InterfaceRegistry
func (app *SimApp) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *SimApp) GetKey(storeKey string) *storetypes.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *SimApp) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *SimApp) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *SimApp) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// SimulationManager implements the SimulationApp interface
func (app *SimApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

// RegisterAPIRoutes registers all application module routes with the provided API server.
func (app *SimApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
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
func (app *SimApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *SimApp) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(
		clientCtx,
		app.BaseApp.GRPCQueryRouter(),
		app.interfaceRegistry,
		app.Query,
	)
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}
	return dupMaccPerms
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
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

	return paramsKeeper
}
