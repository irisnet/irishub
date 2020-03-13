package simulation

import (
	tmkv "github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/cosmos-sdk/codec"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding htlc type
func DecodeStore(cdc *codec.Codec, kvA, kvB tmkv.Pair) string {
	// TODO
	return ""
}
