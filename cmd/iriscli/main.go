package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/irisnet/irishub/address"
	"github.com/irisnet/irishub/app"
)

var (
	encodingConfig = app.MakeEncodingConfig()
	initClientCtx  = client.Context{}.
			WithJSONMarshaler(encodingConfig.Marshaler).
			WithTxGenerator(encodingConfig.TxGenerator).
			WithCodec(encodingConfig.Amino).
			WithInput(os.Stdin).
			WithAccountRetriever(types.NewAccountRetriever(encodingConfig.Marshaler)).
			WithBroadcastMode(flags.BroadcastBlock)
)

func init() {
	authclient.Codec = encodingConfig.Marshaler
}

// TODO: setup keybase, viper object, etc. to be passed into
// the below functions and eliminate global vars, like we do
// with the cdc
func main() {
	cobra.EnableCommandSorting = false

	// Read in the configuration file for the sdk
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(address.Bech32PrefixAccAddr, address.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(address.Bech32PrefixValAddr, address.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(address.Bech32PrefixConsAddr, address.Bech32PrefixConsPub)
	config.Seal()

	rootCmd := &cobra.Command{
		Use:   "iriscli",
		Short: "irishub command line interface",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}
			return server.InterceptConfigsPreRunHandler(cmd)
		},
	}

	rootCmd.AddCommand(
		rpc.StatusCommand(),
		queryCmd(),
		txCmd(),
		flags.LineBreak,
		flags.LineBreak,
		keys.Commands(),
		flags.LineBreak,
		version.NewVersionCommand(),
		cli.NewCompletionCmd(rootCmd, true),
	)

	// Add flags and prefix all env exposed with GA
	executor := cli.PrepareMainCmd(rootCmd, "GA", app.DefaultCLIHome)

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &client.Context{})
	ctx = context.WithValue(ctx, server.ServerContextKey, server.NewDefaultContext())

	if err := executor.ExecuteContext(ctx); err != nil {
		fmt.Printf("failed execution: %s, exiting...\n", err)
		os.Exit(1)
	}
}

func queryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		authcmd.GetAccountCmd(encodingConfig.Amino),
		flags.LineBreak,
		rpc.ValidatorCommand(encodingConfig.Amino),
		rpc.BlockCommand(),
		authcmd.QueryTxsByEventsCmd(encodingConfig.Amino),
		authcmd.QueryTxCmd(encodingConfig.Amino),
		flags.LineBreak,
	)

	app.ModuleBasics.AddQueryCommands(queryCmd, initClientCtx)

	return queryCmd
}

func txCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		bankcmd.NewSendTxCmd(),
		flags.LineBreak,
		authcmd.GetSignCommand(initClientCtx),
		authcmd.GetSignBatchCommand(encodingConfig.Amino),
		authcmd.GetMultiSignCommand(initClientCtx),
		authcmd.GetValidateSignaturesCommand(initClientCtx),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(initClientCtx),
		authcmd.GetEncodeCommand(initClientCtx),
		authcmd.GetDecodeCommand(initClientCtx),
		flags.LineBreak,
	)

	app.ModuleBasics.AddTxCommands(txCmd, initClientCtx)

	return txCmd
}
