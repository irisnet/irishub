package simulation

// DONTCOVER

import (
	"bytes"
	"fmt"

	gogotypes "github.com/cosmos/gogoproto/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/irisnet/irismod/token/types"
	v1 "github.com/irisnet/irismod/token/types/v1"
)

// NewDecodeStore unmarshals the KVPair's Value to the corresponding token type
func NewDecodeStore(cdc codec.BinaryCodec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.PrefixTokenForSymbol):
			var tokenA, tokenB v1.Token
			cdc.MustUnmarshal(kvA.Value, &tokenA)
			cdc.MustUnmarshal(kvB.Value, &tokenB)
			return fmt.Sprintf("%v\n%v", tokenA, tokenB)
		case bytes.Equal(kvA.Key[:1], types.PrefixTokens):
			var symbolA, symbolB gogotypes.Value
			cdc.MustUnmarshal(kvA.Value, &symbolA)
			cdc.MustUnmarshal(kvB.Value, &symbolB)
			return fmt.Sprintf("%v\n%v", symbolA, symbolB)
		case bytes.Equal(kvA.Key[:1], types.PrefixTokenForMinUint):
			var symbolA, symbolB gogotypes.StringValue
			cdc.MustUnmarshal(kvA.Value, &symbolA)
			cdc.MustUnmarshal(kvB.Value, &symbolB)
			return fmt.Sprintf("%v\n%v", symbolA, symbolB)
		case bytes.Equal(kvA.Key[:1], types.PrefixBurnTokenAmt):
			var burntCoinA, burntCoinB sdk.Coin
			cdc.MustUnmarshal(kvA.Value, &burntCoinA)
			cdc.MustUnmarshal(kvB.Value, &burntCoinB)
			return fmt.Sprintf("%v\n%v", burntCoinA, burntCoinB)
		case bytes.Equal(kvA.Key[:1], types.PrefixParamsKey):
			var paramsA, paramsB v1.Params
			cdc.MustUnmarshal(kvA.Value, &paramsA)
			cdc.MustUnmarshal(kvB.Value, &paramsB)
			return fmt.Sprintf("%v\n%v", paramsA, paramsB)
		case bytes.Equal(kvA.Key[:1], types.PrefixTokenForContract):
			var symbolA, symbolB gogotypes.Value
			cdc.MustUnmarshal(kvA.Value, &symbolA)
			cdc.MustUnmarshal(kvB.Value, &symbolB)
			return fmt.Sprintf("%v\n%v", symbolA, symbolB)
		default:
			panic(fmt.Sprintf("invalid %s key prefix %X", types.ModuleName, kvA.Key[:1]))
		}
	}
}
