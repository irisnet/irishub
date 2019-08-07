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
			doubleSwap := keeper.IsDoubleSwap(tc.denom1, tc.denom2)
			require.Equal(t, tc.isDoubleSwap, doubleSwap)
		})
	}
}

func TestGetReservePoolName(t *testing.T) {
	_, keeper, _ := createTestInput(t, sdk.NewInt(0), 0)

	cases := []struct {
		name         string
		denom1       string
		denom2       string
		expectResult string
		expectPass   bool
	}{
		{"denom1 is native", native, "btc", "s-btc", true},
		{"denom2 is native", "btc", native, "s-btc", true},
		{"denom1 equals denom2", "btc", "btc", "s-btc", false},
		{"neither denom is native", "eth", "btc", "s-btc", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reservePoolName, err := keeper.GetReservePoolName(tc.denom1, tc.denom2)
			if tc.expectPass {
				require.Equal(t, tc.expectResult, reservePoolName)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}
