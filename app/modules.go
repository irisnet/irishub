package app

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcclientclient "github.com/cosmos/ibc-go/v7/modules/core/02-client/client"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"

	"github.com/irisnet/irismod/modules/coinswap"
	coinswaptypes "github.com/irisnet/irismod/modules/coinswap/types"
	"github.com/irisnet/irismod/modules/farm"
	farmtypes "github.com/irisnet/irismod/modules/farm/types"
	"github.com/irisnet/irismod/modules/htlc"
	htlctypes "github.com/irisnet/irismod/modules/htlc/types"
	"github.com/irisnet/irismod/modules/mt"
	mttypes "github.com/irisnet/irismod/modules/mt/types"
	nftmodule "github.com/irisnet/irismod/modules/nft/module"
	nfttypes "github.com/irisnet/irismod/modules/nft/types"
	"github.com/irisnet/irismod/modules/oracle"
	oracletypes "github.com/irisnet/irismod/modules/oracle/types"
	"github.com/irisnet/irismod/modules/random"
	randomtypes "github.com/irisnet/irismod/modules/random/types"
	"github.com/irisnet/irismod/modules/record"
	recordtypes "github.com/irisnet/irismod/modules/record/types"
	"github.com/irisnet/irismod/modules/service"
	servicetypes "github.com/irisnet/irismod/modules/service/types"
	"github.com/irisnet/irismod/modules/token"
	tokentypes "github.com/irisnet/irismod/modules/token/types"

	tibcmttransfer "github.com/bianjieai/tibc-go/modules/tibc/apps/mt_transfer"
	tibcmttypes "github.com/bianjieai/tibc-go/modules/tibc/apps/mt_transfer/types"
	tibcnfttransfer "github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer"
	tibcnfttypes "github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"
	tibc "github.com/bianjieai/tibc-go/modules/tibc/core"
	tibchost "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	tibccli "github.com/bianjieai/tibc-go/modules/tibc/core/client/cli"

	"github.com/evmos/ethermint/x/evm"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/evmos/ethermint/x/feemarket"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	nfttransfer "github.com/bianjieai/nft-transfer"
	ibcnfttransfertypes "github.com/bianjieai/nft-transfer/types"

	irisappparams "github.com/irisnet/irishub/v2/app/params"
	irisevm "github.com/irisnet/irishub/v2/modules/evm"
	"github.com/irisnet/irishub/v2/modules/guardian"
	guardiantypes "github.com/irisnet/irishub/v2/modules/guardian/types"
	"github.com/irisnet/irishub/v2/modules/mint"
	minttypes "github.com/irisnet/irishub/v2/modules/mint/types"
)

var (
	legacyProposalHandlers = []govclient.ProposalHandler{
		paramsclient.ProposalHandler,
		//distrclient.ProposalHandler,
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
		authzmodule.AppModuleBasic{},
		genutil.AppModuleBasic{
			GenTxValidator: genutiltypes.DefaultMessageValidator,
		},
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
		ibctm.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		transfer.AppModuleBasic{},
		vesting.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		consensus.AppModuleBasic{},

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
		nfttransfer.AppModuleBasic{},

		evm.AppModuleBasic{},
		feemarket.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     {authtypes.Burner},
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
		evmtypes.ModuleName: {
			authtypes.Minter,
			authtypes.Burner,
		}, // used for secure addition and subtraction of balance using module account
	}
)

