package rand

import (
	"fmt"

	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

var _ params.ParamSet = (*Params)(nil)

const (
	DefaultParamSpace = "rand"
)

// ParamTable for rand module
func ParamTypeTable() params.TypeTable {
	return params.NewTypeTable().RegisterParamSet(&Params{})
}

// rand params
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

// default rand module params
func DefaultParams() Params {
	return Params{}
}

// default rand module params for test
func DefaultParamsForTest() Params {
	return Params{}
}

func validateParams(p Params) error {
	return nil
}

// get rand params from the global param store
func (k Keeper) GetParamSet(ctx sdk.Context) Params {
	var p Params
	k.paramSpace.GetParamSet(ctx, &p)
	return p
}

// set asset params from the global param store
func (k Keeper) SetParamSet(ctx sdk.Context, params Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
