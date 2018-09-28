package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/record/ipfs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tmlibs/cli"
)

func PathExists(path string) (bool, error) {
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

			trustNode := viper.GetBool(client.FlagTrustNode)

			hashHexStr := viper.GetString(FlagTxHash)
			downloadFileName := viper.GetString(FlagFileName)
			home := viper.GetString(cli.HomeFlag)

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			record, err := queryRecordMetadata(cdc, cliCtx, hashHexStr, trustNode)
			if err != nil {
				return err
			}

			if len(record.DataHash) == 0 {
				fmt.Printf("Request file was not found on the blockchain.\n")
				return nil
			}

			// WIP
			record.DataHash = "QmTp2hEo8eXRp6wg7jXv1BLCMh5a4F3B7buAUZNZUu772j"
			filePath := filepath.Join(home, downloadFileName)

			//Begin to download file from ipfs
			sh := ipfs.NewShell("localhost:5001")

			exist, err := PathExists(filePath)
			if exist == true {
				fmt.Printf("%v already exists, please try another file name.\n", filePath)
				os.Exit(1)
			}

			fmt.Printf("Downloading %v ...\n", filePath)

			err = sh.Get(record.DataHash, filePath)
			if err != nil {
				os.Exit(1)
			}

			fmt.Printf("Complete.\n")

			return nil
		},
	}

	cmd.Flags().String(FlagTxHash, "", "tx hash")
	cmd.Flags().String(FlagFileName, "", "file name")
	cmd.Flags().String(FlagTargetPath, "", "target path")

	return cmd
}
