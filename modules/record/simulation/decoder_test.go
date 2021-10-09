package simulation_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/irisnet/irismod/modules/record/simulation"
	"github.com/irisnet/irismod/modules/record/types"
	"github.com/irisnet/irismod/simapp"
)

var (
	creatorPk1   = secp256k1.GenPrivKey().PubKey()
	creatorAddr1 = sdk.AccAddress(creatorPk1.Address())
)

func TestDecodeStore(t *testing.T) {
	cdc := simapp.MakeTestEncodingConfig().Marshaler
	dec := simulation.NewDecodeStore(cdc)

	txHash := make([]byte, 32)
	_, _ = rand.Read(txHash)
	record := types.NewRecord(txHash, nil, creatorAddr1)

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: types.GetRecordKey(txHash), Value: cdc.MustMarshal(&record)},
			{Key: []byte{0x99}, Value: []byte{0x99}},
		},
	}
	tests := []struct {
		name        string
		expectedLog string
	}{
		{"Record", fmt.Sprintf("%v\n%v", record, record)},
		{"other", ""},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			switch i {
			case len(tests) - 1:
				require.Panics(t, func() { dec(kvPairs.Pairs[i], kvPairs.Pairs[i]) }, tt.name)
			default:
				require.Equal(t, tt.expectedLog, dec(kvPairs.Pairs[i], kvPairs.Pairs[i]), tt.name)
			}
		})
	}
}
