package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the rand module
	ModuleName = "rand"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the asset module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the asset module
	RouterKey = ModuleName
)

var (
	KeyDelimiter             = []byte(":")                 // key delimiter
	PrefixRandom             = []byte("rands:")            // key prefix for the random number
	PrefixRandomRequestQueue = []byte("randRequestQueue:") // key prefix for the random number request queue
)

// KeyRandom returns the key for a random number by the specified request id
func KeyRandom(reqID []byte) []byte {
	return append(PrefixRandom, reqID...)
}

// KeyRandomRequestQueue returns the key for the random number request queue by the given height and request id
func KeyRandomRequestQueue(height int64, reqID []byte) []byte {
	prefix := append(PrefixRandomRequestQueue, sdk.Uint64ToBigEndian(uint64(height))...)
	return append(append(prefix, KeyDelimiter...), reqID...)
}

// KeyRandomRequestQueueSubspace returns the key prefix for iterating through all requests at the specified height
func KeyRandomRequestQueueSubspace(height int64) []byte {
	return append(append(PrefixRandomRequestQueue, sdk.Uint64ToBigEndian(uint64(height))...), KeyDelimiter...)
}
