package cli

import (
	"fmt"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v3/oracle"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/codec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetCmdQueryFeed implements the query feed Content definition command
func GetCmdQueryFeed(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-feed [feed-name]",
		Short:   "Query the feed definition",
		Example: "iriscli oracle query-feed <feed-name>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := oracle.QueryFeedParams{
				FeedName: args[0],
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", protocol.OracleRoute, oracle.QueryFeed)
			res, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var feedCtx oracle.FeedContext
			if err := cdc.UnmarshalJSON(res, &feedCtx); err != nil {
				return err
			}

			return cliCtx.PrintOutput(feedCtx)
		},
	}
	cmd.Flags().AddFlagSet(FsQueryFeed)
	return cmd
}

// GetCmdQueryFeed implements the query feed Content definition command
func GetCmdQueryFeeds(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-feeds",
		Short:   "Query a group of feed definition",
		Example: "iriscli oracle query-feeds",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := oracle.QueryFeedsParams{
				State: viper.GetString(FlagFeedState),
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", protocol.OracleRoute, oracle.QueryFeeds)
			res, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var feedCtx oracle.FeedsContext
			if err := cdc.UnmarshalJSON(res, &feedCtx); err != nil {
				return err
			}

			return cliCtx.PrintOutput(feedCtx)
		},
	}
	cmd.Flags().AddFlagSet(FsQueryFeeds)
	return cmd
}

// GetCmdQueryFeedValue implements the query feed value command
func GetCmdQueryFeedValue(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-value [feed-name]",
		Short:   "Query the feed result",
		Example: "iriscli oracle query-value <feed-name>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := oracle.QueryFeedParams{
				FeedName: args[0],
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", protocol.OracleRoute, oracle.QueryFeedValue)
			res, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var feedValue oracle.FeedValues
			if err := cdc.UnmarshalJSON(res, &feedValue); err != nil {
				return err
			}

			return cliCtx.PrintOutput(feedValue)
		},
	}
	cmd.Flags().AddFlagSet(FsQueryFeedValue)
	return cmd
}
