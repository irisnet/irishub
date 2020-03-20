package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v3/rand"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/codec"
)

// GetCmdQueryRand implements the query-rand command.
func GetCmdQueryRand(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-rand",
		Short:   "Query a random number by the request id",
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

			var random rand.Rand
			if err = cdc.UnmarshalJSON(res, &random); err != nil {
				return err
			}

			return cliCtx.PrintOutput(random)
		},
	}

	cmd.Flags().AddFlagSet(FsQueryRand)
	_ = cmd.MarkFlagRequired(FlagReqID)

	return cmd
}

// GetCmdQueryRandRequestQueue implements the query-queue command.
func GetCmdQueryRandRequestQueue(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-queue",
		Short:   "Query the random number request queue with an optional height",
		Example: "iriscli rand query-queue [--queue-height=<height>]",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			height := viper.GetInt64(FlagQueueHeight)
			if height < 0 {
				return fmt.Errorf("the height must not be less than 0: %d", height)
			}

			params := rand.QueryRandRequestQueueParams{
				Height: height,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.RandRoute, rand.QueryRandRequestQueue), bz)
			if err != nil {
				return err
			}

			var requests rand.Requests
			if err = cdc.UnmarshalJSON(res, &requests); err != nil {
				return err
			}

			return cliCtx.PrintOutput(requests)
		},
	}

	cmd.Flags().AddFlagSet(FsQueryQueue)

	return cmd
}
