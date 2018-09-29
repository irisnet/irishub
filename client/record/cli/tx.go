package cli

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// txInfo is used to prepare info to display
type txInfo struct {
	Hash   common.HexBytes        `json:"hash"`
	Height int64                  `json:"height"`
	Tx     sdk.Tx                 `json:"tx"`
	Result abci.ResponseDeliverTx `json:"result"`
}

func parseTx(cdc *wire.Codec, txBytes []byte) (sdk.Tx, error) {
	var tx auth.StdTx
	err := cdc.UnmarshalBinary(txBytes, &tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func formatTxResult(cdc *wire.Codec, res *ctypes.ResultTx) (txInfo, error) {
	// TODO: verify the proof if requested
	tx, err := parseTx(cdc, res.Tx)
	if err != nil {
		return txInfo{}, err
	}

	info := txInfo{
		Hash:   res.Hash,
		Height: res.Height,
		Tx:     tx,
		Result: res.TxResult,
	}
	return info, nil
}
