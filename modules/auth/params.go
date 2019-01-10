package auth

import (
	"fmt"

	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/params"
	sdk "github.com/irisnet/irishub/types"
)

var _ params.ParamSet = (*Params)(nil)

const (
	DefaultParamSpace = "auth"
)
var (
	MinimumGasPrice = sdk.ZeroInt()
	MaximumGasPrice = sdk.NewIntWithDecimal(1, 18) //1iris, 10^18iris-atto
)

//Parameter store key
var (
	// params store for inflation params
	gasPriceThresholdKey = []byte("gasPriceThreshold")
)

// ParamTable for auth module
func ParamTypeTable() params.TypeTable {
	return params.NewTypeTable().RegisterParamSet(&Params{})
}

// auth parameters
type Params struct {
	GasPriceThreshold sdk.Int `json:"gas_price_threshold"` // gas price threshold
}

// Implements params.ParamStruct
func (p *Params) GetParamSpace() string {
	return DefaultParamSpace
}

func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
		{gasPriceThresholdKey, &p.GasPriceThreshold},
	}
}

func (p *Params) Validate(key string, value string) (interface{}, sdk.Error) {
	switch key {
	case string(gasPriceThresholdKey):
		threshold, ok := sdk.NewIntFromString(value)
		if !ok {
			return nil, params.ErrInvalidString(value)
		}
		if !threshold.GT(MinimumGasPrice) || threshold.GT(MaximumGasPrice) {
			return nil, sdk.NewError(params.DefaultCodespace, params.CodeInvalidGasPriceThreshold, fmt.Sprintf("Gas price threshold (%s) should be (0, 10^18iris-atto]", value))
		}
		return threshold, nil
	default:
		return nil, sdk.NewError(params.DefaultCodespace, params.CodeInvalidKey, fmt.Sprintf("%s is not found", key))
	}
}

func (p *Params) StringFromBytes(cdc *codec.Codec, key string, bytes []byte) (string, error) {
	switch key {
	case string(gasPriceThresholdKey):
		err := cdc.UnmarshalJSON(bytes, &p.GasPriceThreshold)
		return p.GasPriceThreshold.String(), err
	default:
		return "", fmt.Errorf("%s is not existed", key)
	}
}

// default auth module parameters
func DefaultParams() Params {
	return Params{
		GasPriceThreshold: sdk.NewIntWithDecimal(2, 10), //20iris-nano, 2*10^10iris-atto
	}
}

func validateParams(p Params) error {
	if sdk.NetworkType != sdk.Mainnet {
		return nil
	}

	if !p.GasPriceThreshold.GT(MinimumGasPrice) || p.GasPriceThreshold.GT(MaximumGasPrice) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidGasPriceThreshold, fmt.Sprintf("Gas price threshold (%s) should be (0, 10^18iris-atto]", p.GasPriceThreshold.String()))
	}
	return nil
}

//______________________________________________________________________


