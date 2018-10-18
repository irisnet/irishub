package cli

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	ipfs "github.com/ipfs/go-ipfs-api"
	"github.com/irisnet/irishub/client/context"
	recordClient "github.com/irisnet/irishub/client/record"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/modules/record"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetCmdSubmitFile implements submitting upload file transaction command.
func GetCmdSubmitFile(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit",
		Short: "Submit the specified file",
		RunE: func(cmd *cobra.Command, args []string) error {
			description := viper.GetString(flagDescription)
			onchainData := viper.GetString(flagOnchainData)
			filePath := viper.GetString(flagFilePath)
			pinedNode := viper.GetString(flagPinedNode)

			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			txCtx := context.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			fromAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			var recordHash string
			var dataSize int64
			// --onchain-data has a high priority over --file-path
			if len(onchainData) != 0 {
				dataSize = int64(binary.Size([]byte(onchainData)))
				if dataSize >= recordClient.UploadLimitOfOnchain {
					fmt.Printf("File %s is too large, upload limit is %d bytes.\n", filePath, recordClient.UploadLimitOfOnchain)
					return err
				}
				sum := sha256.Sum256([]byte(onchainData))
				recordHash = hex.EncodeToString(sum[:recordClient.IpfsHashLength/2])
			} else if len(filePath) != 0 {
				var fileInfo os.FileInfo
				if fileInfo, err = os.Stat(filePath); os.IsNotExist(err) {
					fmt.Printf("File %v doesn't exists, please check correstponding path.\n", filePath)
					return err
				}

				dataSize = fileInfo.Size()
				if dataSize >= recordClient.UploadLimitOfIpfs {
					fmt.Printf("File %s is too large, upload limit is %d bytes.\n", filePath, recordClient.UploadLimitOfIpfs)
					return err
				}

				//upload to ipfs
				sh := ipfs.NewShell(pinedNode)
				f, err := os.Open(filePath)
				if err != nil {
					return err
				}
				ipfsHash, err := sh.Add(f)
				recordHash = ipfsHash
				if err != nil {
					return err
				}
			} else {
				fmt.Printf("--onchain-data and --file-path are both empty and pleae specify one of them")
				return err
			}

			recordID := record.KeyRecord(recordHash)
			res, err := cliCtx.QueryStore([]byte(recordID), storeName)
			if err != nil {
				return err
			}
			if len(res) != 0 {
				// Corresponding record id is already existed, so there is no need to upload file/data
				return fmt.Errorf("Record ID [%s] is already existed", recordID)
			}

			submitTime := time.Now().Unix()
			msg := record.NewMsgSubmitFile(
				description,
				submitTime,
				fromAddr,
				recordHash,
				dataSize,
				onchainData,
			)

			if cliCtx.GenerateOnly {
				return utils.PrintUnsignedStdTx(txCtx, cliCtx, []sdk.Msg{msg})
			}

			// Build and sign the transaction, then broadcast to Tendermint
			cliCtx.PrintResponse = true

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	// common flag
	cmd.Flags().String(flagDescription, "description", "uploaded file description")

	// onchain flag
	cmd.Flags().String(flagOnchainData, "", "on chain data source")

	// ipfs related flag
	cmd.Flags().String(flagFilePath, "", "full path of file (include filename)")
	cmd.Flags().String(flagPinedNode, "localhost:5001", "node to upload file,ip:port")

	return cmd
}
