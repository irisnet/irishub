package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/irisnet/irishub/app"
	bam "github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/server"
	irisInit "github.com/irisnet/irishub/server/init"
	"github.com/irisnet/irishub/version"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	pvm "github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

func main() {
	//	sdk.InitBech32Prefix()
	cdc := app.MakeLatestCodec()
	ctx := server.NewDefaultContext()
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:               "iris",
		Short:             "IRIS Hub Daemon (Server)",
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
		server.ResetCmd(ctx, cdc, resetAppState),
		server.ExportCmd(ctx, cdc, exportAppStateAndTMValidators),
		server.SnapshotCmd(cdc),
		client.LineBreak,
	)

	rootCmd.AddCommand(
		version.ServeVersionCommand(cdc),
	)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "IRIS", app.DefaultNodeHome)
	executor.Execute()
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer, config *cfg.InstrumentationConfig) abci.Application {
	return app.NewIrisApp(logger, db, config, traceStore,
		bam.SetPruning(viper.GetString("pruning")),
		bam.SetMinimumFees(viper.GetString("minimum_fees")),
		bam.SetCheckInvariant(viper.GetBool("check_invariant")),
		bam.SetTrackCoinFlow(viper.GetBool("track_coin_flow")),
	)
}

func exportAppStateAndTMValidators(ctx *server.Context,
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool,
) (int64, json.RawMessage, []tmtypes.GenesisValidator, error) {
	gApp := app.NewIrisApp(logger, db, ctx.Config.Instrumentation, traceStore)
	lastBlockHeight := gApp.LastBlockHeight()
	if height > 0 && height < lastBlockHeight {
		err := gApp.LoadVersion(height, protocol.KeyMain, false)
		if err != nil {
			if strings.Contains(err.Error(), fmt.Sprintf("wanted to load target %v but only found up to", height)) {
				return height, nil, nil, fmt.Errorf("unable to export snapshot height state %v that does not exist. "+
					"If necessary, reset the application state to the specified height using command reset, and then export the state", height)
			}
			return height, nil, nil, err
		}
	} else {
		height = lastBlockHeight
	}
	appState, validators, err := gApp.ExportAppStateAndValidators(forZeroHeight)
	return height, appState, validators, err
}

func resetAppState(ctx *server.Context,
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64) error {
	gApp := app.NewIrisApp(logger, db, ctx.Config.Instrumentation, traceStore)
	if height > 0 {
		if replay, replayHeight := gApp.ResetOrReplay(height); replay {
			_, err := startNodeAndReplay(ctx, gApp, replayHeight)
			if err != nil {
				return err
			}
		}
	}
	if height == 0 {
		return errors.New("No need to reset to zero height, it is always consistent with genesis.json")
	}
	return nil
}

func startNodeAndReplay(ctx *server.Context, app *app.IrisApp, height int64) (n *node.Node, err error) {
	cfg := ctx.Config
	cfg.BaseConfig.ReplayHeight = height

	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		return nil, err
	}
	newNode := func(c chan int) {
		defer func() {
			c <- 0
		}()
		n, err = node.NewNode(
			cfg,
			pvm.LoadOrGenFilePV(cfg.PrivValidatorFile()),
			nodeKey,
			proxy.NewLocalClientCreator(app),
			node.DefaultGenesisDocProviderFunc(cfg),
			node.DefaultDBProvider,
			node.DefaultMetricsProvider(cfg.Instrumentation),
			ctx.Logger.With("module", "node"),
		)
		if err != nil {
			c <- 1
		}
	}
	ch := make(chan int)
	go newNode(ch)
	v := <-ch
	if v == 0 {
		err = nil
	}
	return nil, err
}
