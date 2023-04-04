package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/irisnet/irismod/modules/farm/types"
)

// GetQueryCmd returns the cli query commands for the farm module.
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the farm module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	queryCmd.AddCommand(
		GetCmdQueryFarmPools(),
		GetCmdQueryFarmPool(),
		GetCmdQueryFarmer(),
		GetCmdQueryParams(),
	)
	return queryCmd
}

// GetCmdQueryFarmPools implements the query the farm pool by page.
func GetCmdQueryFarmPools() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pools",
		Example: fmt.Sprintf("$ %s query farm pools", version.AppName),
		Short:   "Query farm pools by page",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			resp, err := queryClient.FarmPools(context.Background(), &types.QueryFarmPoolsRequest{
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "pools")
	return cmd
}

// GetCmdQueryFarmPools implements the query a farm pool.
func GetCmdQueryFarmPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pool",
		Example: fmt.Sprintf("$ %s query farm pool <Farm Pool ID>", version.AppName),
		Short:   "Query a farm pool",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.FarmPool(context.Background(), &types.QueryFarmPoolRequest{
				Id: args[0],
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

// GetCmdQueryFarmer implements the query the farmer reward.
func GetCmdQueryFarmer() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "farmer",
		Example: fmt.Sprintf("$ %s query farm farmer <Farmer Address> --pool-id <Farm Pool Id>", version.AppName),
		Short:   "Query farmer reward",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			poolId, err := cmd.Flags().GetString(FlagFarmPool)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Farmer(context.Background(), &types.QueryFarmerRequest{
				Farmer: args[0],
				PoolId: poolId,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}

	cmd.Flags().AddFlagSet(FsQueryFarmPool)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "params",
		Short:   "Query the current farm parameter values",
		Example: fmt.Sprintf("$ %s query farm params", version.AppName),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
