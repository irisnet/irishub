package main

import (
	"os"

	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/app"
	irisInit "github.com/irisnet/irishub/init"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
	tmlite "github.com/tendermint/tendermint/lite"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	tmliteProxy "github.com/tendermint/tendermint/lite/proxy"
)

func getTrustBasis(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-trust-basis",
		Short: "get the trust basis",
		RunE: func(cmd *cobra.Command, args []string) error {
			chainID := viper.GetString(client.FlagChainID)
			home := viper.GetString(cli.HomeFlag)
			nodeURI := viper.GetString(client.FlagNode)
			node := rpcclient.NewHTTP(nodeURI, "/websocket")
			verifier, err := tmliteProxy.NewVerifier(
				chainID, home,
				node, log.NewNopLogger(), 10,
			)
			return nil
		},
	}
	cmd.Flags().Bool(client.FlagTrustNode, true, "Don't verify proofs for responses")
	cmd.Flags().String(client.FlagChainID, "", "Chain ID of tendermint node")
	cmd.Flags().String(client.FlagNode, "tcp://localhost:26657", "Address of the node to connect to")
	return cmd
}

func verifyState(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify [state file]",
		Short: "verify exported state",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout)
			return nil
		},
	}
	return cmd
}

func main() {
	irisInit.InitBech32Prefix()
	cdc := app.MakeCodec()

	rootCmd := &cobra.Command{
		Use:          "irisverifier",
		Short:        "Iris verifier tool for exported state ",
		SilenceUsage: true,
	}

	rootCmd.AddCommand(getTrustBasis(cdc))
	rootCmd.AddCommand(verifyState(cdc))

	executor := cli.PrepareMainCmd(rootCmd, "IRISVERIFIER", os.ExpandEnv("$HOME/.irisverifier"))
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}
