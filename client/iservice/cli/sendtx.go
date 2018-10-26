package cli

import (
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/spf13/cobra"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/irisnet/irishub/modules/iservice"
	"github.com/spf13/viper"
	"github.com/irisnet/irishub/client"
	cmn "github.com/tendermint/tendermint/libs/common"
	"fmt"
	"strings"
)

func GetCmdScvDef(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "define",
		Short: "create new service definition",
		Example: "iriscli iservice define --chain-id=<chain-id> --from=<key name> --fee=0.004iris " +
			"--name=<service name> --service-description=<service description> --author-description=<author description> " +
			"--tags=\"tag1 tag2\" --idl-content=<interface description content> --broadcast=Broadcast",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			txCtx := context.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			name := viper.GetString(FlagServiceName)
			description := viper.GetString(FlagServiceDescription)
			authorDescription := viper.GetString(FlagAuthorDescription)
			tags := viper.GetStringSlice(FlagTags)
			content := viper.GetString(FlagIdlContent)
			if len(content) > 0 {
				content = strings.Replace(content, `\n`, "\n", -1)
			}
			filePath := viper.GetString(FlagFile)
			if len(filePath) > 0 {
				contentBytes, err := cmn.ReadFile(filePath)
				if err != nil {
					return err
				}
				content = string(contentBytes)
			}
			fmt.Printf("idl condent: \n%s\n", content)
			broadcastStr := viper.GetString(FlagBroadcast)
			chainId := viper.GetString(client.FlagChainID)

			fromAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			broadcast, err := iservice.BroadcastEnumFromString(broadcastStr)
			if err != nil {
				return err
			}

			msg := iservice.NewMsgSvcDef(name, chainId, description, tags, fromAddr, authorDescription, content, broadcast)
			cliCtx.PrintResponse = true
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsServiceName)
	cmd.Flags().AddFlagSet(FsServiceDescription)
	cmd.Flags().AddFlagSet(FsTags)
	cmd.Flags().AddFlagSet(FsAuthorDescription)
	cmd.Flags().AddFlagSet(FsIdlContent)
	cmd.Flags().AddFlagSet(FsBroadcast)
	cmd.Flags().AddFlagSet(FsFile)

	return cmd
}
