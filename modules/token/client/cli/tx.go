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
	v1 "github.com/irisnet/irismod/modules/token/types/v1"
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
		GetCmdSwapFeeToken(),
		GetCmdSwapToErc20(),
		GetCmdSwapFromErc20(),
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

			msg := &v1.MsgIssueToken{
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

			msg := v1.NewMsgEditToken(name, args[0], maxSupply, mintable, owner)
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
		Use:  "mint [coin]",
		Long: "Mint tokens to a specified address.",
		Example: fmt.Sprintf(
			"$ %s tx token mint <coin> "+
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
			addr, err := cmd.Flags().GetString(FlagTo)
			if err != nil {
				return err
			}
			if len(addr) > 0 {
				if _, err = sdk.AccAddressFromBech32(addr); err != nil {
					return err
				}
			}

			coin, token, err := parseMainCoin(clientCtx, args[0])
			if err != nil {
				return err
			}

			msg := &v1.MsgMintToken{
				Coin:  coin,
				To:    addr,
				Owner: owner,
			}

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
				fee, err1 := queryTokenFees(clientCtx, token.GetSymbol())
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
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdBurnToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "burn [coin]",
		Long: "Burn tokens.",
		Example: fmt.Sprintf(
			"$ %s tx token burn <coin> "+
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

			coin, _, err := parseMainCoin(clientCtx, args[0])
			if err != nil {
				return err
			}

			msg := &v1.MsgBurnToken{
				Coin:   coin,
				Sender: clientCtx.GetFromAddress().String(),
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
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

			msg := v1.NewMsgTransferTokenOwner(owner, toAddr, strings.TrimSpace(args[0]))

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

// GetCmdSwapFeeToken implements the swap token command
func GetCmdSwapFeeToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "swap-fee [fee_paid]",
		Long: "Use the input token to exchange for a specified number of other tokens. Note: the exchanged token pair must be registered by the system",
		Example: fmt.Sprintf(
			"$ %s tx token swap-fee <fee_paid> "+
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

			sender := clientCtx.GetFromAddress().String()
			toAddr, err := cmd.Flags().GetString(FlagTo)
			if err != nil {
				return err
			}
			if toAddr != "" {
				if _, err := sdk.AccAddressFromBech32(toAddr); err != nil {
					return err
				}
			}

			coin, _, err := parseMainCoin(clientCtx, args[0])
			if err != nil {
				return err
			}

			msg := &v1.MsgSwapFeeToken{
				FeePaid:   coin,
				Recipient: toAddr,
				Sender:    sender,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsSwapFeeToken)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdSwapToErc20 implements the swap-to-erc20 command
func GetCmdSwapToErc20() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "swap-to-erc20",
		Long: "Swap native token to the corresponding ERC20 at 1:1.",
		Example: fmt.Sprintf(
			"$ %s tx token swap-to-erc20 [paid_amount]"+
				"--to=\"0x0eeb8ec40c6705b669469346ff8f9ce5cad57ed5\" "+
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

			from := clientCtx.GetFromAddress().String()
			paidAmount, token, err := parseMainCoin(clientCtx, args[0])
			if err != nil {
				return err
			}
			if len(token.GetContract()) <= 0 {
				return fmt.Errorf("corresponding erc20 contract of %s does not exist", paidAmount.Denom)
			}

			receiver, err := cmd.Flags().GetString(FlagTo)
			if err != nil {
				return err
			}
			if len(receiver) <= 0 {
				// set default receiver
				receiver = from
			}

			msg := &v1.MsgSwapToERC20{
				Amount:   paidAmount,
				Sender:   from,
				Receiver: receiver,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			var prompt = fmt.Sprintf("Swapping native token to ERC20, sender: %s, receiver: %s, contract: %s, amount: %s", from, receiver, token.GetContract(), paidAmount)

			fmt.Println(prompt)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsSwapToErc20)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdSwapFromErc20 implements the swap-from-erc20 command
func GetCmdSwapFromErc20() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "swap-from-erc20",
		Long: "Swap native token from the corresponding ERC20 at 1:1.",
		Example: fmt.Sprintf(
			"$ %s tx token swap-from-erc20 [wanted_amount]"+
				"--to=\"iaaeeb8ec40c6705b669469346ff8f9ce5cad57ed5\" "+
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

			from := clientCtx.GetFromAddress().String()
			wantedAmount, token, err := parseMainCoin(clientCtx, args[0])
			if err != nil {
				return err
			}
			if len(token.GetContract()) <= 0 {
				return fmt.Errorf("corresponding erc20 contract of %s does not exist", wantedAmount.Denom)
			}

			receiver, err := cmd.Flags().GetString(FlagTo)
			if err != nil {
				return err
			}
			if len(receiver) <= 0 {
				// set default receiver
				receiver = from
			}

			msg := &v1.MsgSwapFromERC20{
				WantedAmount: wantedAmount,
				Sender:       from,
				Receiver:     receiver,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			var prompt = fmt.Sprintf("Swapping native token from ERC20, sender: %s, receiver: %s, contract: %s, amount: %s", from, receiver, token.GetContract(), wantedAmount)

			fmt.Println(prompt)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsSwapFromErc20)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
