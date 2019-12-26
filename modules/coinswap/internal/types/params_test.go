package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValidateParams(t *testing.T) {
	// check that valid case work
	defaultParams := DefaultParams()
	err := defaultParams.Validate()
	require.NoError(t, err)

	require.Panics(t, func() { sdk.NewDecWithPrec(1, 19) }, "should panic")
	require.Panics(t, func() { sdk.NewDecWithPrec(1, -1) }, "should panic")

	// all cases should return an error
	invalidTests := []struct {
		name   string
		params Params
		result bool
	}{
		{"fee == 0 ", NewParams(sdk.ZeroDec(), StandardDenom), false},
		{"fee < 1", NewParams(sdk.NewDecWithPrec(1000, 2), StandardDenom), false},
		{"fee numerator < 0", NewParams(sdk.NewDecWithPrec(-1, 1), StandardDenom), false},
		{"invalid denom", NewParams(sdk.NewDecWithPrec(1, 1), ""), false},
		{"valid", NewParams(sdk.NewDecWithPrec(1, 1), StandardDenom), true},
	}

	for _, tc := range invalidTests {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.params.Validate(); err != nil {
				require.False(t, tc.result)
			} else {
				require.True(t, tc.result)
			}
		})
	}
}
