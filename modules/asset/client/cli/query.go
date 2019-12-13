package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/irisnet/irishub/modules/asset/internal/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetQueryCmd returns the query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                types.ModuleName,
		Short:              "Querying commands for the asset module",
		DisableFlagParsing: true,
	}

	queryCmd.AddCommand(client.GetCommands(
		GetCmdQueryToken(queryRoute, cdc),
		GetCmdQueryTokens(queryRoute, cdc),
		GetCmdQueryFee(queryRoute, cdc),
		GetCmdQueryParams(queryRoute, cdc),
	)...)

	return queryCmd
}

// GetCmdQueryToken implements the query token command.
func GetCmdQueryToken(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "token [token-id]",
		Short:   "Query details of a token",
		Example: fmt.Sprintf("%s query asset token [token-id]", version.ClientName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := types.QueryTokenParams{
				TokenID: args[0],
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryToken), bz)
			if err != nil {
				return err
			}

			var token types.FungibleToken
			err = cdc.UnmarshalJSON(res, &token)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(token)
		},
	}
	return cmd
}

// GetCmdQueryTokens implements the query tokens command.
func GetCmdQueryTokens(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tokens",
		Short: "Query details of a group of tokens",
		Example: fmt.Sprintf("%s asset tokens --source=[native|external] --gateway=[gateway_moniker] "+
			"--owner=[address]", version.ClientName),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := types.QueryTokensParams{
				Source: viper.GetString(FlagSource),
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

// GetCmdQueryFee implements the query asset related fees command.
func GetCmdQueryFee(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "fee [token-id]",
		Short:   "Query the asset related fees",
		Example: fmt.Sprintf("%s query asset fee [token-id]", version.ClientName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// query token fees
			tokenID := args[0]
			if err := types.CheckTokenID(tokenID); err != nil {
				return err
			}

			fees, err := queryTokenFees(cliCtx, queryRoute, tokenID)
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
		Short:   "Query the current asset parameters information",
		Example: fmt.Sprintf("%s query asset params", version.ClientName),
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
