package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v3/asset"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

// getCmdQueryToken implements the query token command.
func getCmdQueryToken(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "token [symbol]",
		Short:   "Query a token by symbol",
		Example: "iriscli asset token token <symbol>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			if err := asset.CheckSymbol(args[0]); err != nil {
				return err
			}

			params := asset.QueryTokenParams{
				Symbol: args[0],
			}

			bz := cdc.MustMarshalJSON(params)
			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.AssetRoute, asset.QueryToken), bz)
			if err != nil {
				return err
			}

			var token asset.TokenOutput
			if err := cdc.UnmarshalJSON(res, &token); err != nil {
				return err
			}

			return cliCtx.PrintOutput(token)
		},
	}

	return cmd
}

// getCmdQueryTokens implements the query tokens command.
func getCmdQueryTokens(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tokens [owner]",
		Short:   "Query tokens by owner",
		Example: "iriscli asset token tokens <owner>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var err error
			var owner sdk.AccAddress

			if len(args) > 0 {
				owner, err = sdk.AccAddressFromBech32(args[0])
				if err != nil {
					return err
				}
			}

			params := asset.QueryTokensParams{
				Owner: owner,
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

	return cmd
}

// getCmdQueryFee implements the query token related fees command.
func getCmdQueryFee(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "fee [symbol]",
		Short:   "Query the token related fees",
		Args:    cobra.ExactArgs(1),
		Example: "iriscli asset token fee <symbol>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			symbol := args[0]
			if err := asset.CheckSymbol(symbol); err != nil {
				return err
			}

			// query token fees
			fees, err := queryTokenFees(cliCtx, symbol)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(fees)
		},
	}

	return cmd
}
