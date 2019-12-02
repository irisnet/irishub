package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

var _ params.ParamSet = (*Params)(nil)

// Parameter store keys
var (
	KeyAssetTaxRate         = []byte("AssetTaxRate")
	KeyIssueTokenBaseFee    = []byte("IssueTokenBaseFee")
	KeyMintTokenFeeRatio    = []byte("MintTokenFeeRatio")
	KeyCreateGatewayBaseFee = []byte("CreateGatewayBaseFee")
	KeyGatewayAssetFeeRatio = []byte("GatewayAssetFeeRatio")
	KeyAssetFeeDenom        = []byte("AssetFeeDenom")
)

// asset parameters
// issuance fee = IssueTokenBaseFee / (ln(len(symbol))/ln3)^4
type Params struct {
	AssetTaxRate         sdk.Dec `json:"asset_tax_rate"`          // e.g., 40%
	IssueTokenBaseFee    sdk.Int `json:"issue_token_base_fee"`    // e.g., 300000*10^18
	MintTokenFeeRatio    sdk.Dec `json:"mint_token_fee_ratio"`    // e.g., 10%
	CreateGatewayBaseFee sdk.Int `json:"create_gateway_base_fee"` // e.g., 600000*10^18
	GatewayAssetFeeRatio sdk.Dec `json:"gateway_asset_fee_ratio"` // e.g., 10%
	AssetFeeDenom        string  `json:"asset_fee_denom"`         // e.g., iris
}

// ParamTable for asset module
func ParamTypeTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(assetTaxRate sdk.Dec, issueTokenBaseFee sdk.Int, mintTokenFeeRatio sdk.Dec,
	createGatewayBaseFee sdk.Int, gatewayAssetFeeRatio sdk.Dec, assetFeeDenom string,
) Params {
	return Params{
		AssetTaxRate:         assetTaxRate,
		IssueTokenBaseFee:    issueTokenBaseFee,
		MintTokenFeeRatio:    mintTokenFeeRatio,
		CreateGatewayBaseFee: createGatewayBaseFee,
		GatewayAssetFeeRatio: gatewayAssetFeeRatio,
		AssetFeeDenom:        assetFeeDenom,
	}
}

// default asset module params
func DefaultParams() Params {
	return Params{
		AssetTaxRate:         sdk.NewDecWithPrec(4, 1), // 0.4 (40%)
		IssueTokenBaseFee:    sdk.NewIntWithDecimal(60000, 18),
		MintTokenFeeRatio:    sdk.NewDecWithPrec(1, 1), // 0.1 (10%)
		CreateGatewayBaseFee: sdk.NewIntWithDecimal(120000, 18),
		GatewayAssetFeeRatio: sdk.NewDecWithPrec(1, 1), // 0.1 (10%)
		AssetFeeDenom:        "iris-atto",
	}
}

func (p Params) String() string {
	return fmt.Sprintf(`Asset Params:
  Asset TaxRate:          %s
  Issue Token BaseFee:     %s
  Mint Token FeeRatio:     %s
  Create Gateway BaseFee:  %s
  Gateway AssetFee Ratio:  %s
  Asset Fee Denom:         %s`,
		p.AssetTaxRate.String(), p.IssueTokenBaseFee.String(), p.MintTokenFeeRatio.String(),
		p.CreateGatewayBaseFee.String(), p.GatewayAssetFeeRatio.String(), p.AssetFeeDenom)
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyAssetTaxRate, Value: &p.AssetTaxRate},
		{Key: KeyIssueTokenBaseFee, Value: &p.IssueTokenBaseFee},
		{Key: KeyMintTokenFeeRatio, Value: &p.MintTokenFeeRatio},
		{Key: KeyCreateGatewayBaseFee, Value: &p.CreateGatewayBaseFee},
		{Key: KeyGatewayAssetFeeRatio, Value: &p.GatewayAssetFeeRatio},
		{Key: KeyAssetFeeDenom, Value: &p.AssetFeeDenom},
	}
}
