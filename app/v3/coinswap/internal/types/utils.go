package types

import (
	"fmt"
	"strings"

	sdk "github.com/irisnet/irishub/types"
)

// GetVoucherCoinName returns the coin name of the voucher for the provided denominations.
// The voucher coin name is in the format of LiquidityVoucherPrefix+token.Symbol, of which the denomination
// is not iris-atto.
func GetVoucherCoinName(denom1, denom2 string) (string, sdk.Error) {
	if denom1 == denom2 {
		return "", ErrEqualDenom("denominations for generating the voucher coin name are equal")
	}

	if denom1 != sdk.IrisAtto && denom2 != sdk.IrisAtto {
		return "", ErrIllegalDenom(fmt.Sprintf("illegal denomnations for generating the voucher coin name, must have one native denom: %s", sdk.IrisAtto))
	}

	denom := denom1
	if denom == sdk.IrisAtto {
		denom = denom2
	}
	coinName, err := sdk.GetCoinNameByDenom(denom)
	if err != nil {
		return "", ErrIllegalDenom(err.Error())
	}

	return fmt.Sprintf(FormatVoucherCoinName, coinName), nil
}

// GetUnderlyingDenom returns the denom of the original token from which the voucher is generated
func GetUnderlyingDenom(voucherDenom string) (string, sdk.Error) {
	err := CheckVoucherDenom(voucherDenom)
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(voucherDenom, LiquidityVoucherPrefix), nil
}

// GetVoucherCoinType returns the voucher coin type
func GetVoucherCoinType(voucherCoinName string) (sdk.CoinType, sdk.Error) {
	voucherDenom, err := GetVoucherDenom(voucherCoinName)
	if err != nil {
		return sdk.CoinType{}, err
	}
	units := make(sdk.Units, 2)
	units[0] = sdk.NewUnit(voucherCoinName, 0)
	units[1] = sdk.NewUnit(voucherDenom, sdk.AttoScale) // the voucher denom has the same decimal with iris-atto
	return sdk.CoinType{
		Name:    voucherCoinName,
		MinUnit: units[1],
		Units:   units,
	}, nil
}

// CheckVoucherDenom returns nil if the voucher denom is valid
func CheckVoucherDenom(voucherDenom string) sdk.Error {
	if !sdk.IsValidCoinDenom(voucherDenom) || !strings.HasPrefix(voucherDenom, LiquidityVoucherPrefix) {
		return ErrIllegalDenom(fmt.Sprintf("illegal voucher denomination: %s", voucherDenom))
	}
	return nil
}

// CheckVoucherCoinName returns nil if the voucher coin name is valid
func CheckVoucherCoinName(voucherCoinName string) sdk.Error {
	if !sdk.IsValidCoinName(voucherCoinName) || !strings.HasPrefix(voucherCoinName, LiquidityVoucherPrefix) {
		return ErrIllegalVoucherCoinName(fmt.Sprintf("illegal voucher coin name: %s", voucherCoinName))
	}
	return nil
}

// GetVoucherDenom returns the voucher denom if the voucher coin name is valid
func GetVoucherDenom(voucherCoinName string) (string, sdk.Error) {
	if err := CheckVoucherCoinName(voucherCoinName); err != nil {
		return "", err
	}

	voucherDenom, err := sdk.GetCoinDenom(voucherCoinName)
	if err != nil {
		return "", ErrIllegalVoucherCoinName(fmt.Sprintf("illegal voucher coin name: %s", voucherCoinName))
	}

	return voucherDenom, nil
}
