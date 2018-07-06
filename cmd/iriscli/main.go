package main

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/tendermint/tmlibs/cli"

	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	ibccmd "github.com/cosmos/cosmos-sdk/x/ibc/client/cli"
	slashingcmd "github.com/cosmos/cosmos-sdk/x/slashing/client/cli"
	stakecmd "github.com/cosmos/cosmos-sdk/x/stake/client/cli"
	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/version"
	"github.com/irisnet/irishub/tools/prometheus"
)

// rootCmd is the entry point for this binary
var (
	rootCmd = &cobra.Command{
		Use:   "iriscli",
		Short: "irishub light-client",
	}
)

func main() {
	cobra.EnableCommandSorting = false
	cdc := app.MakeCodec()

	// TODO: setup keybase, viper object, etc. to be passed into
	// the below functions and eliminate global vars, like we do
	// with the cdc

	// add standard rpc commands
	rpc.AddCommands(rootCmd)

	//Add state commands
	tendermintCmd := &cobra.Command{
		Use:   "tendermint",
		Short: "Tendermint state querying subcommands",
	}
	tendermintCmd.AddCommand(
		rpc.BlockCommand(),
		rpc.ValidatorCommand(),
	)
	tx.AddCommands(tendermintCmd, cdc)

	//Add IBC commands
	ibcCmd := &cobra.Command{
		Use:   "ibc",
		Short: "Inter-Blockchain Communication subcommands",
	}
	ibcCmd.AddCommand(
		client.PostCommands(
			ibccmd.IBCTransferCmd(cdc),
			ibccmd.IBCRelayCmd(cdc),
		)...)

	advancedCmd := &cobra.Command{
		Use:   "advanced",
		Short: "Advanced subcommands",
	}

	advancedCmd.AddCommand(
		tendermintCmd,
		ibcCmd,
		ServeCommand(cdc),
	)
	rootCmd.AddCommand(
		advancedCmd,
		client.LineBreak,
	)

	//Add stake commands
	stakeCmd := &cobra.Command{
		Use:   "stake",
		Short: "Stake and validation subcommands",
	}
	stakeCmd.AddCommand(
		client.GetCommands(
			stakecmd.GetCmdQueryValidator("stake", cdc),
			stakecmd.GetCmdQueryValidators("stake", cdc),
			stakecmd.GetCmdQueryDelegation("stake", cdc),
			stakecmd.GetCmdQueryDelegations("stake", cdc),
			slashingcmd.GetCmdQuerySigningInfo("slashing", cdc),
		)...)
	stakeCmd.AddCommand(
		client.PostCommands(
			stakecmd.GetCmdCreateValidator(cdc),
			stakecmd.GetCmdEditValidator(cdc),
			stakecmd.GetCmdDelegate(cdc),
			stakecmd.GetCmdUnbond(cdc),
			slashingcmd.GetCmdUnrevoke(cdc),
		)...)
	rootCmd.AddCommand(
		stakeCmd,
	)

	//Add auth and bank commands
	rootCmd.AddCommand(
		client.GetCommands(
			authcmd.GetAccountCmd("acc", cdc, authcmd.GetAccountDecoder(cdc)),
		)...)
	rootCmd.AddCommand(
		client.PostCommands(
			bankcmd.SendTxCmd(cdc),
		)...)

	// add proxy, version and key info
	rootCmd.AddCommand(
		keys.Commands(),
		client.LineBreak,
		version.VersionCmd,
	)

	rootCmd.AddCommand(prometheus.MonitorCommand("stake", cdc))

	// prepare and add flags
	executor := cli.PrepareMainCmd(rootCmd, "IRIS", app.DefaultCLIHome)
	executor.Execute()
}
