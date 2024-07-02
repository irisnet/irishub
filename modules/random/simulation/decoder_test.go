package simulation_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"
	"github.com/stretchr/testify/require"

	"mods.irisnet.org/modules/random/simulation"
	"mods.irisnet.org/modules/random/types"
	"mods.irisnet.org/simapp"
)

func TestDecodeStore(t *testing.T) {
	cdc := simapp.MakeTestEncodingConfig().Codec
	dec := simulation.NewDecodeStore(cdc)

	request := types.NewRequest(50, sdk.AccAddress("consumer").String(), hex.EncodeToString([]byte("txHash")), false, nil, "")
	reqID := types.GenerateRequestID(request)
	random := types.NewRandom(hex.EncodeToString([]byte("requestTxHash")), 100, big.NewRat(10, 1000).FloatString(types.RandPrec))

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: types.KeyRandom(reqID), Value: cdc.MustMarshal(&random)},
			{Key: types.KeyRandomRequestQueue(100, reqID), Value: cdc.MustMarshal(&request)},
			{Key: []byte{0x30}, Value: []byte{0x50}},
		},
	}

	tests := []struct {
		pass        bool
		name        string
		expectedLog string
	}{
		{true, "randoms", fmt.Sprintf("randA: %v\nrandB: %v", random, random)},
		{true, "pending requests", fmt.Sprintf("requestA: %v\nrequestB: %v", request, request)},
		{false, "other", ""},
	}

	for i, tt := range tests {
		i, tt := i, tt

		t.Run(tt.name, func(t *testing.T) {
			if tt.pass {
				require.Equal(t, tt.expectedLog, dec(kvPairs.Pairs[i], kvPairs.Pairs[i]), tt.name)
			} else {
				require.Panics(t, func() { dec(kvPairs.Pairs[i], kvPairs.Pairs[i]) }, tt.name)
			}
		})
	}
}
