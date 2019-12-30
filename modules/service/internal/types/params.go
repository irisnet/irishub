package types

import (
	"bytes"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Service params default values
var (
	DefaultMaxRequestTimeout    = int64(100)
	DefaultMinDepositMultiple   = int64(1000)
	DefaultServiceFeeTax        = sdk.NewDecWithPrec(1, 2) // 1%
	DefaultSlashFraction        = sdk.NewDecWithPrec(1, 3) // 0.1%
	DefaultComplaintRetrospect  = 15 * 24 * time.Hour      // 15 days
	DefaultArbitrationTimeLimit = 5 * 24 * time.Hour       // 5 days
	DefaultTxSizeLimit          = uint64(4000)
)

// no lint
var (
	MinRequestTimeout       = int64(2)
	MinDepositMultiple      = int64(500)
	MaxDepositMultiple      = int64(5000)
	MaxServiceFeeTax        = sdk.NewDecWithPrec(2, 1)
	MaxSlashFraction        = sdk.NewDecWithPrec(1, 2)
	MinComplaintRetrospect  = 15 * 24 * time.Hour
	MaxComplaintRetrospect  = 30 * 24 * time.Hour
	MinArbitrationTimeLimit = 5 * 24 * time.Hour
	MaxArbitrationTimeLimit = 10 * 24 * time.Hour
	MinTxSizeLimit          = uint64(2000)
	MaxTxSizeLimit          = uint64(6000)
)

// Keys for parameter access
// nolint
var (
	KeyMaxRequestTimeout    = []byte("MaxRequestTimeout")
	KeyMinDepositMultiple   = []byte("MinDepositMultiple")
	KeyServiceFeeTax        = []byte("ServiceFeeTax")
	KeySlashFraction        = []byte("SlashFraction")
	KeyComplaintRetrospect  = []byte("ComplaintRetrospect")
	KeyArbitrationTimeLimit = []byte("ArbitrationTimeLimit")
	KeyTxSizeLimit          = []byte("TxSizeLimit")
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings for service
type Params struct {
	MaxRequestTimeout    int64         `json:"max_request_timeout" yaml:"max_request_timeout"`
	MinDepositMultiple   int64         `json:"min_deposit_multiple" yaml:"min_deposit_multiple"`
	ServiceFeeTax        sdk.Dec       `json:"service_fee_tax" yaml:"service_fee_tax"`
	SlashFraction        sdk.Dec       `json:"slash_fraction" yaml:"slash_fraction"`
	ComplaintRetrospect  time.Duration `json:"complaint_retrospect" yaml:"complaint_retrospect"`
	ArbitrationTimeLimit time.Duration `json:"arbitration_time_limit" yaml:"arbitration_time_limit"`
	TxSizeLimit          uint64        `json:"tx_size_limit" yaml:"tx_size_limit"`
}

// NewParams creates a new Params instance
func NewParams(maxRequestTimeout, minDepositMultiple int64, serviceFeeTax, slashFraction sdk.Dec,
	complaintRetrospect, arbitrationTimeLimit time.Duration, txSizeLimit uint64) Params {

	return Params{
		MaxRequestTimeout:    maxRequestTimeout,
		MinDepositMultiple:   minDepositMultiple,
		ServiceFeeTax:        serviceFeeTax,
		SlashFraction:        slashFraction,
		ComplaintRetrospect:  complaintRetrospect,
		ArbitrationTimeLimit: arbitrationTimeLimit,
		TxSizeLimit:          txSizeLimit,
	}
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyMaxRequestTimeout, &p.MaxRequestTimeout, validateMaxRequestTimeout),
		params.NewParamSetPair(KeyMinDepositMultiple, &p.MinDepositMultiple, validateMinDepositMultiple),
		params.NewParamSetPair(KeyServiceFeeTax, &p.ServiceFeeTax, validateServiceFeeTax),
		params.NewParamSetPair(KeySlashFraction, &p.SlashFraction, validateSlashFraction),
		params.NewParamSetPair(KeyComplaintRetrospect, &p.ComplaintRetrospect, validateComplaintRetrospect),
		params.NewParamSetPair(KeyArbitrationTimeLimit, &p.ArbitrationTimeLimit, validateArbitrationTimeLimit),
		params.NewParamSetPair(KeyTxSizeLimit, &p.TxSizeLimit, validateTxSizeLimit),
	}
}

