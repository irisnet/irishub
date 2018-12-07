package cli

import (
	"fmt"
	"os"

	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/guardian"
	"github.com/spf13/cobra"
)

func GetCmdQueryProfilers(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "profilers",
		Short:   "Query for all profilers",
		Example: "iriscli guardian profilers",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
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
			output, err := codec.MarshalJSONIndent(cdc, profilers)
			fmt.Println(string(output))
			return nil
		},
	}
	return cmd
}

func GetCmdQueryTrustees(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "trustees",
		Short:   "Query for all trustees",
		Example: "iriscli guardian trustees",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			res, err := cliCtx.QuerySubspace(guardian.GetTrusteesSubspaceKey(), storeName)
			if err != nil {
				return err
			}
			var trustees []guardian.Trustee
			for i := 0; i < len(res); i++ {
				var trustee guardian.Trustee
				cdc.MustUnmarshalBinaryLengthPrefixed(res[i].Value, &trustee)
				trustees = append(trustees, trustee)
			}
			output, err := codec.MarshalJSONIndent(cdc, trustees)
			fmt.Println(string(output))
			return nil
		},
	}
	return cmd
}
