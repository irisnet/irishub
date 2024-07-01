package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"mods.irisnet.org/modules/nft/types"
)

// NewTxCmd returns the transaction commands for this module
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "NFT transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdIssueDenom(),
		GetCmdMintNFT(),
		GetCmdEditNFT(),
		GetCmdTransferNFT(),
		GetCmdBurnNFT(),
		GetCmdTransferDenom(),
	)

	return txCmd
}

// GetCmdIssueDenom is the CLI command for an SaveDenom transaction
func GetCmdIssueDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "issue [denom-id]",
		Long: "Issue a new denom.",
		Example: fmt.Sprintf(
			"$ %s tx nft issue <denom-id> "+
				"--from=<key-name> "+
				"--name=<denom-name> "+
				"--symbol=<denom-symbol> "+
				"--mint-restricted=<mint-restricted> "+
				"--update-restricted=<update-restricted> "+
				"--schema=<schema-content or path to schema.json> "+
				"--description=<description> "+
				"--uri=<uri> "+
				"--uri-hash=<uri-hash> "+
				"--data=<data> "+
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

			denomName, err := cmd.Flags().GetString(FlagDenomName)
			if err != nil {
				return err
			}
			schema, err := cmd.Flags().GetString(FlagSchema)
			if err != nil {
				return err
			}
			symbol, err := cmd.Flags().GetString(FlagSymbol)
			if err != nil {
				return err
			}
			mintRestricted, err := cmd.Flags().GetBool(FlagMintRestricted)
			if err != nil {
				return err
			}
			updateRestricted, err := cmd.Flags().GetBool(FlagUpdateRestricted)
			if err != nil {
				return err
			}
			uri, err := cmd.Flags().GetString(FlagURI)
			if err != nil {
				return err
			}
			uriHash, err := cmd.Flags().GetString(FlagURIHash)
			if err != nil {
				return err
			}
			description, err := cmd.Flags().GetString(FlagDescription)
			if err != nil {
				return err
			}
			data, err := cmd.Flags().GetString(FlagData)
			if err != nil {
				return err
			}
			optionsContent, err := os.ReadFile(schema)
			if err == nil {
				schema = string(optionsContent)
			}

			msg := types.NewMsgIssueDenom(
				args[0],
				denomName,
				schema,
				clientCtx.GetFromAddress().String(),
				symbol,
				mintRestricted,
				updateRestricted,
				description,
				uri,
				uriHash,
				data,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsIssueDenom)
	_ = cmd.MarkFlagRequired(FlagMintRestricted)
	_ = cmd.MarkFlagRequired(FlagUpdateRestricted)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdMintNFT is the CLI command for a SaveNFT transaction
func GetCmdMintNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "mint [denom-id] [nft-id]",
		Long: "Mint an NFT and set the owner to the recipient.",
		Example: fmt.Sprintf(
			"$ %s tx nft mint <denom-id> <nft-id> "+
				"--uri=<uri> "+
				"--uri-hash=<uri-hash> "+
				"--recipient=<recipient> "+
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

			sender := clientCtx.GetFromAddress().String()

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

			tokenName, err := cmd.Flags().GetString(FlagTokenName)
			if err != nil {
				return err
			}
			tokenURI, err := cmd.Flags().GetString(FlagURI)
			if err != nil {
				return err
			}
			tokenURIHash, err := cmd.Flags().GetString(FlagURIHash)
			if err != nil {
				return err
			}
			tokenData, err := cmd.Flags().GetString(FlagData)
			if err != nil {
				return err
			}

			msg := types.NewMsgMintNFT(
				args[1],
				args[0],
				tokenName,
				tokenURI,
				tokenURIHash,
				tokenData,
				sender,
				recipient,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsMintNFT)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdEditNFT is the CLI command for sending an MsgEditNFT transaction
func GetCmdEditNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "edit [denom-id] [nft-id]",
		Long: "Edit the token data of an NFT.",
		Example: fmt.Sprintf(
			"$ %s tx nft edit <denom-id> <nft-id> "+
				"--uri=<uri> "+
				"--uri-hash=<uri-hash> "+
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

			tokenName, err := cmd.Flags().GetString(FlagTokenName)
			if err != nil {
				return err
			}
			tokenURI, err := cmd.Flags().GetString(FlagURI)
			if err != nil {
				return err
			}
			tokenURIHash, err := cmd.Flags().GetString(FlagURIHash)
			if err != nil {
				return err
			}
			tokenData, err := cmd.Flags().GetString(FlagData)
			if err != nil {
				return err
			}
			msg := types.NewMsgEditNFT(
				args[1],
				args[0],
				tokenName,
				tokenURI,
				tokenURIHash,
				tokenData,
				clientCtx.GetFromAddress().String(),
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsEditNFT)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdTransferNFT is the CLI command for sending a TransferNFT transaction
func GetCmdTransferNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "transfer [recipient] [denom-id] [nft-id]",
		Long: "Transfer an NFT to a recipient.",
		Example: fmt.Sprintf(
			"$ %s tx nft transfer <recipient> <denom-id> <nft-id> "+
				"--uri=<uri> "+
				"--uri-hash=<uri-hash> "+
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

			if _, err := sdk.AccAddressFromBech32(args[0]); err != nil {
				return err
			}

			tokenName, err := cmd.Flags().GetString(FlagTokenName)
			if err != nil {
				return err
			}
			tokenURI, err := cmd.Flags().GetString(FlagURI)
			if err != nil {
				return err
			}
			tokenURIHash, err := cmd.Flags().GetString(FlagURIHash)
			if err != nil {
				return err
			}
			tokenData, err := cmd.Flags().GetString(FlagData)
			if err != nil {
				return err
			}
			msg := types.NewMsgTransferNFT(
				args[2],
				args[1],
				tokenName,
				tokenURI,
				tokenURIHash,
				tokenData,
				clientCtx.GetFromAddress().String(),
				args[0],
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsTransferNFT)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdBurnNFT is the CLI command for sending a RemoveNFT transaction
func GetCmdBurnNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "burn [denom-id] [nft-id]",
		Long: "Burn an NFT.",
		Example: fmt.Sprintf(
			"$ %s tx nft burn <denom-id> <nft-id> "+
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

			msg := types.NewMsgBurnNFT(
				clientCtx.GetFromAddress().String(),
				args[1],
				args[0],
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

// GetCmdTransferDenom is the CLI command for sending a TransferDenom transaction
func GetCmdTransferDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "transfer-denom [recipient] [denom-id]",
		Long: "Transfer an Denom to a recipient.",
		Example: fmt.Sprintf(
			"$ %s tx nft transfer-denom <recipient> <denom-id> "+
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

			if _, err := sdk.AccAddressFromBech32(args[0]); err != nil {
				return err
			}

			msg := types.NewMsgTransferDenom(
				args[1],
				clientCtx.GetFromAddress().String(),
				args[0],
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
