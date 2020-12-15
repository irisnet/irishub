package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/irisnet/irismod/modules/token/types"
)

// GetQueryCmd returns the query commands for the token module.
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                types.ModuleName,
		Short:              "Querying commands for the token module",
		DisableFlagParsing: true,
	}

	queryCmd.AddCommand(
		GetCmdQueryToken(),
		GetCmdQueryTokens(),
		GetCmdQueryFee(),
		GetCmdQueryParams(),
	)

	return queryCmd
}

// GetCmdQueryToken implements the query token command.
func GetCmdQueryToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "token [denom]",
		Long:    "Query a token by symbol or minUnit.",
		Example: fmt.Sprintf("$ %s query token token <denom>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())

			if err != nil {
				return err
			}

			if err := types.CheckSymbol(args[0]); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Token(context.Background(), &types.QueryTokenRequest{
				Denom: args[0],
			})

			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res.Token)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryTokens implements the query tokens command.
func GetCmdQueryTokens() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tokens [owner]",
		Long:    "Query token by the owner.",
		Example: fmt.Sprintf("$ %s query token tokens <owner>", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			var owner sdk.AccAddress
			if len(args) > 0 {
				owner, err = sdk.AccAddressFromBech32(args[0])
				if err != nil {
					return err
				}
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			res, err := queryClient.Tokens(
				context.Background(),
				&types.QueryTokensRequest{
					Owner:      owner.String(),
					Pagination: pageReq,
				},
			)
			if err != nil {
				return err
			}

			tokens := make([]types.TokenI, 0, len(res.Tokens))
			for _, eviAny := range res.Tokens {
				var evi types.TokenI
				if err = clientCtx.InterfaceRegistry.UnpackAny(eviAny, &evi); err != nil {
					return err
				}
				tokens = append(tokens, evi)
			}

			return clientCtx.PrintOutputLegacy(tokens)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "all tokens")

	return cmd
}

// GetCmdQueryFee implements the query token related fees command.
func GetCmdQueryFee() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "fee [symbol]",
		Args:    cobra.ExactArgs(1),
		Long:    "Query the token related fees.",
		Example: fmt.Sprintf("$ %s query token fee <symbol>", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			symbol := args[0]
			if err := types.CheckSymbol(symbol); err != nil {
				return err
			}

			// query token fees
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Fees(
				context.Background(),
				&types.QueryFeesRequest{
					Symbol: symbol,
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryParams implements the query token related param command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "params",
		Long:    "Query values set as token parameters.",
		Example: fmt.Sprintf("$ %s query token params", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(&res.Params)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
