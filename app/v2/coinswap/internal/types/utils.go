package types

import (
	"fmt"
	sdk "github.com/irisnet/irishub/types"
	"strings"
)

// GetReservePoolName returns the reserve pool name for the provided denominations.
// The reserve pool name is in the format of 'u-denom' which the denomination
// is not iris-atto.
func GetReservePoolName(denom1, denom2 string) (string, sdk.Error) {
	if denom1 == denom2 {
		return "", ErrEqualDenom("denomnations for forming reserve pool name are equal")
	}

	if denom1 != sdk.IrisAtto && denom2 != sdk.IrisAtto {
		return "", ErrIllegalDenom(fmt.Sprintf("illegal denomnations for forming reserve pool name, must have one native denom: %s", sdk.IrisAtto))
	}

	var denom = denom2
	if denom1 != sdk.IrisAtto {
		denom = denom1
	}
	return fmt.Sprintf(FormatReservePool, denom), nil
}

// GetTokenDenom returns the token denom by uni denom
func GetTokenDenom(uniDenom string) (string, sdk.Error) {
	CheckUniDenom(uniDenom)
	return strings.TrimPrefix(uniDenom, FormatReservePoolPrefix), nil
}

// CheckUniDenom returns nil if the uni denom is valid
func CheckUniDenom(uniDenom string) sdk.Error {
	if !strings.HasPrefix(uniDenom, FormatReservePoolPrefix) {
		return ErrIllegalDenom("illegal uni denomnation")
	}
	return nil
}

// CleanReservePool remove non-pool coins
func CleanReservePool(reservePool sdk.Coins, uniDenom string) (sdk.Coins, sdk.Error) {
	if reservePool == nil {
		return sdk.Coins{}, nil
	}
	tokenDenom, err := GetTokenDenom(uniDenom)
	if err != nil {
		return nil, err
	}

	return sdk.NewCoins(sdk.NewCoin(sdk.IrisAtto, reservePool.AmountOf(sdk.IrisAtto)), sdk.NewCoin(tokenDenom, reservePool.AmountOf(tokenDenom)), sdk.NewCoin(uniDenom, reservePool.AmountOf(uniDenom))), nil
}
