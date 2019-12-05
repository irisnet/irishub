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
	DefaultServiceFeeTax        = sdk.NewDecWithPrec(1, 2)    //1%
	DefaultSlashFraction        = sdk.NewDecWithPrec(1, 3)    //0.1%
	DefaultComplaintRetrospect  = time.Duration(15 * sdk.Day) //15 days
	DefaultArbitrationTimeLimit = time.Duration(5 * sdk.Day)  //5 days
	DefaultTxSizeLimit          = uint64(4000)
)

// nolint - Keys for parameter access
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
	MaxRequestTimeout    int64         `json:"max_request_timeout"`
	MinDepositMultiple   int64         `json:"min_deposit_multiple"`
	ServiceFeeTax        sdk.Dec       `json:"service_fee_tax"`
	SlashFraction        sdk.Dec       `json:"slash_fraction"`
	ComplaintRetrospect  time.Duration `json:"complaint_retrospect"`
	ArbitrationTimeLimit time.Duration `json:"arbitration_time_limit"`
	TxSizeLimit          uint64        `json:"tx_size_limit"`
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

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyMaxRequestTimeout, Value: &p.MaxRequestTimeout},
		{Key: KeyMinDepositMultiple, Value: &p.MinDepositMultiple},
		{Key: KeyServiceFeeTax, Value: &p.ServiceFeeTax},
		{Key: KeySlashFraction, Value: &p.SlashFraction},
		{Key: KeyComplaintRetrospect, Value: &p.ComplaintRetrospect},
		{Key: KeyArbitrationTimeLimit, Value: &p.ArbitrationTimeLimit},
		{Key: KeyTxSizeLimit, Value: &p.TxSizeLimit},
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
	if err != nil {
		return
	}

	return
}

// Validate validates a set of params
func (p Params) Validate() error {
	return validateParams(p)
}

func validateParams(p Params) error {
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
	if err := validateTxSizeLimit(p.TxSizeLimit); err != nil {
		return err
	}
	return nil
}

func validateMaxRequestTimeout(v int64) sdk.Error {
	if sdk.NetworkType == sdk.Mainnet {
		if v < 20 {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMaxRequestTimeout, fmt.Sprintf("Invalid MaxRequestTimeout [%d] should be greater than or equal to 20", v))
		}
	} else if v < 5 {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMaxRequestTimeout, fmt.Sprintf("Invalid MaxRequestTimeout [%d] should be greater than or equal to 5", v))
	}
	return nil
}

func validateMinDepositMultiple(v int64) sdk.Error {
	if sdk.NetworkType == sdk.Mainnet {
		if v < 500 || v > 5000 {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMinDepositMultiple, fmt.Sprintf("Invalid MinDepositMultiple [%d] should be between [500, 5000]", v))
		}
	} else if v < 10 || v > 5000 {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMinDepositMultiple, fmt.Sprintf("Invalid MinDepositMultiple [%d] should be between [10, 5000]", v))
	}
	return nil
}

func validateSlashFraction(v sdk.Dec) sdk.Error {
	if v.LTE(sdk.ZeroDec()) || v.GT(sdk.NewDecWithPrec(1, 2)) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidSlashFraction, fmt.Sprintf("Invalid SlashFraction [%s] should be between (0, 0.01]", v.String()))
	}
	return nil
}

func validateServiceFeeTax(v sdk.Dec) sdk.Error {
	if v.LTE(sdk.ZeroDec()) || v.GT(sdk.NewDecWithPrec(2, 1)) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidServiceFeeTax, fmt.Sprintf("Invalid ServiceFeeTax [%s] should be between (0, 0.2]", v.String()))
	}
	return nil
}

func validateComplaintRetrospect(v time.Duration) sdk.Error {
	if sdk.NetworkType == sdk.Mainnet {
		if v < 15*sdk.Day || v > 30*sdk.Day {
			return sdk.NewError(params.DefaultCodespace, params.CodeComplaintRetrospect, fmt.Sprintf("Invalid ComplaintRetrospect [%s] should be between [15days, 30days]", v.String()))
		}
	} else if v < 20*time.Second {
		return sdk.NewError(params.DefaultCodespace, params.CodeComplaintRetrospect, fmt.Sprintf("Invalid ComplaintRetrospect [%s] should be between [20seconds, )", v.String()))
	}
	return nil
}

func validateArbitrationTimeLimit(v time.Duration) sdk.Error {
	if sdk.NetworkType == sdk.Mainnet {
		if v < 5*sdk.Day || v > 10*sdk.Day {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidArbitrationTimeLimit, fmt.Sprintf("Invalid ArbitrationTimeLimit [%s] should be between [5days, 10days]", v.String()))
		}
	} else if v < 20*time.Second {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidArbitrationTimeLimit, fmt.Sprintf("Invalid ArbitrationTimeLimit [%s] should be between [20seconds, )", v.String()))
	}
	return nil
}

func validateTxSizeLimit(v uint64) sdk.Error {
	if v < uint64(2000) || v > uint64(6000) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidServiceTxSizeLimit, fmt.Sprintf("Invalid ServiceTxSizeLimit [%d] should be between [2000, 6000]", v))
	}
	return nil
}
