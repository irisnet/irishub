package types

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValidateParams(t *testing.T) {
	defaultToken := GetNativeToken()
	tests := []struct {
		testCase string
		Params
		expectPass bool
	}{
		{"Minimum value",
			Params{
				TokenTaxRate:      sdk.ZeroDec(),
				MintTokenFeeRatio: sdk.ZeroDec(),
				IssueTokenBaseFee: sdk.NewCoin(defaultToken.Symbol, sdk.ZeroInt()),
			},
			true,
		},
		{"Maximum value",
			Params{
				TokenTaxRate:      sdk.NewDec(1),
				MintTokenFeeRatio: sdk.NewDec(1),
				IssueTokenBaseFee: sdk.NewCoin(defaultToken.Symbol, sdk.NewInt(math.MaxInt64)),
			},
			true,
		},
		{"TokenTaxRate less than the maximum",
			Params{
				TokenTaxRate:      sdk.NewDecWithPrec(-1, 1),
				MintTokenFeeRatio: sdk.NewDec(0),
				IssueTokenBaseFee: sdk.NewCoin(defaultToken.Symbol, sdk.NewInt(1)),
			},
			false,
		},
		{"MintTokenFeeRatio less than the maximum",
			Params{
				TokenTaxRate:      sdk.NewDec(0),
				MintTokenFeeRatio: sdk.NewDecWithPrec(-1, 1),
				IssueTokenBaseFee: sdk.NewCoin(defaultToken.Symbol, sdk.NewInt(1)),
			},
			false,
		},
		{"TokenTaxRate greater than the maximum",
			Params{
				TokenTaxRate:      sdk.NewDecWithPrec(11, 1),
				MintTokenFeeRatio: sdk.NewDec(1),
				IssueTokenBaseFee: sdk.NewCoin(defaultToken.Symbol, sdk.NewInt(1)),
			},
			false,
		},
		{"MintTokenFeeRatio greater than the maximum",
			Params{
				TokenTaxRate:      sdk.NewDec(1),
				MintTokenFeeRatio: sdk.NewDecWithPrec(11, 1),
				IssueTokenBaseFee: sdk.NewCoin(defaultToken.Symbol, sdk.NewInt(1)),
			},
			false,
		},
		{"IssueTokenBaseFee is negative",
			Params{
				TokenTaxRate:      sdk.NewDec(1),
				MintTokenFeeRatio: sdk.NewDec(1),
				IssueTokenBaseFee: sdk.Coin{Denom: defaultToken.Symbol, Amount: sdk.NewInt(-1)},
			},
			false,
		},
	}

	for _, tc := range tests {
		if tc.expectPass {
			require.Nil(t, ValidateParams(tc.Params), "test: %v", tc.testCase)
		} else {
			require.NotNil(t, ValidateParams(tc.Params), "test: %v", tc.testCase)
		}
	}
}
