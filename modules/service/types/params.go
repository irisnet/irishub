package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"gopkg.in/yaml.v2"
)

// Service params default values
var (
	DefaultMaxRequestTimeout         = int64(100)
	DefaultMinDepositMultiple        = int64(1000)
	DefaultServiceFeeTax             = sdk.NewDecWithPrec(5, 2) // 5%
	DefaultSlashFraction             = sdk.NewDecWithPrec(1, 3) // 0.1%
	DefaultComplaintRetrospect       = 15 * 24 * time.Hour      // 15 days
	DefaultArbitrationTimeLimit      = 5 * 24 * time.Hour       // 5 days
	DefaultTxSizeLimit               = uint64(4000)
	DefaultRestrictedServiceFeeDenom = false
)

// NewParams creates a new Params instance
func NewParams(
	maxRequestTimeout,
	minDepositMultiple int64,
	minDeposit sdk.Coins,
	serviceFeeTax,
	slashFraction sdk.Dec,
	complaintRetrospect,
	arbitrationTimeLimit time.Duration,
	txSizeLimit uint64,
	baseDenom string,
	restrictedServiceFeeDenom bool,
) Params {
	return Params{
		MaxRequestTimeout:         maxRequestTimeout,
		MinDepositMultiple:        minDepositMultiple,
		MinDeposit:                minDeposit,
		ServiceFeeTax:             serviceFeeTax,
		SlashFraction:             slashFraction,
		ComplaintRetrospect:       complaintRetrospect,
		ArbitrationTimeLimit:      arbitrationTimeLimit,
		TxSizeLimit:               txSizeLimit,
		BaseDenom:                 baseDenom,
		RestrictedServiceFeeDenom: restrictedServiceFeeDenom,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		DefaultMaxRequestTimeout,
		DefaultMinDepositMultiple,
		sdk.NewCoins(
			sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(5000)),
		),
		DefaultServiceFeeTax,
		DefaultSlashFraction,
		DefaultComplaintRetrospect,
		DefaultArbitrationTimeLimit,
		DefaultTxSizeLimit,
		sdk.DefaultBondDenom,
		DefaultRestrictedServiceFeeDenom,
	)
}

// String implements stringer
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Validate validates a set of params
func (p Params) Validate() error {
	if err := validateMaxRequestTimeout(p.MaxRequestTimeout); err != nil {
		return err
	}
	if err := validateMinDepositMultiple(p.MinDepositMultiple); err != nil {
		return err
	}
	if err := validateMinDeposit(p.MinDeposit); err != nil {
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
	if err := sdk.ValidateDenom(p.BaseDenom); err != nil {
		return err
	}

	return validateRestrictedServiceFeeDenom(p.RestrictedServiceFeeDenom)
}

func validateMaxRequestTimeout(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("maximum request timeout must be positive: %d", v)
	}

	return nil
}

func validateMinDepositMultiple(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("minimum deposit multiple must be positive: %d", v)
	}

	return nil
}

func validateMinDeposit(i interface{}) error {
	v, ok := i.(sdk.Coins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !v.IsValid() {
		return fmt.Errorf("invalid minimum deposit: %s", v)
	}

	return nil
}

func validateSlashFraction(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.LT(sdk.ZeroDec()) || v.GT(sdk.OneDec()) {
		return fmt.Errorf("slashing fraction must be between [0, 1]: %s", v)
	}

	return nil
}

func validateServiceFeeTax(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.LT(sdk.ZeroDec()) || v.GTE(sdk.OneDec()) {
		return fmt.Errorf("service fee tax must be between [0, 1): %s", v)
	}

	return nil
}

func validateComplaintRetrospect(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("complaint retrospect must be positive: %d", v)
	}

	return nil
}

func validateArbitrationTimeLimit(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("arbitration time limit must be positive: %d", v)
	}

	return nil
}

func validateTxSizeLimit(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("tx size limit must be positive: %d", v)
	}

	return nil
}

func validateBaseDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return sdk.ValidateDenom(v)
}

func validateRestrictedServiceFeeDenom(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
