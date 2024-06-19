package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"irismod.io/mt/types"
)

// NewTxCmd returns the transaction commands for this module
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "MT transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdIssueDenom(),
		GetCmdTransferDenom(),
		GetCmdMintMT(),
		GetCmdEditMT(),
		GetCmdTransferMT(),
		GetCmdBurnMT(),
	)

	return txCmd
}

// GetCmdIssueDenom is the CLI command for an SaveDenom transaction
func GetCmdIssueDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "issue",
		Long: "Issue a new denom.",
		Example: fmt.Sprintf(
			"$ %s tx mt issue "+
				"--name=<denom-name> "+
				"--data=<data> "+
				"--from=<key-name> "+
				"--chain-id=<chain-id> "+
				"--fees=<fee>",
			version.AppName,
		),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var sender = clientCtx.GetFromAddress().String()

			name, err := cmd.Flags().GetString(FlagName)
			if err != nil {
				return err
			}
			data, err := cmd.Flags().GetString(FlagData)
			if err != nil {
				return err
			}

			msg := types.NewMsgIssueDenom(
				name,
				data,
				sender,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsIssueDenom)
	_ = cmd.MarkFlagRequired(FlagName)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdTransferDenom is the CLI command for sending a TransferDenom transaction
func GetCmdTransferDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "transfer-denom [from_key_or_address] [recipient] [denom-id]",
		Long: "Transfer a denom to a recipient.",
		Example: fmt.Sprintf(
			"$ %s tx mt transfer-denom <from_key_or_address> <recipient> <denom-id> "+
				"--chain-id=<chain-id> "+
				"--fees=<fee>",
			version.AppName,
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[0])
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			recipient, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			denomID := args[2]

			msg := types.NewMsgTransferDenom(
				denomID,
				clientCtx.GetFromAddress().String(),
				recipient.String(),
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsTransferDenom)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdMintMT is the CLI command for a MintMT transaction
func GetCmdMintMT() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "mint [denom-id]",
		Long: "Issue or mint MT",
		Example: fmt.Sprintf(
			"$ %s tx mt mint <denom-id> "+
				"--mt-id=<mt-id> "+
				"--amount=<amount> "+
				"--data=<data> "+
				"--recipient=<recipient> "+
				"--from=<key-name> "+
				"--chain-id=<chain-id> "+
				"--fees=<fee>",
			version.AppName,
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			denomID := args[0]

			mtID, err := cmd.Flags().GetString(FlagMTID)
			if err != nil {
				return err
			}

			amountStr, err := cmd.Flags().GetString(FlagAmount)
			if err != nil {
				return err
			}

			amount, err := strconv.ParseUint(amountStr, 10, 64)
			if err != nil {
				return err
			}

			metadata, err := cmd.Flags().GetString(FlagData)
			if err != nil {
				return err
			}

			var sender = clientCtx.GetFromAddress().String()

			recipient, err := cmd.Flags().GetString(FlagRecipient)
			if err != nil {
				return err
			}

			recipientStr := strings.TrimSpace(recipient)
			if len(recipientStr) > 0 {
				if _, err = sdk.AccAddressFromBech32(recipientStr); err != nil {
					return err
				}
			} else {
				recipient = sender
			}

			msg := types.NewMsgMintMT(
				mtID,
				denomID,
				amount,
				metadata,
				sender,
				recipient,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsMintMT)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdEditMT is the CLI command for sending an MsgEditMT transaction
func GetCmdEditMT() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "edit [denom-id] [mt-id]",
		Long: "Edit the metadata of an MT.",
		Example: fmt.Sprintf(
			"$ %s tx mt edit <denom-id> <mt-id> "+
				"--data=<data> "+
				"--from=<key-name> "+
				"--chain-id=<chain-id> "+
				"--fees=<fee>",
			version.AppName,
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			denomID := args[0]
			mtID := args[1]
			metadata, err := cmd.Flags().GetString(FlagData)
			if err != nil {
				return err
			}

			msg := types.NewMsgEditMT(
				mtID,
				denomID,
				metadata,
				clientCtx.GetFromAddress().String(),
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsEditMT)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdTransferMT is the CLI command for sending a TransferMT transaction
func GetCmdTransferMT() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "transfer [from_key_or_address] [recipient] [denom-id] [mt-id] [amount]",
		Long: "Transfer an MT to a recipient.",
		Example: fmt.Sprintf(
			"$ %s tx mt transfer <from_key_or_address> <recipient> <denom-id> <mt-id> <amount> "+
				"--chain-id=<chain-id> "+
				"--fees=<fee>",
			version.AppName,
		),
		Args: cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[0])
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			recipient, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			denomID := args[2]
			mtID := args[3]
			amountStr := args[4]
			amount, err := strconv.ParseUint(amountStr, 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgTransferMT(
				mtID,
				denomID,
				clientCtx.GetFromAddress().String(),
				recipient.String(),
				amount,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsTransferMT)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdBurnMT is the CLI command for sending a BurnMT transaction
func GetCmdBurnMT() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "burn [denom-id] [mt-id] [amount]",
		Long: "Burn amounts of an MT.",
		Example: fmt.Sprintf(
			"$ %s tx mt burn <denom-id> <mt-id> <amount> "+
				"--from=<key-name> "+
				"--chain-id=<chain-id> "+
				"--fees=<fee>",
			version.AppName,
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			denomID := args[0]
			mtID := args[1]
			amountStr := args[2]
			amount, err := strconv.ParseUint(amountStr, 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgBurnMT(
				clientCtx.GetFromAddress().String(),
				mtID,
				denomID,
				amount,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
