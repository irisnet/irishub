package types

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto/tmhash"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

func TestIDGenerate(t *testing.T) {
	txHash := tmbytes.HexBytes(tmhash.Sum([]byte("tx_hash")))
	msgIndex := int64(math.MaxInt64)
	contextID := GenerateRequestContextID(txHash, msgIndex)
	txHash1, msgIndex1, _ := SplitRequestContextID(contextID)
	require.Equal(t, txHash, txHash1)
	require.Equal(t, msgIndex, msgIndex1)

	requestContextBatchCounter := uint64(math.MaxUint64)
	requestHeight := int64(math.MaxInt64)
	batchRequestIndex := int16(math.MaxInt16)
	requestID := GenerateRequestID(contextID, requestContextBatchCounter, requestHeight, batchRequestIndex)
	contextID1, requestContextBatchCounter1, requestHeight1, batchRequestIndex1, _ := SplitRequestID(requestID)
	require.Equal(t, contextID, contextID1)
	require.Equal(t, requestContextBatchCounter, requestContextBatchCounter1)
	require.Equal(t, requestHeight, requestHeight1)
	require.Equal(t, batchRequestIndex, batchRequestIndex1)
}
