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

// GetCmdCreateHtlc implements the create HTLC command
func GetCmdCreateHtlc(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an HTLC",
		Example: "iriscli htlc create --receiver=<receiver> --receiver-on-other-chain=<receiver-on-other-chain> --amount=<amount> --hash-lock=<hash-lock> " +
			"--time-lock=<time-lock> --timestamp=<timestamp>",
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

			amountStr := viper.GetString(FlagAmount)
			coin, err := cliCtx.ParseCoin(amountStr)
			if err != nil {
				return err
			}

			hashLockStr := viper.GetString(FlagHashLock)
			hashLock, err := hex.DecodeString(hashLockStr)
			if err != nil {
				return err
			}

			timestamp := viper.GetInt64(FlagTimestamp)
			timeLock := viper.GetInt64(FlagTimeLock)

			msg := htlc.NewMsgCreateHTLC(
				sender, receiver, receiverOnOtherChain, coin,
				hashLock, uint64(timestamp), uint64(timeLock))

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsCreateHTLC)
	_ = cmd.MarkFlagRequired(FlagReceiver)
	_ = cmd.MarkFlagRequired(FlagReceiverOnOtherChain)
	_ = cmd.MarkFlagRequired(FlagAmount)
	_ = cmd.MarkFlagRequired(FlagHashLock)
	_ = cmd.MarkFlagRequired(FlagTimeLock)

	return cmd
}

// GetCmdClaimHtlc implements the claim HTLC command
func GetCmdClaimHtlc(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "claim",
		Short:   "Claim an opened HTLC",
		Example: "iriscli htlc claim --hash-lock=<hash-lock> --secret=<secret>",
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

			hashLockStr := viper.GetString(FlagHashLock)
			hashLock, err := hex.DecodeString(hashLockStr)
			if err != nil {
				return err
			}

			secretStr := viper.GetString(FlagSecret)
			secret, err := hex.DecodeString(secretStr)
			if err != nil {
				return err
			}

			msg := htlc.NewMsgClaimHTLC(sender, hashLock, secret)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsClaimHTLC)
	_ = cmd.MarkFlagRequired(FlagHashLock)
	_ = cmd.MarkFlagRequired(FlagSecret)

	return cmd
}

// GetCmdRefundHtlc implements the refund HTLC command
func GetCmdRefundHtlc(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "refund",
		Short:   "Refund from an expired HTLC",
		Example: "iriscli htlc refund --hash-lock=<hash-lock>",
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

			hashLockStr := viper.GetString(FlagHashLock)
			hashLock, err := hex.DecodeString(hashLockStr)
			if err != nil {
				return err
			}

			msg := htlc.NewMsgRefundHTLC(
				sender, hashLock)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsRefundHTLC)
	_ = cmd.MarkFlagRequired(FlagHashLock)

	return cmd
}
