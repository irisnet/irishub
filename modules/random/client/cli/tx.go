package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/irisnet/irishub/modules/random/types"
)

// GetTxCmd returns the transaction commands for the rand module.
func GetTxCmd(clientCtx client.Context) *cobra.Command {
	randTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Random transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	randTxCmd.AddCommand(
		GetCmdRequestRandom(clientCtx),
	)
	return randTxCmd
}

// GetCmdRequestRandom implements the request-rand command.
func GetCmdRequestRandom(clientCtx client.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "request-rand",
		Short:   "Request a random number with an optional block interval",
		Example: "iriscli tx rand request-rand [block-interval]",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			var blockInterval uint64

			if len(args) > 0 {
				blockInterval, err = strconv.ParseUint(args[0], 10, 64)
				if err != nil {
					return err
				}
			}

			consumer := clientCtx.GetFromAddress()

			msg := types.NewMsgRequestRandom(consumer, blockInterval)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTx(clientCtx, msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
