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
	"github.com/spf13/viper"
)

func preSignCmd(cmd *cobra.Command, _ []string) {
	// Conditionally mark the account and sequence numbers required as no RPC
	// query will be done.
	if viper.GetString(FlagSource) == "gateway" {
		cmd.MarkFlagRequired(FlagGateway)
	}
}

// GetCmdIssueToken implements the issue asset command
func GetCmdIssueToken(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue-token",
		Short: "Issue a new token",
		Example: "iriscli asset issue-token --family=<family> --source=<source> --gateway=<gateway-moniker> --decimal=<decimal>" +
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
			if !ok || source == asset.EXTERNAL {
				return fmt.Errorf("invalid token source type %s", viper.GetString(FlagSource))
			}

			msg := asset.MsgIssueToken{
				Family:          family,
				Source:          source,
				Gateway:         viper.GetString(FlagGateway),
				Symbol:          viper.GetString(FlagSymbol),
				CanonicalSymbol: viper.GetString(FlagCanonicalSymbol),
				Name:            viper.GetString(FlagName),
				Decimal:         uint8(viper.GetInt(FlagDecimal)),
				MinUnitAlias:    viper.GetString(FlagMinUnitAlias),
				InitialSupply:   uint64(viper.GetInt(FlagInitialSupply)),
				MaxSupply:       uint64(viper.GetInt(FlagMaxSupply)),
				Mintable:        viper.GetBool(FlagMintable),
				Owner:           owner,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			var prompt = "The token issue transaction will consume extra fee"

			if !viper.GetBool(client.FlagGenerateOnly) {
				tokenId, err := asset.GetTokenID(msg.Source, msg.Symbol, msg.Gateway)
				if err != nil {
					return fmt.Errorf("failed to query token issue fee: %s", err.Error())
				}

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
	cmd.MarkFlagRequired(FlagFamily)
	cmd.MarkFlagRequired(FlagSource)
	cmd.MarkFlagRequired(FlagSymbol)
	cmd.MarkFlagRequired(FlagName)
	cmd.MarkFlagRequired(FlagInitialSupply)
	cmd.MarkFlagRequired(FlagDecimal)

	return cmd
}

// GetCmdCreateGateway implements the create gateway command
func GetCmdCreateGateway(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-gateway",
		Short: "Create a gateway",
		Example: "iriscli asset create-gateway --moniker=<moniker> --identity=<identity> --details=<details> " +
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
			identity := viper.GetString(FlagIdentity)
			details := viper.GetString(FlagDetails)
			website := viper.GetString(FlagWebsite)

			var msg sdk.Msg
			msg = asset.NewMsgCreateGateway(
				owner, moniker, identity, details, website,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			var prompt = "The gateway creation transaction will consume extra fee"

			if !viper.GetBool(client.FlagGenerateOnly) {
				// query fee
				creationFee, err := queryGatewayFee(cliCtx, moniker)
				if err != nil {
					return fmt.Errorf("failed to query gateway creation fee: %s", err.Error())
				}

				// append creation fee to prompt
				creationFeeMainUnit := sdk.Coins{creationFee.Fee}.MainUnitString()
				prompt += fmt.Sprintf(": %s", creationFeeMainUnit)
			}

			// a confirmation is needed
			prompt += "\nAre you sure to proceed?"
			confirmed, err := client.GetConfirmation(prompt, bufio.NewReader(os.Stdin))
			if err != nil {
				return err
			}

			if !confirmed {
				return fmt.Errorf("The operation aborted")
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsGatewayCreate)
	cmd.MarkFlagRequired(FlagMoniker)

	return cmd
}

// GetCmdEditGateway implements the edit gateway command
func GetCmdEditGateway(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-gateway",
		Short: "Edit a gateway",
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
			identity := viper.GetString(FlagIdentity)
			details := viper.GetString(FlagDetails)
			website := viper.GetString(FlagWebsite)

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

// GetCmdEditGateway implements the edit asset command
func GetCmdEditAsset(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "edit-token",
		Short:   "Edit a existed token",
		Example: "iriscli asset edit-token <token-id> --name=<name> --canonical-symbol=<canonical-symbol> --min-unit-alias=<min-alias> --max-supply=<max-supply> --mintable=<mintable> --from=<your account name> --chain-id=<chain-id> --fee=0.6iris",
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
			canonicalSymbol := viper.GetString(FlagCanonicalSymbol)
			minUnitAlias := viper.GetString(FlagMinUnitAlias)
			maxSupply := uint64(viper.GetInt(FlagMaxSupply))
			mintable, err := asset.ParseBool(viper.GetString(FlagMintable))
			if err != nil {
				return err
			}
			var msg sdk.Msg
			msg = asset.NewMsgEditToken(name,
				canonicalSymbol, minUnitAlias, tokenId, maxSupply, mintable, owner)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsEditToken)
	return cmd
}

// GetCmdTransferGatewayOwner implements the transfer gateway owner command
func GetCmdTransferGatewayOwner(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-gateway-owner",
		Short: "Transfer the owner of a gateway",
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

			to, err := sdk.AccAddressFromBech32(viper.GetString(FlagTo))
			if err != nil {
				return err
			}

			var msg sdk.Msg
			msg = asset.NewMsgTransferGatewayOwner(
				owner, moniker, to,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsGatewayOwnerTransfer)
	cmd.MarkFlagRequired(FlagMoniker)
	cmd.MarkFlagRequired(FlagTo)

	return cmd
}

func GetCmdMintToken(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "mint-token",
		Short:   "The asset owner and operator can directly mint tokens to a specified address",
		Example: "iriscli asset mint-token <token-id> [flags]",
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
	cmd.MarkFlagRequired(FlagAmount)
	return cmd
}

// GetCmdTransferTokenOwner implements the transfer token owner command
func GetCmdTransferTokenOwner(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "transfer-token-owner",
		Short:   "Transfer the owner of a token to a new owner",
		Example: "iriscli asset transfer-token-owner <token-id> --to=<new owner>",
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
	cmd.MarkFlagRequired(FlagTo)

	return cmd
}
