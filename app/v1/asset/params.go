package asset

import (
	"fmt"
	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

var _ params.ParamSet = (*Params)(nil)

// default paramSpace for asset keeper
const (
	DefaultParamSpace = "asset"
)

//Parameter store key
var (
	// params store for asset params

	// TODO
)

// ParamTable for asset module
func ParamTypeTable() params.TypeTable {
	return params.NewTypeTable().RegisterParamSet(&Params{})
}

// asset params
type Params struct {
	// TODO
}

func (p Params) String() string {
	// TODO
	return fmt.Sprintf(`Asset Params:`)
}

// Implements params.ParamStruct
func (p *Params) GetParamSpace() string {
	return DefaultParamSpace
}

func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
		// TODO
	}
}

func (p *Params) Validate(key string, value string) (interface{}, sdk.Error) {
	switch key {
	// TODO
	default:
		return nil, sdk.NewError(params.DefaultCodespace, params.CodeInvalidKey, fmt.Sprintf("%s is not found", key))
	}
}

func (p *Params) StringFromBytes(cdc *codec.Codec, key string, bytes []byte) (string, error) {
	switch key {
	// TODO
	default:
		return "", fmt.Errorf("%s is not existed", key)
	}
}

// default asset module params
func DefaultParams() Params {
	return Params{
		// TODO
	}
}

// default asset module params for test
func DefaultParamsForTest() Params {
	return Params{
		// TODO
	}
}

func validateParams(p Params) error {
	if sdk.NetworkType != sdk.Mainnet {
		return nil
	}

	// TODO
	return nil
}

//______________________________________________________________________

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

//______________________________________________________________________
