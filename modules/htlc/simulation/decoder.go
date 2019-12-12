package simulation

import (
	cmn "github.com/tendermint/tendermint/libs/common"

	"github.com/cosmos/cosmos-sdk/codec"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding htlc type
func DecodeStore(cdc *codec.Codec, kvA, kvB cmn.KVPair) string {
	// TODO
	return ""
}
