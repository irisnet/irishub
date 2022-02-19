package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		GetCmdQueryDenom(),
		GetCmdQueryDenoms(),
		GetCmdQueryCollection(),
		GetCmdQuerySupply(),
		GetCmdQueryOwner(),
		GetCmdQueryMT(),
	)

	return queryCmd
}

// GetCmdQuerySupply queries the supply of a mt collection
func GetCmdQuerySupply() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "supply [denom-id]",
		Long:    "total supply of a collection or owner of MTs.",
		Example: fmt.Sprintf("$ %s query mt supply <denom-id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var owner sdk.AccAddress
			ownerStr, err := cmd.Flags().GetString(FlagOwner)
			if err != nil {
				return err
			}

			if len(ownerStr) > 0 {
				owner, err = sdk.AccAddressFromBech32(ownerStr)
				if err != nil {
					return err
				}
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Supply(context.Background(), &types.QuerySupplyRequest{
				DenomId: args[0],
				Owner:   owner.String(),
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	cmd.Flags().AddFlagSet(FsQuerySupply)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryOwner queries all the MTs owned by an account
func GetCmdQueryOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "owner [address]",
		Long:    "Get the MTs owned by an account address.",
		Example: fmt.Sprintf("$ %s query mt owner <address> --denom-id=<denom-id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if _, err := sdk.AccAddressFromBech32(args[0]); err != nil {
				return err
			}
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			denomID, err := cmd.Flags().GetString(FlagDenomID)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Owner(context.Background(), &types.QueryOwnerRequest{
				DenomId:    denomID,
				Owner:      args[0],
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	cmd.Flags().AddFlagSet(FsQueryOwner)
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "mts")

	return cmd
}

// GetCmdQueryCollection queries all the MTs from a collection
func GetCmdQueryCollection() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "collection [denom-id]",
		Long:    "Get all the MTs from a given collection.",
		Example: fmt.Sprintf("$ %s query mt collection <denom-id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Collection(
				context.Background(),
				&types.QueryCollectionRequest{
					DenomId:    args[0],
					Pagination: pageReq,
				},
			)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "mts")

	return cmd
}

// GetCmdQueryDenoms queries all denoms
func GetCmdQueryDenoms() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "denoms",
		Long:    "Query all denominations of all collections of MTs.",
		Example: fmt.Sprintf("$ %s query mt denoms", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
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
		Long:    "Query the denom by the specified denom id.",
		Example: fmt.Sprintf("$ %s query mt denom <denom-id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
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

// GetCmdQueryMT queries a single MTs from a collection
func GetCmdQueryMT() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "token [denom-id] [mt-id]",
		Long:    "Query a single MT from a collection.",
		Example: fmt.Sprintf("$ %s query mt token <denom-id> <mt-id>", version.AppName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if err := types.ValidateTokenID(args[1]); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.MT(context.Background(), &types.QueryMTRequest{
				DenomId: args[0],
				TokenId: args[1],
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp.MT)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
