package simulation_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/irisnet/irishub/v3/modules/mint/simulation"
	"github.com/irisnet/irishub/v3/modules/mint/types"
	"github.com/irisnet/irishub/v3/testutil"
)

func TestDecodeStore(t *testing.T) {
	minter := types.NewMinter(time.Now().UTC(), sdkmath.NewIntWithDecimal(2, 9))
	ec := testutil.MakeCodecs()
	dec := simulation.NewDecodeStore(ec.Marshaler)

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: types.MinterKey, Value: ec.Marshaler.MustMarshal(&minter)},
			{Key: []byte{0x99}, Value: []byte{0x99}},
		},
	}
	tests := []struct {
		name        string
		expectedLog string
	}{
		{"Minter", fmt.Sprintf("%v\n%v", minter, minter)},
		{"other", ""},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			switch i {
			case len(tests) - 1:
				require.Panics(t, func() { dec(kvPairs.Pairs[i], kvPairs.Pairs[i]) }, tt.name)
			default:
				require.Equal(t, tt.expectedLog, dec(kvPairs.Pairs[i], kvPairs.Pairs[i]), tt.name)
			}
		})
	}
}
