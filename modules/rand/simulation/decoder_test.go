package simulation

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	cmn "github.com/tendermint/tendermint/libs/common"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/rand/internal/types"
)

func makeTestCodec() (cdc *codec.Codec) {
	cdc = codec.New()
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	types.RegisterCodec(cdc)
	return
}

func TestDecodeStore(t *testing.T) {
	cdc := makeTestCodec()

	request := types.NewRequest(50, sdk.AccAddress([]byte("consumer")), []byte("txHash"))
	reqID := types.GenerateRequestID(request)
	rand := types.NewRand([]byte("requestTxHash"), 100, big.NewRat(10, 1000))

	kvPairs := cmn.KVPairs{
		cmn.KVPair{Key: types.KeyRand(reqID), Value: cdc.MustMarshalBinaryLengthPrefixed(rand)},
		cmn.KVPair{Key: types.KeyRandRequestQueue(100, reqID), Value: cdc.MustMarshalBinaryLengthPrefixed(request)},
		cmn.KVPair{Key: []byte{0x30}, Value: []byte{0x50}},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"rands", fmt.Sprintf("randA: %v\nrandB: %v", rand, rand)},
		{"pending requests", fmt.Sprintf("requestA: %v\nrequestB: %v", request, request)},
		{"other", ""},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			switch i {
			case len(tests) - 1:
				require.Panics(t, func() { DecodeStore(cdc, kvPairs[i], kvPairs[i]) }, tt.name)
			default:
				require.Equal(t, tt.expectedLog, DecodeStore(cdc, kvPairs[i], kvPairs[i]), tt.name)
			}
		})
	}
}
