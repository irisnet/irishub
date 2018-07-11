package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/tendermint/tmlibs/cli"

	basecmd "github.com/cosmos/cosmos-sdk/server/commands"
	"github.com/irisnet/iris-hub/version"
)

// IrisCmd is the entry point for this binary
var (
	IrisCmd = &cobra.Command{
		Use:   "iris",
		Short: "IRIS Hub - a regional Cosmos Hub with a powerful iService infrastructure",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	lineBreak = &cobra.Command{Run: func(*cobra.Command, []string) {}}
)

func main() {
	// disable sorting
	cobra.EnableCommandSorting = false

	// add commands
	prepareNodeCommands()
	prepareRestServerCommands()
	prepareClientCommands()

	IrisCmd.AddCommand(
		nodeCmd,
		restServerCmd,
		clientCmd,

		lineBreak,
		version.VersionCmd,
		//auto.AutoCompleteCmd,
	)

	// prepare and add flags
	basecmd.SetUpRoot(IrisCmd)
	executor := cli.PrepareMainCmd(IrisCmd, "GA", os.ExpandEnv("$HOME/.iris-cli"))
	executor.Execute()
}
