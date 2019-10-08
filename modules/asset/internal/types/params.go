package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
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

func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
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

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
// nolint
func (p *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		{KeyAssetTaxRate, &p.AssetTaxRate},
		{KeyIssueTokenBaseFee, &p.IssueTokenBaseFee},
		{KeyMintTokenFeeRatio, &p.MintTokenFeeRatio},
		{KeyCreateGatewayBaseFee, &p.CreateGatewayBaseFee},
		{KeyGatewayAssetFeeRatio, &p.GatewayAssetFeeRatio},
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
		IssueTokenBaseFee:    sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewIntWithDecimal(60000, int(sdk.AttoScale))),
		MintTokenFeeRatio:    sdk.NewDecWithPrec(1, 1), // 0.1 (10%)
		CreateGatewayBaseFee: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewIntWithDecimal(120000, int(sdk.AttoScale))),
		GatewayAssetFeeRatio: sdk.NewDecWithPrec(1, 1), // 0.1 (10%)
	}
}
