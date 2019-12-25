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
	QueryProfilers = "profilers"
	QueryTrustees  = "trustees"
)

var (
	ProfilerKey = []byte{0x00} // profiler key
	TrusteeKey  = []byte{0x01} // trustee key
)

// GetProfilerKey returns profiler key bytes
func GetProfilerKey(addr sdk.AccAddress) []byte {
	return append(ProfilerKey, addr.Bytes()...)
}

// GetTrusteeKey returns trustee key bytes
func GetTrusteeKey(addr sdk.AccAddress) []byte {
	return append(TrusteeKey, addr.Bytes()...)
}

// GetProfilersSubspaceKey returns the key for getting all profilers from the store
func GetProfilersSubspaceKey() []byte {
	return ProfilerKey
}

// GetProfilersSubspaceKey returns the key for getting all profilers from the store
func GetTrusteesSubspaceKey() []byte {
	return TrusteeKey
}
