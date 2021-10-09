package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/irisnet/irismod/modules/random/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding random type
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.RandomKey):
			var randomA, randomB types.Random
			cdc.MustUnmarshal(kvA.Value, &randomA)
			cdc.MustUnmarshal(kvB.Value, &randomB)
			return fmt.Sprintf("randA: %v\nrandB: %v", randomA, randomB)
		case bytes.Equal(kvA.Key[:1], types.RandomRequestQueueKey):
			var requestA, requestB types.Request
			cdc.MustUnmarshal(kvA.Value, &requestA)
			cdc.MustUnmarshal(kvB.Value, &requestB)
			return fmt.Sprintf("requestA: %v\nrequestB: %v", requestA, requestB)
		default:
			panic(fmt.Sprintf("invalid random key prefix %X", kvA.Key[:1]))
		}
	}
}
