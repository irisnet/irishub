package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/irishub/client/record/ipfs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tmlibs/cli"
)

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetCmdDownload(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download [hash]",
		Short: "download specified file with tx hash",
		RunE: func(cmd *cobra.Command, args []string) error {

			// trustNode := viper.GetBool(client.FlagTrustNode)

			// hashHexStr := viper.GetString(FlagTxHash)
			downloadFileName := viper.GetString(FlagFileName)
			home := viper.GetString(cli.HomeFlag)

			//cliCtx := context.NewCLIContext().WithCodec(cdc)

			var err error
			var record RecordMetadata

			// record, err := queryRecordMetadata(cdc, cliCtx, hashHexStr, trustNode)
			// if err != nil {
			// 	return err
			// }

			// if len(record.DataHash) == 0 {
			// 	fmt.Printf("Request file was not found on the blockchain.\n")
			// 	return nil
			// }

			// WIP
			record.DataHash = "QmTp2hEo8eXRp6wg7jXv1BLCMh5a4F3B7buAUZNZUu772j"
			filePath := filepath.Join(home, downloadFileName)

			//Begin to download file from ipfs

			exist, _ := FileExists(filePath)
			if exist == true {
				fmt.Printf("%v already exists, please try another file name.\n", filePath)
				os.Exit(1)
			}

			fmt.Printf("Downloading %v ...\n", filePath)

			client := ipfs.NewIpfsclient()

			err = client.Get(record.DataHash, filePath)
			if err != nil {
				return err
			}

			fmt.Printf("Complete.\n")

			cid, err := client.Add(strings.NewReader("Upload file!"), true, false)
			fmt.Printf("this is uploadfile hash :%v\n", cid)

			return nil
		},
	}

	cmd.Flags().String(FlagTxHash, "", "tx hash")
	cmd.Flags().String(FlagFileName, "", "download file name")

	return cmd
}
