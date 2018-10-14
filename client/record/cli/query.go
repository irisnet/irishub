package cli

import (
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/irisnet/irishub/client"
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

func GetCmdQureyHash(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query [hash]",
		Short: "query specified file with tx hash",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)
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

func GetDataHash(cdc *wire.Codec, cliCtx context.CLIContext, hashHexStr string, trustNode bool) (string, error) {
	hash, err := hex.DecodeString(hashHexStr)
	if err != nil {
		return "", err
	}

	node, err := cliCtx.GetNode()
	if err != nil {
		return "", err
	}

	res, err := node.Tx(hash, !trustNode)
	if err != nil {
		return "", err
	}

	var tx auth.StdTx
	err = cdc.UnmarshalBinary(res.Tx, &tx)
	if err != nil {
		return "", err
	}

	msgs := tx.GetMsgs()
	if len(msgs) != 1 {
		fmt.Errorf("Record tx format error: there are more than one msg in the tx!\n")
		return "", err
	}

	var ok bool
	var m record.MsgSubmitFile
	if m, ok = msgs[0].(record.MsgSubmitFile); !ok {
		fmt.Errorf("MsgSubmitFile type assertion failed!\n")
		return "", err
	}

	return m.DataHash, nil
}
