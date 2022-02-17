package v150

import (
	"fmt"
	"strings"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// FormatUniABSPrefix defines the prefix of liquidity token
	FormatUniABSPrefix = "swap"
	// FormatUniDenom defines the name of liquidity token
	FormatUniDenom = "swap%s"
)

// GetReservePoolAddr returns the pool address for the provided liquidity denomination.
func GetReservePoolAddr(uniDenom string) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(uniDenom)))
}

// GetUniDenomFromDenom returns the uni denom for the provided denomination.
func GetUniDenomFromDenom(denom string) string {
	return fmt.Sprintf(FormatUniDenom, denom)
}

// GetCoinDenomFromUniDenom returns the token denom by uni denom
func GetCoinDenomFromUniDenom(uniDenom string) (string, error) {
	return strings.TrimPrefix(uniDenom, FormatUniABSPrefix), nil
}
