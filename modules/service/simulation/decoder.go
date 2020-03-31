package simulation

import (
	"fmt"

	tmkv "github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/cosmos-sdk/codec"
)

// DecodeStore unmarshals the Pair's Value to the corresponding service type
func DecodeStore(cdc *codec.Codec, kvA, kvB tmkv.Pair) string {
	switch {
	default:
		panic(fmt.Sprintf("invalid service key prefix %X", kvA.Key[:1]))
	}
}
