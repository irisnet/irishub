package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/irisnet/irishub/modules/random/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding random type
func NewDecodeStore(cdc codec.Marshaler) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:6], types.PrefixRandom):
			var randomA, randomB types.Random
			cdc.MustUnmarshalBinaryBare(kvA.Value, &randomA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &randomB)
			return fmt.Sprintf("randA: %v\nrandB: %v", randomA, randomB)
		case bytes.Equal(kvA.Key[:17], types.PrefixRandomRequestQueue):
			var requestA, requestB types.Request
			cdc.MustUnmarshalBinaryBare(kvA.Value, &requestA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &requestB)
			return fmt.Sprintf("requestA: %v\nrequestB: %v", requestA, requestB)
		default:
			panic(fmt.Sprintf("invalid random key prefix %X", kvA.Key[:1]))
		}
	}
}
