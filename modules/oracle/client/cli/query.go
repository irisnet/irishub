package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/irisnet/irishub/modules/oracle/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	txCmd.AddCommand(flags.GetCommands(
		GetCmdQueryFeed(clientCtx),
		GetCmdQueryFeeds(clientCtx),
		GetCmdQueryFeedValue(clientCtx),
	)...)
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

			params := types.QueryFeedParams{
				FeedName: args[0],
			}

			bz, err := clientCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryFeed)
			res, _, err := clientCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var feedCtx types.FeedContext
			if err := clientCtx.Codec.UnmarshalJSON(res, &feedCtx); err != nil {
				return err
			}

			return clientCtx.PrintOutput(feedCtx)
		},
	}
	cmd.Flags().AddFlagSet(FsQueryFeed)
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

			params := types.QueryFeedsParams{
				State: viper.GetString(FlagFeedState),
			}

			bz, err := clientCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryFeeds)
			res, _, err := clientCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var feedCtx types.FeedsContext
			if err := clientCtx.Codec.UnmarshalJSON(res, &feedCtx); err != nil {
				return err
			}

			return clientCtx.PrintOutput(feedCtx)
		},
	}
	cmd.Flags().AddFlagSet(FsQueryFeeds)
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

			params := types.QueryFeedParams{
				FeedName: args[0],
			}

			bz, err := clientCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryFeedValue)
			res, _, err := clientCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var feedValue types.FeedValues
			if err := clientCtx.Codec.UnmarshalJSON(res, &feedValue); err != nil {
				return err
			}

			return clientCtx.PrintOutput(feedValue)
		},
	}
	cmd.Flags().AddFlagSet(FsQueryFeedValue)
	return cmd
}
