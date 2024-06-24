package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"mods.irisnet.org/record/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding slashing type
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.RecordKey):
			var recordA, recordB types.Record
			cdc.MustUnmarshal(kvA.Value, &recordA)
			cdc.MustUnmarshal(kvB.Value, &recordB)
			return fmt.Sprintf("%v\n%v", recordA, recordB)
		default:
			panic(fmt.Sprintf("invalid record key prefix %X", kvA.Key[:1]))
		}
	}
}
