package cli

import (
	"github.com/irisnet/irishub/modules/iservice"
	"github.com/spf13/cobra"
	"os"
	"github.com/spf13/viper"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/irisnet/irishub/client/context"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"fmt"
	cmn "github.com/irisnet/irishub/client/iservice"
)

func GetCmdQueryScvDef(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "definition",
		Short:   "query service definition",
		Example: "iriscli iservice definition --def-chain-id=<chain-id> --service-name=<service name>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			name := viper.GetString(FlagServiceName)
			defChainId := viper.GetString(FlagDefChainID)

			res, err := cliCtx.QueryStore(iservice.GetServiceDefinitionKey(defChainId, name), storeName)
			if len(res) == 0 || err != nil {
				return fmt.Errorf("chain-id [%s] service [%s] is not existed", defChainId, name)
			}

			var svcDef iservice.SvcDef
			cdc.MustUnmarshalBinary(res, &svcDef)

			res2, err := cliCtx.QuerySubspace(iservice.GetMethodsSubspaceKey(defChainId, name), storeName)
			if err != nil {
				return err
			}

			var methods []iservice.MethodProperty
			for i := 0; i < len(res2); i++ {
				var method iservice.MethodProperty
				cdc.MustUnmarshalBinary(res2[i].Value, &method)
				methods = append(methods, method)
			}

			service := cmn.ServiceOutput{SvcDef: svcDef, Methods: methods}
			output, err := codec.MarshalJSONIndent(cdc, service)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}
	cmd.Flags().AddFlagSet(FsDefChainID)
	cmd.Flags().AddFlagSet(FsServiceName)

	return cmd
}
