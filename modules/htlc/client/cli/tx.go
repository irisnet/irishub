package cli

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/irisnet/irishub/modules/htlc/internal/types"
)

// GetTxCmd returns the transaction commands for the HTLC module.
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	stakingTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "HTLC subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	stakingTxCmd.AddCommand(flags.PostCommands(
		GetCmdCreateHTLC(cdc),
		GetCmdClaimHTLC(cdc),
		GetCmdRefundHTLC(cdc),
	)...)
	return stakingTxCmd
}

// GetCmdCreateHTLC implements the create HTLC command.
func GetCmdCreateHTLC(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an HTLC",
		Example: "iriscli tx htlc create --chain-id=<chain-id> --from=<key-name>" +
			" --fees=0.3iris --to=<to> --receiver-on-other-chain=<receiver-on-other-chain> --amount=<amount>" +
			" --secret=<secret> --timestamp=<timestamp> --time-lock=<time-lock>",
		PreRunE: preCheckCmd,
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			sender := cliCtx.GetFromAddress()

			toAddrStr := viper.GetString(FlagTo)
			toAddr, err := sdk.AccAddressFromBech32(toAddrStr)
			if err != nil {
				return err
			}

			receiverOnOtherChain := viper.GetString(FlagReceiverOnOtherChain)

			amountStr := viper.GetString(FlagAmount)
			amount, err := sdk.ParseCoins(amountStr)
			if err != nil {
				return err
			}

			timestamp := viper.GetInt64(FlagTimestamp)
			timeLock := viper.GetInt64(FlagTimeLock)

			secret := make(types.HTLCSecret, 32)
			var hashLock types.HTLCHashLock

			flags := cmd.Flags()
			if flags.Changed(FlagHashLock) {
				hashLockStr := strings.TrimSpace(viper.GetString(FlagHashLock))
				hashLock, err = hex.DecodeString(hashLockStr)
				if err != nil {
					return err
				}
			} else {
				secretStr := strings.TrimSpace(viper.GetString(FlagSecret))
				if len(secretStr) > 0 {
					if len(secretStr) != 2*types.SecretLength {
						return fmt.Errorf("length of the secret str must be %d", 2*types.SecretLength)
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
				sender,
				toAddr,
				receiverOnOtherChain,
				amount,
				hashLock,
				uint64(timestamp),
				uint64(timeLock),
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			if err := utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}); err == nil && !flags.Changed(FlagHashLock) {
				fmt.Println("**Important** save this secret, hashLock in a safe place.")
				fmt.Println("It is the only way to claim or refund the locked coins from an HTLC")
				fmt.Println()
				fmt.Printf("Secret:      %s\nHashLock:    %s\n", secret.String(), hashLock.String())
			}
			return err
		},
	}

	cmd.Flags().AddFlagSet(FsCreateHTLC)
	_ = cmd.MarkFlagRequired(FlagTo)
	_ = cmd.MarkFlagRequired(FlagAmount)
	_ = cmd.MarkFlagRequired(FlagTimeLock)

	return cmd
}

// GetCmdClaimHTLC implements the claim HTLC command.
func GetCmdClaimHTLC(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "claim",
		Short:   "Claim an opened HTLC",
		Example: "iriscli tx htlc claim --chain-id=<chain-id> --from=<key-name> --fees=0.3iris --hash-lock=<hash-lock> --secret=<secret>",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			sender := cliCtx.GetFromAddress()
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

			msg := types.NewMsgClaimHTLC(sender, hashLock, secret)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsClaimHTLC)
	_ = cmd.MarkFlagRequired(FlagHashLock)
	_ = cmd.MarkFlagRequired(FlagSecret)

	return cmd
}

// GetCmdRefundHTLC implements the refund HTLC command.
func GetCmdRefundHTLC(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "refund",
		Short:   "Refund from an expired HTLC",
		Example: "iriscli tx htlc refund --chain-id=<chain-id> --from=<key-name> --fees=0.3iris --hash-lock=<hash-lock>",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			sender := cliCtx.GetFromAddress()
			hashLockStr := viper.GetString(FlagHashLock)
			hashLock, err := hex.DecodeString(hashLockStr)
			if err != nil {
				return err
			}

			msg := types.NewMsgRefundHTLC(sender, hashLock)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsRefundHTLC)
	_ = cmd.MarkFlagRequired(FlagHashLock)

	return cmd
}

func preCheckCmd(cmd *cobra.Command, _ []string) error {
	// make sure either the secret or hash lock is provided
	flags := cmd.Flags()
	if flags.Changed(FlagSecret) && flags.Changed(FlagHashLock) {
		return fmt.Errorf("only one flag is allowed among the secret and hash lock")
	}

	return nil
}
