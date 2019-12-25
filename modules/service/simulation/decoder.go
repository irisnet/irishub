package simulation

import (
	"fmt"

	cmn "github.com/tendermint/tendermint/libs/common"

	"github.com/cosmos/cosmos-sdk/codec"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding service type
func DecodeStore(cdc *codec.Codec, kvA, kvB cmn.KVPair) string {
	switch {
	default:
		panic(fmt.Sprintf("invalid service key prefix %X", kvA.Key[:1]))
	}
}
