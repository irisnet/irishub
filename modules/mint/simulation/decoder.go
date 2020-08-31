package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/irisnet/irishub/modules/mint/types"
)

// NewDecodeStore returns a function closure that unmarshals the KVPair's values
// to the corresponding types.
func NewDecodeStore(cdc codec.Marshaler) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key, types.MinterKey):
			var minterA, minterB types.Minter
			cdc.MustUnmarshalBinaryBare(kvA.Value, &minterA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &minterB)
			return fmt.Sprintf("%v\n%v", minterA, minterB)
		default:
			panic(fmt.Sprintf("invalid mint key %X", kvA.Key))
		}
	}
}
