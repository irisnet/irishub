package simulation

// DONTCOVER

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	cmn "github.com/tendermint/tendermint/libs/common"
)

// TODO
// DecodeStore unmarshals the KVPair's Value to the corresponding gov type
func DecodeStore(cdc *codec.Codec, kvA, kvB cmn.KVPair) string {
	switch {
	case true:
		return ""

	default:
		panic(fmt.Sprintf("invalid asset key %X", kvA.Key))
	}
}
