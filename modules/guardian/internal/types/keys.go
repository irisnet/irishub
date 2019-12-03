package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	profilerKey = []byte{0x00}
	trusteeKey  = []byte{0x01}
)

// nolint
const (
	// module name
	ModuleName = "guardian"

	// default paramspace for params keeper
	DefaultParamspace = ModuleName

	// StoreKey is the default store key for mint
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute is the querier route for the minting store.
	QuerierRoute = StoreKey

	// Query endpoints supported by the guardian querier
	QueryProfilers = "profilers"
	QueryTrustees  = "trustees"
)

func GetProfilerKey(addr sdk.AccAddress) []byte {
	return append(profilerKey, addr.Bytes()...)
}

func GetTrusteeKey(addr sdk.AccAddress) []byte {
	return append(trusteeKey, addr.Bytes()...)
}

// Key for getting all profilers from the store
func GetProfilersSubspaceKey() []byte {
	return profilerKey
}

// Key for getting all profilers from the store
func GetTrusteesSubspaceKey() []byte {
	return trusteeKey
}
