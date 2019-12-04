package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValidateParams(t *testing.T) {
	// check that valid case work
	defaultParams := DefaultParams()
	err := ValidateParams(defaultParams)
	require.Nil(t, err)

	// all cases should return an error
	invalidTests := []struct {
		name   string
		params Params
		result bool
	}{
		{"fee == 0 ", NewParams(sdk.ZeroDec()), false},
		{"fee < 1", NewParams(sdk.NewDecWithPrec(1000, 2)), false},
		{"fee numerator < 0", NewParams(sdk.NewDecWithPrec(-1, 1)), false},
		{"fee denominator < 0", NewParams(sdk.NewDecWithPrec(1, -1)), false},
	}

	for _, tc := range invalidTests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateParams(tc.params)
			if err != nil {
				require.False(t, tc.result)
			} else {
				require.True(t, tc.result)
			}
		})
	}
}
