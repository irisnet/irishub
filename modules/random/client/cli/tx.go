package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/irisnet/irismod/modules/random/types"
)

// NewTxCmd returns the transaction commands for the random module.
func NewTxCmd() *cobra.Command {
	randTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Random transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	randTxCmd.AddCommand(
		GetCmdRequestRandom(),
	)
	return randTxCmd
}

// GetCmdRequestRandom implements the request-random command.
func GetCmdRequestRandom() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request",
		Short: "Request a random number with an optional block interval",
		Example: fmt.Sprintf(
			"%s tx random request [--block-interval=10] [--oracle=true --service-fee-cap=1iris]",
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			consumer := clientCtx.GetFromAddress().String()

			var serviceFeeCap sdk.Coins
			oracle, err := cmd.Flags().GetBool(FlagOracle)
			if err != nil {
				return err
			}
			rawServiceFeeCap, err := cmd.Flags().GetString(FlagServiceFeeCap)
			if err != nil {
				return err
			}
			if oracle {
				if serviceFeeCap, err = sdk.ParseCoins(rawServiceFeeCap); err != nil {
					return err
				}
			}

			blockInterval, err := cmd.Flags().GetUint64(FlagBlockInterval)
			if err != nil {
				return err
			}
			msg := types.NewMsgRequestRandom(consumer, blockInterval, oracle, serviceFeeCap)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsRequestRand)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
