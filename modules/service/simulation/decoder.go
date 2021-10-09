package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/irisnet/irismod/modules/service/types"
)

// NewDecodeStore unmarshals the KVPair's Value to the corresponding service type
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.ServiceDefinitionKey):
			var definition1, definition2 types.ServiceDefinition
			cdc.MustUnmarshal(kvA.Value, &definition1)
			cdc.MustUnmarshal(kvB.Value, &definition2)
			return fmt.Sprintf("%v\n%v", definition1, definition2)
		case bytes.Equal(kvA.Key[:1], types.ServiceBindingKey):
			var binding1, binding2 types.ServiceBinding
			cdc.MustUnmarshal(kvA.Value, &binding1)
			cdc.MustUnmarshal(kvB.Value, &binding2)
			return fmt.Sprintf("%v\n%v", binding1, binding2)
		case bytes.Equal(kvA.Key[:1], types.RequestContextKey):
			var requestContext1, requestContext2 types.RequestContext
			cdc.MustUnmarshal(kvA.Value, &requestContext1)
			cdc.MustUnmarshal(kvB.Value, &requestContext2)
			return fmt.Sprintf("%v\n%v", requestContext1, requestContext2)
		case bytes.Equal(kvA.Key[:1], types.RequestKey):
			var request1, request2 types.Request
			cdc.MustUnmarshal(kvA.Value, &request1)
			cdc.MustUnmarshal(kvB.Value, &request2)
			return fmt.Sprintf("%v\n%v", request1, request2)
		case bytes.Equal(kvA.Key[:1], types.ResponseKey):
			var response1, response2 types.Response
			cdc.MustUnmarshal(kvA.Value, &response1)
			cdc.MustUnmarshal(kvB.Value, &response2)
			return fmt.Sprintf("%v\n%v", response1, response2)
		default:
			panic(fmt.Sprintf("invalid service key prefix %X", kvA.Key[:1]))
		}
	}
}
