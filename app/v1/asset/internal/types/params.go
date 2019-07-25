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
  asset/AssetTaxRate:          %s
  asset/IssueTokenBaseFee:     %s
  asset/MintTokenFeeRatio:     %s
  asset/CreateGatewayBaseFee:  %s
  asset/GatewayAssetFeeRatio:  %s`,
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
		if err != nil || fee.Denom != sdk.IrisAtto {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateIssueTokenBaseFee(fee); err != nil {
			return nil, err
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
		if err != nil || fee.Denom != sdk.IrisAtto {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateCreateGatewayBaseFee(fee); err != nil {
			return nil, err
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

func (p *Params) ReadOnly() bool {
	return false
}

// default asset module params
func DefaultParams() Params {
	return Params{
		AssetTaxRate:         sdk.NewDecWithPrec(4, 1), // 0.4 (40%)
		IssueTokenBaseFee:    sdk.NewCoin(sdk.IrisAtto, sdk.NewIntWithDecimal(60000, int(sdk.AttoScale))),
		MintTokenFeeRatio:    sdk.NewDecWithPrec(1, 1), // 0.1 (10%)
		CreateGatewayBaseFee: sdk.NewCoin(sdk.IrisAtto, sdk.NewIntWithDecimal(120000, int(sdk.AttoScale))),
		GatewayAssetFeeRatio: sdk.NewDecWithPrec(1, 1), // 0.1 (10%)
	}
}

// default asset module params for test
func DefaultParamsForTest() Params {
	return Params{
		AssetTaxRate:         sdk.NewDecWithPrec(4, 1), // 0.4 (40%)
		IssueTokenBaseFee:    sdk.NewCoin(sdk.IrisAtto, sdk.NewIntWithDecimal(30, int(sdk.AttoScale))),
		MintTokenFeeRatio:    sdk.NewDecWithPrec(1, 1), // 0.1 (10%)
		CreateGatewayBaseFee: sdk.NewCoin(sdk.IrisAtto, sdk.NewIntWithDecimal(60, int(sdk.AttoScale))),
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
	if err := validateIssueTokenBaseFee(p.IssueTokenBaseFee); err != nil {
		return err
	}
	if err := validateCreateGatewayBaseFee(p.CreateGatewayBaseFee); err != nil {
		return err
	}

	return nil
}

func validateAssetTaxRate(v sdk.Dec) sdk.Error {
	if v.GT(sdk.NewDec(1)) || v.LT(sdk.ZeroDec()) {
		return sdk.NewError(
			params.DefaultCodespace,
			params.CodeInvalidAssetTaxRate,
			fmt.Sprintf("Asset tax rate [%s] should be between [0, 1]", v.String()),
		)
	}
	return nil
}

func validateMintTokenFeeRatio(v sdk.Dec) sdk.Error {
	if v.GT(sdk.NewDec(1)) || v.LT(sdk.ZeroDec()) {
		return sdk.NewError(
			params.DefaultCodespace,
			params.CodeInvalidMintTokenFeeRatio,
			fmt.Sprintf("Fee ratio for minting tokens [%s] should be between [0, 1]", v.String()),
		)
	}
	return nil
}

func validateGatewayAssetFeeRatio(v sdk.Dec) sdk.Error {
	if v.GT(sdk.NewDec(1)) || v.LT(sdk.ZeroDec()) {
		return sdk.NewError(
			params.DefaultCodespace,
			params.CodeInvalidGatewayAssetFeeRatio,
			fmt.Sprintf("Fee ratio for gateway tokens [%s] should be between [0, 1]", v.String()),
		)
	}
	return nil
}

func validateIssueTokenBaseFee(coin sdk.Coin) sdk.Error {
	if !coin.IsNotNegative() {
		return sdk.NewError(
			params.DefaultCodespace,
			params.CodeInvalidIssueTokenBaseFee,
			fmt.Sprintf("Base fee for issuing token should not be negative"),
		)
	}
	return nil
}

func validateCreateGatewayBaseFee(coin sdk.Coin) sdk.Error {
	if !coin.IsNotNegative() {
		return sdk.NewError(
			params.DefaultCodespace,
			params.CodeInvalidCreateGatewayBaseFee,
			fmt.Sprintf("Base fee for creating gateway should not be negative"),
		)
	}
	return nil
}
