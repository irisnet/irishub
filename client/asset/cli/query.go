package cli

import (
	"fmt"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v3/asset"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/codec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmdQueryTokens implements the query tokens command.
func getCmdQueryTokens(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tokens",
		Short:   "Query details of a group of tokens",
		Example: "iriscli asset token tokens --token-id=<token-id> --owner=<address>",
		PreRunE: preQueryTokenCmd,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := asset.QueryTokensParams{
				Owner:   viper.GetString(FlagOwner),
				TokenID: viper.GetString(FlagTokenID),
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.AssetRoute, asset.QueryTokens), bz)
			if err != nil {
				return err
			}

			var tokens asset.TokensOutput
			err = cdc.UnmarshalJSON(res, &tokens)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(tokens)
		},
	}

	cmd.Flags().AddFlagSet(FsTokensQuery)

	return cmd
}

// getCmdQueryFee implements the query asset related fees command.
func getCmdQueryFee(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "fee",
		Short:   "Query the asset related fees",
		Example: "iriscli asset token fee --symbol=<symbol>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			// query token fees
			tokenID := viper.GetString(FlagSymbol)
			if err := asset.CheckTokenID(tokenID); err != nil {
				return err
			}

			fees, err := queryTokenFees(cliCtx, tokenID)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(fees)
		},
	}

	cmd.Flags().AddFlagSet(FsFeeQuery)

	return cmd
}

// preQueryTokenCmd is used to check if the specified flags are valid
func preQueryTokenCmd(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()

	if flags.Changed(FlagOwner) && flags.Changed(FlagTokenID) {
		return fmt.Errorf("only one flag is allowed among the owner and token-id")
	}
	return nil
}
