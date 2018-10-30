package record

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestAddRecord(t *testing.T) {
	mapp, keeper, _, addrs, _, _ := getMockApp(t, 2)
	mapp.BeginBlock(abci.RequestBeginBlock{})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	data := "record data"
	sum := sha256.Sum256([]byte(data))
	recordHash := hex.EncodeToString(sum[:])
	dataSize := int64(binary.Size([]byte(data)))

	record1 := NewMsgSubmitRecord(
		"record description",
		time.Now().Unix(),
		addrs[0],
		recordHash,
		dataSize,
		data,
	)
	keeper.AddRecord(ctx, record1)

	err, record2 := getRecord(ctx, keeper, recordHash)
	require.Nil(t, err)
	require.True(t, recordEqual(record1, record2))

}
