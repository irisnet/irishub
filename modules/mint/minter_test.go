package mint

import (
	"testing"
	"time"

	stakeTypes "github.com/irisnet/irishub/modules/stake/types"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
)

func TestNextInflation(t *testing.T) {
	minter := NewMinter(time.Now(), stakeTypes.StakeDenom, sdk.NewIntWithDecimal(100, 18))
	tests := []struct {
		params Params
	}{
		{DefaultParams()},
		{DefaultParams()},
		{DefaultParams()},
		{DefaultParams()},
		{Params{sdk.NewDecWithPrec(20, 2)}},
		{Params{sdk.NewDecWithPrec(10, 2)}},
		{Params{sdk.NewDecWithPrec(5, 2)}},
	}
	for _, tc := range tests {
		annualProvisions := minter.NextAnnualProvisions(tc.params)
		mintCoin := minter.BlockProvision(annualProvisions)

		blockProvision := annualProvisions.QuoInt(sdk.NewInt(12 * 60 * 8766))

		require.True(t, mintCoin.Amount.Equal(blockProvision.TruncateInt()), "mint amount:"+mintCoin.Amount.String()+", block provision amount: "+blockProvision.TruncateInt().String())
	}
}
