package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/modules/record"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type RecordMetadata struct {
	OwnerAddress string
	SubmitTime   string
	DataHash     string
	DataSize     string
	PinedNode    string
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

	cmd.Flags().StringP(client.FlagNode, "n", "tcp://localhost:26657", "Node to connect to")
	cmd.Flags().Bool(client.FlagTrustNode, false, "Trust connected full node (don't verify proofs for responses)")
	cmd.Flags().String(client.FlagChainID, "", "Chain ID of Tendermint node")
	cmd.Flags().String(FlagTxHash, "", "tx hash")
	cmd.Flags().String(FlagTargetPath, "", "target path")
	cmd.Flags().String(FlagFileName, "", "file name")

	return cmd
}

func queryRecordMetadata(cdc *wire.Codec, cliCtx context.CLIContext, hashHexStr string, trustNode bool) (RecordMetadata, error) {

	tx, err := queryTx(cdc, cliCtx, hashHexStr, trustNode)
	if err != nil {
		return RecordMetadata{}, err
	}

	msgs := tx.GetMsgs()

	if len(msgs) != 1 {
		return RecordMetadata{}, nil
	}

	// WIP
	var ok bool
	var m record.MsgSubmitFile
	if m, ok = msgs[0].(record.MsgSubmitFile); ok {
		return RecordMetadata{}, nil
	}

	return GetMetadata(m)
}

func GetMetadata(msg record.MsgSubmitFile) (RecordMetadata, error) {
	// Get record msg from record type msg (TO DO)
	var metadata RecordMetadata
	metadata.OwnerAddress = "address from record type msg"
	metadata.DataHash = "data hash from record type msg"
	metadata.DataSize = "data size from record type msg"
	metadata.PinedNode = "pined node from record type msg"
	metadata.SubmitTime = "submit time  from record type msg"

	return metadata, nil
}
