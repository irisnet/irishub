package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// nolint
const (
	// module name
	ModuleName = "guardian"

	// StoreKey is the default store key for guardian
	StoreKey = ModuleName

	// RouterKey is the message route for guardian
	RouterKey = ModuleName

	// QuerierRoute is the querier route for the guardian store.
	QuerierRoute = StoreKey

	// Query endpoints supported by the guardian querier
	QuerySupers = "supers"
)

var (
	SuperKey = []byte{0x00} // super key
)

// GetSuperKey returns super key bytes
func GetSuperKey(addr sdk.AccAddress) []byte {
	return append(SuperKey, addr.Bytes()...)
}

// GetSupersSubspaceKey returns the key for getting all supers from the store
func GetSupersSubspaceKey() []byte {
	return SuperKey
}
