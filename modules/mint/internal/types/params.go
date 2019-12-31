package types

import (
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/params"
)

var _ params.ParamSet = (*Params)(nil)

// default paramspace for params keeper
const (
	DefaultParamSpace = "mint"
	MintDenom         = sdk.DefaultBondDenom
)

//Parameter store key
var (
	// params store for inflation params
	KeyInflation = []byte("Inflation")
	KeyMintDenom = []byte("MintDenom")
)

// ParamTable for mint module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// Params defines mint parameters
type Params struct {
	Inflation sdk.Dec `json:"inflation" yaml:"inflation"`   // inflation rate
	MintDenom string  `json:"mint_denom" yaml:"mint_denom"` // type of coin to mint
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyInflation, &p.Inflation, validateInflation),
		params.NewParamSetPair(KeyMintDenom, &p.MintDenom, validateMintDenom),
	}
}

// GetParamSpace implements params.ParamStruct
func (p *Params) GetParamSpace() string {
	return DefaultParamSpace
}

// DefaultParams returns default minting module parameters
func DefaultParams() Params {
	return Params{
		Inflation: sdk.NewDecWithPrec(4, 2),
		MintDenom: MintDenom,
	}
}

// Validate returns err if the Params is invalid
func (p Params) Validate() error {
	if p.Inflation.GT(sdk.NewDecWithPrec(2, 1)) || p.Inflation.LT(sdk.ZeroDec()) {
		return sdkerrors.Wrapf(ErrInvalidMintInflation, "Mint inflation [%s] should be between [0, 0.2] ", p.Inflation.String())
	}
	if len(p.MintDenom) == 0 {
		return sdkerrors.Wrapf(ErrInvalidMintDenom, "Mint denom [%s] should not be empty", p.MintDenom)
	}
	return nil
}

func validateInflation(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.GT(sdk.NewDecWithPrec(2, 1)) || v.LT(sdk.ZeroDec()) {
		return fmt.Errorf("Mint inflation [%s] should be between [0, 0.2] ", v.String())
	}

	return nil
}

func validateMintDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("mint denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}
