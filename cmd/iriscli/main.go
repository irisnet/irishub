package main

import (
	"fmt"
	"os"
	"path"

	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/client"
	bankcmd "github.com/irisnet/irishub/client/bank/cli"
	distributioncmd "github.com/irisnet/irishub/client/distribution/cli"
	govcmd "github.com/irisnet/irishub/client/gov/cli"
	guardiancmd "github.com/irisnet/irishub/client/guardian/cli"
	keyscmd "github.com/irisnet/irishub/client/keys/cli"
	servicecmd "github.com/irisnet/irishub/client/service/cli"
	slashingcmd "github.com/irisnet/irishub/client/slashing/cli"
	stakecmd "github.com/irisnet/irishub/client/stake/cli"
	tendermintrpccmd "github.com/irisnet/irishub/client/tendermint/rpc"
	tenderminttxcmd "github.com/irisnet/irishub/client/tendermint/tx"
	upgradecmd "github.com/irisnet/irishub/client/upgrade/cli"
	"github.com/irisnet/irishub/client/utils"
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
	defer func() {
		if r := recover(); r != nil {
			switch rType := r.(type) {
			case string:
				println(rType)
			default:
				panic(r)
			}
		}
	}()
	//	sdk.InitBech32Prefix()
	cobra.EnableCommandSorting = false
	cdc := app.MakeLatestCodec()

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
		tendermintrpccmd.BlockResultCommand(),
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
			bankcmd.BurnTxCmd(cdc),
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
			distributioncmd.GetWithdrawAddress(cdc),
			distributioncmd.GetDelegationDistInfo(cdc),
			distributioncmd.GetValidatorDistInfo(cdc),
			distributioncmd.GetAllDelegationDistInfo(cdc),
			distributioncmd.GetRewards(cdc),
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
			stakecmd.GetCmdQueryValidator(cdc),
			stakecmd.GetCmdQueryValidators(cdc),
			stakecmd.GetCmdQueryDelegation(cdc),
			stakecmd.GetCmdQueryDelegations(cdc),
			stakecmd.GetCmdQueryUnbondingDelegation(cdc),
			stakecmd.GetCmdQueryUnbondingDelegations(cdc),
			stakecmd.GetCmdQueryValidatorDelegations(cdc),
			stakecmd.GetCmdQueryValidatorUnbondingDelegations(cdc),
			stakecmd.GetCmdQueryValidatorRedelegations(cdc),
			stakecmd.GetCmdQueryRedelegation(cdc),
			stakecmd.GetCmdQueryRedelegations(cdc),
			stakecmd.GetCmdQueryPool(cdc),
			stakecmd.GetCmdQueryParams(cdc),
			slashingcmd.GetCmdQuerySigningInfo(protocol.SlashingRoute, cdc),
		)...)
	stakeCmd.AddCommand(
		client.PostCommands(
			stakecmd.GetCmdCreateValidator(cdc),
			stakecmd.GetCmdEditValidator(cdc),
			stakecmd.GetCmdDelegate(cdc),
			stakecmd.GetCmdUnbond(cdc),
			stakecmd.GetCmdRedelegate(cdc),
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
			guardiancmd.GetCmdQueryProfilers(cdc),
			guardiancmd.GetCmdQueryTrustees(cdc),
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

	if err := executor.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
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
