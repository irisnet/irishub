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
)

// asset parameters
// issuance fee = IssueTokenBaseFee / (ln(len(symbol))/ln3)^4
type Params struct {
	AssetTaxRate      sdk.Dec  `json:"asset_tax_rate" yaml:"asset_tax_rate"`             // e.g., 40%
	IssueTokenBaseFee sdk.Coin `json:"issue_token_base_fee" yaml:"issue_token_base_fee"` // e.g., 60000*10^18iris-atto
	MintTokenFeeRatio sdk.Dec  `json:"mint_token_fee_ratio" yaml:"mint_token_fee_ratio"` // e.g., 10%
}

// NewParams asset params constructor
func NewParams(assetTaxRate sdk.Dec, issueTokenBaseFee sdk.Coin,
	mintTokenFeeRatio sdk.Dec,
) Params {
	return Params{
		AssetTaxRate:      assetTaxRate,
		IssueTokenBaseFee: issueTokenBaseFee,
		MintTokenFeeRatio: mintTokenFeeRatio,
	}
}

// DefaultParams returns default asset module params
func DefaultParams() Params {
	return Params{
		AssetTaxRate:      sdk.NewDecWithPrec(4, 1), // 0.4 (40%)
		IssueTokenBaseFee: sdk.NewCoin(IrisToken().MinUnit, sdk.NewIntWithDecimal(60000, int(IrisToken().Scale))),
		MintTokenFeeRatio: sdk.NewDecWithPrec(1, 1), // 0.1 (10%)
	}
}

// ParamSetPairs Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyAssetTaxRate, Value: &p.AssetTaxRate},
		{Key: KeyIssueTokenBaseFee, Value: &p.IssueTokenBaseFee},
		{Key: KeyMintTokenFeeRatio, Value: &p.MintTokenFeeRatio},
	}
}

// Validate validates a set of params
func (p Params) Validate() error {
	// TODO should validate Params
	return nil
}
