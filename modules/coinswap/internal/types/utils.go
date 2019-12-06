package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/config"
)

const (
	MinDenomSuffix = "-min"
)

// GetUniId returns the unique uni id for the provided denominations.
// The uni id is in the format of 'u-coin-name' which the denomination
// is not iris-atto.
func GetUniId(denom1, denom2 string) (string, sdk.Error) {
	if denom1 == denom2 {
		return "", ErrEqualDenom("denomnations for forming uni id are equal")
	}

	if denom1 != config.IrisAtto && denom2 != config.IrisAtto {
		return "", ErrIllegalDenom(fmt.Sprintf("illegal denomnations for forming uni id, must have one native denom: %s", config.IrisAtto))
	}

	denom := denom1
	if denom == config.IrisAtto {
		denom = denom2
	}
	coinName, err := GetCoinNameByDenom(denom)
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

// CheckUniDenom returns nil if the uni denom is valid
func CheckUniDenom(uniDenom string) sdk.Error {
	if _, isValid := sdk.GetDenomUnit(uniDenom); !isValid || !strings.HasPrefix(uniDenom, FormatUniABSPrefix) {
		return ErrIllegalDenom(fmt.Sprintf("illegal liquidity denomnation: %s", uniDenom))
	}
	return nil
}

// CheckUniId returns nil if the uni id is valid
func CheckUniId(uniId string) sdk.Error {
	if _, isValid := sdk.GetDenomUnit(uniId); !isValid || !strings.HasPrefix(uniId, FormatUniABSPrefix) {
		return ErrIllegalUniId(fmt.Sprintf("illegal liquidity id: %s", uniId))
	}
	return nil
}

// GetUniDenom returns uni denom if the uni id is valid
func GetUniDenom(uniId string) (string, sdk.Error) {
	if err := CheckUniId(uniId); err != nil {
		return "", err
	}

	uniDenom, err := GetCoinMinDenom(uniId)
	if err != nil {
		return "", ErrIllegalUniId(fmt.Sprintf("illegal liquidity id: %s", uniId))
	}
	return uniDenom, nil
}

func GetCoinNameByDenom(denom string) (coinName string, err error) {
	denom = strings.ToLower(denom)
	if strings.HasPrefix(denom, config.Iris+"-") {
		if _, isValid := sdk.GetDenomUnit(denom); !isValid {
			return "", fmt.Errorf("invalid denom for getting coin name: %s", denom)
		}
		return config.Iris, nil
	}
	if _, isValid := sdk.GetDenomUnit(denom); !isValid {
		return "", fmt.Errorf("invalid denom for getting coin name: %s", denom)
	}
	coinName = strings.TrimSuffix(denom, MinDenomSuffix)
	if coinName == "" {
		return coinName, fmt.Errorf("coin name is empty")
	}
	return coinName, nil
}

func GetCoinMinDenom(coinName string) (denom string, err error) {
	coinName = strings.ToLower(strings.TrimSpace(coinName))

	if coinName == config.Iris {
		return config.IrisAtto, nil
	}

	return fmt.Sprintf("%s%s", coinName, MinDenomSuffix), nil
}
