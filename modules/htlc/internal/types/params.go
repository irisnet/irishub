package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/params"
)

var _ params.ParamSet = (*Params)(nil)

const (
	DefaultParamSpace = "htlc"
)

// ParamTable for HTLC module
func ParamTypeTable() params.TypeTable {
	return params.NewTypeTable().RegisterParamSet(&Params{})
}

// HTLC params
type Params struct {
}

func (p Params) String() string {
	return ""
}

// Implements params.ParamSet
func (p *Params) GetParamSpace() string {
	return DefaultParamSpace
}

func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{}
}

func (p *Params) Validate(key string, value string) (interface{}, sdk.Error) {
	switch key {
	default:
		return nil, nil
	}
}

func (p *Params) StringFromBytes(cdc *codec.Codec, key string, bytes []byte) (string, error) {
	return "", fmt.Errorf("this method is not implemented")
}

func (p *Params) ReadOnly() bool {
	return false
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
