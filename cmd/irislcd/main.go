package main

import (
	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/lite"
	_ "github.com/irisnet/irishub/lite/statik"
	"github.com/irisnet/irishub/version"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
)

// rootCmd is the entry point for this binary
var (
	rootCmd = &cobra.Command{
		Use:   "irislcd",
		Short: "IRIS Hub API Server (Lite Client Daemon)",
	}
)

func main() {
	// sdk.InitBech32Prefix()
	cobra.EnableCommandSorting = false
	cdc := app.MakeLatestCodec()

	rootCmd.AddCommand(
		lite.ServeLCDStartCommand(cdc),
		version.ServeVersionCommand(cdc),
	)

	// prepare and add flags
	executor := cli.PrepareMainCmd(rootCmd, "IRISLCD", app.DefaultLCDHome)
	err := executor.Execute()
	if err != nil {
		// handle with #870
		panic(err)
	}
}
