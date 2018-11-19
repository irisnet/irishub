package main

import (
	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/tools/prometheus"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"
	irisInit "github.com/irisnet/irishub/init"
)

func init() {
	cobra.EnableCommandSorting = false
	cdc := app.MakeCodec()
	rootCmd = prometheus.MonitorCommand(cdc)
	rootCmd.SilenceUsage = true
}

var rootCmd *cobra.Command

func main() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(irisInit.Bech32PrefixAccAddr, irisInit.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(irisInit.Bech32PrefixValAddr, irisInit.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(irisInit.Bech32PrefixConsAddr, irisInit.Bech32PrefixConsPub)
	config.Seal()

	executor := cli.PrepareMainCmd(rootCmd, "IRIS", app.DefaultNodeHome)
	err := executor.Execute()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
