package types

import (
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/params"
	"github.com/irisnet/irishub/codec"
	"time"
	"strconv"
	"fmt"
)

const DefaultParamSpace = "distr"

var _ params.ParamSet = (*Params)(nil)

var (
	KeyBaseProposerReward  = []byte("BaseProposerReward")
	KeyBonusProposerReward = []byte("BonusProposerReward")
	KeyCommunityTax        = []byte("CommunityTax")
)

// Params defines the high level settings for distribution
type Params struct {
	CommunityTax        sdk.Dec `json:"community_tax"`
	BaseProposerReward  sdk.Dec `json:"base_proposer_reward"`
	BonusProposerReward sdk.Dec `json:"bonus_proposer_reward"`
}

// Implements params.Params
func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
		{KeyCommunityTax, &p.CommunityTax},
		{KeyBaseProposerReward, &p.BaseProposerReward},
		{KeyBonusProposerReward, &p.BonusProposerReward},
	}
}

func (p *Params) Validate(key string, value string) (interface{}, sdk.Error) {
	switch key {
	case string(KeyCommunityTax):
		communityTax, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateCommunityTax(communityTax); err != nil {
			return nil, err
		}
		return communityTax, nil
	case string(KeyBaseProposerReward):
		baseProposerReward, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateBaseProposerReward(baseProposerReward); err != nil {
			return nil, err
		}
		return baseProposerReward, nil
	case string(KeyBonusProposerReward):
		bonusProposerReward, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateBonusProposerReward(bonusProposerReward); err != nil {
			return nil, err
		}
		return bonusProposerReward, nil
	default:
		return nil, sdk.NewError(params.DefaultCodespace, params.CodeInvalidKey, fmt.Sprintf("%s is not found", key))
	}
}

func (p *Params) GetParamSpace() string {
	return DefaultParamSpace
}

func (p *Params) StringFromBytes(cdc *codec.Codec, key string, bytes []byte) (string, error) {
	switch key {
	case string(KeyCommunityTax):
		var communityTax time.Duration
		err := cdc.UnmarshalJSON(bytes, &communityTax)
		return communityTax.String(), err
	case string(KeyBaseProposerReward):
		var unBondingTime time.Duration
		err := cdc.UnmarshalJSON(bytes, &unBondingTime)
		return unBondingTime.String(), err
	case string(KeyBonusProposerReward):
		var maxValidators uint16
		err := cdc.UnmarshalJSON(bytes, &maxValidators)
		return strconv.Itoa(int(maxValidators)), err
	default:
		return "", fmt.Errorf("%s is not existed", key)
	}
}

// default distribution module params
func DefaultParams() Params {
	return Params{
		CommunityTax:        sdk.NewDecWithPrec(2, 2), // 2%
		BaseProposerReward:  sdk.NewDecWithPrec(1, 2), // 1%
		BonusProposerReward: sdk.NewDecWithPrec(4, 2), // 4%
	}
}

func ValidateParams(p Params) error {
	if err := validateCommunityTax(p.BaseProposerReward); err != nil {
		return err
	}
	if err := validateBaseProposerReward(p.BaseProposerReward); err != nil {
		return err
	}
	if err := validateBonusProposerReward(p.BonusProposerReward); err != nil {
		return err
	}
	return nil
}

//______________________________________________________________________

func validateCommunityTax(v sdk.Dec) sdk.Error {
	if v.LTE(sdk.ZeroDec()) || v.GTE(sdk.NewDec(1)) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidCommunityTax, fmt.Sprintf("Invalid CommunityTax [%s] should be between 0 and 1", v.String()))
	}
	return nil
}

func validateBaseProposerReward(v sdk.Dec) sdk.Error {
	if v.LTE(sdk.ZeroDec()) || v.GTE(sdk.NewDec(1)) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidBaseProposerReward, fmt.Sprintf("Invalid BaseProposerReward [%s] should be between 0 and 1", v.String()))
	}
	return nil
}

func validateBonusProposerReward(v sdk.Dec) sdk.Error {
	if v.LTE(sdk.ZeroDec()) || v.GTE(sdk.NewDec(1)) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidBonusProposerReward, fmt.Sprintf("Invalid BonusProposerReward [%s] should be between 0 and 1", v.String()))
	}
	return nil
}
