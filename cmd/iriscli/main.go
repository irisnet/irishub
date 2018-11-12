package main

import (
	"os"
	"path"

	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/client"
	bankcmd "github.com/irisnet/irishub/client/bank/cli"
	distributioncmd "github.com/irisnet/irishub/client/distribution/cli"
	govcmd "github.com/irisnet/irishub/client/gov/cli"
	servicecmd "github.com/irisnet/irishub/client/service/cli"
	keyscmd "github.com/irisnet/irishub/client/keys/cli"
	recordcmd "github.com/irisnet/irishub/client/record/cli"
	slashingcmd "github.com/irisnet/irishub/client/slashing/cli"
	stakecmd "github.com/irisnet/irishub/client/stake/cli"
	tendermintrpccmd "github.com/irisnet/irishub/client/tendermint/rpc"
	tenderminttxcmd "github.com/irisnet/irishub/client/tendermint/tx"
	upgradecmd "github.com/irisnet/irishub/client/upgrade/cli"
	"github.com/irisnet/irishub/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
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

	rootCmd.AddCommand(client.ConfigCmd())
	rootCmd.AddCommand(client.LineBreak)

	rootCmd.AddCommand(tendermintrpccmd.StatusCommand())
	//Add state commands
	tendermintCmd := &cobra.Command{
		Use:   "tendermint",
		Short: "Tendermint state querying subcommands",
	}
	tendermintCmd.AddCommand(
		tenderminttxcmd.QueryTxCmd(cdc),
		tenderminttxcmd.SearchTxCmd(cdc),
		tendermintrpccmd.BlockCommand(),
		tendermintrpccmd.ValidatorCommand(),
	)
	rootCmd.AddCommand(tendermintCmd)
	rootCmd.AddCommand(client.LineBreak)

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
			bankcmd.GetSignCommand(cdc, authcmd.GetAccountDecoder(cdc)),
		)...)
	rootCmd.AddCommand(
		bankCmd,
	)

	//Add distribution commands
	distributionCmd := &cobra.Command{
		Use:   "distribution",
		Short: "Distribution subcommands",
	}
	distributionCmd.AddCommand(
		client.GetCommands(
			distributioncmd.GetWithdrawAddress("distr", cdc),
			distributioncmd.GetDelegationDistInfo("distr", cdc),
			distributioncmd.GetValidatorDistInfo("distr", cdc),
			distributioncmd.GetAllDelegationDistInfo("distr", cdc),
		)...)
	distributionCmd.AddCommand(
		client.PostCommands(
			distributioncmd.GetCmdSetWithdrawAddr(cdc),
			distributioncmd.GetCmdWithdrawRewards(cdc),
		)...)
	rootCmd.AddCommand(
		distributionCmd,
	)

	//Add gov commands
	govCmd := &cobra.Command{
		Use:   "gov",
		Short: "Governance and voting subcommands",
	}
	govCmd.AddCommand(
		client.GetCommands(
			govcmd.GetCmdQueryProposal("gov", cdc),
			govcmd.GetCmdQueryProposals("gov", cdc),
			govcmd.GetCmdQueryVote("gov", cdc),
			govcmd.GetCmdQueryVotes("gov", cdc),
			govcmd.GetCmdQueryDeposit("gov", cdc),
			govcmd.GetCmdQueryDeposits("gov", cdc),
			govcmd.GetCmdQueryTally("gov", cdc),
			govcmd.GetCmdQueryGovConfig("params", cdc),
			govcmd.GetCmdPullGovConfig("params", cdc),
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
			stakecmd.GetCmdQueryPool("stake", cdc),
			stakecmd.GetCmdQueryParams("stake", cdc),
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
			upgradecmd.GetInfoCmd("upgrade", cdc),
			upgradecmd.GetCmdQuerySwitch("upgrade", cdc),
		)...)
	upgradeCmd.AddCommand(
		client.PostCommands(
			upgradecmd.GetCmdSubmitSwitch(cdc),
		)...)
	rootCmd.AddCommand(
		upgradeCmd,
	)

	//Add service commands
	serviceCmd := &cobra.Command{
		Use:   "service",
		Short: "Service subcommands",
	}
	serviceCmd.AddCommand(
		client.GetCommands(
			servicecmd.GetCmdQueryScvDef("service", cdc),
			servicecmd.GetCmdQueryScvBind("service", cdc),
			servicecmd.GetCmdQueryScvBinds("service", cdc),
		)...)
	serviceCmd.AddCommand(client.PostCommands(
		servicecmd.GetCmdScvDef(cdc),
		servicecmd.GetCmdScvBind(cdc),
		servicecmd.GetCmdScvBindUpdate(cdc),
		servicecmd.GetCmdScvDisable(cdc),
		servicecmd.GetCmdScvEnable(cdc),
		servicecmd.GetCmdScvRefundDeposit(cdc),
	)...)

	rootCmd.AddCommand(
		serviceCmd,
	)

	//add record command
	recordCmd := &cobra.Command{
		Use:   "record",
		Short: "Record subcommands",
	}

	recordCmd.AddCommand(
		client.GetCommands(
			recordcmd.GetCmdQureyRecord("record", cdc),
			recordcmd.GetCmdDownload("record", cdc),
		)...)

	recordCmd.AddCommand(
		client.PostCommands(
			recordcmd.GetCmdSubmitRecord("record", cdc),
		)...)
	rootCmd.AddCommand(
		recordCmd,
	)

	//Add keys and version commands
	rootCmd.AddCommand(
		client.LineBreak,
		keyscmd.Commands(),
		client.LineBreak,
		version.ServeVersionCommand(cdc),
	)

	// prepare and add flags
	executor := cli.PrepareMainCmd(rootCmd, "IRISCLI", app.DefaultCLIHome)
	err := initConfig(rootCmd)
	if err != nil {
		panic(err)
	}

	err = executor.Execute()
	if err != nil {
		// handle with #870
		panic(err)
	}
}

func initConfig(cmd *cobra.Command) error {
	home, err := cmd.PersistentFlags().GetString(cli.HomeFlag)
	if err != nil {
		return err
	}

	cfgFile := path.Join(home, "config", "config.toml")
	if _, err := os.Stat(cfgFile); err == nil {
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}

	if err := viper.BindPFlag(cli.EncodingFlag, cmd.PersistentFlags().Lookup(cli.EncodingFlag)); err != nil {
		return err
	}
	return viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag))
}
