package types

import (
	"math"
	"testing"

	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
)

func TestValidateParams(t *testing.T) {
	tests := []struct {
		testCase string
		Params
		expectPass bool
	}{
		{"Minimum value",
			Params{
				AssetTaxRate:         sdk.ZeroDec(),
				MintTokenFeeRatio:    sdk.ZeroDec(),
				GatewayAssetFeeRatio: sdk.ZeroDec(),
				IssueTokenBaseFee:    sdk.NewCoin(sdk.IrisAtto, sdk.ZeroInt()),
				CreateGatewayBaseFee: sdk.NewCoin(sdk.IrisAtto, sdk.ZeroInt()),
			},
			true,
		},
		{"Maximum value",
			Params{
				AssetTaxRate:         sdk.NewDec(1),
				MintTokenFeeRatio:    sdk.NewDec(1),
				GatewayAssetFeeRatio: sdk.NewDec(1),
				IssueTokenBaseFee:    sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(math.MaxInt64)),
				CreateGatewayBaseFee: sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(math.MaxInt64)),
			},
			true,
		},
		{"AssetTaxRate less than the maximum",
			Params{
				AssetTaxRate:         sdk.NewDecWithPrec(-1, 1),
				MintTokenFeeRatio:    sdk.NewDec(0),
				GatewayAssetFeeRatio: sdk.NewDec(0),
				IssueTokenBaseFee:    sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(1)),
				CreateGatewayBaseFee: sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(1)),
			},
			false,
		},
		{"MintTokenFeeRatio less than the maximum",
			Params{
				AssetTaxRate:         sdk.NewDec(0),
				MintTokenFeeRatio:    sdk.NewDecWithPrec(-1, 1),
				GatewayAssetFeeRatio: sdk.NewDec(0),
				IssueTokenBaseFee:    sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(1)),
				CreateGatewayBaseFee: sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(1)),
			},
			false,
		},
		{"GatewayAssetFeeRatio less than the maximum",
			Params{
				AssetTaxRate:         sdk.NewDec(0),
				MintTokenFeeRatio:    sdk.NewDec(0),
				GatewayAssetFeeRatio: sdk.NewDecWithPrec(-1, 1),
				IssueTokenBaseFee:    sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(1)),
				CreateGatewayBaseFee: sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(1)),
			},
			false,
		},
		{"AssetTaxRate greater than the maximum",
			Params{
				AssetTaxRate:         sdk.NewDecWithPrec(11, 1),
				MintTokenFeeRatio:    sdk.NewDec(1),
				GatewayAssetFeeRatio: sdk.NewDec(1),
				IssueTokenBaseFee:    sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(1)),
				CreateGatewayBaseFee: sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(1)),
			},
			false,
		},
		{"MintTokenFeeRatio greater than the maximum",
			Params{
				AssetTaxRate:         sdk.NewDec(1),
				MintTokenFeeRatio:    sdk.NewDecWithPrec(11, 1),
				GatewayAssetFeeRatio: sdk.NewDec(1),
				IssueTokenBaseFee:    sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(1)),
				CreateGatewayBaseFee: sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(1)),
			},
			false,
		},
		{"GatewayAssetFeeRatio greater than the maximum",
			Params{
				AssetTaxRate:         sdk.NewDec(1),
				MintTokenFeeRatio:    sdk.NewDec(1),
				GatewayAssetFeeRatio: sdk.NewDecWithPrec(11, 1),
				IssueTokenBaseFee:    sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(1)),
				CreateGatewayBaseFee: sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(1)),
			},
			false,
		},
		{"IssueTokenBaseFee is negative",
			Params{
				AssetTaxRate:         sdk.NewDec(1),
				MintTokenFeeRatio:    sdk.NewDec(1),
				GatewayAssetFeeRatio: sdk.NewDec(1),
				IssueTokenBaseFee:    sdk.Coin{Denom: sdk.IrisAtto, Amount: sdk.NewInt(-1)},
				CreateGatewayBaseFee: sdk.Coin{Denom: sdk.IrisAtto, Amount: sdk.NewInt(1)},
			},
			false,
		},
		{"CreateGatewayBaseFee is Negative",
			Params{
				AssetTaxRate:         sdk.NewDec(1),
				MintTokenFeeRatio:    sdk.NewDec(1),
				GatewayAssetFeeRatio: sdk.NewDec(1),
				IssueTokenBaseFee:    sdk.Coin{Denom: sdk.IrisAtto, Amount: sdk.NewInt(1)},
				CreateGatewayBaseFee: sdk.Coin{Denom: sdk.IrisAtto, Amount: sdk.NewInt(-1)},
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
