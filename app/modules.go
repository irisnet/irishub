package app

import (
	"cosmossdk.io/x/evidence"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	feegrantmodule "cosmossdk.io/x/feegrant/module"
	"cosmossdk.io/x/upgrade"
	upgradeclient "cosmossdk.io/x/upgrade/client/cli"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
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
	"github.com/cosmos/ibc-go/modules/capability"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	iristypes "github.com/irisnet/irishub/v4/types"
	"github.com/spf13/cobra"

	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v8/modules/core"

	// ibcclientclient "github.com/cosmos/ibc-go/v8/modules/core/02-client/client"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"

	"mods.irisnet.org/modules/coinswap"
	coinswaptypes "mods.irisnet.org/modules/coinswap/types"
	"mods.irisnet.org/modules/farm"
	farmtypes "mods.irisnet.org/modules/farm/types"
	"mods.irisnet.org/modules/htlc"
	htlctypes "mods.irisnet.org/modules/htlc/types"
	"mods.irisnet.org/modules/mt"
	mttypes "mods.irisnet.org/modules/mt/types"
	"mods.irisnet.org/modules/nft"
	nfttypes "mods.irisnet.org/modules/nft/types"
	"mods.irisnet.org/modules/oracle"
	oracletypes "mods.irisnet.org/modules/oracle/types"
	"mods.irisnet.org/modules/random"
	randomtypes "mods.irisnet.org/modules/random/types"
	"mods.irisnet.org/modules/record"
	recordtypes "mods.irisnet.org/modules/record/types"
	"mods.irisnet.org/modules/service"
	servicetypes "mods.irisnet.org/modules/service/types"
	"mods.irisnet.org/modules/token"
	tokentypes "mods.irisnet.org/modules/token/types"

	tibcmttypes "github.com/bianjieai/tibc-go/modules/tibc/apps/mt_transfer/types"
	tibcnfttypes "github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"
	tibc "github.com/bianjieai/tibc-go/modules/tibc/core"
	tibchost "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	tibccli "github.com/bianjieai/tibc-go/modules/tibc/core/client/cli"

	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/evmos/ethermint/x/feemarket"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	ibcnfttransfertypes "github.com/bianjieai/nft-transfer/types"

	irisappparams "github.com/irisnet/irishub/v4/app/params"
	irisevm "github.com/irisnet/irishub/v4/modules/evm"
	"github.com/irisnet/irishub/v4/modules/guardian"
	guardiantypes "github.com/irisnet/irishub/v4/modules/guardian/types"
	"github.com/irisnet/irishub/v4/modules/mint"
	minttypes "github.com/irisnet/irishub/v4/modules/mint/types"
)

var (
	legacyProposalHandlers = []govclient.ProposalHandler{
		paramsclient.ProposalHandler,
		govclient.NewProposalHandler(func() *cobra.Command {
			return upgradeclient.NewCmdSubmitUpgradeProposal(addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()))
		}),
		govclient.NewProposalHandler(func() *cobra.Command {
			return upgradeclient.NewCmdSubmitCancelUpgradeProposal(addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()))
		}),
	}

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
		icatypes.ModuleName:            nil,
		evmtypes.ModuleName: {
			authtypes.Minter,
			authtypes.Burner,
		}, // used for secure addition and subtraction of balance using module account
	}
)

// ModuleBasics defines the module BasicManager that is in charge of setting up basic,
// non-dependant module elements, such as codec registration
// and genesis verification.
func newBasicManagerFromManager(app *IrisApp) module.BasicManager {
	basicManager := module.NewBasicManagerFromManager(
		app.mm,
		map[string]module.AppModuleBasic{
			genutiltypes.ModuleName: genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
			govtypes.ModuleName:     gov.NewAppModuleBasic(append(legacyProposalHandlers, tibccli.GovHandlers...)),
		})
	basicManager.RegisterLegacyAminoCodec(app.legacyAmino)
	basicManager.RegisterInterfaces(app.interfaceRegistry)
	return basicManager
}

func appModules(
	app *IrisApp,
	encodingConfig irisappparams.EncodingConfig,
	skipGenesisInvariants bool,
) []module.AppModule {
	appCodec := encodingConfig.Codec

	return []module.AppModule{
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app.BaseApp,
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
			app.interfaceRegistry,
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
		upgrade.NewAppModule(app.UpgradeKeeper, addresscodec.NewBech32Codec(iristypes.Bech32PrefixAccAddr)),
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
		app.TransferModule,
		app.IBCNftTransferModule,
		app.ICAModule,
		app.NftTransferModule,
		app.MtTransferModule,
		guardian.NewAppModule(appCodec, app.GuardianKeeper),
		token.NewAppModule(
			appCodec,
			app.TokenKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.GetSubspace(tokentypes.ModuleName),
		),
		record.NewAppModule(appCodec, app.RecordKeeper, app.AccountKeeper, app.BankKeeper),
		nft.NewAppModule(appCodec, app.NFTKeeper, app.AccountKeeper, app.BankKeeper),
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
	appCodec := encodingConfig.Codec

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
			app.interfaceRegistry,
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
		app.TransferModule,
		app.IBCNftTransferModule,
		guardian.NewAppModule(appCodec, app.GuardianKeeper),
		token.NewAppModule(
			appCodec,
			app.TokenKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.GetSubspace(tokentypes.ModuleName),
		),
		record.NewAppModule(appCodec, app.RecordKeeper, app.AccountKeeper, app.BankKeeper),
		nft.NewAppModule(appCodec, app.NFTKeeper, app.AccountKeeper, app.BankKeeper),
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
		icatypes.ModuleName,
		consensustypes.ModuleName,

		// self module
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
		icatypes.ModuleName,
		consensustypes.ModuleName,

		// self module
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
		icatypes.ModuleName,
		consensustypes.ModuleName,

		// self module
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
