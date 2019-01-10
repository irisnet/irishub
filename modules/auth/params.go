package auth

import (
	"fmt"

	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/params"
	sdk "github.com/irisnet/irishub/types"
	"strconv"
)

var _ params.ParamSet = (*Params)(nil)

const (
	DefaultParamSpace = "auth"
)

var (
	MinimumGasPrice    = sdk.ZeroInt()
	MaximumGasPrice    = sdk.NewIntWithDecimal(1, 18) //1iris, 10^18iris-atto
	MinimumTxSizeLimit = uint32(500)
	MaximumTxSizeLimit = uint32(1500)
)

//Parameter store key
var (
	// params store for inflation params
	gasPriceThresholdKey = []byte("gasPriceThreshold")
	TxSizeLimitKey       = []byte("txSizeLimit")
)

// ParamTable for auth module
func ParamTypeTable() params.TypeTable {
	return params.NewTypeTable().RegisterParamSet(&Params{})
}

// auth parameters
type Params struct {
	GasPriceThreshold sdk.Int `json:"gas_price_threshold"` // gas price threshold
	TxSizeLimit       uint32  `json:"tx_size"`             // tx size limit
}

// Implements params.ParamStruct
func (p *Params) GetParamSpace() string {
	return DefaultParamSpace
}

func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
		{gasPriceThresholdKey, &p.GasPriceThreshold},
		{TxSizeLimitKey, &p.TxSizeLimit},
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
	case string(TxSizeLimitKey):
		txsize, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if uint32(txsize) < MinimumTxSizeLimit || uint32(txsize) > MaximumTxSizeLimit {
			return nil, sdk.NewError(params.DefaultCodespace, params.CodeInvalidTxSizeLimit, fmt.Sprintf("Tx size limit (%s) should be [500, 1500]", value))
		}
		return uint32(txsize), nil
	default:
		return nil, sdk.NewError(params.DefaultCodespace, params.CodeInvalidKey, fmt.Sprintf("%s is not found", key))
	}
}

func (p *Params) StringFromBytes(cdc *codec.Codec, key string, bytes []byte) (string, error) {
	switch key {
	case string(gasPriceThresholdKey):
		err := cdc.UnmarshalJSON(bytes, &p.GasPriceThreshold)
		return p.GasPriceThreshold.String(), err
	case string(TxSizeLimitKey):
		err := cdc.UnmarshalJSON(bytes, &p.TxSizeLimit)
		return strconv.FormatUint(uint64(p.TxSizeLimit), 10), err
	default:
		return "", fmt.Errorf("%s is not existed", key)
	}
}

// default auth module parameters
func DefaultParams() Params {
	return Params{
		GasPriceThreshold: sdk.NewIntWithDecimal(2, 10), //20iris-nano, 2*10^10iris-atto
		TxSizeLimit:       1000,
	}
}

func validateParams(p Params) error {
	if !p.GasPriceThreshold.GT(MinimumGasPrice) || p.GasPriceThreshold.GT(MaximumGasPrice) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidGasPriceThreshold, fmt.Sprintf("Gas price threshold (%s) should be (0, 10^18iris-atto]", p.GasPriceThreshold.String()))
	}
	if p.TxSizeLimit < MinimumTxSizeLimit || p.TxSizeLimit > MaximumTxSizeLimit {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidTxSizeLimit, fmt.Sprintf("Tx size limit (%s) should be [500, 1500]", strconv.FormatUint(uint64(p.TxSizeLimit), 10)))
	}
	return nil
}

//______________________________________________________________________
