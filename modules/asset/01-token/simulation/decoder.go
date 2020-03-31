package simulation

// DONTCOVER

import (
	"fmt"
	"strings"

	tmkv "github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/irisnet/irishub/modules/asset/01-token/internal/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding gov type
func DecodeStore(cdc *codec.Codec, kvA, kvB tmkv.Pair) string {
	switch {
	case strings.HasPrefix(string(kvA.Key), string(types.KeyTokenPrefix())):
		var tokenA, tokenB types.Tokens
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &tokenA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &tokenB)
		return fmt.Sprintf("%v\n%v", tokenA, tokenB)

	default:
		panic(fmt.Sprintf("invalid token key %X", kvA.Key))
	}
}
