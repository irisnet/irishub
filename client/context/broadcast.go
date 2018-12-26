package context

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/pkg/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// TODO: This should get deleted eventually, and perhaps
// ctypes.ResultBroadcastTx be stripped of unused fields, and
// ctypes.ResultBroadcastTxCommit returned for tendermint RPC BroadcastTxSync.
//
// The motivation is that we want a unified type to return, and the better
// option is the one that can hold CheckTx/DeliverTx responses optionally.
func resultBroadcastTxToCommit(res *ctypes.ResultBroadcastTx) *ctypes.ResultBroadcastTxCommit {
	return &ctypes.ResultBroadcastTxCommit{
		Hash: res.Hash,
		// NOTE: other fields are unused for async.
	}
}

// BroadcastTx broadcasts a transactions either synchronously or asynchronously
// based on the context parameters. The result of the broadcast is parsed into
// an intermediate structure which is logged if the context has a logger
// defined.
func (cliCtx CLIContext) BroadcastTx(txBytes []byte) (*ctypes.ResultBroadcastTxCommit, error) {
	// commit
	if cliCtx.WaitForCommit {
		return cliCtx.broadcastTxCommit(txBytes)
	}
	// sync
	if !cliCtx.Async {
		return cliCtx.broadcastTxSync(txBytes)
	}
	// async
	res, err := cliCtx.broadcastTxAsync(txBytes)
	if err != nil {
		return nil, err
	}

	resCommit := resultBroadcastTxToCommit(res)
	return resCommit, err

}

// BroadcastTxAndAwaitCommit broadcasts transaction bytes to a Tendermint node
// and waits for a commit.
func (cliCtx CLIContext) BroadcastTxAndAwaitCommit(tx []byte) (*ctypes.ResultBroadcastTxCommit, error) {
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	res, err := node.BroadcastTxCommit(tx)
	if err != nil {
		return res, err
	}

	if !res.CheckTx.IsOK() {
		return res, errors.Errorf(res.CheckTx.Log)
	}

	if !res.DeliverTx.IsOK() {
		return res, errors.Errorf(res.DeliverTx.Log)
	}

	return res, err
}

// BroadcastTxSync broadcasts transaction bytes to a Tendermint node
// synchronously.
func (cliCtx CLIContext) BroadcastTxSync(tx []byte) (*ctypes.ResultBroadcastTxCommit, error) {
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	res, err := node.BroadcastTxSync(tx)

	result := &ctypes.ResultBroadcastTxCommit{
		Hash: res.Hash,
		CheckTx: abci.ResponseCheckTx{
			Code: res.Code,
			Data: res.Data,
			Log:  res.Log,
		},
	}
	if err != nil {
		return result, err
	}

	return result, err
}

// BroadcastTxAsync broadcasts transaction bytes to a Tendermint node
// asynchronously.
func (cliCtx CLIContext) BroadcastTxAsync(tx []byte) (*ctypes.ResultBroadcastTx, error) {
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	res, err := node.BroadcastTxAsync(tx)
	if err != nil {
		return res, err
	}

	return res, err
}

func (cliCtx CLIContext) broadcastTxSync(txBytes []byte) (*ctypes.ResultBroadcastTxCommit, error) {
	res, err := cliCtx.BroadcastTxSync(txBytes)
	if err != nil {
		return res, err
	}

	if cliCtx.Logger != nil {
		if cliCtx.JSON {
			// Since JSON is intended for automated scripts, always include response in
			// JSON mode.
			type toJSON struct {
				TxHash   string
				Response abci.ResponseCheckTx
			}

			if cliCtx.Logger != nil {
				resJSON := toJSON{res.Hash.String(), res.CheckTx}
				bz, err := cliCtx.Codec.MarshalJSON(resJSON)
				if err != nil {
					return res, err
				}

				cliCtx.Logger.Write(bz)
				io.WriteString(cliCtx.Logger, "\n")
			}

			return res, nil
		} else {
			resStr := fmt.Sprintf("sync tx sent (tx hash: %s)\n", res.Hash)
			if cliCtx.PrintResponse {
				jsonStr, _ := checkTxMarshalIndentJSON(res.CheckTx)
				resStr = fmt.Sprintf("Sync committed (tx hash: %s, response:\n %+v)\n", res.Hash, string(jsonStr))
			}

			io.WriteString(cliCtx.Logger, resStr)
		}
	}

	return res, nil
}

