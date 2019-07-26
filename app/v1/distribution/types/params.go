package types

import (
	"fmt"

	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
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

func (p Params) String() string {
	return fmt.Sprintf(`Distribution Params:
  distr/CommunityTax:        %s
  distr/BaseProposerReward:  %s
  distr/BonusProposerReward:        %s`,
		p.CommunityTax.String(), p.BaseProposerReward.String(), p.BonusProposerReward.String())
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
		err := cdc.UnmarshalJSON(bytes, &p.CommunityTax)
		return p.CommunityTax.String(), err
	case string(KeyBaseProposerReward):
		err := cdc.UnmarshalJSON(bytes, &p.BaseProposerReward)
		return p.BaseProposerReward.String(), err
	case string(KeyBonusProposerReward):
		err := cdc.UnmarshalJSON(bytes, &p.BonusProposerReward)
		return p.BonusProposerReward.String(), err
	default:
		return "", fmt.Errorf("%s is not existed", key)
	}
}

func (p *Params) ReadOnly() bool {
	return false
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
	if err := validateCommunityTax(p.CommunityTax); err != nil {
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
	if v.LTE(sdk.ZeroDec()) || v.GT(sdk.NewDecWithPrec(2, 1)) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidCommunityTax, fmt.Sprintf("Invalid CommunityTax [%s] should be between (0, 0.2]", v.String()))
	}
	return nil
}

func validateBaseProposerReward(v sdk.Dec) sdk.Error {
	if v.LTE(sdk.ZeroDec()) || v.GT(sdk.NewDecWithPrec(2, 2)) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidBaseProposerReward, fmt.Sprintf("Invalid BaseProposerReward [%s] should be between (0, 0.02]", v.String()))
	}
	return nil
}

func validateBonusProposerReward(v sdk.Dec) sdk.Error {
	if v.LTE(sdk.ZeroDec()) || v.GT(sdk.NewDecWithPrec(8, 2)) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidBonusProposerReward, fmt.Sprintf("Invalid BonusProposerReward [%s] should be between (0, 0.08]", v.String()))
	}
	return nil
}
