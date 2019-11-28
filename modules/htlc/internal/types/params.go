package types

import (
	"github.com/cosmos/cosmos-sdk/x/params"
)

var _ params.ParamSet = (*Params)(nil)

const (
	DefaultParamSpace = "htlc"
)

// ParamTable for HTLC module
func ParamTypeTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// HTLC params
type Params struct {
}

func (p Params) String() string {
	return ""
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{}
}

// default HTLC module params
func DefaultParams() Params {
	return Params{}
}

// default HTLC module params for test
func DefaultParamsForTest() Params {
	return Params{}
}

func ValidateParams(p Params) error {
	return nil
}
