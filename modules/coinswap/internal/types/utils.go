package types

import (
	"fmt"
	"strings"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetReservePoolAddr returns the poor address for the provided provided liquidity denomination.
func GetReservePoolAddr(uniDenom string) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(uniDenom)))
}

// GetTokenPairByDenom return the token pair for the provided denominations
func GetTokenPairByDenom(inputDenom, outputDenom string) string {
	return fmt.Sprintf("%s-%s", outputDenom, inputDenom)
}

// GetUniDenomFromDenoms returns the uni denom for the provided denominations.
func GetUniDenomFromDenoms(denom1, denom2 string) (string, sdk.Error) {
	if denom1 == denom2 {
		return "", ErrEqualDenom("denomnations for forming uni id are equal")
	}

	if denom1 != StandardDenom && denom2 != StandardDenom {
		return "", ErrIllegalDenom(fmt.Sprintf("illegal denomnations for forming uni id, must have one native denom: %s", StandardDenom))
	}

	if denom1 == StandardDenom {
		return fmt.Sprintf(FormatUniDenom, denom2), nil
	}

	return fmt.Sprintf(FormatUniDenom, denom1), nil
}

// GetUniDenomFromDenom returns the uni denom for the provided denomination.
func GetUniDenomFromDenom(denom string) (string, sdk.Error) {
	if denom == StandardDenom {
		return "", ErrIllegalDenom("illegal denomnation for forming uni denom, must not be NativeDenom")
	}
	return fmt.Sprintf(FormatUniDenom, denom), nil
}

// GetCoinDenomFromUniDenom returns the token denom by uni denom
func GetCoinDenomFromUniDenom(uniDenom string) (string, sdk.Error) {
	if err := CheckUniDenom(uniDenom); err != nil {
		return "", err
	}
	return strings.TrimPrefix(uniDenom, FormatUniABSPrefix), nil
}

// CheckUniDenom returns nil if the uni denom is valid
func CheckUniDenom(uniDenom string) sdk.Error {
	if !strings.HasPrefix(uniDenom, FormatUniABSPrefix) {
		return ErrIllegalDenom(fmt.Sprintf("illegal liquidity denomnation: %s", uniDenom))
	}
	return nil
}
