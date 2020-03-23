package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v3/asset"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/codec"
)

// getCmdQueryTokens implements the query tokens command.
func getCmdQueryTokens(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tokens",
		Short:   "Query details of a group of tokens",
		Example: "iriscli asset token tokens --symbol=<symbol> --owner=<address>",
		PreRunE: preQueryTokenCmd,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := asset.QueryTokensParams{
				Owner:  viper.GetString(FlagOwner),
				Symbol: viper.GetString(FlagSymbol),
			}

			bz := cdc.MustMarshalJSON(params)
			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.AssetRoute, asset.QueryTokens), bz)
			if err != nil {
				return err
			}

			var tokens asset.TokensOutput
			if err := cdc.UnmarshalJSON(res, &tokens); err != nil {
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
		Use:     "fee [symbol]",
		Short:   "Query the asset related fees",
		Args:    cobra.ExactArgs(1),
		Example: "iriscli asset token fee [symbol]",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// query token fees
			symbol := args[0]
			if err := asset.CheckSymbol(symbol); err != nil {
				return err
			}

			fees, err := queryTokenFees(cliCtx, symbol)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(fees)
		},
	}

	return cmd
}

// preQueryTokenCmd is used to check if the specified flags are valid
func preQueryTokenCmd(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()

	if flags.Changed(FlagOwner) && flags.Changed(FlagSymbol) {
		return fmt.Errorf("only one flag is allowed among the owner and symbol")
	}

	return nil
}
