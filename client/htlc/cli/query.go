package cli

import (
	"fmt"
	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v2/htlc"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/codec"
	"github.com/spf13/cobra"
)

// GetCmdQueryHtlc implements the query htlc command.
func GetCmdQueryHtlc(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-htlc",
		Short:   "Query details of a htlc",
		Example: "iriscli htlc query-htlc <hash-lock>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := htlc.QueryHTLCParams{
				SecretHashLock: args[0],
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.HtlcRoute, htlc.QueryHTLC), bz)
			if err != nil {
				return err
			}

			var htlc htlc.HTLC
			err = cdc.UnmarshalJSON(res, &htlc)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(htlc)
		},
	}

	return cmd
}
