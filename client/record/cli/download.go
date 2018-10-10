package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/context"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tmlibs/cli"

	shell "github.com/ipfs/go-ipfs-api"
)

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

			filePath := filepath.Join(home, downloadFileName)
			sh := shell.NewShell("localhost:5001")

			//Begin to download file from ipfs
			if _, err := os.Stat("/path/to/whatever"); !os.IsNotExist(err) {
				fmt.Printf("%v already exists, please try another file name.\n", filePath)
				return err
			}

			fmt.Printf("Downloading %v ...\n", filePath)
			err = sh.Get(record.DataHash, filePath)
			if err != nil {
				return err
			}
			fmt.Println("Download file complete.")

			return nil
		},
	}

	cmd.Flags().String(FlagTxHash, "", "tx hash")
	cmd.Flags().String(FlagFileName, "", "download file name")

	return cmd
}
