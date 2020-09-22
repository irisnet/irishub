package simulation

// DONTCOVER

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"
	gogotypes "github.com/gogo/protobuf/types"

	"github.com/irisnet/irismod/modules/token/types"
)

// NewDecodeStore unmarshals the KVPair's Value to the corresponding token type
func NewDecodeStore(cdc codec.Marshaler) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.PrefixTokenForSymbol):
			var tokenA, tokenB types.Token
			cdc.MustUnmarshalBinaryBare(kvA.Value, &tokenA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &tokenB)
			return fmt.Sprintf("%v\n%v", tokenA, tokenB)
		case bytes.Equal(kvA.Key[:1], types.PrefixTokens):
			var symbolA, symbolB gogotypes.Value
			cdc.MustUnmarshalBinaryBare(kvA.Value, &symbolA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &symbolB)
			return fmt.Sprintf("%v\n%v", symbolA, symbolB)
		case bytes.Equal(kvA.Key[:1], types.PrefixTokenForMinUint):
			var symbolA, symbolB gogotypes.StringValue
			cdc.MustUnmarshalBinaryBare(kvA.Value, &symbolA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &symbolB)
			return fmt.Sprintf("%v\n%v", symbolA, symbolB)
		default:
			panic(fmt.Sprintf("invalid %s key prefix %X", types.ModuleName, kvA.Key[:1]))
		}
	}
}
