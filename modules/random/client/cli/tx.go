package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/irisnet/irishub/modules/random/types"
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
		Use:   "request-random",
		Short: "Request a random number with an optional block interval",
		Example: fmt.Sprintf(
			"%s tx random request-random [--block-interval=10] [--oracle=true --service-fee-cap=1iris]",
			version.AppName,
		),
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

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsRequestRand)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
