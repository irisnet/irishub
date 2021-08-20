package simulation

import (
	"bytes"
	"fmt"

	"github.com/irisnet/irismod/modules/farm/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"
)

// NewDecodeStore unmarshals the KVPair's Value to the corresponding slashing type
func NewDecodeStore(cdc codec.Marshaler) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.FarmPoolKey):
			var farmPoolA, farmPoolB types.FarmPool
			cdc.MustUnmarshalBinaryBare(kvA.Value, &farmPoolA)
			cdc.MustUnmarshalBinaryBare(kvA.Value, &farmPoolB)
			return fmt.Sprintf("%v\n%v", farmPoolA, farmPoolB)

		case bytes.Equal(kvA.Key[:1], types.FarmPoolRuleKey):
			var farmPoolRuleA, farmPoolRuleB types.RewardRule
			cdc.MustUnmarshalBinaryBare(kvA.Value, &farmPoolRuleA)
			cdc.MustUnmarshalBinaryBare(kvA.Value, &farmPoolRuleB)
			return fmt.Sprintf("%v\n%v", farmPoolRuleA, farmPoolRuleB)

		case bytes.Equal(kvA.Key[:1], types.FarmerKey):
			var farmerA, farmerB types.FarmInfo
			cdc.MustUnmarshalBinaryBare(kvA.Value, &farmerA)
			cdc.MustUnmarshalBinaryBare(kvA.Value, &farmerB)
			return fmt.Sprintf("%v\n%v", farmerA, farmerB)

		case bytes.Equal(kvA.Key[:1], types.ActiveFarmPoolKey):
			var ActiveFarmPoolA, ActiveFarmPoolB types.FarmPool
			cdc.MustUnmarshalBinaryBare(kvA.Value, &ActiveFarmPoolA)
			cdc.MustUnmarshalBinaryBare(kvA.Value, &ActiveFarmPoolB)
			return fmt.Sprintf("%v\n%v", ActiveFarmPoolA, ActiveFarmPoolB)

		default:
			panic(fmt.Sprintf("invalid farm key prefix %X", kvA.Key[:1]))
		}
	}

}
