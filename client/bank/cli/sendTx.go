package cli

import (
	"fmt"
	"os"

	"github.com/irisnet/irishub/client/bank"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagTo     = "to"
	flagAmount = "amount"
	flagHolder = "holder"
)

// SendTxCmd will create a send tx and sign it with the given key.
func SendTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "send",
		Short:   "Create and sign a send tx",
		Example: "iriscli bank send --to=<account address> --from <key name> --fee=0.4iris --chain-id=<chain-id> --amount=10iris",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).WithCliCtx(cliCtx)

			toStr := viper.GetString(flagTo)

			to, err := sdk.AccAddressFromBech32(toStr)
			if err != nil {
				return err
			}

			// parse coins trying to be sent
			amount := viper.GetString(flagAmount)
			coins, err := cliCtx.ParseCoins(amount)
			if err != nil {
				return err
			}

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := bank.BuildBankSendMsg(from, to, coins)
			if cliCtx.GenerateOnly {
				return utils.PrintUnsignedStdTx(txCtx, cliCtx, []sdk.Msg{msg}, true)
			}

			account, err := cliCtx.GetAccount(from)
			if err != nil {
				return err
			}

			// ensure account has enough coins
			if !account.GetCoins().IsAllGTE(coins) {
				return fmt.Errorf("Address %s doesn't have enough coins to pay for this transaction.", from)
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagTo, "", "Bech32 encoding address to receive coins")
	cmd.Flags().String(flagAmount, "", "Amount of coins to send, for instance: 10iris")
	cmd.MarkFlagRequired(flagTo)
	cmd.MarkFlagRequired(flagAmount)

	return cmd
}

// BurnTxCmd will create a burn token tx and sign it with the given key.
func BurnTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "burn",
		Short:   "Create and sign a tx to burn coins",
		Example: "iriscli bank burn --from <key name> --fee=0.4iris --chain-id=<chain-id> --amount=10iris",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).WithCliCtx(cliCtx)

			// parse coins trying to be sent
			amount := viper.GetString(flagAmount)
			coins, err := cliCtx.ParseCoins(amount)
			if err != nil {
				return err
			}

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}
			// build and sign the transaction, then broadcast to Tendermint
			msg := bank.BuildBankBurnMsg(from, coins)
			if cliCtx.GenerateOnly {
				return utils.PrintUnsignedStdTx(txCtx, cliCtx, []sdk.Msg{msg}, true)
			}

			//account, err := cliCtx.GetAccount(from)
			//if err != nil {
			//	return err
			//}
			//
			//// ensure account has enough coins
			//if !account.GetCoins().IsAllGTE(coins) {
			//	return fmt.Errorf("Address %s doesn't have enough coins to pay for this transaction.", from)
			//}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagAmount, "", "Amount of coins to burn, for instance: 10iris")
	cmd.MarkFlagRequired(flagAmount)

	return cmd
}

// FreezeTxCmd will freeze specific token.
func FreezeTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "freeze",
		Short:   "Create and sign a tx to freeze coins",
		Example: "iriscli bank freeze --amount=10iris --holder=<account address>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).WithCliCtx(cliCtx)

			// parse coins trying to be sent
			amount := viper.GetString(flagAmount)
			coin, err := cliCtx.ParseCoin(amount)
			if err != nil {
				return err
			}

			owner, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := bank.BuildBankFreezeMsg(owner, coin)
			if cliCtx.GenerateOnly {
				return utils.PrintUnsignedStdTx(txCtx, cliCtx, []sdk.Msg{msg}, true)
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagAmount, "", "Amount of coins to freeeze, for instance: 10iris")
	cmd.MarkFlagRequired(flagAmount)

	return cmd
}

// UnfreezeTxCmd will unfreeze specific token.
func UnfreezeTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "unfreeze",
		Short:   "Create and sign a tx to unfreeze coins",
		Example: "iriscli bank unfreeze --amount=10iris",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).WithCliCtx(cliCtx)

			// parse coins trying to be sent
			amount := viper.GetString(flagAmount)
			coin, err := cliCtx.ParseCoin(amount)
			if err != nil {
				return err
			}
			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}
			// build and sign the transaction, then broadcast to Tendermint
			msg := bank.BuildBankUnfreezeMsg(from,coin)
			if cliCtx.GenerateOnly {
				return utils.PrintUnsignedStdTx(txCtx, cliCtx, []sdk.Msg{msg}, true)
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagAmount, "", "Amount of coins to freeeze, for instance: 10iris")
	cmd.Flags().String(flagHolder, "", "address of token-holder, only asset-owner can pass this param")
	cmd.MarkFlagRequired(flagAmount)

	return cmd
}
