package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestNextInflation(t *testing.T) {
	minter := NewMinter(time.Now(), sdkmath.NewIntWithDecimal(100, 18))
	tests := []struct{ params Params }{
		{Params{Inflation: sdkmath.LegacyNewDecWithPrec(20, 2), MintDenom: sdk.DefaultBondDenom}},
		{Params{Inflation: sdkmath.LegacyNewDecWithPrec(10, 2), MintDenom: sdk.DefaultBondDenom}},
		{Params{Inflation: sdkmath.LegacyNewDecWithPrec(5, 2), MintDenom: sdk.DefaultBondDenom}},
	}
	for _, tc := range tests {
		annualProvisions := minter.NextAnnualProvisions(tc.params)
		mintCoin := minter.BlockProvision(tc.params)
		blockProvision := annualProvisions.QuoInt(sdkmath.NewInt(12 * 60 * 8766))
		require.True(t, mintCoin.Amount.Equal(blockProvision.TruncateInt()), "mint amount:"+mintCoin.Amount.String()+", block provision amount: "+blockProvision.TruncateInt().String())
	}
}

func TestDefaultMinter(t *testing.T) {
	err := ValidateMinter(DefaultMinter())
	require.NoError(t, err)
}

func TestMinterValidate(t *testing.T) {
	tests := []struct {
		expectPass    bool
		LastUpdate    time.Time
		InflationBase sdkmath.Int
	}{
		{false, time.Unix(-1, -1), initialIssue.Mul(sdkmath.NewIntWithDecimal(1, 18))},
		{false, time.Unix(0, 0), initialIssue.Mul(sdkmath.NewIntWithDecimal(0, 0))},
		{true, time.Unix(0, 0), initialIssue.Mul(sdkmath.NewIntWithDecimal(1, 18))},
	}
	for i, tc := range tests {
		minter := NewMinter(tc.LastUpdate, tc.InflationBase)
		err := ValidateMinter(minter)
		if tc.expectPass {
			require.NoError(t, err, "%d: %+v", i, err)
		} else {
			require.Error(t, err, "%d: %+v", i, err)
		}
	}
}
