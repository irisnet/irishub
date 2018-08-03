package main

import (
	"encoding/json"
	"io"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/irisnet/irishub/examples/basecoin1/app"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/irisnet/irishub/examples/basecoin1/version"

	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/irisnet/irishub/tools/prometheus"
)

func main() {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:               "basecoind",
		Short:             "basecoin Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	rootCmd.PersistentFlags().String("log_level", ctx.Config.LogLevel, "Log level")

	tendermintCmd := &cobra.Command{
		Use:   "tendermint",
		Short: "Tendermint subcommands",
	}

	tendermintCmd.AddCommand(
		server.ShowNodeIDCmd(ctx),
		server.ShowValidatorCmd(ctx),
	)

	startCmd := server.StartCmd(ctx, server.ConstructAppCreator(newApp, "basecoin"))
	startCmd.Flags().Bool(app.FlagReplay, false, "Replay the last block")
	rootCmd.AddCommand(
		server.InitCmd(ctx, cdc, app.BasecoinAppInit()),
		startCmd,
		server.TestnetFilesCmd(ctx, cdc, app.BasecoinAppInit()),
		server.UnsafeResetAllCmd(ctx),
		client.LineBreak,
		tendermintCmd,
		server.ExportCmd(ctx, cdc, server.ConstructAppExporter(exportAppStateAndTMValidators, "basecoin")),
		client.LineBreak,
		version.VersionCmd,
	)

	rootCmd.AddCommand(prometheus.MonitorCommand(cdc))

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "BASECOIN", app.DefaultNodeHome)
	executor.Execute()
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	return app.NewBasecoinApp(logger, db, traceStore, bam.SetPruning(viper.GetString("pruning")))
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	gApp := app.NewBasecoinApp(logger, db, traceStore)
	return gApp.ExportAppStateAndValidators()
}
