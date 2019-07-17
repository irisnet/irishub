package cli

import (
	"fmt"
	"github.com/irisnet/irishub/app/v1/coinswap"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"time"
)

const (
	flagDepositToken = "deposit-token"
	flagNativeToken  = "native-token"
)

// SendTxCmd will create a send tx and sign it with the given key.
func AddLiquidity(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add-liquidity",
		Short:   "Create and sign a send tx",
		Example: "iriscli bank send --to=<account address> --from <key name> --fee=0.4iris --chain-id=<chain-id> --amount=10iris",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).WithCliCtx(cliCtx)

			// parse coins trying to be sent
			depositTokenStr := viper.GetString(flagDepositToken)
			depositToken, err := cliCtx.ParseCoin(depositTokenStr)
			if err != nil {
				return err
			}

			nativeTokenStr := viper.GetString(flagNativeToken)
			nativeToken, err := cliCtx.ParseCoin(nativeTokenStr)
			if err != nil {
				return err
			}

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			// TODO deadline
			msg := coinswap.NewMsgAddLiquidity(depositToken, nativeToken.Amount, sdk.ZeroInt(), time.Now(), from)
			if cliCtx.GenerateOnly {
				return utils.PrintUnsignedStdTx(txCtx, cliCtx, []sdk.Msg{msg}, true)
			}

			account, err := cliCtx.GetAccount(from)
			if err != nil {
				return err
			}

			// ensure account has enough coins
			if !account.GetCoins().IsAllGTE(sdk.Coins{depositToken, nativeToken}) {
				return fmt.Errorf("Address %s doesn't have enough coins to pay for this transaction.", from)
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagDepositToken, "", "Amount of coins to send, for instance: 10iris")
	cmd.Flags().String(flagNativeToken, "", "Amount of coins to send, for instance: 10iris")
	cmd.MarkFlagRequired(flagDepositToken)
	cmd.MarkFlagRequired(flagNativeToken)

	return cmd
}
