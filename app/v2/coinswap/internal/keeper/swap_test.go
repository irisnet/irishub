package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/irisnet/irishub/types"
)

var (
	native = sdk.IrisAtto
)

func TestIsDoubleSwap(t *testing.T) {
	ctx, keeper, _ := createTestInput(t, sdk.NewInt(0), 0)

	cases := []struct {
		name         string
		denom1       string
		denom2       string
		isDoubleSwap bool
	}{
		{"denom1 is native", native, "btc", false},
		{"denom2 is native", "btc", native, false},
		{"neither denom is native", "eth", "btc", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			doubleSwap := keeper.IsDoubleSwap(ctx, tc.denom1, tc.denom2)
			require.Equal(t, tc.isDoubleSwap, doubleSwap)
		})
	}
}
