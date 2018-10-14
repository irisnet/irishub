package cli

import (
	"fmt"
	"os"
	"path/filepath"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/modules/record"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tmlibs/cli"

	shell "github.com/ipfs/go-ipfs-api"
)

func GetCmdDownload(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download [hash]",
		Short: "download specified file with tx hash",
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			downloadFileName := viper.GetString(FlagFileName)
			home := viper.GetString(cli.HomeFlag)
			hashHexStr := viper.GetString(FlagTxHash)
			trustNode := viper.GetBool(client.FlagTrustNode)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			ipfsHash, err := GetDataHash(cdc, cliCtx, hashHexStr, trustNode)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryStore(record.KeyRecord(addr, ipfsHash), storeName)
			if len(res) == 0 || err != nil {
				return fmt.Errorf("Record hash [%s] is not existed", hashHexStr)
			}

			var submitFile record.MsgSubmitFile
			cdc.MustUnmarshalBinary(res, &submitFile)

			if len(submitFile.DataHash) == 0 {
				fmt.Errorf("Request file was not found on the blockchain.\n")
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
			err = sh.Get(submitFile.DataHash, filePath)
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
