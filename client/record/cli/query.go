package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/modules/record"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	recordClient "github.com/irisnet/irishub/client/record"
	cmn "github.com/tendermint/tendermint/libs/common"
)

type RecordMetadata struct {
	OwnerAddress sdk.AccAddress
	SubmitTime   int64
	DataHash     string
	DataSize     int64
	//PinedNode    string
}

func GetCmdQureyHash(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query [hash]",
		Short: "query specified file with tx hash",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			hashHexStr := viper.GetString(FlagTxHash)

			var tmpkey = cmn.HexBytes{}
			res, err := cliCtx.QueryStore(tmpkey /*record.KeyProposal(hashHexStr)*/, storeName)
			if len(res) == 0 || err != nil {
				return fmt.Errorf("Record hash [%s] is not existed", hashHexStr)
			}

			var submitFile record.MsgSubmitFile
			cdc.MustUnmarshalBinary(res, &submitFile)

			recordResponse, err := recordClient.ConvertRecordToRecordOutput(cliCtx, submitFile)
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
