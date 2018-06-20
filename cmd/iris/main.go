package main

import (
	"encoding/json"

	"github.com/spf13/cobra"

	abci "github.com/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/tendermint/tmlibs/cli"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/irisnet/iris-hub/app"
)

func main() {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:               "iris",
		Short:             "iris Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	server.AddCommands(ctx, cdc, rootCmd, app.GaiaAppInit(),
		server.ConstructAppCreator(newApp, "iris"),
		server.ConstructAppExporter(exportAppStateAndTMValidators, "iris"))

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "IRIS", app.DefaultNodeHome)
	executor.Execute()
}

func newApp(logger log.Logger, db dbm.DB) abci.Application {
	return app.NewIrisApp(logger, db)
}

func exportAppStateAndTMValidators(logger log.Logger, db dbm.DB) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	irisApp := app.NewIrisApp(logger, db)
	return irisApp.ExportAppStateAndValidators()
}