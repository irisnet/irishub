package cli

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/irisnet/irismod/modules/mt/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                types.ModuleName,
		Short:              "Querying commands for the MT module",
		DisableFlagParsing: true,
	}

	queryCmd.AddCommand(
		GetCmdQueryDenoms(),
		GetCmdQueryDenom(),
		GetCmdQueryMTSupply(),
		GetCmdQueryMTs(),
		GetCmdQueryMT(),
		GetCmdQueryBalances(),
	)

	return queryCmd
}

// GetCmdQueryDenoms queries all denoms
func GetCmdQueryDenoms() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "denoms",
		Long:    "Query all denoms.",
		Example: fmt.Sprintf("$ %s query mt denoms", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Denoms(context.Background(), &types.QueryDenomsRequest{Pagination: pageReq})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "all denoms")
	return cmd
}

// GetCmdQueryDenom queries the specified denom
func GetCmdQueryDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "denom [denom-id]",
		Long:    "Query denom by ID.",
		Example: fmt.Sprintf("$ %s query mt denom <denom-id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Denom(
				context.Background(),
				&types.QueryDenomRequest{DenomId: args[0]},
			)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp.Denom)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryMTSupply queries the total supply of given denom and mt ID
func GetCmdQueryMTSupply() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "supply [denom-id] [mt-id]",
		Long:    "Query total supply of an MT.",
		Example: fmt.Sprintf("$ %s query mt supply <denom-id> <mt-id>", version.AppName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.MTSupply(context.Background(), &types.QueryMTSupplyRequest{
				DenomId: args[0],
				MtId:    args[1],
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryMTs queries all MTs of a denom
func GetCmdQueryMTs() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tokens [denom-id]",
		Long:    "Query all MTs of a denom.",
		Example: fmt.Sprintf("$ %s query mt tokens <denom-id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.MTs(context.Background(), &types.QueryMTsRequest{
				DenomId:    args[0],
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "all tokens")

	return cmd
}

// GetCmdQueryMT queries MT by ID
func GetCmdQueryMT() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "token [denom-id] [mt-id]",
		Long:    "Query MT by ID.",
		Example: fmt.Sprintf("$ %s query mt token <denom-id> <mt-id>", version.AppName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.MT(context.Background(), &types.QueryMTRequest{
				DenomId: args[0],
				MtId:    args[1],
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp.Mt)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryBalances queries the MT balances of a specified owner
func GetCmdQueryBalances() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "balances [owner] [denom-id]",
		Long:    "Query balances of an owner.",
		Example: fmt.Sprintf("$ %s query mt balances <owner> <denom-id>", version.AppName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Balances(context.Background(), &types.QueryBalancesRequest{
				Owner:      args[0],
				DenomId:    args[1],
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "all balances")

	return cmd
}
