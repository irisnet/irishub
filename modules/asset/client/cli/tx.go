package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/irisnet/irishub/modules/asset/types"
	iristypes "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Asset transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		GetCmdIssueToken(queryRoute, cdc),
		GetCmdTransferTokenOwner(cdc),
		GetCmdEditToken(cdc),
		GetCmdMintToken(queryRoute, cdc),
	)
	return txCmd
}

func preSignCmd(cmd *cobra.Command, _ []string) {
	// Conditionally mark the account and sequence numbers required as no RPC
	// query will be done.
	if viper.GetString(FlagSource) == "gateway" {
		cmd.MarkFlagRequired(FlagGateway)
	}
}

// GetCmdIssueToken implements the issue asset command
func GetCmdIssueToken(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue-token",
		Short: "Issue a new token",
		Example: "iriscli asset issue-token --family=<family> --source=<source> --gateway=<gateway-moniker> --decimal=<decimal>" +
			" --symbol=<symbol> --name=<token-name> --initial-supply=<initial-supply> --from=<key-name> --chain-id=<chain-id> --fee=0.6iris",
		PreRun: preSignCmd,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := authtypes.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			owner := cliCtx.GetFromAddress()

			family, ok := types.StringToAssetFamilyMap[strings.ToLower(viper.GetString(FlagFamily))]
			if !ok {
				return fmt.Errorf("invalid token family type %s", viper.GetString(FlagFamily))
			}

			source, ok := types.StringToAssetSourceMap[strings.ToLower(viper.GetString(FlagSource))]
			if !ok || source == types.EXTERNAL {
				return fmt.Errorf("invalid token source type %s", viper.GetString(FlagSource))
			}

			msg := types.MsgIssueToken{
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
				tokenId, err := types.GetTokenID(msg.Source, msg.Symbol, msg.Gateway)
				if err != nil {
					return fmt.Errorf("failed to query token issue fee: %s", err.Error())
				}

				// query fee
				fee, err1 := queryTokenFees(cliCtx, queryRoute, tokenId)
				if err1 != nil {
					return fmt.Errorf("failed to query token issue fee: %s", err1.Error())
				}

				// append issue fee to prompt
				issueFeeMainUnit := sdk.Coins{fee.IssueFee}.String()
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

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
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

// GetCmdEditToken implements the edit asset command
func GetCmdEditToken(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "edit-token",
		Short:   "Edit a existed token",
		Example: "iriscli asset edit-token <token-id> --name=<name> --canonical-symbol=<canonical-symbol> --min-unit-alias=<min-alias> --max-supply=<max-supply> --mintable=<mintable> --from=<your account name> --chain-id=<chain-id> --fee=0.6iris",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := authtypes.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			owner := cliCtx.GetFromAddress()

			tokenId := args[0]
			name := viper.GetString(FlagName)
			canonicalSymbol := viper.GetString(FlagCanonicalSymbol)
			minUnitAlias := viper.GetString(FlagMinUnitAlias)
			maxSupply := uint64(viper.GetInt(FlagMaxSupply))
			mintable, err := types.ParseBool(viper.GetString(FlagMintable))
			if err != nil {
				return err
			}
			var msg sdk.Msg
			msg = types.NewMsgEditToken(name,
				canonicalSymbol, minUnitAlias, tokenId, maxSupply, mintable, owner)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsEditToken)
	return cmd
}

func GetCmdMintToken(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "mint-token",
		Short:   "The asset owner and operator can directly mint tokens to a specified address",
		Example: "iriscli asset mint-token <token-id> [flags]",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := authtypes.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			owner := cliCtx.GetFromAddress()

			amount := uint64(viper.GetInt64(FlagAmount))
			var to sdk.AccAddress
			var err error
			addr := viper.GetString(FlagTo)
			if len(strings.TrimSpace(addr)) > 0 {
				to, err = sdk.AccAddressFromBech32(addr)
				if err != nil {
					return err
				}
			}

			var msg sdk.Msg
			msg = types.NewMsgMintToken(
				args[0], owner, to, amount,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			var prompt = "The token mint transaction will consume extra fee"

			if !viper.GetBool(client.FlagGenerateOnly) {
				tokenId, _ := iristypes.ConvertIdToTokenKeyId(args[0])
				// query fee
				fee, err1 := queryTokenFees(cliCtx, queryRoute, tokenId)
				if err1 != nil {
					return fmt.Errorf("failed to query token mint fee: %s", err1.Error())
				}

				// append mint fee to prompt
				mintFeeMainUnit := sdk.Coins{fee.MintFee}.String()
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

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
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
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := authtypes.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			owner := cliCtx.GetFromAddress()

			to, err := sdk.AccAddressFromBech32(viper.GetString(FlagTo))
			if err != nil {
				return err
			}

			var msg sdk.Msg
			msg = types.NewMsgTransferTokenOwner(owner, to, args[0])

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsTransferTokenOwner)
	cmd.MarkFlagRequired(FlagTo)

	return cmd
}