func appModules(
	app *IrisApp,
	encodingConfig irisappparams.EncodingConfig,
	skipGenesisInvariants bool,
) []module.AppModule {
	appCodec := encodingConfig.Marshaler

	return []module.AppModule{
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(
			appCodec,
			app.AccountKeeper,
			authsims.RandomGenesisAccounts,
			app.GetSubspace(authtypes.ModuleName),
		),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(
			appCodec,
			app.BankKeeper,
			app.AccountKeeper,
			app.GetSubspace(banktypes.ModuleName),
		),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		crisis.NewAppModule(
			app.CrisisKeeper,
			skipGenesisInvariants,
			app.GetSubspace(crisistypes.ModuleName),
		),
		gov.NewAppModule(
			appCodec,
			app.GovKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.GetSubspace(govtypes.ModuleName),
		),
		mint.NewAppModule(
			appCodec,
			app.MintKeeper,
			app.GetSubspace(minttypes.ModuleName),
		),
		slashing.NewAppModule(
			appCodec,
			app.SlashingKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.StakingKeeper,
			app.GetSubspace(slashingtypes.ModuleName),
		),
		distr.NewAppModule(
			appCodec,
			app.DistrKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.StakingKeeper,
			app.GetSubspace(distrtypes.ModuleName),
		),
		staking.NewAppModule(
			appCodec,
			app.StakingKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.GetSubspace(stakingtypes.ModuleName),
		),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(*app.EvidenceKeeper),
		feegrantmodule.NewAppModule(
			appCodec,
			app.AccountKeeper,
			app.BankKeeper,
			app.FeeGrantKeeper,
			app.interfaceRegistry,
		),
		authzmodule.NewAppModule(
			appCodec,
			app.AuthzKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.interfaceRegistry,
		),
		consensus.NewAppModule(appCodec, app.ConsensusParamsKeeper),
		ibc.NewAppModule(app.IBCKeeper), tibc.NewAppModule(app.TIBCKeeper),
		params.NewAppModule(app.ParamsKeeper),
		app.transferModule,
		app.ibcnfttransferModule,
		app.nfttransferModule,
		app.mttransferModule,
		guardian.NewAppModule(appCodec, app.GuardianKeeper),
		token.NewAppModule(
			appCodec,
			app.TokenKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.GetSubspace(tokentypes.ModuleName),
		),
		record.NewAppModule(appCodec, app.RecordKeeper, app.AccountKeeper, app.BankKeeper),
		nftmodule.NewAppModule(appCodec, app.NFTKeeper, app.AccountKeeper, app.BankKeeper),
		mt.NewAppModule(appCodec, app.MTKeeper, app.AccountKeeper, app.BankKeeper),
		htlc.NewAppModule(
			appCodec,
			app.HTLCKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.GetSubspace(htlctypes.ModuleName),
		),
		coinswap.NewAppModule(
			appCodec,
			app.CoinswapKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.GetSubspace(coinswaptypes.ModuleName),
		),
		service.NewAppModule(
			appCodec,
			app.ServiceKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.GetSubspace(servicetypes.ModuleName),
		),
		oracle.NewAppModule(appCodec, app.OracleKeeper, app.AccountKeeper, app.BankKeeper),
		random.NewAppModule(appCodec, app.RandomKeeper, app.AccountKeeper, app.BankKeeper),
		farm.NewAppModule(
			appCodec,
			app.FarmKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.GetSubspace(farmtypes.ModuleName),
		),

		// Ethermint app modules
		irisevm.NewAppModule(
			app.EvmKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.GetSubspace(evmtypes.ModuleName),
		),
		feemarket.NewAppModule(app.FeeMarketKeeper, app.GetSubspace(feemarkettypes.ModuleName)),
	}
}

// simulationModules returns modules for simulation manager
// define the order of the modules for deterministic simulations
func simulationModules(
	app *IrisApp,
	encodingConfig irisappparams.EncodingConfig,
	_ bool,
) []module.AppModuleSimulation {
	appCodec := encodingConfig.Marshaler

	return []module.AppModuleSimulation{
		auth.NewAppModule(
			appCodec,
			app.AccountKeeper,
			authsims.RandomGenesisAccounts,
			app.GetSubspace(authtypes.ModuleName),
		),
		bank.NewAppModule(
			appCodec,
			app.BankKeeper,
			app.AccountKeeper,
			app.GetSubspace(banktypes.ModuleName),
		),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		gov.NewAppModule(
			appCodec,
			app.GovKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.GetSubspace(govtypes.ModuleName),
		),
		mint.NewAppModule(appCodec, app.MintKeeper, app.GetSubspace(minttypes.ModuleName)),
		feegrantmodule.NewAppModule(
			appCodec,
			app.AccountKeeper,
			app.BankKeeper,
			app.FeeGrantKeeper,
			app.interfaceRegistry,
		),
		staking.NewAppModule(
			appCodec,
			app.StakingKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.GetSubspace(stakingtypes.ModuleName),
		),
		distr.NewAppModule(
			appCodec,
			app.DistrKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.StakingKeeper,
			app.GetSubspace(distrtypes.ModuleName),
		),
		slashing.NewAppModule(
			appCodec,
			app.SlashingKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.StakingKeeper,
			app.GetSubspace(slashingtypes.ModuleName),
		),
		params.NewAppModule(app.ParamsKeeper),
		evidence.NewAppModule(*app.EvidenceKeeper),
		authzmodule.NewAppModule(
			appCodec,
			app.AuthzKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.interfaceRegistry,
		),
		ibc.NewAppModule(app.IBCKeeper),
		app.transferModule,
		app.ibcnfttransferModule,
		guardian.NewAppModule(appCodec, app.GuardianKeeper),
		token.NewAppModule(
			appCodec,
			app.TokenKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.GetSubspace(tokentypes.ModuleName),
		),
		record.NewAppModule(appCodec, app.RecordKeeper, app.AccountKeeper, app.BankKeeper),
		nftmodule.NewAppModule(appCodec, app.NFTKeeper, app.AccountKeeper, app.BankKeeper),
		htlc.NewAppModule(
			appCodec,
			app.HTLCKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.GetSubspace(htlctypes.ModuleName),
		),
		coinswap.NewAppModule(
			appCodec,
			app.CoinswapKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.GetSubspace(coinswaptypes.ModuleName),
		),
		service.NewAppModule(
			appCodec,
			app.ServiceKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.GetSubspace(servicetypes.ModuleName),
		),
		oracle.NewAppModule(appCodec, app.OracleKeeper, app.AccountKeeper, app.BankKeeper),
		random.NewAppModule(appCodec, app.RandomKeeper, app.AccountKeeper, app.BankKeeper),
		farm.NewAppModule(
			appCodec,
			app.FarmKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.GetSubspace(farmtypes.ModuleName),
		),
		tibc.NewAppModule(app.TIBCKeeper),

		// Ethermint app modules
		irisevm.NewAppModule(
			app.EvmKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.GetSubspace(evmtypes.ModuleName),
		),
		feemarket.NewAppModule(app.FeeMarketKeeper, app.GetSubspace(feemarkettypes.ModuleName)),
	}
}

