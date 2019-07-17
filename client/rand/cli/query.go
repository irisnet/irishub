package cli

import (
	"fmt"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v1/rand"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetCmdQueryRand implements the query-rand command.
func GetCmdQueryRand(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-rand",
		Short:   "query a random number by the request id",
		Example: "iriscli rand query-rand --request-id=<request id>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			reqID := viper.GetString(FlagReqID)
			if err := rand.CheckReqID(reqID); err != nil {
				return err
			}

			params := rand.QueryRandParams{
				ReqID: reqID,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.RandRoute, rand.QueryRand), bz)
			if err != nil {
				return err
			}

			var rand rand.Rand
			err = cdc.UnmarshalJSON(res, &rand)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(rand)
		},
	}

	cmd.Flags().AddFlagSet(FsQueryRand)
	cmd.MarkFlagRequired(FlagReqID)

	return cmd
}

// GetCmdQueryRands implements the query-rands command.
func GetCmdQueryRands(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-rands",
		Short:   "Query all random numbers with an optional consumer",
		Example: "iriscli rand query-rands [--consumer=<consumer>]",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var (
				consumer sdk.AccAddress
				err      error
			)

			consumerStr := viper.GetString(FlagConsumer)
			if consumerStr != "" {
				consumer, err = sdk.AccAddressFromBech32(consumerStr)
				if err != nil {
					return err
				}
			}

			params := rand.QueryRandsParams{
				Consumer: consumer,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.RandRoute, rand.QueryRands), bz)
			if err != nil {
				return err
			}

			var rands []rand.Rand
			err = cdc.UnmarshalJSON(res, &rands)
			if err != nil {
				return err
			}

			// TODO
			// return cliCtx.PrintOutput(rands)
			return nil
		},
	}

	cmd.Flags().AddFlagSet(FsQueryRands)

	return cmd
}

// GetCmdQueryRequest implements the query-request command.
func GetCmdQueryRequest(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-request",
		Short:   "Query a random number request by the given request id",
		Example: "iriscli rand query-request --request-id=<request id>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			reqID := viper.GetString(FlagReqID)
			if err := rand.CheckReqID(reqID); err != nil {
				return err
			}

			params := rand.QueryRandRequestParams{
				ReqID: reqID,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.RandRoute, rand.QueryRandRequest), bz)
			if err != nil {
				return err
			}

			var request rand.Request
			err = cdc.UnmarshalJSON(res, &request)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(request)
		},
	}

	cmd.Flags().AddFlagSet(FsQueryRequest)
	cmd.MarkFlagRequired(FlagReqID)

	return cmd
}

// GetCmdQueryRequests implements the query-requests command.
func GetCmdQueryRequests(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-requests",
		Short:   "Query all random number requests with an optional consumer",
		Example: "iriscli rand query-requests [--consumer=<consumer>]",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var (
				consumer sdk.AccAddress
				err      error
			)

			consumerStr := viper.GetString(FlagConsumer)
			if consumerStr != "" {
				consumer, err = sdk.AccAddressFromBech32(consumerStr)
				if err != nil {
					return err
				}
			}

			params := rand.QueryRandRequestsParams{
				Consumer: consumer,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.RandRoute, rand.QueryRandRequests), bz)
			if err != nil {
				return err
			}

			var requests []rand.Request
			err = cdc.UnmarshalJSON(res, &requests)
			if err != nil {
				return err
			}

			// TODO
			// return cliCtx.PrintOutput(requests)
			return nil
		},
	}

	cmd.Flags().AddFlagSet(FsQueryRequests)

	return cmd
}

// GetCmdQueryQueue implements the query-queue command.
func GetCmdQueryQueue(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-queue",
		Short:   "Query the random number request queue with an optional height",
		Example: "iriscli rand query-queue [--height=<height>]",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := rand.QueryRandRequestQueueParams{
				Height: viper.GetInt64(FlagHeight),
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.RandRoute, rand.QueryRandRequestQueue), bz)
			if err != nil {
				return err
			}

			var requests []rand.Request
			err = cdc.UnmarshalJSON(res, &requests)
			if err != nil {
				return err
			}

			// TODO
			// return cliCtx.PrintOutput(requests)
			return nil
		},
	}

	cmd.Flags().AddFlagSet(FsQueryQueue)

	return cmd
}
