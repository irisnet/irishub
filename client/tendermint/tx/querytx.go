package tx

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/tendermint"
	"github.com/irisnet/irishub/client/tendermint/rpc"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/auth"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// QueryTxCmd implements the default command for a tx query.
func QueryTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tx [hash]",
		Short:   "Matches this txhash over all committed blocks",
		Example: "iriscli tendermint tx <transaction hash>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// find the key to look up the account
			hashHexStr := args[0]

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			output, err := queryTx(cdc, cliCtx, hashHexStr)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}
	cmd.Flags().Bool(client.FlagIndentResponse, true, "Add indent to JSON response")
	cmd.Flags().StringP(client.FlagNode, "n", "tcp://localhost:26657", "Node to connect to")
	cmd.Flags().Bool(client.FlagTrustNode, false, "Trust connected full node (don't verify proofs for responses)")
	cmd.Flags().String(client.FlagChainID, "", "Chain ID of Tendermint node")
	return cmd
}

// QueryTxCmd implements the default command for a tx query.
func QueryTxWithCoinFlowCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tx-coin-flow [hash]",
		Short:   "Matches this txhash over all committed blocks, all coin flow record triggered by this tx will be included",
		Example: "iriscli tendermint tx-coin-flow <transaction hash>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// find the key to look up the account
			hashHexStr := args[0]

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			output, err := queryTxWithCoinFlow(cdc, cliCtx, hashHexStr)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}
	cmd.Flags().Bool(client.FlagIndentResponse, true, "Add indent to JSON response")
	cmd.Flags().StringP(client.FlagNode, "n", "tcp://localhost:26657", "Node to connect to")
	cmd.Flags().Bool(client.FlagTrustNode, false, "Trust connected full node (don't verify proofs for responses)")
	cmd.Flags().String(client.FlagChainID, "", "Chain ID of Tendermint node")
	return cmd
}

func queryTx(cdc *codec.Codec, cliCtx context.CLIContext, hashHexStr string) ([]byte, error) {
	hash, err := hex.DecodeString(hashHexStr)
	if err != nil {
		return nil, err
	}

	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	res, err := node.Tx(hash, !cliCtx.TrustNode)
	if err != nil {
		return nil, err
	}

	if !cliCtx.TrustNode {
		err := ValidateTxResult(cliCtx, res)
		if err != nil {
			return nil, err
		}
	}

	info, err := formatTxResult(cdc, res)
	if err != nil {
		return nil, err
	}

	if cliCtx.Indent {
		return cdc.MarshalJSONIndent(info, "", "  ")
	}
	return cdc.MarshalJSON(info)
}

// ValidateTxResult performs transaction verification
func ValidateTxResult(cliCtx context.CLIContext, res *ctypes.ResultTx) error {
	check, err := cliCtx.Verify(res.Height)
	if err != nil {
		return err
	}

	err = res.Proof.Validate(check.Header.DataHash)
	if err != nil {
		return err
	}
	return nil
}

type ResponseDeliverTx struct {
	Code                 uint32
	Data                 []byte
	Log                  string
	Info                 string
	GasWanted            int64
	GasUsed              int64
	Tags                 []tendermint.ReadableTag
	Codespace            string
	XXX_NoUnkeyedLiteral struct{}
	XXX_unrecognized     []byte
	XXX_sizecache        int32
}

func MakeResponseHumanReadable(dtx abci.ResponseDeliverTx) ResponseDeliverTx {
	return ResponseDeliverTx{
		Code:      dtx.Code,
		Data:      dtx.Data,
		Log:       dtx.Log,
		Info:      dtx.Info,
		GasWanted: dtx.GasWanted,
		GasUsed:   dtx.GasUsed,
		Codespace: dtx.Codespace,
		Tags:      tendermint.MakeTagsHumanReadable(dtx.Tags),
	}
}

func formatTxResult(cdc *codec.Codec, res *ctypes.ResultTx) (Info, error) {
	tx, err := parseTx(cdc, res.Tx)
	if err != nil {
		return Info{}, err
	}

	return Info{
		Hash:   res.Hash,
		Height: res.Height,
		Tx:     tx,
		Result: MakeResponseHumanReadable(res.TxResult),
	}, nil
}

// Info is used to prepare info to display
type Info struct {
	Hash   common.HexBytes   `json:"hash"`
	Height int64             `json:"height"`
	Tx     sdk.Tx            `json:"tx"`
	Result ResponseDeliverTx `json:"result"`
}

func parseTx(cdc *codec.Codec, txBytes []byte) (sdk.Tx, error) {
	var tx auth.StdTx

	err := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// transaction query REST handler
func QueryTxRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hashHexStr := vars["hash"]

		output, err := queryTx(cdc, cliCtx, hashHexStr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		utils.PostProcessResponse(w, cdc, output, cliCtx.Indent)
	}
}

// Info is used to prepare info to display
type InfoCoinFlow struct {
	Hash     common.HexBytes          `json:"hash"`
	Height   int64                    `json:"height"`
	Tx       sdk.Tx                   `json:"tx"`
	Result   ResponseDeliverTx        `json:"result"`
	CoinFlow []tendermint.ReadableTag `json:"coin_flow"`
}

func queryTxWithCoinFlow(cdc *codec.Codec, cliCtx context.CLIContext, hashHexStr string) ([]byte, error) {
	hash, err := hex.DecodeString(hashHexStr)
	if err != nil {
		return nil, err
	}

	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	res, err := node.Tx(hash, !cliCtx.TrustNode)
	if err != nil {
		return nil, err
	}

	if !cliCtx.TrustNode {
		err := ValidateTxResult(cliCtx, res)
		if err != nil {
			return nil, err
		}
	}

	coinFlow, err := rpc.GetTxCoinFlow(cliCtx, &res.Height, strings.ToLower(hashHexStr))
	if err != nil {
		return nil, err
	}

	tx, err := parseTx(cdc, res.Tx)
	if err != nil {
		return nil, err
	}

	info := InfoCoinFlow{
		Hash:     res.Hash,
		Height:   res.Height,
		Tx:       tx,
		Result:   MakeResponseHumanReadable(res.TxResult),
		CoinFlow: coinFlow,
	}

	if cliCtx.Indent {
		return cdc.MarshalJSONIndent(info, "", "  ")
	}
	return cdc.MarshalJSON(info)
}

func QueryTxWithCoinFlowRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hashHexStr := vars["hash"]

		output, err := queryTxWithCoinFlow(cdc, cliCtx, hashHexStr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		utils.PostProcessResponse(w, cdc, output, cliCtx.Indent)
	}
}
