package simulation

import (
	"fmt"
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmkv "github.com/tendermint/tendermint/libs/kv"

	"github.com/irisnet/irishub/modules/mint/simulation"
	"github.com/irisnet/irishub/modules/random/types"
	"github.com/irisnet/irishub/simapp"
)

func TestDecodeStore(t *testing.T) {
	cdc, _ := simapp.MakeCodecs()
	dec := simulation.NewDecodeStore(cdc)

	request := types.NewRequest(50, sdk.AccAddress("consumer"), []byte("txHash"), false, nil, nil)
	reqID := types.GenerateRequestID(request)
	rand := types.NewRandom([]byte("requestTxHash"), 100, big.NewRat(10, 1000).FloatString(types.RandPrec))

	kvPairs := tmkv.Pairs{
		tmkv.Pair{Key: types.KeyRandom(reqID), Value: cdc.MustMarshalBinaryBare(&rand)},
		tmkv.Pair{Key: types.KeyRandomRequestQueue(100, reqID), Value: cdc.MustMarshalBinaryBare(request)},
		tmkv.Pair{Key: []byte{0x30}, Value: []byte{0x50}},
	}

	tests := []struct {
		pass        bool
		name        string
		expectedLog string
	}{
		{true, "rands", fmt.Sprintf("randA: %v\nrandB: %v", rand, rand)},
		{true, "pending requests", fmt.Sprintf("requestA: %v\nrequestB: %v", request, request)},
		{false, "other", ""},
	}

	for i, tt := range tests {
		i, tt := i, tt

		t.Run(tt.name, func(t *testing.T) {
			if tt.pass {
				require.Equal(t, tt.expectedLog, NewDecodeStore(cdc), tt.name)
			} else {
				require.Panics(t, func() { dec(kvPairs[i], kvPairs[i]) }, tt.name)
			}
		})
	}
}
