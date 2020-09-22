package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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
		Use: "supply [denomID]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`total supply of a collection or owner of NFTs.
Example:
$ %s query nft supply [denom]`, version.AppName)),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			var owner sdk.AccAddress

			ownerStr := strings.TrimSpace(viper.GetString(FlagOwner))
			if len(ownerStr) > 0 {
				owner, err = sdk.AccAddressFromBech32(ownerStr)
				if err != nil {
					return err
				}
			}

			denom := strings.TrimSpace(args[0])
			if err := types.ValidateDenomID(denom); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Supply(context.Background(), &types.QuerySupplyRequest{
				Denom: denom,
				Owner: owner,
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
		Use: "owner [address]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get the NFTs owned by an account address
Example:
$ %s query nft owner <address> --denom=<denom>`, version.AppName)),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			denom := viper.GetString(FlagDenom)
			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Owner(context.Background(), &types.QueryOwnerRequest{
				Denom: denom,
				Owner: address,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(resp.Owner)
		},
	}
	cmd.Flags().AddFlagSet(FsQueryOwner)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryCollection queries all the NFTs from a collection
func GetCmdQueryCollection() *cobra.Command {
	cmd := &cobra.Command{
		Use: "collection [denomID]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get all the NFTs from a given collection
Example:
$ %s query nft collection <denom>`, version.AppName)),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			denom := strings.TrimSpace(args[0])
			if err := types.ValidateDenomID(denom); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Collection(context.Background(), &types.QueryCollectionRequest{
				Denom: denom,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintOutput(resp.Collection)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryDenoms queries all denoms
func GetCmdQueryDenoms() *cobra.Command {
	cmd := &cobra.Command{
		Use: "denoms",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all denominations of all collections of NFTs
Example:
$ %s query nft denoms`, version.AppName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Denoms(context.Background(), &types.QueryDenomsRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintOutput(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryDenoms queries the specified denoms
func GetCmdQueryDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use: "denom [denomID]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the denominations by the specified denmo name
Example:
$ %s query nft denom <denom>`, version.AppName)),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			denom := strings.TrimSpace(args[0])
			if err := types.ValidateDenomID(denom); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Denom(context.Background(), &types.QueryDenomRequest{
				Denom: denom,
			})
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
		Use: "token [denomID] [tokenID]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query a single NFT from a collection
Example:
$ %s query nft token <denom> <tokenID>`, version.AppName)),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			denom := strings.TrimSpace(args[0])
			if err := types.ValidateDenomID(denom); err != nil {
				return err
			}

			tokenID := strings.TrimSpace(args[1])
			if err := types.ValidateTokenID(tokenID); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.NFT(context.Background(), &types.QueryNFTRequest{
				Denom: denom,
				Id:    tokenID,
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
