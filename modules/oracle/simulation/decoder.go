package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/irisnet/irismod/modules/oracle/types"
)

// NewDecodeStore unmarshals the KVPair's Value to the corresponding slashing type
func NewDecodeStore(cdc codec.Marshaler) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.PrefixFeedKey):
			var feedA, feedB types.Feed
			cdc.MustUnmarshalBinaryBare(kvA.Value, &feedA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &feedB)
			return fmt.Sprintf("%v\n%v", feedA, feedB)
		case bytes.Equal(kvA.Key[:1], types.PrefixFeedValueKey):
			var feedValueA, feedValueB types.FeedValue
			cdc.MustUnmarshalBinaryBare(kvA.Value, &feedValueA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &feedValueB)
			return fmt.Sprintf("%v\n%v", feedValueA, feedValueB)
		default:
			panic(fmt.Sprintf("invalid record key prefix %X", kvA.Key[:1]))
		}
	}
}
