package cli

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/codec"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"github.com/irisnet/irishub/client/bank"
)

const (
	flagTo     = "to"
	flagAmount = "amount"
)

// SendTxCmd will create a send tx and sign it with the given key.
func SendTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "Create and sign a send tx",
		Example: "iriscli bank send --to=<account address> --from <key name> --fee=0.004iris --chain-id=<chain-id> --amount=10iris",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			txCtx := context.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			toStr := viper.GetString(flagTo)

			to, err := sdk.AccAddressFromBech32(toStr)
			if err != nil {
				return err
			}

			// parse coins trying to be sent
			amountString := viper.GetString(flagAmount)

			coins, err := cliCtx.ParseCoins(amountString)
			if err != nil {
				return err
			}

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			account, err := cliCtx.GetAccount(from)
			if err != nil {
				return err
			}

			// ensure account has enough coins
			if !account.GetCoins().IsGTE(coins) {
				return fmt.Errorf("address %s doesn't have enough coins to pay for this transaction", from)
			}

			msg := bank.BuildMsg(from, to, coins)

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagTo, "", "Bech32 encoding address to receive coins")
	cmd.Flags().String(flagAmount, "", "Amount of coins to send, for instance: 10iris")

	return cmd
}