func (cliCtx CLIContext) broadcastTxAsync(txBytes []byte) (*ctypes.ResultBroadcastTx, error) {
	res, err := cliCtx.BroadcastTxAsync(txBytes)
	if err != nil {
		return res, err
	}

	if cliCtx.Logger != nil {
		if cliCtx.JSON {
			type toJSON struct {
				TxHash string
			}

			resJSON := toJSON{res.Hash.String()}
			bz, err := cliCtx.Codec.MarshalJSON(resJSON)
			if err != nil {
				return res, err
			}

			cliCtx.Logger.Write(bz)
			io.WriteString(cliCtx.Logger, "\n")
		} else {
			io.WriteString(cliCtx.Logger, fmt.Sprintf("async tx sent (tx hash: %s)\n", res.Hash))
		}
	}

	return res, nil
}

func (cliCtx CLIContext) broadcastTxCommit(txBytes []byte) (*ctypes.ResultBroadcastTxCommit, error) {
	res, err := cliCtx.BroadcastTxAndAwaitCommit(txBytes)
	if err != nil {
		return res, err
	}

	if cliCtx.JSON {
		// Since JSON is intended for automated scripts, always include response in
		// JSON mode.
		type toJSON struct {
			Height   int64
			TxHash   string
			Response abci.ResponseDeliverTx
		}

		if cliCtx.Logger != nil {
			resJSON := toJSON{res.Height, res.Hash.String(), res.DeliverTx}
			bz, err := cliCtx.Codec.MarshalJSON(resJSON)
			if err != nil {
				return res, err
			}

			cliCtx.Logger.Write(bz)
			io.WriteString(cliCtx.Logger, "\n")
		}

		return res, nil
	}

	if cliCtx.Logger != nil {
		resStr := fmt.Sprintf("Committed at block %d (tx hash: %s)\n", res.Height, res.Hash.String())

		if cliCtx.PrintResponse {
			jsonStr, _ := deliverTxMarshalIndentJSON(res.DeliverTx)
			resStr = fmt.Sprintf("Committed at block %d (tx hash: %s, response:\n %+v)\n",
				res.Height, res.Hash.String(), string(jsonStr),
			)
		}

		io.WriteString(cliCtx.Logger, resStr)
	}

	return res, nil
}

type ReadableTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func deliverTxMarshalIndentJSON(dtx abci.ResponseDeliverTx) ([]byte, error) {
	tags := make([]ReadableTag, len(dtx.Tags))
	for i, kv := range dtx.Tags {
		tags[i] = ReadableTag{
			Key:   string(kv.Key),
			Value: string(kv.Value),
		}
	}
	return json.MarshalIndent(&struct {
		Code      uint32        `json:"code"`
		Data      []byte        `json:"data"`
		Log       string        `json:"log"`
		Info      string        `json:"info"`
		GasWanted int64         `json:"gas_wanted"`
		GasUsed   int64         `json:"gas_used"`
		Codespace string        `json:"codespace"`
		Tags      []ReadableTag `json:"tags,omitempty"`
	}{
		Code:      dtx.Code,
		Data:      dtx.Data,
		Log:       dtx.Log,
		Info:      dtx.Info,
		GasWanted: dtx.GasWanted,
		GasUsed:   dtx.GasUsed,
		Codespace: dtx.Codespace,
		Tags:      tags,
	}, " ", "  ")
}

func checkTxMarshalIndentJSON(ctx abci.ResponseCheckTx) ([]byte, error) {
	tags := make([]ReadableTag, len(ctx.Tags))
	for i, kv := range ctx.Tags {
		tags[i] = ReadableTag{
			Key:   string(kv.Key),
			Value: string(kv.Value),
		}
	}

	return json.MarshalIndent(&struct {
		Code      uint32        `json:"code"`
		Data      []byte        `json:"data"`
		Log       string        `json:"log"`
		Info      string        `json:"info"`
		GasWanted int64         `json:"gas_wanted"`
		GasUsed   int64         `json:"gas_used"`
		Tags      []ReadableTag `json:"tags,omitempty"`
		Codespace string        `json:"codespace"`
	}{
		Code:      ctx.Code,
		Data:      ctx.Data,
		Log:       ctx.Log,
		Info:      ctx.Info,
		GasWanted: ctx.GasWanted,
		GasUsed:   ctx.GasUsed,
		Codespace: ctx.Codespace,
		Tags:      tags,
	}, " ", "  ")
}
