package mint

import (
	"fmt"

	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/params"
	sdk "github.com/irisnet/irishub/types"
)

var _ params.ParamSet = (*Params)(nil)

// default paramspace for params keeper
const (
	DefaultParamspace = "mint"
)

//Parameter store key
var (
	// params store for inflation params
	KeyInflation = []byte("Inflation")
)

// ParamTable for mint module
func ParamTypeTable() params.TypeTable {
	return params.NewTypeTable().RegisterParamSet(&Params{})
}

// mint parameters
type Params struct {
	Inflation sdk.Dec `json:"inflation"` // inflation rate
}

// Implements params.ParamStruct
func (p *Params) GetParamsKey() string {
	return DefaultParamspace
}

func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
		{KeyInflation, &p.Inflation},
	}
}

func (p *Params) Validate(key string, value string) (interface{},sdk.Error) {
	switch key {
	case string(KeyInflation):
		inflation, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil , params.ErrInvalidString(value)
		}
		if inflation.GT(sdk.NewDecWithPrec(2,1)) || inflation.LT(sdk.NewDecWithPrec(0,0)){
            return nil, sdk.NewError(params.DefaultCodespace,params.CodeInvalidMintInflation,fmt.Sprintf("Mint Inflation should be between 0 and 20% "))
		}
		return inflation,nil
	default:
		return nil,nil
	}
}

func (p *Params) StringFromBytes(cdc *codec.Codec, key string, bytes []byte) (string, error) {
	switch key {
	case string(KeyInflation):
		cdc.MustUnmarshalJSON(bytes, &p.Inflation)
		return p.Inflation.String(), nil
	default:
		return "", fmt.Errorf("%s is not existed", key)
	}
}

// default minting module parameters
func DefaultParams() Params {
	return Params{
		Inflation: sdk.NewDecWithPrec(4, 2),
	}
}

func validateParams(params Params) error {
	if params.Inflation.LT(sdk.ZeroDec()) {
		return fmt.Errorf("mint parameter Max inflation must be greater than or equal to min inflation")
	}
	return nil
}

//______________________________________________________________________

// get inflation params from the global param store
func (k Keeper) GetParamSet(ctx sdk.Context) Params {
	var params Params
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// set inflation params from the global param store
func (k Keeper) SetParamSet(ctx sdk.Context, params Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
