package types

import (
	"fmt"

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
		params.NewParamSetPair(KeyAssetTaxRate, &p.AssetTaxRate, validateAssetTaxRate),
		params.NewParamSetPair(KeyIssueTokenBaseFee, &p.IssueTokenBaseFee, validateIssueTokenBaseFee),
		params.NewParamSetPair(KeyMintTokenFeeRatio, &p.MintTokenFeeRatio, validateMintTokenFeeRatio),
	}
}

// Validate validates a set of params
func (p Params) Validate() error {
	// TODO should validate Params
	return nil
}

func validateAssetTaxRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("inflation rate change cannot be negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("inflation rate change too large: %s", v)
	}

	return nil
}

func validateIssueTokenBaseFee(i interface{}) error {
	_, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	//if v.IsValid() {
	//	return fmt.Errorf("issue token base fee change cannot be invalid: %s", v.String())
	//}

	return nil
}

func validateMintTokenFeeRatio(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("mint token fee ratio change cannot be negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("mint token fee ratio change too large: %s", v)
	}

	return nil
}
