package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/asset/types"
)

const (
	// SubModuleName is the name of the module
	SubModuleName = "token"

	// StoreKey is the store key string for gov
	StoreKey = SubModuleName

	// RouterKey is the message route for gov
	RouterKey = SubModuleName

	// QuerierRoute is the querier route for gov
	QuerierRoute = SubModuleName

	// DefaultParamspace default name for parameter store
	DefaultParamspace = SubModuleName
)

func keyPath(prefix int) []byte {
	return types.KeyPrefixBytes(SubModuleName, prefix)
}

// KeyTokenPrefix return the base prefix of the token module
func KeyTokenPrefix() []byte {
	return keyPath(types.KeyTokenPrefix)
}

// KeyToken returns the key of the specified token
func KeyToken(symbol string) []byte {
	return append(
		keyPath(types.KeyTokenPrefix),
		[]byte(symbol)...,
	)
}

// KeyTokens returns the key of the specified owner . Intended for querying all tokens of an owner
func KeyTokens(owner sdk.AccAddress, symbol string) []byte {
	return append(
		keyPath(types.KeyTokenSymbolPrefix),
		[]byte(fmt.Sprintf("%s/%s", owner, symbol))...,
	)
}

// KeyMinUnit returns the key of the specified minUnit
func KeyMinUnit(minUnit string) []byte {
	return append(
		keyPath(types.KeyTokenMinUnitPrefix),
		[]byte(minUnit)...,
	)
}
