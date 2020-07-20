package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/irisnet/irishub/modules/oracle/types"
)

// GetQueryCmd returns the cli query commands for the guardian module.
func GetQueryCmd(clientCtx client.Context) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the oracle module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		GetCmdQueryFeed(clientCtx),
		GetCmdQueryFeeds(clientCtx),
		GetCmdQueryFeedValue(clientCtx),
	)
	return txCmd
}

// GetCmdQueryFeed implements the query feed Content definition command
func GetCmdQueryFeed(clientCtx client.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-feed [feed-name]",
		Short:   "Query the feed definition",
		Example: "iriscli oracle query-feed <feed-name>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Feed(context.Background(), &types.QueryFeedRequest{FeedName: args[0]})
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res.Feed)
		},
	}
	cmd.Flags().AddFlagSet(FsQueryFeed)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryFeed implements the query feed Content definition command
func GetCmdQueryFeeds(clientCtx client.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-feeds",
		Short:   "Query a group of feed definition",
		Example: "iriscli oracle query-feeds",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Feeds(context.Background(), &types.QueryFeedsRequest{State: viper.GetString(FlagFeedState)})
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res.Feeds)
		},
	}
	cmd.Flags().AddFlagSet(FsQueryFeeds)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryFeedValue implements the query feed value command
func GetCmdQueryFeedValue(clientCtx client.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-value [feed-name]",
		Short:   "Query the feed result",
		Example: "iriscli oracle query-value <feed-name>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.FeedValue(context.Background(), &types.QueryFeedValueRequest{FeedName: args[0]})
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res.FeedValues)
		},
	}
	cmd.Flags().AddFlagSet(FsQueryFeedValue)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
