package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/irisnet/irismod/modules/token/types"
)

// NewTxCmd returns the transaction commands for the token module.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Token transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdIssueToken(),
		GetCmdEditToken(),
		GetCmdMintToken(),
		GetCmdBurnToken(),
		GetCmdTransferTokenOwner(),
	)

	return txCmd
}

// GetCmdIssueToken implements the issue token command
func GetCmdIssueToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "issue",
		Long: "Issue a new token.",
		Example: fmt.Sprintf(
			"$ %s tx token issue "+
				"--name=\"Kitty Token\" "+
				"--symbol=\"kitty\" "+
				"--min-unit=\"kitty\" "+
				"--scale=0 "+
				"--initial-supply=100000000000 "+
				"--max-supply=1000000000000 "+
				"--mintable=true "+
				"--from=<key-name> "+
				"--chain-id=<chain-id> "+
				"--fees=<fee>",
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			owner := clientCtx.GetFromAddress()
			symbol, err := cmd.Flags().GetString(FlagSymbol)
			if err != nil {
				return err
			}
			name, err := cmd.Flags().GetString(FlagName)
			if err != nil {
				return err
			}
			minUnit, err := cmd.Flags().GetString(FlagMinUnit)
			if err != nil {
				return err
			}
			scale, err := cmd.Flags().GetUint32(FlagScale)
			if err != nil {
				return err
			}
			initialSupply, err := cmd.Flags().GetUint64(FlagInitialSupply)
			if err != nil {
				return err
			}
			maxSupply, err := cmd.Flags().GetUint64(FlagMaxSupply)
			if err != nil {
				return err
			}
			mintable, err := cmd.Flags().GetBool(FlagMintable)
			if err != nil {
				return err
			}

			msg := &types.MsgIssueToken{
				Symbol:        symbol,
				Name:          name,
				MinUnit:       minUnit,
				Scale:         scale,
				InitialSupply: initialSupply,
				MaxSupply:     maxSupply,
				Mintable:      mintable,
				Owner:         owner.String(),
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			var prompt = "The token issuance transaction will consume extra fee"

			generateOnly, err := cmd.Flags().GetBool(flags.FlagGenerateOnly)
			if err != nil {
				return err
			}
			if !generateOnly {
				// query fee
				fee, err1 := queryTokenFees(clientCtx, msg.Symbol)
				if err1 != nil {
					return fmt.Errorf("failed to query token issuance fee: %s", err1.Error())
				}

				// append issuance fee to prompt
				issueFeeMainUnit := sdk.Coins{fee.IssueFee}.String()
				prompt += fmt.Sprintf(": %s", issueFeeMainUnit)
			}

			fmt.Println(prompt)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsIssueToken)
	_ = cmd.MarkFlagRequired(FlagSymbol)
	_ = cmd.MarkFlagRequired(FlagName)
	_ = cmd.MarkFlagRequired(FlagInitialSupply)
	_ = cmd.MarkFlagRequired(FlagScale)
	_ = cmd.MarkFlagRequired(FlagMinUnit)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdEditToken implements the edit token command
func GetCmdEditToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "edit [symbol]",
		Long: "Edit an existing token.",
		Example: fmt.Sprintf(
			"$ %s tx token edit <symbol> "+
				"--name=\"Cat Token\" "+
				"--max-supply=100000000000 "+
				"--mintable=true "+
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

			owner := clientCtx.GetFromAddress().String()

			name, err := cmd.Flags().GetString(FlagName)
			if err != nil {
				return err
			}
			maxSupply, err := cmd.Flags().GetUint64(FlagMaxSupply)
			if err != nil {
				return err
			}
			rawMintable, err := cmd.Flags().GetString(FlagMintable)
			if err != nil {
				return err
			}
			mintable, err := types.ParseBool(rawMintable)
			if err != nil {
				return err
			}

			msg := types.NewMsgEditToken(name, args[0], maxSupply, mintable, owner)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsEditToken)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdMintToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "mint [symbol]",
		Long: "Mint tokens to a specified address.",
		Example: fmt.Sprintf(
			"$ %s tx token mint <symbol> "+
				"--amount=<amount> "+
				"--to=<to> "+
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

			owner := clientCtx.GetFromAddress().String()

			amount, err := cmd.Flags().GetUint64(FlagAmount)
			if err != nil {
				return err
			}

			addr, err := cmd.Flags().GetString(FlagTo)
			if err != nil {
				return err
			}
			if len(addr) > 0 {
				if _, err = sdk.AccAddressFromBech32(addr); err != nil {
					return err
				}
			}

			msg := types.NewMsgMintToken(
				strings.TrimSpace(args[0]), owner, addr, amount,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			var prompt = "The token minting transaction will consume extra fee"

			generateOnly, err := cmd.Flags().GetBool(flags.FlagGenerateOnly)
			if err != nil {
				return err
			}
			if !generateOnly {
				// query fee
				fee, err1 := queryTokenFees(clientCtx, args[0])
				if err1 != nil {
					return fmt.Errorf("failed to query token minting fee: %s", err1.Error())
				}

				// append mint fee to prompt
				mintFeeMainUnit := sdk.Coins{fee.MintFee}.String()
				prompt += fmt.Sprintf(": %s", mintFeeMainUnit)
			}

			fmt.Println(prompt)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsMintToken)
	_ = cmd.MarkFlagRequired(FlagAmount)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdBurnToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "burn [symbol]",
		Long: "Burn tokens.",
		Example: fmt.Sprintf(
			"$ %s tx token burn <symbol> "+
				"--amount=<amount> "+
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

			owner := clientCtx.GetFromAddress().String()

			amount, err := cmd.Flags().GetUint64(FlagAmount)
			if err != nil {
				return err
			}

			addr, err := cmd.Flags().GetString(FlagTo)
			if err != nil {
				return err
			}
			if len(addr) > 0 {
				if _, err = sdk.AccAddressFromBech32(addr); err != nil {
					return err
				}
			}

			msg := types.NewMsgBurnToken(
				strings.TrimSpace(args[0]), owner, amount,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsMintToken)
	_ = cmd.MarkFlagRequired(FlagAmount)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdTransferTokenOwner implements the transfer token owner command
func GetCmdTransferTokenOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "transfer [symbol]",
		Long: "Transfer the owner of a token to a new owner.",
		Example: fmt.Sprintf(
			"$ %s tx token transfer <symbol> "+
				"--to=<to> "+
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

			owner := clientCtx.GetFromAddress().String()

			toAddr, err := cmd.Flags().GetString(FlagTo)
			if err != nil {
				return err
			}
			if _, err := sdk.AccAddressFromBech32(toAddr); err != nil {
				return err
			}

			msg := types.NewMsgTransferTokenOwner(owner, toAddr, strings.TrimSpace(args[0]))

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsTransferTokenOwner)
	_ = cmd.MarkFlagRequired(FlagTo)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
