package cli

import (
	"fmt"
	"os"

	"github.com/irisnet/irishub/client/context"
	upgcli "github.com/irisnet/irishub/client/upgrade"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/upgrade"
	"github.com/spf13/cobra"
	protocol "github.com/irisnet/irishub/app/protocol/keeper"
)

func GetInfoCmd(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "info",
		Short:   "query the information of upgrade module",
		Example: "iriscli upgrade info",
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))

			res_protocolVersion, _ := cliCtx.QueryStore(protocol.CurrentProtocolVersionKey, "protocol")
			var protocolVersion uint64
			cdc.MustUnmarshalBinaryLengthPrefixed(res_protocolVersion, &protocolVersion)

			res_upgradeConfig, _ := cliCtx.QueryStore(protocol.UpgradeConfigkey, "protocol")
			var upgradeConfig protocol.UpgradeConfig
			cdc.MustUnmarshalBinaryLengthPrefixed(res_upgradeConfig, &upgradeConfig)


			res_appVersion, _ := cliCtx.QueryStore(upgrade.GetAppVersionKey(protocolVersion), storeName)
			var appVersion upgrade.AppVersion
			cdc.MustUnmarshalBinaryLengthPrefixed(res_appVersion, &appVersion)

			upgradeInfoOutput := upgcli.ConvertUpgradeInfoToUpgradeOutput(appVersion,upgradeConfig)

			output, err := codec.MarshalJSONIndent(cdc, upgradeInfoOutput)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}
	return cmd
}

