package types

import (
	"fmt"
	sdk "github.com/irisnet/irishub/types"
	"strings"
)

// GetUniId returns the unique uni id for the provided denominations.
// The uni id is in the format of 'u-coin-name' which the denomination
// is not iris-atto.
func GetUniId(denom1, denom2 string) (string, sdk.Error) {
	if denom1 == denom2 {
		return "", ErrEqualDenom("denomnations for forming uni id are equal")
	}

	if denom1 != sdk.IrisAtto && denom2 != sdk.IrisAtto {
		return "", ErrIllegalDenom(fmt.Sprintf("illegal denomnations for forming uni id, must have one native denom: %s", sdk.IrisAtto))
	}

	denom := denom1
	if denom == sdk.IrisAtto {
		denom = denom2
	}
	coinName, err := sdk.GetCoinNameByDenom(denom)
	if err != nil {
		return "", ErrIllegalDenom(err.Error())
	}
	return fmt.Sprintf(FormatUniId, coinName), nil
}

// GetCoinMinDenomFromUniDenom returns the token denom by uni denom
func GetCoinMinDenomFromUniDenom(uniDenom string) (string, sdk.Error) {
	err := CheckUniDenom(uniDenom)
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(uniDenom, FormatUniABSPrefix), nil
}

// GetUniCoinType returns the uni coin type
func GetUniCoinType(uniId string) (sdk.CoinType, sdk.Error) {
	uniDenom, err := GetUniDenom(uniId)
	if err != nil {
		return sdk.CoinType{}, err
	}
	units := make(sdk.Units, 2)
	units[0] = sdk.NewUnit(uniId, 0)
	units[1] = sdk.NewUnit(uniDenom, sdk.AttoScale) // the uni denom has the same decimal with iris-atto
	return sdk.CoinType{
		Name:    uniId,
		MinUnit: units[1],
		Units:   units,
	}, nil
}

// CheckUniDenom returns nil if the uni denom is valid
func CheckUniDenom(uniDenom string) sdk.Error {
	if !sdk.IsCoinMinDenomValid(uniDenom) || !strings.HasPrefix(uniDenom, FormatUniABSPrefix) {
		return ErrIllegalDenom(fmt.Sprintf("illegal liquidity denomnation: %s", uniDenom))
	}
	return nil
}

// CheckUniId returns nil if the uni id is valid
func CheckUniId(uniId string) sdk.Error {
	if !sdk.IsCoinNameValid(uniId) || !strings.HasPrefix(uniId, FormatUniABSPrefix) {
		return ErrIllegalUniId(fmt.Sprintf("illegal liquidity id: %s", uniId))
	}
	return nil
}

// GetUniDenom returns uni denom if the uni id is valid
func GetUniDenom(uniId string) (string, sdk.Error) {
	if err := CheckUniId(uniId); err != nil {
		return "", err
	}

	uniDenom, err := sdk.GetCoinMinDenom(uniId)
	if err != nil {
		return "", ErrIllegalUniId(fmt.Sprintf("illegal liquidity id: %s", uniId))
	}
	return uniDenom, nil
}
