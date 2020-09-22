package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/irisnet/irismod/modules/token/types"
)

// NewTxCmd returns the transaction commands for the token module.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Asset transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		getCmdIssueToken(),
		getCmdEditToken(),
		getCmdMintToken(),
		getCmdTransferTokenOwner(),
	)

	return txCmd
}

// getCmdIssueToken implements the issue token command
func getCmdIssueToken() *cobra.Command {
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
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			owner := clientCtx.GetFromAddress()

			msg := &types.MsgIssueToken{
				Symbol:        viper.GetString(FlagSymbol),
				Name:          viper.GetString(FlagName),
				MinUnit:       viper.GetString(FlagMinUnit),
				Scale:         uint32(viper.GetInt(FlagScale)),
				InitialSupply: uint64(viper.GetInt(FlagInitialSupply)),
				MaxSupply:     uint64(viper.GetInt(FlagMaxSupply)),
				Mintable:      viper.GetBool(FlagMintable),
				Owner:         owner,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			var prompt = "The token issue transaction will consume extra fee"

			if !viper.GetBool(flags.FlagGenerateOnly) {
				// query fee
				fee, err1 := queryTokenFees(clientCtx, msg.Symbol)
				if err1 != nil {
					return fmt.Errorf("failed to query token issue fee: %s", err1.Error())
				}

				// append issue fee to prompt
				issueFeeMainUnit := sdk.Coins{fee.IssueFee}.String()
				prompt += fmt.Sprintf(": %s", issueFeeMainUnit)
			}

			// a confirmation is needed
			prompt += "\nAre you sure to proceed?"
			confirmed, err := input.GetConfirmation(prompt, bufio.NewReader(os.Stdin), cmd.ErrOrStderr())
			if err != nil {
				return err
			}

			if !confirmed {
				return fmt.Errorf("operation aborted")
			}

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

// getCmdEditToken implements the edit token command
func getCmdEditToken() *cobra.Command {
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
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			owner := clientCtx.GetFromAddress()

			name := viper.GetString(FlagName)
			maxSupply := uint64(viper.GetInt(FlagMaxSupply))

			mintable, err := types.ParseBool(viper.GetString(FlagMintable))
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

func getCmdMintToken() *cobra.Command {
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
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			owner := clientCtx.GetFromAddress()

			amount := uint64(viper.GetInt64(FlagAmount))

			var to sdk.AccAddress
			addr := viper.GetString(FlagTo)
			if len(strings.TrimSpace(addr)) > 0 {
				to, err = sdk.AccAddressFromBech32(addr)
				if err != nil {
					return err
				}
			}

			msg := types.NewMsgMintToken(
				args[0], owner, to, amount,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			var prompt = "The token mint transaction will consume extra fee"

			if !viper.GetBool(flags.FlagGenerateOnly) {
				// query fee
				fee, err1 := queryTokenFees(clientCtx, args[0])
				if err1 != nil {
					return fmt.Errorf("failed to query token mint fee: %s", err1.Error())
				}

				// append mint fee to prompt
				mintFeeMainUnit := sdk.Coins{fee.MintFee}.String()
				prompt += fmt.Sprintf(": %s", mintFeeMainUnit)
			}

			// a confirmation is needed
			prompt += "\nAre you sure to proceed?"
			confirmed, err := input.GetConfirmation(prompt, bufio.NewReader(os.Stdin), cmd.ErrOrStderr())
			if err != nil {
				return err
			}

			if !confirmed {
				return fmt.Errorf("operation aborted")
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsMintToken)
	_ = cmd.MarkFlagRequired(FlagAmount)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// getCmdTransferTokenOwner implements the transfer token owner command
func getCmdTransferTokenOwner() *cobra.Command {
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
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			owner := clientCtx.GetFromAddress()

			to, err := sdk.AccAddressFromBech32(viper.GetString(FlagTo))
			if err != nil {
				return err
			}

			msg := types.NewMsgTransferTokenOwner(owner, to, args[0])

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
