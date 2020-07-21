package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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
		Example: fmt.Sprintf("%s tx rand request-rand [--block-interval=10] [--oracle=true --service-fee-cap=1iris]", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			consumer := clientCtx.GetFromAddress()

			oracle := viper.GetBool(FlagOracle)

			var serviceFeeCap sdk.Coins
			if oracle {
				if serviceFeeCap, err = sdk.ParseCoins(viper.GetString(FlagServiceFeeCap)); err != nil {
					return err
				}
			}

			msg := types.NewMsgRequestRandom(consumer, uint64(viper.GetInt64(FlagBlockInterval)), oracle, serviceFeeCap)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTx(clientCtx, msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
