package types

import (
	"fmt"
	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

var _ params.ParamSet = (*Params)(nil)

const (
	DefaultParamSpace = "asset"
)

// parameter keys
var (
	KeyAssetTaxRate         = []byte("AssetTaxRate")
	KeyIssueTokenBaseFee    = []byte("IssueTokenBaseFee")
	KeyMintTokenFeeRatio    = []byte("MintTokenFeeRatio")
	KeyCreateGatewayBaseFee = []byte("CreateGatewayBaseFee")
	KeyGatewayAssetFeeRatio = []byte("GatewayAssetFeeRatio")
)

// ParamTable for asset module
func ParamTypeTable() params.TypeTable {
	return params.NewTypeTable().RegisterParamSet(&Params{})
}

// asset params
type Params struct {
	AssetTaxRate         sdk.Dec  `json:"asset_tax_rate"`          // e.g., 40%
	IssueTokenBaseFee    sdk.Coin `json:"issue_token_base_fee"`    // e.g., 300000*10^18iris-atto
	MintTokenFeeRatio    sdk.Dec  `json:"mint_token_fee_ratio"`    // e.g., 10%
	CreateGatewayBaseFee sdk.Coin `json:"create_gateway_base_fee"` // e.g., 600000*10^18iris-atto
	GatewayAssetFeeRatio sdk.Dec  `json:"gateway_asset_fee_ratio"` // e.g., 10%
} // issuance fee = IssueTokenBaseFee / (ln(len(symbol))/ln3)^4

func (p Params) String() string {
	return fmt.Sprintf(`Asset Params:
  Asset Tax Rate:                              %s
  Base Fee for Issuing Token:                  %s
  Fee Ratio for Minting (vs Issuing) Token:    %s
  Base Fee for Creating Gateway:               %s
  Fee Ratio for Gateway (vs Native) Token:     %s`,
		p.AssetTaxRate.String(), p.IssueTokenBaseFee.String(), p.MintTokenFeeRatio.String(), p.CreateGatewayBaseFee.String(), p.GatewayAssetFeeRatio.String())
}

// Implements params.ParamSet
func (p *Params) GetParamSpace() string {
	return DefaultParamSpace
}

func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
		{KeyAssetTaxRate, &p.AssetTaxRate},
		{KeyIssueTokenBaseFee, &p.IssueTokenBaseFee},
		{KeyMintTokenFeeRatio, &p.MintTokenFeeRatio},
		{KeyCreateGatewayBaseFee, &p.CreateGatewayBaseFee},
		{KeyGatewayAssetFeeRatio, &p.GatewayAssetFeeRatio},
	}
}

func (p *Params) Validate(key string, value string) (interface{}, sdk.Error) {
	switch key {
	case string(KeyAssetTaxRate):
		rate, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateAssetTaxRate(rate); err != nil {
			return nil, err
		}
		return rate, nil
	case string(KeyIssueTokenBaseFee):
		fee, err := sdk.ParseCoin(value)
		if err != nil || fee.Denom != sdk.NativeTokenMinDenom {
			return nil, params.ErrInvalidString(value)
		}
		return fee, nil
	case string(KeyMintTokenFeeRatio):
		ratio, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateMintTokenFeeRatio(ratio); err != nil {
			return nil, err
		}
		return ratio, nil
	case string(KeyCreateGatewayBaseFee):
		fee, err := sdk.ParseCoin(value)
		if err != nil || fee.Denom != sdk.NativeTokenMinDenom {
			return nil, params.ErrInvalidString(value)
		}

		return fee, nil
	case string(KeyGatewayAssetFeeRatio):
		ratio, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateGatewayAssetFeeRatio(ratio); err != nil {
			return nil, err
		}
		return ratio, nil
	default:
		return nil, sdk.NewError(params.DefaultCodespace, params.CodeInvalidKey, fmt.Sprintf("%s is an invalid key", key))
	}
}

func (p *Params) StringFromBytes(cdc *codec.Codec, key string, bytes []byte) (string, error) {
	return "", fmt.Errorf("this method is not implemented")
}

// default asset module params
func DefaultParams() Params {
	return Params{
		AssetTaxRate:         sdk.NewDecWithPrec(4, 1), // 0.4 (40%)
		IssueTokenBaseFee:    sdk.NewCoin(sdk.NativeTokenMinDenom, sdk.NewIntWithDecimal(300000, 18)),
		MintTokenFeeRatio:    sdk.NewDecWithPrec(1, 1), // 0.1 (10%)
		CreateGatewayBaseFee: sdk.NewCoin(sdk.NativeTokenMinDenom, sdk.NewIntWithDecimal(600000, 18)),
		GatewayAssetFeeRatio: sdk.NewDecWithPrec(1, 1), // 0.1 (10%)
	}
}

// default asset module params for test
func DefaultParamsForTest() Params {
	return Params{
		AssetTaxRate:         sdk.NewDecWithPrec(4, 1), // 0.4 (40%)
		IssueTokenBaseFee:    sdk.NewCoin(sdk.NativeTokenMinDenom, sdk.NewIntWithDecimal(30, 18)),
		MintTokenFeeRatio:    sdk.NewDecWithPrec(1, 1), // 0.1 (10%)
		CreateGatewayBaseFee: sdk.NewCoin(sdk.NativeTokenMinDenom, sdk.NewIntWithDecimal(60, 18)),
		GatewayAssetFeeRatio: sdk.NewDecWithPrec(1, 1), // 0.1 (10%)
	}
}

func ValidateParams(p Params) error {
	if err := validateAssetTaxRate(p.AssetTaxRate); err != nil {
		return err
	}
	if err := validateMintTokenFeeRatio(p.MintTokenFeeRatio); err != nil {
		return err
	}
	if err := validateGatewayAssetFeeRatio(p.GatewayAssetFeeRatio); err != nil {
		return err
	}

	return nil
}

func validateAssetTaxRate(v sdk.Dec) sdk.Error {
	if v.GT(sdk.NewDec(1)) || v.LTE(sdk.ZeroDec()) {
		return sdk.NewError(
			params.DefaultCodespace,
			params.CodeInvalidAssetTaxRate,
			fmt.Sprintf("Asset Tax Rate [%s] should be between (0, 1]", v.String()),
		)
	}
	return nil
}

func validateMintTokenFeeRatio(v sdk.Dec) sdk.Error {
	if v.GTE(sdk.NewDec(1)) || v.LTE(sdk.ZeroDec()) {
		return sdk.NewError(
			params.DefaultCodespace,
			params.CodeInvalidMintTokenFeeRatio,
			fmt.Sprintf("Fee Ratio for Minting Tokens [%s] should be between (0, 1)", v.String()),
		)
	}
	return nil
}

func validateGatewayAssetFeeRatio(v sdk.Dec) sdk.Error {
	if v.GTE(sdk.NewDec(1)) || v.LTE(sdk.ZeroDec()) {
		return sdk.NewError(
			params.DefaultCodespace,
			params.CodeInvalidGatewayAssetFeeRatio,
			fmt.Sprintf("Fee Ratio for Gateway Tokens [%s] should be between (0, 1)", v.String()),
		)
	}
	return nil
}
