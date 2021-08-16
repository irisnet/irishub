package types

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetReservePoolAddr returns the pool address for the provided liquidity denomination.
func GetReservePoolAddr(lptDenom string) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(lptDenom)))
}

// GetTokenPairByDenom returns the token pair for the provided denominations
func GetTokenPairByDenom(inputDenom, outputDenom string) string {
	return fmt.Sprintf("%s-%s", outputDenom, inputDenom)
}

// GetPoolId returns the pool coin denom by specified sequence.
func GetPoolId(counterpartyDenom string) string {
	return fmt.Sprintf("pool-%s", counterpartyDenom)
}

// GetLptDenom returns the pool coin denom by specified sequence.
func GetLptDenom(sequence uint64) string {
	return fmt.Sprintf(LptTokenFormat, sequence)
}

func ParseLptDenom(lptDenom string) (uint64, error) {
	result := strings.Split(lptDenom, "-")
	if len(result) != 2 {
		return 0, fmt.Errorf("invalid lpt denom: %s", lptDenom)
	}
	return strconv.ParseUint(result[1], 10, 64)
}
