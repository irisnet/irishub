package cli

import (
	"os"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/client/context"
	authcmd "github.com/irisnet/irishub/client/auth/cli"
	"github.com/irisnet/irishub/modules/profiling"
)

func GetCmdQueryProfilers(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "profilers",
		Short:   "Query for all profilers",
		Example: "iriscli profiling profilers",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			res, err := cliCtx.QuerySubspace(profiling.GetProfilersSubspaceKey(), storeName)
			if err != nil {
				return err
			}
			var profilers []profiling.Profiler
			for i := 0; i < len(res); i++ {
				var profiler profiling.Profiler
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
