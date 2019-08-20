package types

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/params"
	sdk "github.com/irisnet/irishub/types"
)

var _ params.ParamSet = (*Params)(nil)

const (
	// Default parameter namespace
	DefaultParamSpace = "stake"

	// Delay, in blocks, between when validator updates are returned to Tendermint and when they are applied
	// For example, if this is 0, the validator set at the end of a block will sign the next block, or
	// if this is 1, the validator set at the end of a block will sign the block after the next.
	// Constant as this should not change without a hard fork.
	ValidatorUpdateDelay int64 = 1

	// Stake token denomination "iris-atto"
	StakeDenom = sdk.IrisAtto
)

// nolint - Keys for parameter access
var (
	KeyUnbondingTime = []byte("UnbondingTime")
	KeyMaxValidators = []byte("MaxValidators")
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings for staking
type Params struct {
	UnbondingTime time.Duration `json:"unbonding_time"`
	MaxValidators uint16        `json:"max_validators"` // maximum number of validators
}

func (p Params) String() string {
	return fmt.Sprintf(`Stake Params:
  Unbonding Time:         %s
  Max Validators:         %d`,
		p.UnbondingTime, p.MaxValidators)
}

// Implements params.Params
func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
		{KeyUnbondingTime, &p.UnbondingTime},
		{KeyMaxValidators, &p.MaxValidators},
	}
}

func (p *Params) Validate(key string, value string) (interface{}, sdk.Error) {
	switch key {
	case string(KeyUnbondingTime):
		unbondingTime, err := time.ParseDuration(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateUnbondingTime(unbondingTime); err != nil {
			return nil, err
		}
		return unbondingTime, nil
	case string(KeyMaxValidators):
		maxValidators, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateMaxValidators(uint16(maxValidators)); err != nil {
			return nil, err
		}
		return uint16(maxValidators), nil
	default:
		return nil, sdk.NewError(params.DefaultCodespace, params.CodeInvalidKey, fmt.Sprintf("%s is not found", key))
	}
}

func (p *Params) GetParamSpace() string {
	return DefaultParamSpace
}

func (p *Params) StringFromBytes(cdc *codec.Codec, key string, bytes []byte) (string, error) {
	switch key {
	case string(KeyUnbondingTime):
		err := cdc.UnmarshalJSON(bytes, &p.UnbondingTime)
		return p.UnbondingTime.String(), err
	case string(KeyMaxValidators):
		err := cdc.UnmarshalJSON(bytes, &p.MaxValidators)
		return strconv.Itoa(int(p.MaxValidators)), err
	default:
		return "", fmt.Errorf("%s is not existed", key)
	}
}

// Equal returns a boolean determining if two Param types are identical.
func (p Params) Equal(p2 Params) bool {
	bz1 := MsgCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := MsgCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// default stake module params
func DefaultParams() Params {
	return Params{
		UnbondingTime: 3 * sdk.Week,
		MaxValidators: 100,
	}
}

func ValidateParams(p Params) error {
	if err := validateUnbondingTime(p.UnbondingTime); err != nil {
		return err
	}
	if err := validateMaxValidators(p.MaxValidators); err != nil {
		return err
	}
	return nil
}

// HumanReadableString returns a human readable string representation of the
// parameters.
func (p Params) HumanReadableString() string {
	resp := "Params \n"
	resp += fmt.Sprintf("Unbonding Time: %s\n", p.UnbondingTime)
	resp += fmt.Sprintf("Max Validators: %d: \n", p.MaxValidators)
	return resp
}

//______________________________________________________________________

func validateUnbondingTime(v time.Duration) sdk.Error {
	if sdk.NetworkType == sdk.Mainnet {
		if v < 2*sdk.Week {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidUnbondingTime, fmt.Sprintf("Invalid UnbondingTime [%s] should be greater than or equal to 2 weeks", v.String()))
		}
	} else if v < 2*time.Minute {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidUnbondingTime, fmt.Sprintf("Invalid UnbondingTime [%s] should be greater than or equal to 2 minutes", v.String()))
	}
	return nil
}

func validateMaxValidators(v uint16) sdk.Error {
	if sdk.NetworkType == sdk.Mainnet {
		if v < 100 || v > 200 {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMaxValidators, fmt.Sprintf("Invalid MaxValidators [%d] should be between [100, 200]", v))
		}
	} else if v == 0 || v > 200 {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMaxValidators, fmt.Sprintf("Invalid MaxValidators [%d] should be between [1, 200]", v))
	}
	return nil
}
