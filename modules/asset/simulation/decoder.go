package simulation

// DONTCOVER

import (
	"fmt"

	tmkv "github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/cosmos-sdk/codec"
)

// TODO
// DecodeStore unmarshals the KVPair's Value to the corresponding gov type
func DecodeStore(cdc *codec.Codec, kvA, kvB tmkv.Pair) string {
	switch {
	case true:
		return ""

	default:
		panic(fmt.Sprintf("invalid asset key %X", kvA.Key))
	}
}
