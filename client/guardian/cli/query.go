package cli

import (
	"os"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/client/context"
	authcmd "github.com/irisnet/irishub/client/auth/cli"
	"github.com/irisnet/irishub/modules/guardian"
)

func GetCmdQueryProfilers(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "guardian",
		Short:   "Query for all profilers",
		Example: "iriscli guardian profilers",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			res, err := cliCtx.QuerySubspace(guardian.GetProfilersSubspaceKey(), storeName)
			if err != nil {
				return err
			}
			var profilers []guardian.Profiler
			for i := 0; i < len(res); i++ {
				var profiler guardian.Profiler
				cdc.MustUnmarshalBinaryLengthPrefixed(res[i].Value, &profiler)
				profilers = append(profilers, profiler)
			}
			output, err := codec.MarshalJSONIndent(cdc,profilers)
			fmt.Println(string(output))
			return nil
		},
	}
	return cmd
}

func GetCmdQueryTrustees(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "guardian",
		Short:   "Query for all trustees",
		Example: "iriscli guardian trustees",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			res, err := cliCtx.QuerySubspace(guardian.GetTrusteesSubspaceKey(), storeName)
			if err != nil {
				return err
			}
			var profilers []guardian.Profiler
			for i := 0; i < len(res); i++ {
				var profiler guardian.Profiler
				cdc.MustUnmarshalBinaryLengthPrefixed(res[i].Value, &profiler)
				profilers = append(profilers, profiler)
			}
			output, err := codec.MarshalJSONIndent(cdc,profilers)
			fmt.Println(string(output))
			return nil
		},
	}
	return cmd
}
