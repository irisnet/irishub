package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/modules/record"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	recordClient "github.com/irisnet/irishub/client/record"
)

type RecordMetadata struct {
	OwnerAddress sdk.AccAddress
	SubmitTime   int64
	DataHash     string
	DataSize     int64
	//PinedNode    string
}

func GetCmdQureyRecord(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query",
		Short:   "query specified record",
		Example: "iriscli record query --chain-id=<chain-id> --record-id=<record-id>",
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			recordID := viper.GetString(flagRecordID)

			res, err := cliCtx.QueryStore([]byte(recordID), storeName)
			if len(res) == 0 || err != nil {
				return fmt.Errorf("Record ID [%s] doesn't exist", recordID)
			}

			var submitRecord record.MsgSubmitRecord
			cdc.MustUnmarshalBinary(res, &submitRecord)

			recordResponse, err := recordClient.ConvertRecordToRecordOutput(cliCtx, submitRecord)
			if err != nil {
				return err
			}

			output, err := wire.MarshalJSONIndent(cdc, recordResponse)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil

		},
	}

	cmd.Flags().String(flagRecordID, "", "record ID for query")

	return cmd
}
