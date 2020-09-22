package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the Token module
	ModuleName = "token"

	// StoreKey is the string store representation
	StoreKey string = ModuleName

	// QuerierRoute is the querier route for the Asset module
	QuerierRoute string = ModuleName

	// RouterKey is the msg router key for the Asset module
	RouterKey string = ModuleName

	// DefaultParamspace default name for parameter store
	DefaultParamspace = ModuleName
)

var (
	PrefixTokenForSymbol  = []byte{0x1} // symbol prefix for the token
	PrefixTokenForMinUint = []byte{0x2} // min_unit prefix for the token
	PrefixTokens          = []byte{0x3} // prefix for the tokens
)

// KeySymbol returns the key of the token with the specified symbol
func KeySymbol(symbol string) []byte {
	symbol = strings.ToLower(strings.TrimSpace(symbol))
	return append(PrefixTokenForSymbol, []byte(symbol)...)
}

// KeyMinUint returns the key of the token with the specified min_unit
func KeyMinUint(minUnit string) []byte {
	minUnit = strings.ToLower(strings.TrimSpace(minUnit))
	return append(PrefixTokenForMinUint, []byte(minUnit)...)
}

// KeyTokens returns the key of the specified owner and symbol. Intended for querying all tokens of an owner
func KeyTokens(owner sdk.AccAddress, symbol string) []byte {
	symbol = strings.ToLower(strings.TrimSpace(symbol))
	return append(append(PrefixTokens, owner.Bytes()...), []byte(symbol)...)
}
