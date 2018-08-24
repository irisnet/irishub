package bank

import (
	"os"

	"github.com/irisnet/irishub/client/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/bank/client"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/irisnet/irishub/app"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := app.NewContext().WithCodeC(cdc)

			cliCtx := ctx.WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			ctx = ctx.WithCLIContext(cliCtx)

			if err := ctx.EnsureAccountExists(); err != nil {
				return err
			}

			toStr := viper.GetString(flagTo)

			to, err := sdk.AccAddressFromBech32(toStr)
			if err != nil {
				return err
			}

			// parse coins trying to be sent
			amount := viper.GetString(flagAmount)
			coins, err := ctx.ParseCoins(amount)
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
				return errors.Errorf("Address %s doesn't have enough coins to pay for this transaction.", from)
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := client.BuildMsg(from, to, coins)

			return utils.SendTx(ctx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagTo, "", "Address to send coins")
	cmd.Flags().String(flagAmount, "", "Amount of coins to send")

	return cmd
}
