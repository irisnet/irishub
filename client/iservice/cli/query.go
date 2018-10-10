package cli

import (
	"github.com/irisnet/irishub/modules/iservice"
	"github.com/spf13/cobra"
	"os"
	"github.com/spf13/viper"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"fmt"
)

func GetCmdQuerySevDef(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "definition [chain-id] [name]",
		Short: "create new service definition",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			name := viper.GetString(FlagServiceName)
			chainId := viper.GetString(client.FlagChainID)
			res, err := cliCtx.QueryStore(iservice.GetServiceDefinitionKey(chainId, name), storeName)
			if err != nil {
				return err
			}

			var msgSvcDef iservice.MsgSvcDef
			cdc.MustUnmarshalBinary(res, &msgSvcDef)

			output, err := wire.MarshalJSONIndent(cdc, msgSvcDef)
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
