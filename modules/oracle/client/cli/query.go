package cli

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"mods.irisnet.org/modules/oracle/types"
)

// GetQueryCmd returns the cli query commands for the oracle module.
func GetQueryCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the oracle module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		GetCmdQueryFeed(),
		GetCmdQueryFeeds(),
		GetCmdQueryFeedValue(),
	)
	return txCmd
}

// GetCmdQueryFeed implements the query feed definition command
func GetCmdQueryFeed() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "feed [feed-name]",
		Short:   "Query the feed definition.",
		Example: fmt.Sprintf("%s query oracle feed <feed-name>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Feed(context.Background(), &types.QueryFeedRequest{FeedName: args[0]})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Feed)
		},
	}
	cmd.Flags().AddFlagSet(FsQueryFeed)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryFeed implements the query feed definitions command
func GetCmdQueryFeeds() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "feeds",
		Short:   "Query a group of feed definitions.",
		Example: fmt.Sprintf("%s query oracle feeds", version.AppName),
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
			feedState, err := cmd.Flags().GetString(FlagFeedState)
			if err != nil {
				return err
			}

			res, err := queryClient.Feeds(context.Background(), &types.QueryFeedsRequest{State: feedState, Pagination: pageReq})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	cmd.Flags().AddFlagSet(FsQueryFeeds)
	flags.AddPaginationFlagsToCmd(cmd, "all feeds")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryFeedValue implements the query feed value command
func GetCmdQueryFeedValue() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "value [feed-name]",
		Short:   "Query the feed result.",
		Example: fmt.Sprintf("%s query oracle value <feed-name>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.FeedValue(context.Background(), &types.QueryFeedValueRequest{FeedName: args[0]})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	cmd.Flags().AddFlagSet(FsQueryFeedValue)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
