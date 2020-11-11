package cli

import (
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/irisnet/irismod/modules/htlc/types"
)

// NewTxCmd returns the transaction commands for this module
func NewTxCmd() *cobra.Command {
	htlcTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "HTLC transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	htlcTxCmd.AddCommand(
		GetCmdCreateHTLC(),
		GetCmdClaimHTLC(),
		GetCmdRefundHTLC(),
	)

	return htlcTxCmd
}

// GetCmdCreateHTLC implements creating an HTLC command
func GetCmdCreateHTLC() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an HTLC",
		Long:  "Create an HTLC.",
		Example: fmt.Sprintf(
			"$ %s tx htlc create "+
				"--to=<recipient> "+
				"--receiver-on-other-chain=<receiver-on-other-chain> "+
				"--amount=<amount> "+
				"--hash-lock=<hash-lock> "+
				"--timestamp=<timestamp> "+
				"--time-lock=<time-lock> "+
				"--from=mykey",
			version.AppName,
		),
		PreRunE: preCheckCmd,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			sender := clientCtx.GetFromAddress()

			toAddr, err := cmd.Flags().GetString(FlagTo)
			if err != nil {
				return err
			}

			if _, err := sdk.AccAddressFromBech32(toAddr); err != nil {
				return err
			}

			receiverOnOtherChain, err := cmd.Flags().GetString(FlagReceiverOnOtherChain)
			if err != nil {
				return err
			}
			amountStr, err := cmd.Flags().GetString(FlagAmount)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoins(amountStr)
			if err != nil {
				return err
			}

			timestamp, err := cmd.Flags().GetUint64(FlagTimestamp)
			if err != nil {
				return err
			}
			timeLock, err := cmd.Flags().GetUint64(FlagTimeLock)
			if err != nil {
				return err
			}

			secret := make([]byte, 32)
			var hashLock []byte

			flags := cmd.Flags()
			if flags.Changed(FlagHashLock) {
				hashLock, err = cmd.Flags().GetBytesHex(FlagHashLock)
				if err != nil {
					return err
				}
			} else {
				secret, err = cmd.Flags().GetBytesHex(FlagSecret)
				if err != nil {
					return err
				}
				hashLock = types.GetHashLock(secret, timestamp)
			}

			msg := types.NewMsgCreateHTLC(
				sender.String(), toAddr, receiverOnOtherChain, amount,
				hex.EncodeToString(hashLock), timestamp, timeLock,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			if err = tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg); err == nil && !flags.Changed(FlagHashLock) {
				fmt.Println("**Important** save this secret, hashLock in a safe place.")
				fmt.Println("It is the only way to claim or refund the locked coins from an HTLC")
				fmt.Println()
				fmt.Printf("Secret:      %s\nHashLock:    %s\n",
					hex.EncodeToString(secret), hex.EncodeToString(hashLock),
				)
			}

			return err
		},
	}

	cmd.Flags().AddFlagSet(FsCreateHTLC)
	_ = cmd.MarkFlagRequired(FlagTo)
	_ = cmd.MarkFlagRequired(FlagAmount)
	_ = cmd.MarkFlagRequired(FlagTimeLock)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdClaimHTLC implements claiming an HTLC command
func GetCmdClaimHTLC() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "claim [hash-lock] [secret]",
		Short:   "Claim an HTLC",
		Long:    "Claim an open HTLC with a secret.",
		Example: fmt.Sprintf("$ %s tx htlc claim <hash-lock> <secret> --from mykey", version.AppName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			sender := clientCtx.GetFromAddress().String()

			if _, err := hex.DecodeString(args[0]); err != nil {
				return err
			}

			if _, err := hex.DecodeString(args[1]); err != nil {
				return err
			}

			msg := types.NewMsgClaimHTLC(sender, args[0], args[1])
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdRefundHTLC implements refunding an HTLC command
func GetCmdRefundHTLC() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "refund [hash-lock]",
		Short:   "Refund an HTLC",
		Long:    "Refund from an expired HTLC.",
		Example: fmt.Sprintf("$ %s tx htlc refund <hash-lock> --from mykey", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			sender := clientCtx.GetFromAddress().String()

			if _, err := hex.DecodeString(args[0]); err != nil {
				return err
			}

			msg := types.NewMsgRefundHTLC(sender, args[0])
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func preCheckCmd(cmd *cobra.Command, _ []string) error {
	flags := cmd.Flags()

	if flags.Changed(FlagSecret) && flags.Changed(FlagHashLock) {
		return fmt.Errorf("can not provide both the secret and hash lock")
	}

	return nil
}
