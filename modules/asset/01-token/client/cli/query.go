package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/irisnet/irishub/modules/asset/01-token/internal/types"
)

// GetQueryCmd returns the query commands for the asset module.
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                types.SubModuleName,
		Short:              "Querying commands for the token module",
		DisableFlagParsing: true,
	}

	queryCmd.AddCommand(flags.GetCommands(
		GetCmdQueryTokens(queryRoute, cdc),
		GetCmdQueryFee(queryRoute, cdc),
		GetCmdQueryParams(queryRoute, cdc),
	)...)

	return queryCmd
}

// GetCmdQueryTokens implements the query tokens command.
func GetCmdQueryTokens(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tokens",
		Short:   "Query details of a group of tokens",
		Example: fmt.Sprintf("%s asset token tokens --symbol=[symbol] \n%s asset token tokens --owner=[address]", version.ClientName, version.ClientName),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := types.QueryTokensParams{
				Symbol: viper.GetString(FlagSymbol),
				Owner:  viper.GetString(FlagOwner),
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryTokens), bz)
			if err != nil {
				return err
			}

			var tokens types.Tokens
			if err := cdc.UnmarshalJSON(res, &tokens); err != nil {
				return err
			}

			return cliCtx.PrintOutput(tokens)
		},
	}
	cmd.Flags().AddFlagSet(FsTokensQuery)
	return cmd
}

// GetCmdQueryFee implements the query asset related fees command.
func GetCmdQueryFee(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "fee [symbol]",
		Short:   "Query the token related fees",
		Example: fmt.Sprintf("%s query asset token fee [symbol]", version.ClientName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// query token fees
			symbol := args[0]
			if err := types.ValidateSymbol(symbol); err != nil {
				return err
			}

			fees, err := queryTokenFees(cliCtx, queryRoute, symbol)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(fees)
		},
	}
	return cmd
}

// GetCmdQueryParams implements the params query command.
func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "params",
		Short:   "Query the current token parameters information",
		Example: fmt.Sprintf("%s query asset token params", version.ClientName),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryParameters)
			bz, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			cdc.MustUnmarshalJSON(bz, &params)
			return cliCtx.PrintOutput(params)
		},
	}
	return cmd
}
