package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/irisnet/irismod/modules/service/types"
)

// NewDecodeStore unmarshals the KVPair's Value to the corresponding service type
func NewDecodeStore(cdc codec.Marshaler) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.ServiceDefinitionKey):
			var definition1, definition2 types.ServiceDefinition
			cdc.MustUnmarshalBinaryBare(kvA.Value, &definition1)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &definition2)
			return fmt.Sprintf("%v\n%v", definition1, definition2)

		default:
			panic(fmt.Sprintf("invalid service key prefix %X", kvA.Key[:1]))
		}
	}
}
