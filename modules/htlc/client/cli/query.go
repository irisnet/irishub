package cli

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/irisnet/irismod/modules/htlc/types"
)

// GetQueryCmd returns the cli query commands for the module.
func GetQueryCmd() *cobra.Command {
	htlcQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the HTLC module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	htlcQueryCmd.AddCommand(
		GetCmdQueryHTLC(),
		GetCmdQueryAssetSupply(),
		GetCmdQueryAssetSupplies(),
		GetCmdQueryParams(),
	)

	return htlcQueryCmd
}

// GetCmdQueryHTLC implements the query HTLC command.
func GetCmdQueryHTLC() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "htlc [id]",
		Short:   "Query an HTLC",
		Long:    "Query details of an HTLC with the specified id.",
		Example: fmt.Sprintf("$ %s query htlc htlc <id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			id, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			param := types.QueryHTLCRequest{Id: tmbytes.HexBytes(id).String()}
			response, err := queryClient.HTLC(context.Background(), &param)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(response.Htlc)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryAssetSupply queries as asset's current in swap supply, active, supply, and supply limit
func GetCmdQueryAssetSupply() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "supply [denom]",
		Short:   "Query information about an asset's supply",
		Example: fmt.Sprintf("$ %s query htlc supply <htltbcbnb>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			param := types.QueryAssetSupplyRequest{Denom: args[0]}
			response, err := queryClient.AssetSupply(context.Background(), &param)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(response.AssetSupply)

		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryAssetSupplies queries AssetSupplies in the store
func GetCmdQueryAssetSupplies() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "supplies",
		Short:   "Query the list of all asset supplies",
		Example: fmt.Sprintf("$ %s query htlc supplies", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			param := types.QueryAssetSuppliesRequest{}
			response, err := queryClient.AssetSupplies(context.Background(), &param)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(response)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "params",
		Short:   "Query the current htlc parameter values",
		Long:    "Query values set as htlc parameters.",
		Example: fmt.Sprintf("$ %s query htlc params", version.AppName),
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
