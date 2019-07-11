package types

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/irisnet/irishub/app/v1/params"
	sdk "github.com/irisnet/irishub/types"
)

const (
	DefaultParamSpace = ModuleName
)

// Parameter store keys
var (
	KeyNativeDenom = []byte("nativeDenom")
	KeyFee         = []byte("fee")
)

// Params defines the fee and native denomination for coinswap
type Params struct {
	NativeDenom string   `json:"native_denom"`
	Fee         FeeParam `json:"fee"`
}

func NewParams(nativeDenom string, fee FeeParam) Params {
	return Params{
		NativeDenom: nativeDenom,
		Fee:         fee,
	}
}

// FeeParam defines the numerator and denominator used in calculating the
// amount to be reserved as a liquidity fee.
// TODO: come up with a more descriptive name than Numerator/Denominator
// Fee = 1 - (Numerator / Denominator) TODO: move this to spec
type FeeParam struct {
	Numerator   sdk.Int `json:"fee_numerator"`
	Denominator sdk.Int `json:"fee_denominator"`
}

func NewFeeParam(numerator, denominator sdk.Int) FeeParam {
	return FeeParam{
		Numerator:   numerator,
		Denominator: denominator,
	}
}

// ParamKeyTable returns the KeyTable for coinswap module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// Implements params.ParamSet.
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyNativeDenom, &p.NativeDenom},
		{KeyFee, &p.Fee},
	}
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  Native Denom:	%s
  Fee:			%s`, p.NativeDenom, p.Fee,
	)
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
		txsize, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if txsize < MinimumTxSizeLimit || txsize > MaximumTxSizeLimit {
			return nil, sdk.NewError(params.DefaultCodespace, params.CodeInvalidTxSizeLimit, fmt.Sprintf("Tx size limit (%s) should be [500, 1500]", value))
		}
		return txsize, nil
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
