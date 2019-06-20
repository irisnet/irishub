package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/irisnet/irishub/app/v1/asset"
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func preSignCmd(cmd *cobra.Command, _ []string) {
	// Conditionally mark the account and sequence numbers required as no RPC
	// query will be done.
	if viper.GetString(FlagSource) == "gateway" {
		cmd.MarkFlagRequired(FlagGateway)
	}
}

// GetCmdIssueAsset implements the issue asset command
func GetCmdIssueAsset(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue-token",
		Short: "issue a new token",
		Example: "iriscli asset issue-token --family=<family> --source=<source> --gateway=<gateway-moniker>" +
			" --symbol=<symbol> --name=<token-name> --initial-supply=<initial-supply> --from=<key-name> --chain-id=<chain-id> --fee=0.6iris",
		PreRun: preSignCmd,
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

			family, ok := asset.StringToAssetFamilyMap[strings.ToLower(viper.GetString(FlagFamily))]
			if !ok {
				return fmt.Errorf("invalid token family type %s", viper.GetString(FlagFamily))
			}

			source, ok := asset.StringToAssetSourceMap[strings.ToLower(viper.GetString(FlagSource))]
			if !ok {
				return fmt.Errorf("invalid token source type %s", viper.GetString(FlagSource))
			}

			msg := asset.MsgIssueToken{
				Family:         family,
				Source:         source,
				Gateway:        viper.GetString(FlagGateway),
				Symbol:         viper.GetString(FlagSymbol),
				SymbolAtSource: viper.GetString(FlagSymbolAtSource),
				Name:           viper.GetString(FlagName),
				Decimal:        uint8(viper.GetInt(FlagDecimal)),
				SymbolMinAlias: viper.GetString(FlagSymbolMinAlias),
				InitialSupply:  uint64(viper.GetInt(FlagInitialSupply)),
				MaxSupply:      uint64(viper.GetInt(FlagMaxSupply)),
				Mintable:       viper.GetBool(FlagMintable),
				Owner:          owner,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsAssetIssue)
	cmd.MarkFlagRequired(FlagFamily)
	cmd.MarkFlagRequired(FlagSource)
	cmd.MarkFlagRequired(FlagSymbol)
	cmd.MarkFlagRequired(FlagName)
	cmd.MarkFlagRequired(FlagInitialSupply)

	return cmd
}

// GetCmdCreateGateway implements the create gateway command
func GetCmdCreateGateway(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-gateway",
		Short: "create a gateway",
		Example: "iriscli asset create-gateway --moniker=<moniker> --identity=<identity> --details=<details> " +
			"--website=<website> --create-fee=<gateway create fee>",
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

			moniker := viper.GetString(FlagMoniker)
			identity := viper.GetString(FlagIdentity)
			details := viper.GetString(FlagDetails)
			website := viper.GetString(FlagWebsite)
			createFee := viper.GetString(FlagCreateFee)

			createFeeCoin, err := sdk.ParseCoin(createFee)
			if err != nil {
				return err
			}

			if createFeeCoin.Denom == sdk.NativeTokenName {
				createFeeCoin = ConvertToNativeTokenMin(createFeeCoin)
			}

			var msg sdk.Msg
			msg = asset.NewMsgCreateGateway(
				owner, moniker, identity, details, website, createFeeCoin,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			var prompt = "The gateway creation transaction will consume extra fee"

			if !viper.GetBool(client.FlagGenerateOnly) {
				// query fee
				actualFee, err := queryGatewayFee(cliCtx, moniker)
				if err != nil {
					return fmt.Errorf("failed to query gateway creation fee: %s", err.Error())
				}

				// check if the provided fee is enough
				if createFeeCoin.IsLT(actualFee.Fee) {
					return fmt.Errorf("insufficient creation fee: expected %s, got %s", actualFee.Fee, createFeeCoin)
				}

				// append actual fee to prompt
				actualNativeTokenFee := ConvertToNativeToken(actualFee.Fee)
				prompt += fmt.Sprintf(": %s", actualNativeTokenFee)
			}

			// a confirmation is needed
			prompt += "\nAre you sure to proceed? (y/n)"
			confirmed, err := client.GetConfirmation(prompt, bufio.NewReader(os.Stdin))
			if err != nil {
				return err
			}

			if !confirmed {
				return fmt.Errorf("transaction aborted")
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsGatewayCreate)
	cmd.MarkFlagRequired(FlagMoniker)
	cmd.MarkFlagRequired(FlagCreateFee)

	return cmd
}

// GetCmdEditGateway implements the edit gateway command
func GetCmdEditGateway(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-gateway",
		Short: "edit a gateway",
		Example: "iriscli asset edit-gateway --moniker=<moniker> --identity=<identity> --details=<details> " +
			"--website=<website>",
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

			moniker := viper.GetString(FlagMoniker)
			identity := (*string)(nil)
			details := (*string)(nil)
			website := (*string)(nil)

			flags := cmd.Flags()
			flags.Visit(func(f *pflag.Flag) {
				if f.Name == FlagIdentity {
					value := f.Value.String()
					identity = &value
				}

				if f.Name == FlagDetails {
					value := f.Value.String()
					details = &value
				}

				if f.Name == FlagWebsite {
					value := f.Value.String()
					website = &value
				}
			})

			var msg sdk.Msg
			msg = asset.NewMsgEditGateway(
				owner, moniker, identity, details, website,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsGatewayEdit)
	cmd.MarkFlagRequired(FlagMoniker)

	return cmd
}
