package types

import (
	"fmt"
	"strconv"

	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

var _ params.ParamSet = (*Params)(nil)

// default paramSpace for oracle keeper
const (
	DefaultParamSpace = "oracle"
)

//Parameter store key
var (
	// params store for oracle params
	KeyMaxHistory = []byte("MaxHistory")
)

// ParamTable for oracle module
func ParamTypeTable() params.TypeTable {
	return params.NewTypeTable().RegisterParamSet(&Params{})
}

// oracle params
type Params struct {
	MaxHistory uint64 `json:"max_history"`
}

func (p Params) String() string {
	return fmt.Sprintf(`Oracle Params:
  oracle/MaxHistory:     %d`,
		p.MaxHistory)
}

// Implements params.ParamStruct
func (p *Params) GetParamSpace() string {
	return DefaultParamSpace
}

func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
		{KeyMaxHistory, &p.MaxHistory},
	}
}

func (p *Params) Validate(key string, value string) (interface{}, sdk.Error) {
	switch key {
	case string(KeyMaxHistory):
		maxHistory, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateMaxHistory(maxHistory); err != nil {
			return nil, err
		}
		return maxHistory, nil
	default:
		return nil, sdk.NewError(params.DefaultCodespace, params.CodeInvalidKey, fmt.Sprintf("%s is not found", key))
	}
}

func (p *Params) StringFromBytes(cdc *codec.Codec, key string, bytes []byte) (string, error) {
	switch key {
	case string(KeyMaxHistory):
		err := cdc.UnmarshalJSON(bytes, &p.MaxHistory)
		return strconv.FormatUint(p.MaxHistory, 10), err
	default:
		return "", fmt.Errorf("%s is not existed", key)
	}
}

func (p *Params) ReadOnly() bool {
	return false
}

// default oracle module params
func DefaultParams() Params {
	return Params{
		MaxHistory: MaxHistory,
	}
}

// default oracle module params for test
func DefaultParamsForTest() Params {
	return Params{
		MaxHistory: MaxHistory,
	}
}

func validateParams(p Params) error {
	return validateMaxHistory(p.MaxHistory)
}

func validateMaxHistory(v uint64) sdk.Error {
	if v < 1 {
		return ErrInvalidMaxHistory(DefaultCodespace)
	}
	return nil
}
