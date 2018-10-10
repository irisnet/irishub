package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/modules/record"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type RecordMetadata struct {
	OwnerAddress sdk.AccAddress
	SubmitTime   int64
	DataHash     string
	DataSize     int64
	//PinedNode    string
}

func GetCmdQureyHash(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query [hash]",
		Short: "query specified file with tx hash",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			trustNode := viper.GetBool(client.FlagTrustNode)

			hashHexStr := viper.GetString(FlagTxHash)

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			record, err := queryRecordMetadata(cdc, cliCtx, hashHexStr, trustNode)
			if err != nil {
				return err
			}

			fmt.Println("Record metadata %v", record)
			return nil
		},
	}

	cmd.Flags().String(FlagTxHash, "", "tx hash for query")

	return cmd
}

func queryRecordMetadata(cdc *wire.Codec, cliCtx context.CLIContext, hashHexStr string, trustNode bool) (RecordMetadata, error) {

	tx, err := QueryTx(cdc, cliCtx, hashHexStr, trustNode)

	if err != nil {
		return RecordMetadata{}, err
	}

	msgs := tx.GetMsgs()

	for i := 0; i < len(msgs); i++ {
		if msgs[i].Type() == record.MsgType {
			var ok bool
			var m record.MsgSubmitFile
			if m, ok = msgs[i].(record.MsgSubmitFile); ok {

				var metadata RecordMetadata
				metadata.OwnerAddress = m.OwnerAddress
				metadata.DataHash = m.DataHash
				metadata.DataSize = m.DataSize
				metadata.SubmitTime = m.SubmitTime

				return metadata, nil
			}
			return RecordMetadata{}, nil
		}
	}

	return RecordMetadata{}, nil
}

func QueryTx(cdc *wire.Codec, cliCtx context.CLIContext, hashHexStr string, trustNode bool) (sdk.Tx, error) {

	return nil, nil
}
