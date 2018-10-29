package cli

import (
	"github.com/irisnet/irishub/modules/iservice"
	"github.com/spf13/cobra"
	"os"
	"github.com/spf13/viper"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"fmt"
	cmn "github.com/irisnet/irishub/client/iservice"
)

func GetCmdQueryScvDef(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "definition",
		Short: "query service definition",
		Example: "iriscli iservice definition --name=<service name> --chain-id=<chain-id>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			name := viper.GetString(FlagServiceName)
			chainId := viper.GetString(client.FlagChainID)
			res, err := cliCtx.QueryStore(iservice.GetServiceDefinitionKey(chainId, name), storeName)
			if len(res) == 0 || err != nil {
				return fmt.Errorf("service [%s] is not existed", name)
			}

			var msgSvcDef iservice.MsgSvcDef
			cdc.MustUnmarshalBinary(res, &msgSvcDef)

			res2, err := cliCtx.QuerySubspace(iservice.GetMethodsSubspaceKey(chainId, name), storeName)
			if err != nil {
				return err
			}

			var methods []iservice.MethodProperty
			for i := 0; i < len(res2); i++ {
				var method iservice.MethodProperty
				cdc.MustUnmarshalBinary(res2[i].Value, &method)
				methods = append(methods, method)
			}

			service := cmn.ServiceOutput{MsgSvcDef: msgSvcDef, Methods: methods}
			output, err := wire.MarshalJSONIndent(cdc, service)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}
	cmd.Flags().AddFlagSet(FsServiceName)

	return cmd
}
