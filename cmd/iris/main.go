package main

import (
	"encoding/json"
	"io"

	"github.com/irisnet/irishub/app"
	bam "github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/server"
	irisInit "github.com/irisnet/irishub/server/init"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

func main() {

	sdk.InitBech32Prefix()

	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:               "iris",
		Short:             "iris Daemon (server)",
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
		server.ShowAddressCmd(ctx),
	)

	startCmd := server.StartCmd(ctx, newApp)
	startCmd.Flags().Bool(app.FlagReplay, false, "Replay the last block")
	rootCmd.AddCommand(
		irisInit.InitCmd(ctx, cdc),
		irisInit.GenTxCmd(ctx, cdc),
		irisInit.AddGenesisAccountCmd(ctx, cdc),
		irisInit.TestnetFilesCmd(ctx, cdc),
		irisInit.CollectGenTxsCmd(ctx, cdc),
		startCmd,
		//server.TestnetFilesCmd(ctx, cdc, app.IrisAppInit()),
		server.UnsafeResetAllCmd(ctx),
		client.LineBreak,
		tendermintCmd,
		server.ExportCmd(ctx, cdc, exportAppStateAndTMValidators),
		client.LineBreak,
	)

	rootCmd.AddCommand(
		version.ServeVersionCommand(cdc),
	)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "IRIS", app.DefaultNodeHome)
	executor.Execute()
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	return app.NewIrisApp(logger, db, traceStore,
		bam.SetPruning(viper.GetString("pruning")),
		bam.SetMinimumFees(viper.GetString("minimum_fees")),
		bam.SetInvariantLevel(viper.GetString("invariant_level")),
	)
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	gApp := app.NewIrisApp(logger, db, traceStore)
	if height != -1 {
		err := gApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
	}
	return gApp.ExportAppStateAndValidators(forZeroHeight)
}
