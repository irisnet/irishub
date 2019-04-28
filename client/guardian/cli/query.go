package cli

import (
	"fmt"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/guardian"
	"github.com/spf13/cobra"
)

func GetCmdQueryProfilers(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "profilers",
		Short:   "Query for all profilers",
		Example: "iriscli guardian profilers",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.GuardianRoute, guardian.QueryProfilers), nil)

			if err != nil {
				return err
			}

			var profilers guardian.Profilers
			err = cdc.UnmarshalJSON(res, &profilers)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(profilers)
		},
	}
	return cmd
}

func GetCmdQueryTrustees(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "trustees",
		Short:   "Query for all trustees",
		Example: "iriscli guardian trustees",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.GuardianRoute, guardian.QueryTrustees), nil)

			if err != nil {
				return err
			}

			var trustees guardian.Trustees
			err = cdc.UnmarshalJSON(res, &trustees)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(trustees)
		},
	}
	return cmd
}
