package main

import (
	"os"
	"path"

	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/client"
	bankcmd "github.com/irisnet/irishub/client/bank/cli"
	distributioncmd "github.com/irisnet/irishub/client/distribution/cli"
	govcmd "github.com/irisnet/irishub/client/gov/cli"
	guardiancmd "github.com/irisnet/irishub/client/guardian/cli"
	keyscmd "github.com/irisnet/irishub/client/keys/cli"
	recordcmd "github.com/irisnet/irishub/client/record/cli"
	servicecmd "github.com/irisnet/irishub/client/service/cli"
	slashingcmd "github.com/irisnet/irishub/client/slashing/cli"
	stakecmd "github.com/irisnet/irishub/client/stake/cli"
	tendermintrpccmd "github.com/irisnet/irishub/client/tendermint/rpc"
	tenderminttxcmd "github.com/irisnet/irishub/client/tendermint/tx"
	upgradecmd "github.com/irisnet/irishub/client/upgrade/cli"
	"github.com/irisnet/irishub/client/utils"
	sdk "github.com/irisnet/irishub/types"
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

	sdk.InitBech32Prefix()

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
			bankcmd.GetAccountCmd("acc", cdc, utils.GetAccountDecoder(cdc)),
			bankcmd.GetCmdQueryTokenStats(cdc, "acc", "stake"),
		)...)
	bankCmd.AddCommand(
		client.PostCommands(
			bankcmd.SendTxCmd(cdc),
			bankcmd.GetSignCommand(cdc, utils.GetAccountDecoder(cdc)),
			bankcmd.GetBroadcastCommand(cdc),
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
			distributioncmd.GetDelegationDistInfo("distr", cdc),
			distributioncmd.GetValidatorDistInfo("distr", cdc),
			distributioncmd.GetAllDelegationDistInfo("distr", cdc),
		)...)
	distributionCmd.AddCommand(
		client.PostCommands(
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
			stakecmd.GetCmdQueryValidatorDelegations("stake", cdc),
			stakecmd.GetCmdQueryValidatorUnbondingDelegations("stake", cdc),
			stakecmd.GetCmdQueryValidatorRedelegations("stake", cdc),
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
			upgradecmd.GetCmdQuerySignals("upgrade", cdc),
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
			servicecmd.GetCmdQuerySvcDef("service", cdc),
			servicecmd.GetCmdQuerySvcBind("service", cdc),
			servicecmd.GetCmdQuerySvcBinds("service", cdc),
			servicecmd.GetCmdQuerySvcRequests("service", cdc),
			servicecmd.GetCmdQuerySvcResponse("service", cdc),
			servicecmd.GetCmdQuerySvcFees("service", cdc),
		)...)
	serviceCmd.AddCommand(client.PostCommands(
		servicecmd.GetCmdSvcDef(cdc),
		servicecmd.GetCmdSvcBind(cdc),
		servicecmd.GetCmdSvcBindUpdate(cdc),
		servicecmd.GetCmdSvcDisable(cdc),
		servicecmd.GetCmdSvcEnable(cdc),
		servicecmd.GetCmdSvcRefundDeposit(cdc),
		servicecmd.GetCmdSvcCall(cdc),
		servicecmd.GetCmdSvcRespond(cdc),
		servicecmd.GetCmdSvcRefundFees(cdc),
		servicecmd.GetCmdSvcWithdrawFees(cdc),
		servicecmd.GetCmdSvcWithdrawTax(cdc),
	)...)

	rootCmd.AddCommand(
		serviceCmd,
	)

	//add guardian command
	guardianCmd := &cobra.Command{
		Use:   "guardian",
		Short: "Guardian subcommands",
	}
	guardianCmd.AddCommand(
		client.GetCommands(
			guardiancmd.GetCmdQueryProfilers("guardian", cdc),
			guardiancmd.GetCmdQueryTrustees("guardian", cdc),
		)...)

	guardianCmd.AddCommand(
		client.PostCommands(
			guardiancmd.GetCmdCreateProfiler(cdc),
			guardiancmd.GetCmdDeleteProfiler(cdc),
			guardiancmd.GetCmdCreateTrustee(cdc),
			guardiancmd.GetCmdDeleteTrustee(cdc),
		)...)
	rootCmd.AddCommand(
		guardianCmd,
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
