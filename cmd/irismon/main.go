package main

import (
	"os"
	"github.com/spf13/cobra"
	"github.com/irisnet/irishub/tools/prometheus"
	"github.com/irisnet/irishub/app"
	"github.com/tendermint/tendermint/libs/cli"
)

func init() {
	cobra.EnableCommandSorting = false
	cdc := app.MakeCodec()
	rootCmd = prometheus.MonitorCommand(cdc)
	rootCmd.SilenceUsage = true
}

var rootCmd *cobra.Command

func main() {
	executor := cli.PrepareMainCmd(rootCmd, "IRIS", app.DefaultNodeHome)
	err := executor.Execute()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
