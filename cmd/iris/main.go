package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"golang.org/x/net/context"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/irisnet/irishub/app"

	"github.com/irisnet/irishub/address"
)

var (
	encodingConfig = app.MakeEncodingConfig()
	initClientCtx  = client.Context{}.
			WithJSONMarshaler(encodingConfig.Marshaler).
			WithTxGenerator(encodingConfig.TxGenerator).
			WithCodec(encodingConfig.Amino).
			WithInput(os.Stdin).
			WithAccountRetriever(types.NewAccountRetriever(encodingConfig.Marshaler)).
			WithBroadcastMode(flags.BroadcastBlock)
)

func main() {
	appCodec, cdc := app.MakeCodecs()

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(address.Bech32PrefixAccAddr, address.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(address.Bech32PrefixValAddr, address.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(address.Bech32PrefixConsAddr, address.Bech32PrefixConsPub)
	config.Seal()

	ctx := server.NewDefaultContext()
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:   "iris",
		Short: "IRIS Hub Daemon (Server)",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}
			return server.InterceptConfigsPreRunHandler(cmd)
		},
	}

	rootCmd.AddCommand(
		genutilcli.InitCmd(app.ModuleBasics, app.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}),
		genutilcli.MigrateGenesisCmd(),
		genutilcli.GenTxCmd(
			app.ModuleBasics, banktypes.GenesisBalancesIterator{},
		),
		genutilcli.ValidateGenesisCmd(app.ModuleBasics),
		AddGenesisAccountCmd(ctx, cdc, appCodec, app.DefaultCLIHome),
		cli.NewCompletionCmd(rootCmd, true),
		testnetCmd(ctx, cdc, app.ModuleBasics, banktypes.GenesisBalancesIterator{}),
		debug.Cmd(),
	)

	server.AddCommands(rootCmd, newApp, exportAppStateAndTMValidators)

	executorCtx := context.Background()
	executorCtx = context.WithValue(executorCtx, client.ClientContextKey, &client.Context{})
	executorCtx = context.WithValue(executorCtx, server.ServerContextKey, server.NewDefaultContext())

	executor := cli.PrepareBaseCmd(rootCmd, "", app.DefaultNodeHome)
	if err := executor.ExecuteContext(executorCtx); err != nil {
		fmt.Printf("failed execution: %s, exiting...\n", err)
		os.Exit(1)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer, appOpts server.AppOptions) server.Application {
	var cache sdk.MultiStorePersistentCache

	if cast.ToBool(appOpts.Get(server.FlagInterBlockCache)) {
		cache = store.NewCommitKVStoreCacheManager()
	}

	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	pruningOpts, err := server.GetPruningOptionsFromFlags(appOpts)
	if err != nil {
		panic(err)
	}

	return app.NewIrisApp(
		logger, db, traceStore, true, skipUpgradeHeights,
		cast.ToString(appOpts.Get(flags.FlagHome)),
		cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)),
		baseapp.SetPruning(pruningOpts),
		baseapp.SetMinGasPrices(cast.ToString(appOpts.Get(server.FlagMinGasPrices))),
		baseapp.SetHaltHeight(cast.ToUint64(appOpts.Get(server.FlagHaltHeight))),
		baseapp.SetHaltTime(cast.ToUint64(appOpts.Get(server.FlagHaltTime))),
		baseapp.SetInterBlockCache(cache),
		baseapp.SetTrace(cast.ToBool(appOpts.Get(server.FlagTrace))),
	)
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, *abci.ConsensusParams, error) {

	var irisApp *app.IrisApp
	if height != -1 {
		irisApp = app.NewIrisApp(logger, db, traceStore, false, map[int64]bool{}, "", uint(1))

		if err := irisApp.LoadHeight(height); err != nil {
			return nil, nil, nil, err
		}
	} else {
		irisApp = app.NewIrisApp(logger, db, traceStore, true, map[int64]bool{}, "", uint(1))
	}

	return irisApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}
