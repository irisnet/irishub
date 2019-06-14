package asset

import (
	"fmt"
	"strconv"

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
	KeyIssueFTBaseFee       = []byte("IssueFTBaseFee")
	KeyMintFTFeeRatio       = []byte("MintFTBaseRatio")
	KeyCreateGatewayBaseFee = []byte("CreateGatewayBaseFee")
	KeyGatewayAssetFeeRatio = []byte("GatewayAssetFeeRatio")
)

// ParamTable for asset module
func ParamTypeTable() params.TypeTable {
	return params.NewTypeTable().RegisterParamSet(&Params{})
}

// asset params
type Params struct {
	AssetTaxRate         sdk.Dec `json:"asset_tax_rate"`          // e.g., 40%
	IssueFTBaseFee       uint32  `json:"issue_ft_base_fee"`       // e.g., 300000
	MintFTFeeRatio       sdk.Dec `json:"mint_ft_fee_ratio"`       // e.g., 10%
	CreateGatewayBaseFee uint32  `json:"create_gateway_base_fee"` // e.g., 600000
	GatewayAssetFeeRatio sdk.Dec `json:"gateway_asset_fee_ratio"` // e.g., 10%
} // issuance fee = IssueFTBaseFee / (ln(len(symbol))/ln3)^4

func (p Params) String() string {
	return fmt.Sprintf(`Asset Params:
  Asset Tax Rate:                                           %s
  Base Fee for Issuing Fungible Token:                      %d
  Fee Ratio for Minting (vs Issuing) Fungible Token:        %s
  Base Fee for Creating Gateway:                            %d
  Fee Ratio for Gateway (vs Native) Assets:                 %s`,
		p.AssetTaxRate.String(), p.IssueFTBaseFee, p.MintFTFeeRatio.String(), p.CreateGatewayBaseFee, p.GatewayAssetFeeRatio.String())
}

// Implements params.ParamSet
func (p *Params) GetParamSpace() string {
	return DefaultParamSpace
}

func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
		{KeyAssetTaxRate, &p.AssetTaxRate},
		{KeyIssueFTBaseFee, &p.IssueFTBaseFee},
		{KeyMintFTFeeRatio, &p.MintFTFeeRatio},
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
	case string(KeyIssueFTBaseFee):
		fee, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		return uint32(fee), nil
	case string(KeyMintFTFeeRatio):
		ratio, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateMintFTBaseFeeRatio(ratio); err != nil {
			return nil, err
		}
		return ratio, nil
	case string(KeyCreateGatewayBaseFee):
		fee, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		return uint32(fee), nil
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
	return "", fmt.Errorf("This method is not implemented!")
}

// default asset module params
func DefaultParams() Params {
	return Params{
		AssetTaxRate:         sdk.NewDecWithPrec(4, 1), // 0.4 (40%)
		IssueFTBaseFee:       300000,
		MintFTFeeRatio:       sdk.NewDecWithPrec(1, 1), // 0.1 (10%)
		CreateGatewayBaseFee: 600000,
		GatewayAssetFeeRatio: sdk.NewDecWithPrec(1, 1), // 0.1 (10%)
	}
}

// default asset module params for test
func DefaultParamsForTest() Params {
	return Params{
		AssetTaxRate:         sdk.NewDecWithPrec(4, 1), // 0.4 (40%)
		IssueFTBaseFee:       300000,
		MintFTFeeRatio:       sdk.NewDecWithPrec(1, 1), // 0.1 (10%)
		CreateGatewayBaseFee: 600000,
		GatewayAssetFeeRatio: sdk.NewDecWithPrec(1, 1), // 0.1 (10%)
	}
}

func validateParams(p Params) error {
	if sdk.NetworkType != sdk.Mainnet {
		return nil
	}

	if err := validateAssetTaxRate(p.AssetTaxRate); err != nil {
		return err
	}
	if err := validateMintFTBaseFeeRatio(p.MintFTFeeRatio); err != nil {
		return err
	}
	if err := validateGatewayAssetFeeRatio(p.GatewayAssetFeeRatio); err != nil {
		return err
	}

	return nil
}

func validateAssetTaxRate(v sdk.Dec) sdk.Error {
	if v.GT(sdk.NewDecWithPrec(1, 0)) || v.LT(sdk.NewDecWithPrec(0, 0)) {
		return sdk.NewError(
			params.DefaultCodespace,
			params.CodeInvalidAssetTaxRate,
			fmt.Sprintf("Asset Tax Rate [%s] should be between [0, 1]", v.String()),
		)
	}
	return nil
}

func validateMintFTBaseFeeRatio(v sdk.Dec) sdk.Error {
	if v.GT(sdk.NewDecWithPrec(1, 0)) || v.LT(sdk.NewDecWithPrec(0, 0)) {
		return sdk.NewError(
			params.DefaultCodespace,
			params.CodeInvalidMintFTBaseFeeRatio,
			fmt.Sprintf("Base Fee Ratio for Minting FTs [%s] should be between [0, 1]", v.String()),
		)
	}
	return nil
}

func validateGatewayAssetFeeRatio(v sdk.Dec) sdk.Error {
	if v.GT(sdk.NewDecWithPrec(1, 0)) || v.LT(sdk.NewDecWithPrec(0, 0)) {
		return sdk.NewError(
			params.DefaultCodespace,
			params.CodeInvalidGatewayAssetFeeRatio,
			fmt.Sprintf("Fee Ratio for Gateway Assets [%s] should be between [0, 1]", v.String()),
		)
	}
	return nil
}

// get asset params from the global param store
func (k Keeper) GetParamSet(ctx sdk.Context) Params {
	var params Params
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// set asset params from the global param store
func (k Keeper) SetParamSet(ctx sdk.Context, params Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
