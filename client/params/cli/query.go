package cli

import (
	"fmt"
	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/codec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

const flagModule = "module"

func Commands(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "params",
		Short:   "Query parameter",
		Example: "iriscli params --module=<module name>",
		RunE: func(cmd *cobra.Command, args []string) error {
			moduleStr := strings.TrimSpace(viper.GetString(flagModule))

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			param := params.QueryModuleParams{
				Module: moduleStr,
			}
			bz, err := cdc.MarshalJSON(param)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/module", protocol.ParamsRoute), bz)
			if err != nil {
				return err
			}

			// query all
			if len(moduleStr) == 0 {
				var paramSets params.ParamSets
				if err := cdc.UnmarshalJSON(res, &paramSets); err != nil {
					return err
				}
				return cliCtx.PrintOutput(paramSets)
			}
			// query by module
			var paramSet params.ParamSet
			if err := cdc.UnmarshalJSON(res, &paramSet); err != nil {
				return err
			}
			return cliCtx.PrintOutput(paramSet)
		},
	}

	cmd.Flags().String(flagModule, "", "module name can be stake/mint/distr/slashing/service/asset/auth")
	return cmd
}
