package simulation

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	tmkv "github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/random/types"
)

func makeTestCodec() (clientCtx client.Context) {
	cdc = codec.New()
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	types.RegisterCodec(cdc)
	return
}

func TestDecodeStore(t *testing.T) {
	cdc := makeTestCodec()

	request := types.NewRequest(50, sdk.AccAddress("consumer"), []byte("txHash"))
	reqID := types.GenerateRequestID(request)
	rand := types.NewRandom([]byte("requestTxHash"), 100, big.NewRat(10, 1000).FloatString(types.RandomPrec))

	kvPairs := tmkv.Pairs{
		tmkv.Pair{Key: types.KeyRandom(reqID), Value: cdc.MustMarshalBinaryBare(rand)},
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
				require.Equal(t, tt.expectedLog, DecodeStore(cdc, kvPairs[i], kvPairs[i]), tt.name)
			} else {
				require.Panics(t, func() { DecodeStore(cdc, kvPairs[i], kvPairs[i]) }, tt.name)
			}
		})
	}
}
