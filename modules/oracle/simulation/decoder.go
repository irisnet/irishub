package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/irisnet/irismod/modules/oracle/types"
)

// NewDecodeStore unmarshals the KVPair's Value to the corresponding slashing type
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.PrefixFeedKey):
			var feedA, feedB types.Feed
			cdc.MustUnmarshal(kvA.Value, &feedA)
			cdc.MustUnmarshal(kvB.Value, &feedB)
			return fmt.Sprintf("%v\n%v", feedA, feedB)
		case bytes.Equal(kvA.Key[:1], types.PrefixFeedValueKey):
			var feedValueA, feedValueB types.FeedValue
			cdc.MustUnmarshal(kvA.Value, &feedValueA)
			cdc.MustUnmarshal(kvB.Value, &feedValueB)
			return fmt.Sprintf("%v\n%v", feedValueA, feedValueB)
		default:
			panic(fmt.Sprintf("invalid record key prefix %X", kvA.Key[:1]))
		}
	}
}
