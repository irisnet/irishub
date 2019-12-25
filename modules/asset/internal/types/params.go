package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

var _ params.ParamSet = (*Params)(nil)

// Parameter store keys
var (
	KeyAssetTaxRate      = []byte("AssetTaxRate")      //
	KeyIssueTokenBaseFee = []byte("IssueTokenBaseFee") //
	KeyMintTokenFeeRatio = []byte("MintTokenFeeRatio") //
	KeyAssetFeeDenom     = []byte("AssetFeeDenom")     //
)

// asset parameters
// issuance fee = IssueTokenBaseFee / (ln(len(symbol))/ln3)^4
type Params struct {
	AssetTaxRate      sdk.Dec `json:"asset_tax_rate" yaml:"asset_tax_rate"`             // e.g., 40%
	IssueTokenBaseFee sdk.Int `json:"issue_token_base_fee" yaml:"issue_token_base_fee"` // e.g., 300000*10^18
	MintTokenFeeRatio sdk.Dec `json:"mint_token_fee_ratio" yaml:"mint_token_fee_ratio"` // e.g., 10%
	AssetFeeDenom     string  `json:"asset_fee_denom" yaml:"asset_fee_denom"`           // e.g., iris
}

// NewParams asset params constructor
func NewParams(assetTaxRate sdk.Dec, issueTokenBaseFee sdk.Int,
	mintTokenFeeRatio sdk.Dec, assetFeeDenom string,
) Params {
	return Params{
		AssetTaxRate:      assetTaxRate,
		IssueTokenBaseFee: issueTokenBaseFee,
		MintTokenFeeRatio: mintTokenFeeRatio,
		AssetFeeDenom:     assetFeeDenom,
	}
}

// DefaultParams returns default asset module params
func DefaultParams() Params {
	return Params{
		AssetTaxRate:      sdk.NewDecWithPrec(4, 1), // 0.4 (40%)
		IssueTokenBaseFee: sdk.NewIntWithDecimal(60000, 18),
		MintTokenFeeRatio: sdk.NewDecWithPrec(1, 1), // 0.1 (10%)
		AssetFeeDenom:     sdk.DefaultBondDenom,
	}
}

// ParamSetPairs Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyAssetTaxRate, Value: &p.AssetTaxRate},
		{Key: KeyIssueTokenBaseFee, Value: &p.IssueTokenBaseFee},
		{Key: KeyMintTokenFeeRatio, Value: &p.MintTokenFeeRatio},
		{Key: KeyAssetFeeDenom, Value: &p.AssetFeeDenom},
	}
}

// Validate validates a set of params
func (p Params) Validate() error {
	// TODO should validate Params
	return nil
}
