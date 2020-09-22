package cli

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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
				"--secret=<secret> "+
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

			to, err := sdk.AccAddressFromBech32(viper.GetString(FlagTo))
			if err != nil {
				return err
			}

			receiverOnOtherChain := viper.GetString(FlagReceiverOnOtherChain)

			amount, err := sdk.ParseCoins(viper.GetString(FlagAmount))
			if err != nil {
				return err
			}

			timestamp := viper.GetInt64(FlagTimestamp)
			timeLock := viper.GetInt64(FlagTimeLock)

			secret := make([]byte, 32)
			var hashLock []byte

			flags := cmd.Flags()
			if flags.Changed(FlagHashLock) {
				hashLockStr := strings.TrimSpace(viper.GetString(FlagHashLock))
				if hashLock, err = hex.DecodeString(hashLockStr); err != nil {
					return err
				}
			} else {
				secretStr := strings.TrimSpace(viper.GetString(FlagSecret))
				if len(secretStr) > 0 {
					if len(secretStr) != 2*types.SecretLength {
						return fmt.Errorf("length of the secret must be %d in bytes", types.SecretLength)
					}

					if secret, err = hex.DecodeString(secretStr); err != nil {
						return err
					}
				} else {
					if _, err := rand.Read(secret); err != nil {
						return err
					}
				}

				hashLock = types.GetHashLock(secret, uint64(timestamp))
			}

			msg := types.NewMsgCreateHTLC(
				sender, to, receiverOnOtherChain, amount,
				hashLock, uint64(timestamp), uint64(timeLock),
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

			sender := clientCtx.GetFromAddress()

			hashLock, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			secret, err := hex.DecodeString(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgClaimHTLC(sender, hashLock, secret)
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

			sender := clientCtx.GetFromAddress()

			hashLock, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgRefundHTLC(sender, hashLock)
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
