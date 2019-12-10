package types

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestNextInflation(t *testing.T) {
	minter := NewMinter(time.Now(), sdk.NewIntWithDecimal(100, 18))
	tests := []struct {
		params Params
	}{
		{Params{Inflation: sdk.NewDecWithPrec(20, 2), MintDenom: "iris"}},
		{Params{Inflation: sdk.NewDecWithPrec(10, 2), MintDenom: "iris"}},
		{Params{Inflation: sdk.NewDecWithPrec(5, 2), MintDenom: "iris"}},
	}
	for _, tc := range tests {
		annualProvisions := minter.NextAnnualProvisions(tc.params)
		mintCoin := minter.BlockProvision(tc.params)

		blockProvision := annualProvisions.QuoInt(sdk.NewInt(12 * 60 * 8766))

		require.True(t, mintCoin.Amount.Equal(blockProvision.TruncateInt()), "mint amount:"+mintCoin.Amount.String()+", block provision amount: "+blockProvision.TruncateInt().String())
	}
}
