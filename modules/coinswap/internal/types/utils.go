package types

import (
	"fmt"
	"strings"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
func GetUniDenomFromDenoms(denom1, denom2 string) (string, error) {
	if denom1 == denom2 {
		return "", ErrEqualDenom
	}
	if denom1 != StandardDenom && denom2 != StandardDenom {
		return "", sdkerrors.Wrap(ErrNotContainStandardDenom, fmt.Sprintf("standard denom: %s,denom1: %s,denom2: %s", StandardDenom, denom1, denom2))
	}
	if denom1 == StandardDenom {
		return fmt.Sprintf(FormatUniDenom, denom2), nil
	}
	return fmt.Sprintf(FormatUniDenom, denom1), nil
}

// GetUniDenomFromDenom returns the uni denom for the provided denomination.
func GetUniDenomFromDenom(denom string) (string, error) {
	if denom == StandardDenom {
		return "", ErrMustStandardDenom
	}
	return fmt.Sprintf(FormatUniDenom, denom), nil
}

// GetCoinDenomFromUniDenom returns the token denom by uni denom
func GetCoinDenomFromUniDenom(uniDenom string) (string, error) {
	if err := CheckUniDenom(uniDenom); err != nil {
		return "", err
	}
	return strings.TrimPrefix(uniDenom, FormatUniABSPrefix), nil
}

// CheckUniDenom returns nil if the uni denom is valid
func CheckUniDenom(uniDenom string) error {
	if !strings.HasPrefix(uniDenom, FormatUniABSPrefix) {
		return sdkerrors.Wrap(ErrInvalidDenom, uniDenom)
	}
	return nil
}
