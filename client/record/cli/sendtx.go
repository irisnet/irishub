package cli

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/record"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetCmdSubmitFile implements submitting upload file transaction command.
func GetCmdSubmitRecord(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "submit",
		Short:   "Submit a new record",
		Example: "iriscli record submit --chain-id=<chain-id> --description=<record description> --onchain-data=<record data> --from=<key name> --fee=0.004iris",
		RunE: func(cmd *cobra.Command, args []string) error {
			description := viper.GetString(flagDescription)
			onchainData := viper.GetString(flagOnchainData)

			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).WithCliCtx(cliCtx)

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			fromAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			sum := sha256.Sum256([]byte(onchainData))
			recordHash := hex.EncodeToString(sum[:])

			recordID := record.KeyRecord(recordHash)
			res, err := cliCtx.QueryStore([]byte(recordID), storeName)
			if err != nil {
				return err
			}
			if len(res) != 0 {
				// Corresponding record id already exists, so there is no need to upload file/data
				fmt.Printf("Warning: Record ID %v already exists.\n", string(recordID))
				return nil
			}

			submitTime := time.Now().Unix()
			dataSize := int64(binary.Size([]byte(onchainData)))
			msg := record.NewMsgSubmitRecord(
				description,
				submitTime,
				fromAddr,
				recordHash,
				dataSize,
				onchainData,
			)

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	// common flag
	cmd.Flags().String(flagDescription, "description", "uploaded file description")

	// onchain flag
	cmd.Flags().String(flagOnchainData, "", "on chain data source")

	return cmd
}
