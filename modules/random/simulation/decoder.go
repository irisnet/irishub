package simulation

import (
	"bytes"
	"fmt"

	tmkv "github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/irisnet/irishub/modules/random/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding rand type
func NewDecodeStore(cdc codec.Marshaler) func(kvA, kvB tmkv.Pair) string {
	return func(kvA, kvB tmkv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:6], types.PrefixRandom):
			var randA, randB types.Random
			cdc.MustMarshalBinaryBare(kvA.Value, &randA)
			cdc.MustMarshalBinaryBare(kvB.Value, &randB)
			return fmt.Sprintf("randA: %v\nrandB: %v", randA, randB)
		case bytes.Equal(kvA.Key[:17], types.PrefixRandomRequestQueue):
			var requestA, requestB types.Request
			cdc.MustMarshalBinaryBare(kvA.Value, &requestA)
			cdc.MustMarshalBinaryBare(kvB.Value, &requestB)
			return fmt.Sprintf("requestA: %v\nrequestB: %v", requestA, requestB)
		default:
			panic(fmt.Sprintf("invalid rand key prefix %X", kvA.Key[:1]))
		}
	}
}
