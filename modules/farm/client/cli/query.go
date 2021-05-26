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
		GetCmdQueryFarmer(),
		GetCmdQueryParams(),
	)
	return queryCmd
}

// GetCmdQueryFarmPools implements the query the farm pool.
func GetCmdQueryFarmPools() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pools",
		Example: fmt.Sprintf("$ %s query farm pools --pool-name <Farm Pool Name>", version.AppName),
		Short:   "Query a farm",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			var poolName string
			if len(args) > 0 {
				poolName = args[0]
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Pools(context.Background(), &types.QueryPoolsRequest{
				Name: poolName,
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

// GetCmdQueryFarmer implements the query the farmer reward.
func GetCmdQueryFarmer() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "farmer",
		Example: fmt.Sprintf("$ %s query farm farmer <Farmer Address> --pool-name <Farm Pool Name>", version.AppName),
		Short:   "Query farmer reward",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			poolName, err := cmd.Flags().GetString(FlagFarmPool)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Farmer(context.Background(), &types.QueryFarmerRequest{
				Farmer:   args[0],
				PoolName: poolName,
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
