package cmd

import (
	"errors"
	"fmt"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/v4/testutil"
	"io"
	"os"

	"github.com/spf13/cobra"

	tmcfg "github.com/cometbft/cometbft/config"
	tmcli "github.com/cometbft/cometbft/libs/cli"

	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/log"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/pruning"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"

	ethermintdebug "github.com/evmos/ethermint/client/debug"
	etherminthd "github.com/evmos/ethermint/crypto/hd"
	servercfg "github.com/evmos/ethermint/server/config"

	"github.com/irisnet/irishub/v4/app"
	"github.com/irisnet/irishub/v4/app/params"
	iristypes "github.com/irisnet/irishub/v4/types"
)

// NewRootCmd creates a new root command for simd. It is called once in the
// main function.
func NewRootCmd() (*cobra.Command, params.EncodingConfig) {

	tempApplication := app.NewIrisApp(log.NewNopLogger(), dbm.NewMemDB(), nil, true, app.RegisterEncodingConfig(), testutil.EmptyAppOptions{})

	encodingConfig := app.RegisterEncodingConfig()
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastSync).
		WithHomeDir(iristypes.DefaultNodeHome).
		WithKeyringOptions(etherminthd.EthSecp256k1Option()).
		WithViper("")

	rootCmd := &cobra.Command{
		Use:   "iris",
		Short: "IRIS Hub app command",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			if err = client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			converter.handlePreRun(cmd, args)

			customTemplate, customIRISHubConfig := initAppConfig()
			customTMConfig := initTendermintConfig()
			return server.InterceptConfigsPreRunHandler(
				cmd,
				customTemplate,
				customIRISHubConfig,
				customTMConfig,
			)
		},
		PersistentPostRun: func(cmd *cobra.Command, _ []string) {
			converter.handlePostRun(cmd)
		},
	}

	initRootCmd(rootCmd, encodingConfig)

	autoCliOpts := enrichAutoCliOpts(tempApplication.AutoCliOpts(), initClientCtx)
	if err := autoCliOpts.EnhanceRootCommand(rootCmd); err != nil {
		panic(err)
	}

	return rootCmd, encodingConfig
}

func enrichAutoCliOpts(autoCliOpts autocli.AppOptions, clientCtx client.Context) autocli.AppOptions {
	autoCliOpts.AddressCodec = addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
	autoCliOpts.ValidatorAddressCodec = addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix())
	autoCliOpts.ConsensusAddressCodec = addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix())

	autoCliOpts.ClientCtx = clientCtx

	return autoCliOpts
}

// initTendermintConfig helps to override default Tendermint Config values.
// return tmcfg.DefaultConfig if no custom configuration is required for the application.
func initTendermintConfig() *tmcfg.Config {
	cfg := tmcfg.DefaultConfig()

	// these values put a higher strain on node memory
	// cfg.P2P.MaxNumInboundPeers = 100
	// cfg.P2P.MaxNumOutboundPeers = 40

	return cfg
}

// initAppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func initAppConfig() (string, interface{}) {
	customAppTemplate, customAppConfig := servercfg.AppConfig(iristypes.NativeToken.MinUnit)
	srvCfg, ok := customAppConfig.(servercfg.Config)
	if !ok {
		panic(fmt.Errorf("unknown app config type %T", customAppConfig))
	}

	srvCfg.StateSync.SnapshotInterval = 1000
	srvCfg.StateSync.SnapshotKeepRecent = 10

	return customAppTemplate, srvCfg
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig params.EncodingConfig) {
	ac := appCreator{
		encCfg: encodingConfig,
	}

	rootCmd.AddCommand(
		genutilcli.InitCmd(app.ModuleBasics, iristypes.DefaultNodeHome),
		genutilcli.ValidateGenesisCmd(app.ModuleBasics),
		AddGenesisAccountCmd(iristypes.DefaultNodeHome),
		tmcli.NewCompletionCmd(rootCmd, true),
		testnetCmd(app.ModuleBasics, banktypes.GenesisBalancesIterator{}),
		ethermintdebug.Cmd(),
		mergeGenesisCmd(encodingConfig),
		pruning.Cmd(ac.newApp, iristypes.DefaultNodeHome),
	)

	//ethermintserver.AddCommands(
	//	rootCmd,
	//	ethermintserver.NewDefaultStartOptions(ac.newApp, iristypes.DefaultNodeHome),
	//	ac.appExport,
	//	addModuleInitFlags,
	//)
	server.AddCommands(
		rootCmd,
		iristypes.DefaultNodeHome,
		ac.newApp,
		ac.appExport,
		addModuleInitFlags,
	)

	// add keybase, auxiliary RPC, query, and tx child commands
	rootCmd.AddCommand(
		server.StatusCommand(),
		genesisCommand(encodingConfig),
		queryCommand(),
		txCommand(),
		Commands(iristypes.DefaultNodeHome),
	)
}

// genesisCommand builds genesis-related `simd genesis` command. Users may provide application specific commands as a parameter
func genesisCommand(encodingConfig params.EncodingConfig, cmds ...*cobra.Command) *cobra.Command {
	cmd := genutilcli.GenesisCoreCommand(
		encodingConfig.TxConfig,
		app.ModuleBasics,
		iristypes.DefaultNodeHome,
	)

	for _, subCmd := range cmds {
		cmd.AddCommand(subCmd)
	}
	return cmd
}

func addModuleInitFlags(rootCmd *cobra.Command) {
	crisis.AddModuleInitFlags(rootCmd)
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		rpc.ValidatorCommand(),
		authcmd.QueryTxsByEventsCmd(),
		authcmd.QueryTxCmd(),
	)

	app.ModuleBasics.AddQueryCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetValidateSignaturesCommand(),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
	)

	//app.ModuleBasics.AddTxCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

type appCreator struct {
	encCfg params.EncodingConfig
}

func (ac appCreator) newApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	appOpts servertypes.AppOptions,
) servertypes.Application {
	baseappOptions := server.DefaultBaseappOptions(appOpts)
	return app.NewIrisApp(
		logger,
		db,
		traceStore,
		true,
		ac.encCfg,
		appOpts,
		baseappOptions...,
	)
}

// createIrisappAndExport creates a new irisapp (optionally at a given height) and exports state.
func (ac appCreator) appExport(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
	modulesToExport []string,
) (servertypes.ExportedApp, error) {
	homePath, ok := appOpts.Get(flags.FlagHome).(string)
	if !ok || homePath == "" {
		return servertypes.ExportedApp{}, errors.New("application home is not set")
	}

	var loadLatest bool
	if height == -1 {
		loadLatest = true
	}

	irisApp := app.NewIrisApp(
		logger,
		db,
		traceStore,
		loadLatest,
		ac.encCfg,
		appOpts,
	)

	if height != -1 {
		if err := irisApp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	}

	return irisApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs, modulesToExport)
}
