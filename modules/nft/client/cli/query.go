package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/irisnet/irismod/modules/nft/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                types.ModuleName,
		Short:              "Querying commands for the NFT module",
		DisableFlagParsing: true,
	}

	queryCmd.AddCommand(
		GetCmdQueryDenom(),
		GetCmdQueryDenoms(),
		GetCmdQueryCollection(),
		GetCmdQuerySupply(),
		GetCmdQueryOwner(),
		GetCmdQueryNFT(),
	)

	return queryCmd
}

// GetCmdQuerySupply queries the supply of a nft collection
func GetCmdQuerySupply() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "supply [denom-id]",
		Long:    "total supply of a collection or owner of NFTs.",
		Example: fmt.Sprintf("$ %s query nft supply <denom-id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			var owner sdk.AccAddress
			rawOwner, err := cmd.Flags().GetString(FlagOwner)
			if err != nil {
				return err
			}
			ownerStr := strings.TrimSpace(rawOwner)
			if len(ownerStr) > 0 {
				owner, err = sdk.AccAddressFromBech32(ownerStr)
				if err != nil {
					return err
				}
			}

			denomID := strings.TrimSpace(args[0])
			if err := types.ValidateDenomID(denomID); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Supply(context.Background(), &types.QuerySupplyRequest{
				DenomId: denomID,
				Owner:   owner.String(),
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintOutput(resp)
		},
	}
	cmd.Flags().AddFlagSet(FsQuerySupply)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryOwner queries all the NFTs owned by an account
func GetCmdQueryOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "owner [address]",
		Long:    "Get the NFTs owned by an account address.",
		Example: fmt.Sprintf("$ %s query nft owner <address> --denom-id=<denom-id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
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
			return clientCtx.PrintOutput(resp)
		},
	}
	cmd.Flags().AddFlagSet(FsQueryOwner)
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "nfts")

	return cmd
}

// GetCmdQueryCollection queries all the NFTs from a collection
func GetCmdQueryCollection() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "collection [denom-id]",
		Long:    "Get all the NFTs from a given collection.",
		Example: fmt.Sprintf("$ %s query nft collection <denom-id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			denomID := strings.TrimSpace(args[0])
			if err := types.ValidateDenomID(denomID); err != nil {
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
					DenomId:    denomID,
					Pagination: pageReq,
				},
			)
			if err != nil {
				return err
			}
			return clientCtx.PrintOutput(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "nfts")

	return cmd
}

// GetCmdQueryDenoms queries all denoms
func GetCmdQueryDenoms() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "denoms",
		Long:    "Query all denominations of all collections of NFTs.",
		Example: fmt.Sprintf("$ %s query nft denoms", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
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
			return clientCtx.PrintOutput(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "all denoms")
	return cmd
}

// GetCmdQueryDenoms queries the specified denoms
func GetCmdQueryDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "denom [denom-id]",
		Long:    "Query the denominations by the specified denmo name.",
		Example: fmt.Sprintf("$ %s query nft denom <denom-id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			denomID := strings.TrimSpace(args[0])
			if err := types.ValidateDenomID(denomID); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Denom(
				context.Background(),
				&types.QueryDenomRequest{DenomId: denomID},
			)
			if err != nil {
				return err
			}
			return clientCtx.PrintOutput(resp.Denom)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryNFT queries a single NFTs from a collection
func GetCmdQueryNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "token [denom-id] [token-id]",
		Long:    "Query a single NFT from a collection.",
		Example: fmt.Sprintf("$ %s query nft token <denom-id> <token-id>", version.AppName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			denomID := strings.TrimSpace(args[0])
			if err := types.ValidateDenomID(denomID); err != nil {
				return err
			}

			tokenID := strings.TrimSpace(args[1])
			if err := types.ValidateTokenID(tokenID); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.NFT(context.Background(), &types.QueryNFTRequest{
				DenomId: denomID,
				TokenId: tokenID,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintOutput(resp.NFT)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