// Equal returns a boolean determining if two Param types are identical.
// TODO: This is slower than comparing struct fields directly
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		DefaultMaxRequestTimeout,
		DefaultMinDepositMultiple,
		DefaultServiceFeeTax,
		DefaultSlashFraction,
		DefaultComplaintRetrospect,
		DefaultArbitrationTimeLimit,
		DefaultTxSizeLimit,
	)
}

// String implements stringer
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  Max Request Timeout:     %d
  Min Deposit Multiple:    %d
  Service Fee Tax:         %s
  Slash Fraction:          %s
  Complaint Retrospect:    %s
  Arbitration Time Limit:  %s
  Tx Size Limit:           %d`,
		p.MaxRequestTimeout, p.MinDepositMultiple, p.ServiceFeeTax.String(), p.SlashFraction.String(),
		p.ComplaintRetrospect, p.ArbitrationTimeLimit, p.TxSizeLimit)
}

// MustUnmarshalParams unmarshals the current service params value from store key or panic
func MustUnmarshalParams(cdc *codec.Codec, value []byte) Params {
	params, err := UnmarshalParams(cdc, value)
	if err != nil {
		panic(err)
	}
	return params
}

// UnmarshalParams unmarshals the current service params value from store key
func UnmarshalParams(cdc *codec.Codec, value []byte) (params Params, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &params)
	return
}

// Validate validates a set of params
func (p Params) Validate() error {
	if err := validateMaxRequestTimeout(p.MaxRequestTimeout); err != nil {
		return err
	}
	if err := validateMinDepositMultiple(p.MinDepositMultiple); err != nil {
		return err
	}
	if err := validateSlashFraction(p.SlashFraction); err != nil {
		return err
	}
	if err := validateServiceFeeTax(p.ServiceFeeTax); err != nil {
		return err
	}
	if err := validateComplaintRetrospect(p.ComplaintRetrospect); err != nil {
		return err
	}
	if err := validateArbitrationTimeLimit(p.ArbitrationTimeLimit); err != nil {
		return err
	}
	return validateTxSizeLimit(p.TxSizeLimit)
}

func validateMaxRequestTimeout(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < MinRequestTimeout {
		return fmt.Errorf("MaxRequestTimeout [%d] should be greater than or equal to %d", v, MinRequestTimeout)
	}
	return nil
}

func validateMinDepositMultiple(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < MinDepositMultiple || v > MaxDepositMultiple {
		return fmt.Errorf("MinDepositMultiple [%d] should be between [%d, %d]", v, MinDepositMultiple, MaxDepositMultiple)
	}
	return nil
}

func validateSlashFraction(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.LTE(sdk.ZeroDec()) || v.GT(MaxSlashFraction) {
		return fmt.Errorf("SlashFraction [%s] should be between (0, %s]", v, MaxSlashFraction)
	}
	return nil
}

func validateServiceFeeTax(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.LTE(sdk.ZeroDec()) || v.GT(MaxServiceFeeTax) {
		return fmt.Errorf("ServiceFeeTax [%s] should be between (0, %s]", v, MaxServiceFeeTax)
	}
	return nil
}

func validateComplaintRetrospect(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < MinComplaintRetrospect || v > MaxComplaintRetrospect {
		return fmt.Errorf("ComplaintRetrospect [%s] should be between [%s, %s]", v, MinComplaintRetrospect, MaxComplaintRetrospect)
	}
	return nil
}

func validateArbitrationTimeLimit(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < MinArbitrationTimeLimit || v > MaxArbitrationTimeLimit {
		return fmt.Errorf("ArbitrationTimeLimit [%s] should be between [%s, %s]", v, MinArbitrationTimeLimit, MaxArbitrationTimeLimit)
	}
	return nil
}

func validateTxSizeLimit(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < MinTxSizeLimit || v > MaxTxSizeLimit {
		return fmt.Errorf("TxSizeLimit [%d] should be between [%d, %d]", v, MinTxSizeLimit, MaxTxSizeLimit)
	}
	return nil
}
