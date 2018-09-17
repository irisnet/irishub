package main

import (
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/client"
	bankcmd "github.com/irisnet/irishub/client/bank/cli"
	govcmd "github.com/irisnet/irishub/client/gov/cli"
	keyscmd "github.com/irisnet/irishub/client/keys/cli"
	slashingcmd "github.com/irisnet/irishub/client/slashing/cli"
	stakecmd "github.com/irisnet/irishub/client/stake/cli"
	tendermintcmd "github.com/irisnet/irishub/client/tendermint/cli"
	upgradecmd "github.com/irisnet/irishub/client/upgrade/cli"
	"github.com/irisnet/irishub/version"
)

// rootCmd is the entry point for this binary
var (
	rootCmd = &cobra.Command{
		Use:   "iriscli",
		Short: "irishub command line interface",
	}
)

func main() {
	cobra.EnableCommandSorting = false
	cdc := app.MakeCodec()

	//Add state commands
	tendermintCmd := &cobra.Command{
		Use:   "tendermint",
		Short: "Tendermint state querying subcommands",
	}
	tendermintCmd.AddCommand(
		tendermintcmd.QueryTxCmd(cdc),
		tendermintcmd.SearchTxCmd(cdc),
		tendermintcmd.BlockCommand(),
		tendermintcmd.ValidatorCommand(),
	)
	rootCmd.AddCommand(tendermintCmd)

	//Add bank commands
	bankCmd := &cobra.Command{
		Use:   "bank",
		Short: "Bank subcommands",
	}
	bankCmd.AddCommand(
		client.GetCommands(
			bankcmd.GetCmdQueryCoinType(cdc),
			bankcmd.GetAccountCmd("acc", cdc, authcmd.GetAccountDecoder(cdc)),
		)...)
	bankCmd.AddCommand(
		client.PostCommands(
			bankcmd.SendTxCmd(cdc),
		)...)
	rootCmd.AddCommand(
		bankCmd,
	)

	//Add stake commands
	govCmd := &cobra.Command{
		Use:   "gov",
		Short: "Governance and voting subcommands",
	}
	govCmd.AddCommand(
		client.GetCommands(
			govcmd.GetCmdQueryProposal("gov", cdc),
			govcmd.GetCmdQueryVote("gov", cdc),
			govcmd.GetCmdQueryVotes("gov", cdc),
			govcmd.GetCmdQueryProposals("gov", cdc),
			govcmd.GetCmdQueryConfig("gov", cdc),
		)...)
	govCmd.AddCommand(
		client.PostCommands(
			govcmd.GetCmdSubmitProposal(cdc),
			govcmd.GetCmdDeposit(cdc),
			govcmd.GetCmdVote(cdc),
		)...)
	rootCmd.AddCommand(
		govCmd,
	)

	//Add staking and slashing commands
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
			stakecmd.GetCmdQueryUnbondingDelegation("stake", cdc),
			stakecmd.GetCmdQueryUnbondingDelegations("stake", cdc),
			stakecmd.GetCmdQueryRedelegation("stake", cdc),
			stakecmd.GetCmdQueryRedelegations("stake", cdc),
			slashingcmd.GetCmdQuerySigningInfo("slashing", cdc),
		)...)
	stakeCmd.AddCommand(
		client.PostCommands(
			stakecmd.GetCmdCreateValidator(cdc),
			stakecmd.GetCmdEditValidator(cdc),
			stakecmd.GetCmdDelegate(cdc),
			stakecmd.GetCmdUnbond("stake", cdc),
			stakecmd.GetCmdRedelegate("stake", cdc),
			slashingcmd.GetCmdUnrevoke(cdc),
		)...)
	rootCmd.AddCommand(
		stakeCmd,
	)

	//Add upgrade commands
	upgradeCmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Software Upgrade subcommands",
	}
	upgradeCmd.AddCommand(
		client.GetCommands(
			upgradecmd.GetCmdVersion("upgrade", cdc),
		)...)
	rootCmd.AddCommand(
		upgradeCmd,
	)

	//Add keys and version commands
	rootCmd.AddCommand(
		client.LineBreak,
		keyscmd.Commands(),
		version.ServeVersionCommand(cdc),
	)

	// prepare and add flags
	executor := cli.PrepareMainCmd(rootCmd, "IRISCLI", app.DefaultCLIHome)
	err := executor.Execute()
	if err != nil {
		// handle with #870
		panic(err)
	}
}