/*
orderBeginBlockers tells the app's module manager how to set the order of
BeginBlockers, which are run at the beginning of every block.
Interchain Security Requirements:
During begin block slashing happens after distr.BeginBlocker so that
there is nothing left over in the validator fee pool, so as to keep the
CanWithdrawInvariant invariant.
NOTE: staking module is required if HistoricalEntries param > 0
NOTE: capability module's beginblocker must come before any modules using capabilities (e.g. IBC)
*/

func orderBeginBlockers() []string {
	return []string{
		// upgrades should be run first
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		minttypes.ModuleName,
		feemarkettypes.ModuleName,
		evmtypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibcexported.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
		consensustypes.ModuleName,

		//self module
		tokentypes.ModuleName,
		tibchost.ModuleName,
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

		ibcnfttransfertypes.ModuleName,
	}
}

/*
Interchain Security Requirements:
- provider.EndBlock gets validator updates from the staking module;
thus, staking.EndBlock must be executed before provider.EndBlock;
- creating a new consumer chain requires the following order,
CreateChildClient(), staking.EndBlock, provider.EndBlock;
thus, gov.EndBlock must be executed before staking.EndBlock
*/
func orderEndBlockers() []string {
	return []string{
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		evmtypes.ModuleName,
		feemarkettypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibcexported.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		consensustypes.ModuleName,

		//self module
		tokentypes.ModuleName,
		tibchost.ModuleName,
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

		ibcnfttransfertypes.ModuleName,
	}
}

/*
NOTE: The genutils module must occur after staking so that pools are
properly initialized with tokens from genesis accounts.
NOTE: The genutils module must also occur after auth so that it can access the params from auth.
NOTE: Capability module must occur first so that it can initialize any capabilities
so that other modules that want to create or claim capabilities afterwards in InitChain
can do so safely.
*/
func orderInitBlockers() []string {
	return []string{
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibcexported.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		evmtypes.ModuleName,
		// NOTE: feemarket module needs to be initialized before genutil module:
		// gentx transactions use MinGasPriceDecorator.AnteHandle
		feemarkettypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		consensustypes.ModuleName,

		//self module
		tokentypes.ModuleName,
		tibchost.ModuleName,
		nfttypes.ModuleName,
		htlctypes.ModuleName,
		recordtypes.ModuleName,
		// NOTE: coinswap module needs to be initialized before farm module:
		coinswaptypes.ModuleName,
		farmtypes.ModuleName,
		randomtypes.ModuleName,
		servicetypes.ModuleName,
		oracletypes.ModuleName,
		mttypes.ModuleName,
		tibcnfttypes.ModuleName,
		tibcmttypes.ModuleName,
		guardiantypes.ModuleName,
		// NOTE: crisis module must go at the end to check for invariants on each module
		crisistypes.ModuleName,

		ibcnfttransfertypes.ModuleName,
	}
}
