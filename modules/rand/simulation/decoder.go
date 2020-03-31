package simulation

import (
	"bytes"
	"fmt"

	tmkv "github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/irisnet/irishub/modules/rand/internal/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding rand type
func DecodeStore(cdc *codec.Codec, kvA, kvB tmkv.Pair) string {
	switch {
	case bytes.Equal(kvA.Key[:6], types.PrefixRand):
		var randA, randB types.Rand
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &randA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &randB)
		return fmt.Sprintf("randA: %v\nrandB: %v", randA, randB)
	case bytes.Equal(kvA.Key[:17], types.PrefixRandRequestQueue):
		var requestA, requestB types.Request
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &requestA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &requestB)
		return fmt.Sprintf("requestA: %v\nrequestB: %v", requestA, requestB)
	default:
		panic(fmt.Sprintf("invalid rand key prefix %X", kvA.Key[:1]))
	}
}
