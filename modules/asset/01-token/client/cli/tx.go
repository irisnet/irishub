package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/irisnet/irishub/modules/asset/01-token/internal/types"
)

// GetTxCmd returns the transaction commands for the asset module.
func GetTxCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.SubModuleName,
		Short:                      "Asset transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(flags.PostCommands(
		GetCmdIssueToken(queryRoute, cdc),
		GetCmdEditToken(cdc),
		GetCmdMintToken(queryRoute, cdc),
		GetCmdTransferToken(cdc),
		GetCmdBurnToken(cdc),
	)...)

	return txCmd
}

// GetCmdIssueToken implements the issue token command.
func GetCmdIssueToken(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue",
		Short: "Issue a new token",
		Example: fmt.Sprintf("%s tx asset token issue --symbol=<symbol> --name=<token-name>"+
			" --scale=<token-scale> --min-unit=<token-min-unit> --initial-supply=<initial-supply> --from=<key-name>"+
			" --chain-id=<chain-id> --fees=0.6iris", version.ClientName),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			owner := cliCtx.GetFromAddress()

			msg := types.MsgIssueToken{
				Symbol:        viper.GetString(FlagSymbol),
				Name:          viper.GetString(FlagName),
				Scale:         uint8(viper.GetInt(FlagScale)),
				MinUnit:       viper.GetString(FlagMinUnit),
				InitialSupply: uint64(viper.GetInt(FlagInitialSupply)),
				MaxSupply:     uint64(viper.GetInt(FlagMaxSupply)),
				Mintable:      viper.GetBool(FlagMintable),
				Owner:         owner,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			if !viper.GetBool(flags.FlagGenerateOnly) {
				// query fee
				fee, err1 := queryTokenFees(cliCtx, queryRoute, msg.Symbol)
				if err1 != nil {
					return fmt.Errorf("failed to query token issue fee: %s", err1.Error())
				}

				// append issue fee to prompt
				issueFeeMainUnit := sdk.Coins{fee.IssueFee}.String()
				fmt.Printf("The token issue transaction will consume extra fee: %s\n", issueFeeMainUnit)
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsTokenIssue)
	_ = cmd.MarkFlagRequired(FlagSymbol)
	_ = cmd.MarkFlagRequired(FlagName)
	_ = cmd.MarkFlagRequired(FlagScale)
	_ = cmd.MarkFlagRequired(FlagMinUnit)
	_ = cmd.MarkFlagRequired(FlagInitialSupply)

	return cmd
}

// GetCmdEditToken implements the edit token command.
func GetCmdEditToken(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit [symbol]",
		Short: "Edit a existed token",
		Example: fmt.Sprintf("%s tx asset edit-token <symbol> --name=<name>"+
			" --max-supply=<max-supply> --mintable=<mintable> --from=<your account name>"+
			" --chain-id=<chain-id> --fees=0.6iris", version.ClientName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			owner := cliCtx.GetFromAddress()
			symbol := args[0]
			name := viper.GetString(FlagName)
			maxSupply := uint64(viper.GetInt(FlagMaxSupply))
			mintable, err := types.ParseBool(viper.GetString(FlagMintable))
			if err != nil {
				return err
			}

			msg := types.NewMsgEditToken(name, symbol, maxSupply, mintable, owner)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsEditToken)
	return cmd
}

// GetCmdMintToken implements the mint token command.
func GetCmdMintToken(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint [symbol]",
		Short: "The asset owner and operator can directly mint tokens to a specified address",
		Example: fmt.Sprintf("%s tx asset mint [symbol] --recipient=[address] --amount=[amount] --from=<your account name>"+
			" --chain-id=<chain-id> --fees=0.6iris", version.ClientName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var (
				to  sdk.AccAddress
				err error
			)

			if addr := viper.GetString(FlagRecipient); len(strings.TrimSpace(addr)) > 0 {
				if to, err = sdk.AccAddressFromBech32(addr); err != nil {
					return err
				}
			}

			owner := cliCtx.GetFromAddress()
			amount := uint64(viper.GetInt64(FlagAmount))
			msg := types.NewMsgMintToken(
				args[0], owner, to, amount,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			if !viper.GetBool(flags.FlagGenerateOnly) {
				// query fee
				var fee types.TokenFeesOutput
				if fee, err = queryTokenFees(cliCtx, queryRoute, msg.Symbol); err != nil {
					return fmt.Errorf("failed to query token mint fee: %s", err.Error())
				}

				// append mint fee to prompt
				mintFeeMainUnit := sdk.Coins{fee.MintFee}.String()
				fmt.Printf("The token mint transaction will consume extra fee: %s\n", mintFeeMainUnit)
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsMintToken)
	_ = cmd.MarkFlagRequired(FlagAmount)
	return cmd
}

// GetCmdTransferToken implements the transfer token owner command.
func GetCmdTransferToken(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer [symbol]",
		Short: "Transfer the owner of a token to a new owner",
		Example: fmt.Sprintf("%s tx asset transfer [symbol] --recipient=<new owner> --from=<your account name>"+
			" --chain-id=<chain-id> --fees=0.6iris", version.ClientName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			owner := cliCtx.GetFromAddress()

			recipient, err := sdk.AccAddressFromBech32(viper.GetString(FlagRecipient))
			if err != nil {
				return err
			}

			msg := types.NewMsgTransferToken(owner, recipient, args[0])

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsTransferToken)
	_ = cmd.MarkFlagRequired(FlagRecipient)

	return cmd
}

// GetCmdBurnToken implements the transfer token owner command.
func GetCmdBurnToken(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "burn [amount]",
		Short: "burn the specified amount token from the account of the sender",
		Example: fmt.Sprintf("%s tx asset burn [amount] --from=<your account name>"+
			" --chain-id=<chain-id> --fees=0.6iris", version.ClientName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// parse coins trying to be burn
			coins, err := sdk.ParseCoins(args[0])
			if err != nil {
				return err
			}

			owner := cliCtx.GetFromAddress()
			msg := types.NewMsgBurnToken(owner, coins)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
