package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/irisnet/irismod/token/types"
	v1 "github.com/irisnet/irismod/token/types/v1"
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
		GetCmdQueryTotalBurn(),
		GetCmdQueryParams(),
		GetCmdQueryBalances(),
	)

	return queryCmd
}

// GetCmdQueryToken implements the query token command.
func GetCmdQueryToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "token [denom]",
		Long:    "Query a token by symbol or min unit.",
		Example: fmt.Sprintf("$ %s query token token <denom>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			if err := types.ValidateSymbol(args[0]); err != nil {
				return err
			}

			queryClient := v1.NewQueryClient(clientCtx)

			res, err := queryClient.Token(context.Background(), &v1.QueryTokenRequest{
				Denom: args[0],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res.Token)
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
			clientCtx, err := client.GetClientQueryContext(cmd)
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

			queryClient := v1.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			res, err := queryClient.Tokens(
				context.Background(),
				&v1.QueryTokensRequest{
					Owner:      owner.String(),
					Pagination: pageReq,
				},
			)
			if err != nil {
				return err
			}

			tokens := make([]v1.TokenI, 0, len(res.Tokens))
			for _, eviAny := range res.Tokens {
				var evi v1.TokenI
				if err = clientCtx.InterfaceRegistry.UnpackAny(eviAny, &evi); err != nil {
					return err
				}
				tokens = append(tokens, evi)
			}

			return clientCtx.PrintObjectLegacy(tokens)
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
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			symbol := args[0]
			if err := types.ValidateSymbol(symbol); err != nil {
				return err
			}

			// query token fees
			queryClient := v1.NewQueryClient(clientCtx)
			res, err := queryClient.Fees(
				context.Background(),
				&v1.QueryFeesRequest{
					Symbol: symbol,
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
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
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := v1.NewQueryClient(clientCtx)
			res, err := queryClient.Params(context.Background(), &v1.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryTotalBurn return the total amount of all burned tokens
func GetCmdQueryTotalBurn() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "total-burn",
		Long:    "Query the total amount of all burned tokens.",
		Example: fmt.Sprintf("$ %s query token params", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := v1.NewQueryClient(clientCtx)
			res, err := queryClient.TotalBurn(context.Background(), &v1.QueryTotalBurnRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryBalances return all the balances of an owner in special denom
func GetCmdQueryBalances() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "balances [addr] [denom]",
		Args:    cobra.ExactArgs(2),
		Long:    "Query all the balances of an owner in special denom.",
		Example: fmt.Sprintf("$ %s query token balances <addr> <denom>", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := v1.NewQueryClient(clientCtx)
			res, err := queryClient.Balances(context.Background(), &v1.QueryBalancesRequest{
				Address: args[0],
				Denom:   args[1],
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
