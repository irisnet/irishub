package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/irisnet/irishub/app/v3/asset"
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetTokenCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "token",
		Short:                      "token transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
	}
	cmd.AddCommand(client.PostCommands(
		getCmdIssueToken(cdc),
		getCmdEditAsset(cdc),
		getCmdMintToken(cdc),
		getCmdTransferTokenOwner(cdc),
	)...)

	cmd.AddCommand(client.GetCommands(
		getCmdQueryTokens(cdc),
		getCmdQueryFee(cdc),
	)...)

	return cmd
}

// getCmdIssueToken implements the issue asset command
func getCmdIssueToken(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "issue",
		Short:   "Issue a new token",
		Example: `iriscli asset token issue --symbol="kitty" --name="Kitty Token" --initial-supply=100000000000 --max-supply=1000000000000 --from=<key-name> --chain-id=<chain-id> --fee=0.6iris`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			owner, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := asset.MsgIssueToken{
				Symbol:        viper.GetString(FlagSymbol),
				Name:          viper.GetString(FlagName),
				Decimal:       uint8(viper.GetInt(FlagDecimal)),
				InitialSupply: uint64(viper.GetInt(FlagInitialSupply)),
				MaxSupply:     uint64(viper.GetInt(FlagMaxSupply)),
				Mintable:      viper.GetBool(FlagMintable),
				Owner:         owner,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			var prompt = "The token issue transaction will consume extra fee"

			if !viper.GetBool(client.FlagGenerateOnly) {
				tokenId := asset.GetTokenID(msg.Symbol)

				// query fee
				fee, err1 := queryTokenFees(cliCtx, tokenId)
				if err1 != nil {
					return fmt.Errorf("failed to query token issue fee: %s", err1.Error())
				}

				// append issue fee to prompt
				issueFeeMainUnit := sdk.Coins{fee.IssueFee}.MainUnitString()
				prompt += fmt.Sprintf(": %s", issueFeeMainUnit)
			}

			// a confirmation is needed
			prompt += "\nAre you sure to proceed?"
			confirmed, err := client.GetConfirmation(prompt, bufio.NewReader(os.Stdin))
			if err != nil {
				return err
			}

			if !confirmed {
				return fmt.Errorf("operation aborted")
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsTokenIssue)
	_ = cmd.MarkFlagRequired(FlagSymbol)
	_ = cmd.MarkFlagRequired(FlagName)
	_ = cmd.MarkFlagRequired(FlagInitialSupply)
	_ = cmd.MarkFlagRequired(FlagDecimal)
	return cmd
}

// getCmdEditGateway implements the edit asset command
func getCmdEditAsset(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "edit",
		Short:   "Edit a existed token",
		Example: `iriscli asset token edit <token-id> --name="Cat Token" --max-supply=100000000000 --mintable=true --from=<your account name> --chain-id=<chain-id> --fee=0.6iris`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			owner, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			tokenId := args[0]
			name := viper.GetString(FlagName)
			maxSupply := uint64(viper.GetInt(FlagMaxSupply))
			mintable, err := asset.ParseBool(viper.GetString(FlagMintable))
			if err != nil {
				return err
			}
			var msg sdk.Msg
			msg = asset.NewMsgEditToken(name, tokenId, maxSupply, mintable, owner)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsEditToken)
	return cmd
}

func getCmdMintToken(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "mint",
		Short:   "The asset owner and operator can directly mint tokens to a specified address",
		Example: `iriscli tx asset token mint [token-id] --amount="" --to="" --from=<key-name> --chain-id=irishub --fee=0.3iris`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			owner, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			amount := uint64(viper.GetInt64(FlagAmount))
			var to sdk.AccAddress
			addr := viper.GetString(FlagTo)
			if len(strings.TrimSpace(addr)) > 0 {
				to, err = sdk.AccAddressFromBech32(addr)
				if err != nil {
					return err
				}
			}

			var msg sdk.Msg
			msg = asset.NewMsgMintToken(
				args[0], owner, to, amount,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			var prompt = "The token mint transaction will consume extra fee"

			if !viper.GetBool(client.FlagGenerateOnly) {
				tokenId, _ := sdk.ConvertIdToTokenKeyId(args[0])
				// query fee
				fee, err1 := queryTokenFees(cliCtx, tokenId)
				if err1 != nil {
					return fmt.Errorf("failed to query token mint fee: %s", err1.Error())
				}

				// append mint fee to prompt
				mintFeeMainUnit := sdk.Coins{fee.MintFee}.MainUnitString()
				prompt += fmt.Sprintf(": %s", mintFeeMainUnit)
			}

			// a confirmation is needed
			prompt += "\nAre you sure to proceed?"
			confirmed, err := client.GetConfirmation(prompt, bufio.NewReader(os.Stdin))
			if err != nil {
				return err
			}

			if !confirmed {
				return fmt.Errorf("operation aborted")
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsMintToken)
	_ = cmd.MarkFlagRequired(FlagAmount)
	return cmd
}

// getCmdTransferTokenOwner implements the transfer token owner command
func getCmdTransferTokenOwner(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "transfer",
		Short:   "Transfer the owner of a token to a new owner",
		Example: `iriscli tx asset token transfer [token-id] --to="" --from=<key-name> --chain-id=irishub --fee=0.3iris`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			owner, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			to, err := sdk.AccAddressFromBech32(viper.GetString(FlagTo))
			if err != nil {
				return err
			}

			var msg sdk.Msg
			msg = asset.NewMsgTransferTokenOwner(owner, to, args[0])

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsTransferTokenOwner)
	_ = cmd.MarkFlagRequired(FlagTo)
	return cmd
}
