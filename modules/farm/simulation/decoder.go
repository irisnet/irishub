package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/irisnet/irismod/modules/farm/types"
)

// NewDecodeStore unmarshals the KVPair's Value to the corresponding slashing type
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.FarmPoolKey):
			var farmPoolA, farmPoolB types.FarmPool
			cdc.MustUnmarshal(kvA.Value, &farmPoolA)
			cdc.MustUnmarshal(kvA.Value, &farmPoolB)
			return fmt.Sprintf("%v\n%v", farmPoolA, farmPoolB)

		case bytes.Equal(kvA.Key[:1], types.FarmPoolRuleKey):
			var farmPoolRuleA, farmPoolRuleB types.RewardRule
			cdc.MustUnmarshal(kvA.Value, &farmPoolRuleA)
			cdc.MustUnmarshal(kvA.Value, &farmPoolRuleB)
			return fmt.Sprintf("%v\n%v", farmPoolRuleA, farmPoolRuleB)

		case bytes.Equal(kvA.Key[:1], types.FarmerKey):
			var farmerA, farmerB types.FarmInfo
			cdc.MustUnmarshal(kvA.Value, &farmerA)
			cdc.MustUnmarshal(kvA.Value, &farmerB)
			return fmt.Sprintf("%v\n%v", farmerA, farmerB)

		case bytes.Equal(kvA.Key[:1], types.ActiveFarmPoolKey):
			var ActiveFarmPoolA, ActiveFarmPoolB types.FarmPool
			cdc.MustUnmarshal(kvA.Value, &ActiveFarmPoolA)
			cdc.MustUnmarshal(kvA.Value, &ActiveFarmPoolB)
			return fmt.Sprintf("%v\n%v", ActiveFarmPoolA, ActiveFarmPoolB)

		default:
			panic(fmt.Sprintf("invalid farm key prefix %X", kvA.Key[:1]))
		}
	}

}
