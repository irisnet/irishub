package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/modules/record"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tmlibs/cli"

	shell "github.com/ipfs/go-ipfs-api"
)

func GetCmdDownload(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download [record ID]",
		Short: "download specified file with record ID",
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			strPinedNode := viper.GetString(flagPinedNode)
			downloadFileName := viper.GetString(FlagFileName)
			home := viper.GetString(cli.HomeFlag)
			recordID := viper.GetString(FlagRecordID)

			res, err := cliCtx.QueryStore([]byte(recordID), storeName)
			if len(res) == 0 || err != nil {
				return fmt.Errorf("Record id [%s] is not existed", recordID)
			}

			var submitFile record.MsgSubmitFile
			cdc.MustUnmarshalBinary(res, &submitFile)

			if len(submitFile.DataHash) == 0 {
				fmt.Errorf("Request file was not found on the blockchain.\n")
				return nil
			}

			filePath := filepath.Join(home, downloadFileName)
			sh := shell.NewShell(strPinedNode)

			//Begin to download file from ipfs
			if _, err := os.Stat("/path/to/whatever"); !os.IsNotExist(err) {
				fmt.Printf("%v already exists, please try another file name.\n", filePath)
				return err
			}

			fmt.Printf("Downloading %v ...\n", filePath)
			err = sh.Get(submitFile.DataHash, filePath)
			if err != nil {
				return err
			}
			fmt.Println("Download file complete.")

			return nil
		},
	}

	cmd.Flags().String(flagPinedNode, "localhost:5001", "node to download file,ip:port")
	cmd.Flags().String(FlagRecordID, "", "record ID")
	cmd.Flags().String(FlagFileName, "", "download file name")

	return cmd
}
