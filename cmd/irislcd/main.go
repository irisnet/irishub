package main

import (
	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/lite"
	_ "github.com/irisnet/irishub/lite/statik"
	irisInit "github.com/irisnet/irishub/server/init"
	"github.com/irisnet/irishub/version"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
)

// rootCmd is the entry point for this binary
var (
	rootCmd = &cobra.Command{
		Use:   "irislcd",
		Short: "irishub lite server interface",
	}
)

func main() {

	irisInit.InitBech32Prefix()

	cobra.EnableCommandSorting = false
	cdc := app.MakeCodec()

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
