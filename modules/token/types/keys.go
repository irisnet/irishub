package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the token module
	ModuleName = "token"

	// StoreKey is the string store representation
	StoreKey string = ModuleName

	// QuerierRoute is the querier route for the token module
	QuerierRoute string = ModuleName

	// RouterKey is the msg router key for the token module
	RouterKey string = ModuleName

	// DefaultParamspace is the default name for parameter store
	DefaultParamspace = ModuleName
)

var (
	// PrefixTokenForSymbol defines a symbol prefix for the token
	PrefixTokenForSymbol = []byte{0x01}
	// PrefixTokenForMinUint defines the min unit prefix for the token
	PrefixTokenForMinUint = []byte{0x02}
	// PrefixTokens defines a prefix for the tokens
	PrefixTokens = []byte{0x03}
	// PeffixBurnTokenAmt defines a prefix for the amount of token burnt
	PeffixBurnTokenAmt = []byte{0x04}
)

// KeySymbol returns the key of the token with the specified symbol
func KeySymbol(symbol string) []byte {
	return append(PrefixTokenForSymbol, []byte(symbol)...)
}

// KeyMinUint returns the key of the token with the specified min unit
func KeyMinUint(minUnit string) []byte {
	return append(PrefixTokenForMinUint, []byte(minUnit)...)
}

// KeyTokens returns the key of the specified owner and symbol. Intended for querying all tokens of an owner
func KeyTokens(owner sdk.AccAddress, symbol string) []byte {
	return append(append(PrefixTokens, owner.Bytes()...), []byte(symbol)...)
}

// KeyBurnTokenAmt returns the key of the specified min unit.
func KeyBurnTokenAmt(minUint string) []byte {
	return append(PeffixBurnTokenAmt, []byte(minUint)...)
}
