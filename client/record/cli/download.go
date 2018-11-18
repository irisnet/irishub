package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/modules/record"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tmlibs/cli"

	shell "github.com/ipfs/go-ipfs-api"
)

func GetCmdDownload(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "download [record ID]",
		Short:   "download related data with unique record ID to specified file",
		Example: "iriscli record download --chain-id=<chain-id> --file-name=<file name> --record-id=<record-id>",
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			pinedNode := viper.GetString(flagPinedNode)
			downloadFileName := viper.GetString(flagFileName)
			home := viper.GetString(cli.HomeFlag)
			recordID := viper.GetString(flagRecordID)

			res, err := cliCtx.QueryStore([]byte(recordID), storeName)
			if len(res) == 0 || err != nil {
				return fmt.Errorf("Record id [%s] doesn't exist", recordID)
			}

			var submitFile record.MsgSubmitRecord
			cdc.MustUnmarshalBinaryLengthPrefixed(res, &submitFile)

			filePath := filepath.Join(home, downloadFileName)
			if _, err := os.Stat(filePath); !os.IsNotExist(err) {
				fmt.Printf("Warning: %v already exists, please try another file name.\n", filePath)
				return err
			}

			if len(submitFile.RecordID) == 0 {
				fmt.Errorf("Request file was not found on the blockchain.\n")
				return nil
			}
			if len(submitFile.Data) != 0 {
				//Begin to download file from blockchain directly
				fmt.Printf("[ONCHAIN] Downloading %v from blockchain directly...\n", filePath)
				fh, err := os.Create(filePath)
				if err != nil {
					return err
				}

				defer func() {
					if err := fh.Close(); err != nil {
						panic(err)
					}
				}()

				if _, err := fh.Write([]byte(submitFile.Data)); err != nil {
					return err
				}
				fmt.Println("[ONCHAIN] Download file from blockchain complete.")
			} else {
				//Begin to download file from ipfs
				fmt.Printf("[IPFS] Downloading %v from ipfs...\n", filePath)
				sh := shell.NewShell(pinedNode)
				err = sh.Get(submitFile.DataHash, filePath)
				if err != nil {
					return err
				}
				fmt.Println("[IPFS] Download file from ipfs complete.")
			}
			return nil
		},
	}

	cmd.Flags().String(flagRecordID, "", "record ID")
	cmd.Flags().String(flagFileName, "", "download file name")

	return cmd
}
