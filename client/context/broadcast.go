package context

import (
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
	if cliCtx.Async {
		res, err := cliCtx.broadcastTxAsync(txBytes)
		if err != nil {
			return nil, err
		}

		resCommit := resultBroadcastTxToCommit(res)
		return resCommit, err
	}

	return cliCtx.broadcastTxCommit(txBytes)
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
func (cliCtx CLIContext) BroadcastTxSync(tx []byte) (*ctypes.ResultBroadcastTx, error) {
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	res, err := node.BroadcastTxSync(tx)
	if err != nil {
		return res, err
	}

	return res, err
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
			resStr = fmt.Sprintf("Committed at block %d (tx hash: %s, response: %+v)\n",
				res.Height, res.Hash.String(), res.DeliverTx,
			)
		}

		io.WriteString(cliCtx.Logger, resStr)
	}

	return res, nil
}
