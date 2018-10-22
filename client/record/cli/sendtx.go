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
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/modules/record"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetCmdSubmitFile implements submitting upload file transaction command.
func GetCmdSubmitFile(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit",
		Short: "Submit the specified file/data",
		RunE: func(cmd *cobra.Command, args []string) error {
			description := viper.GetString(flagDescription)
			onchainData := viper.GetString(flagOnchainData)

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
				if dataSize >= record.UploadLimitOfOnchain {
					fmt.Printf("Onchain data is too large, upload limit is %d bytes.\n", record.UploadLimitOfOnchain)
					return err
				}
				sum := sha256.Sum256([]byte(onchainData))
				recordHash = hex.EncodeToString(sum[:])
			} else {
				fmt.Println("--onchain-data is empty and pleae double check this option")
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
			msg := record.NewMsgSubmitRecord(
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

	return cmd
}
