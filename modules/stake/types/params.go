package types

import (
	"bytes"
	"fmt"
	"time"

	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/params"
	sdk "github.com/irisnet/irishub/types"
	"strconv"
)

var _ params.ParamSet = (*Params)(nil)

const (
	// Default parameter namespace
	DefaultParamSpace = "stake"

	// defaultUnbondingTime reflects three weeks in seconds as the default
	// unbonding time.
	defaultUnbondingTime time.Duration = 60 * 60 * 24 * 3 * time.Second

	// Delay, in blocks, between when validator updates are returned to Tendermint and when they are applied
	// For example, if this is 0, the validator set at the end of a block will sign the next block, or
	// if this is 1, the validator set at the end of a block will sign the block after the next.
	// Constant as this should not change without a hard fork.
	ValidatorUpdateDelay int64 = 1

	// Stake token denomination "iris-atto"
	StakeDenom = sdk.NativeTokenMinDenom

	// Stake token name "iris"
	StakeTokenName = sdk.NativeTokenName
)

// nolint - Keys for parameter access
var (
	KeyUnbondingTime = []byte("UnbondingTime")
	KeyMaxValidators = []byte("MaxValidators")
	KeyBondDenom     = []byte("BondDenom")
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings for staking
type Params struct {
	UnbondingTime time.Duration `json:"unbonding_time"`

	MaxValidators uint16 `json:"max_validators"` // maximum number of validators
	BondDenom     string `json:"bond_denom"`     // bondable coin denomination
}

// Implements params.Params
func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
		{KeyUnbondingTime, &p.UnbondingTime},
		{KeyMaxValidators, &p.MaxValidators},
		{KeyBondDenom, &p.BondDenom},
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
		return maxValidators, nil
	case string(KeyBondDenom):
		bondDenom := string(value)
		if err := validateBondDenom(bondDenom); err != nil {
			return nil, err
		}
		return bondDenom, nil
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
	case string(KeyBondDenom):
		err := cdc.UnmarshalJSON(bytes, &p.BondDenom)
		return p.BondDenom, err
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
		UnbondingTime: defaultUnbondingTime,
		MaxValidators: 100,
		BondDenom:     StakeDenom,
	}
}

func ValidateParams(p Params) error {
	if err := validateUnbondingTime(p.UnbondingTime); err != nil {
		return err
	}
	if err := validateMaxValidators(p.MaxValidators); err != nil {
		return err
	}
	if err := validateBondDenom(p.BondDenom); err != nil {
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
	resp += fmt.Sprintf("Bonded Coin Denomination: %s\n", p.BondDenom)
	return resp
}

//______________________________________________________________________

func validateUnbondingTime(v time.Duration) sdk.Error {
	if v < time.Minute || v > sdk.EightWeeks {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidUnbondingTime, fmt.Sprintf("Invalid UnbondingTime [%d] should be between [10min, 8week]", v))
	}
	return nil
}

func validateMaxValidators(v uint16) sdk.Error {
	if v < 100 || v > 200 {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMaxValidators, fmt.Sprintf("Invalid MaxValidators [%d] should be between [100, 200]", v))
	}
	return nil
}

func validateBondDenom(v string) sdk.Error {
	if len(v) == 0 {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidBondDenom, "staking parameter BondDenom can't be an empty string")
	}
	return nil
}
