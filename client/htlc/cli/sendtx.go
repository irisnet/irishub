package cli

import (
	"encoding/hex"
	"os"

	"github.com/irisnet/irishub/app/v2/htlc"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetCmdCreateHtlc implements the create htlc command
func GetCmdCreateHtlc(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a HTLC",
		Example: "iriscli htlc create --receiver=<receiver> --receiver-on-other-chain=<receiver-on-other-chain> --hash-lock=<hash-lock> " +
			"--in-amount=<in-amount> --amount=<amount> --time-lock=<time-lock>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			sender, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			receiverStr := viper.GetString(FlagReceiver)
			receiver, err := sdk.AccAddressFromBech32(receiverStr)
			if err != nil {
				return err
			}

			receiverOnOtherChainStr := viper.GetString(FlagReceiverOnOtherChain)

			receiverOnOtherChain, err := hex.DecodeString(receiverOnOtherChainStr)
			if err != nil {
				return err
			}

			inAmount := viper.GetInt64(FlagInAmount)
			amountStr := viper.GetString(FlagAmount)
			coin, err := cliCtx.ParseCoin(amountStr)
			if err != nil {
				return err
			}

			hashLock := viper.GetString(FlagHashLock)
			timestamp := viper.GetInt64(FlagTimestamp)
			timeLock := viper.GetInt64(FlagTimeLock)

			var msg sdk.Msg
			msg = htlc.NewMsgCreateHTLC(
				sender, receiver, receiverOnOtherChain, coin, uint64(inAmount),
				hashLock, uint64(timestamp), uint64(timeLock))

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsCreateHTLC)
	cmd.MarkFlagRequired(FlagReceiver)
	cmd.MarkFlagRequired(FlagReceiverOnOtherChain)
	cmd.MarkFlagRequired(FlagHashLock)
	cmd.MarkFlagRequired(FlagInAmount)
	cmd.MarkFlagRequired(FlagAmount)
	cmd.MarkFlagRequired(FlagTimeLock)
	cmd.MarkFlagRequired(FlagTimestamp)
	cmd.MarkFlagRequired(FlagSecret)

	return cmd
}
