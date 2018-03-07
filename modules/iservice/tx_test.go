package iservice

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cosmos/cosmos-sdk"
	wire "github.com/tendermint/go-wire"
)

func TestAllAreTx(t *testing.T) {
	assert := assert.New(t)

	txDefineSvc := NewTxDefineService("", "")
	_, ok := txDefineSvc.Unwrap().(TxDefineService)
	assert.True(ok, "%#v", txDefineSvc)
}

func TestSerializeTx(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		tx sdk.Tx
	}{
		{NewTxDefineService("", "")},
	}

	for i, tc := range cases {
		var tx sdk.Tx
		bs := wire.BinaryBytes(tc.tx)
		err := wire.ReadBinaryBytes(bs, &tx)
		if assert.NoError(err, "%d", i) {
			assert.Equal(tc.tx, tx, "%d", i)
		}
	}
}
