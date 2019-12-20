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
	ProfilerKey = []byte{0x00}
	TrusteeKey  = []byte{0x01}
)

func GetProfilerKey(addr sdk.AccAddress) []byte {
	return append(ProfilerKey, addr.Bytes()...)
}

func GetTrusteeKey(addr sdk.AccAddress) []byte {
	return append(TrusteeKey, addr.Bytes()...)
}

// Key for getting all profilers from the store
func GetProfilersSubspaceKey() []byte {
	return ProfilerKey
}

// Key for getting all profilers from the store
func GetTrusteesSubspaceKey() []byte {
	return TrusteeKey
}
